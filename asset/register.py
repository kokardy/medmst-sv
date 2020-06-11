#encoding:utf8

import psycopg2 as pg
import csv
import os, re
import sys
import os.path 
import codecs

ASSET_DIR = "/asset"

if "MEDMST_SAVE" in os.environ:
    SAVE_DIR = os.environ["MEDMST_SAVE"]
else:
    SAVE_DIR = "/bootstrap/save"

PARAM = dict(
            host = os.environ.get("PG_HOST", "postgres"),
            port = os.environ.get("PG_PORT", 5432),
            user = os.environ.get("PG_USER", "postgres"),
            password = os.environ.get("PG_PASSWORD", "allyourbaseisbelongtous"),
            database = os.environ.get("PG_DATABASE", "medmst"),
)

def connection():
    return _connection(PARAM)

def _connection(param):
    return pg.connect(**param)

#get_files
#@param save_dir 探すディレクトリ
#@return {medis:[filelist], y:[filelist]}
#ディレクトリからinsertするためのファイルのリスト探してくる
def get_files(save_dir=SAVE_DIR):
    hot = "hot"
    y = "y"
    medis = "medis"
    generic = "generic"
    result = dict(medis=[], y=[])
    for path, _dirs, filenames in os.walk(save_dir):
        for filename in filenames:
            fullpath = os.path.join(path, filename)
            if os.path.basename(path) == hot:
                r = result[medis]
                if re.match("^MEDIS\\d{8}.TXT$", filename):
                    r.append(fullpath)
                elif re.match("^\\d{8}.txt", filename):
                    r.append(fullpath)
                elif re.match("^h\\d{8}del.txt", filename):
                    r.append(fullpath)
            if os.path.basename(path) == y:
                if filename == "y.csv":
                    result[y] = [fullpath]
            if os.path.basename(path) == generic:
                if re.match("^tp\\d{8}-\\d{2}_\\d{2}.xlsx", filename):
                    result[generic] = [fullpath]

    print("result", result)

    ##medisの古いデータをinsert deleteしないようにする
    ##MDIS********.TXTの日付が最新なのでそれより古いものは省く
    r = result[medis]

    #filenameにMEDIS*********.TXTが入る
    for filename in r:
        if filename.startswith("MEDIS"):
            break

    oldest = filename.lstrip("MEDIS")

    #_ok 最終的にMEDIS始まりのファイルとそのファイルより新しいのだけにする
    def _ok(filename): 
        if filename.find("MEDIS") > -1:
            return True
        elif oldest < filename:
            return True
        else:
            return False

    r = filter(_ok, r)
    result[medis] = r

    return result

#想定使用方法
#*.sqlのファイルを読み込んでsql文字列を得る
def _sql_from_file(filepath):
    with codecs.open(filepath, "r", "utf8") as f:
        lines = [line for line in f]
        sql = "\n".join(lines)
    return sql


#CREATE TABLEを実行する
def create(con):
    filepathlist = [
        os.path.join(ASSET_DIR, "medis_def.sql"),
        os.path.join(ASSET_DIR, "y_def.sql"),
        os.path.join(ASSET_DIR, "generic_def.sql"),
        os.path.join(ASSET_DIR, "available_def.sql"),
    ]

    for f in filepathlist:
        print("execute:", f)
        sql = _sql_from_file(f)
        cur = con.cursor()
        try:
            cur.execute(sql)
            print("sql success")
        except Exception as e:
            print("error occured in creating tables.")
            print(e)

#@param infiles={medis:[filelist], y:[filelist]}
#INSERT INTO TABLEを実行
def insert(con, infiles):
    
    insert_list = []
    for (table, skip) in [("medis", True), ("y", False)]:
        sql_template = os.path.join(ASSET_DIR, "{0}_insert.sql").format(table)
        insert_data =  infiles[table]
        insert_list.append([sql_template, insert_data, skip])
    for (sql_template, insert_data, skip) in insert_list:
        _insert(con, sql_template, insert_data, skip)

    import pandas as pd
    import sqlalchemy as sa
    engine = sa.create_engine(
        f'postgresql://{PARAM["user"]}:{PARAM["password"]}@{PARAM["host"]}:{PARAM["port"]}/{PARAM["database"]}'
    )

    gfile = infiles["generic"]
    insert_data = pd.read_excel(gfile[0], dtype={4:str}) #日付が数値に解釈されるのを防止する
    index = "薬価基準収載医薬品コード"
    columns = [
        "成分名",
        "品名",
        "後発情報",
        "収載年月日",
        "経過措置による使用期限" ,
        "備考"]   
    insert_data = insert_data.set_index(index)
    insert_data.columns = columns

    print("insert: generic table")
    insert_data.to_sql("generic", engine, if_exists="append" )#if_exists=append テーブルがあったら追加


#insert(con, infles)から呼ぶ用
#sqlファイルと入力用ファイルを一行目をスキップするかどうかを与えて
#INSERT INTO TABLE 実行する
def _insert(con, sql_file, insert_files, line1skip):
    try:
        sql = _sql_from_file(sql_file)
    except IOError as e:
        print(e)
        return
    for insert_file in insert_files:
        print(f'insert: {insert_file} table')
        with open(insert_file, "r", encoding="cp932") as f:
            r = csv.reader(f)
            r = [line for line in r]
            if line1skip:
                r = r[1:]
            cur = con.cursor()
            try: 
                cur.executemany(sql, r)
                con.commit()
            except Exception as e:
                print("Error occured in executing SQL")
                print(cur.query)
                raise e

#テーブル全消し
def delete(con):
    sqls = [
        """DELETE FROM "medis";""",
        """DELETE FROM "y";""",
        """DELETE FROM "generic";""",
    ]
    cur = con.cursor()
    for sql in sqls:
        cur.execute(sql)
    con.commit()




def main():
    options = sys.argv[1].lstrip("-")
    
    if len(options) == 0:
        print("OPTION must be -[C][D][I]")
        print("C: create table")
        print("D: delete table data")
        print("I: insert data to table")
        return
    
    #OPTIONの分だけcreate delete insert関数登録
    con = connection()
    if "C" in options:
        C(con)
    if "D" in options:
        con = connection()
        con = D(con, commit=False)
    if "I" in options:
        con = I(con, commit=False)
    con.commit()


#option -Cで実行する
#テーブル作成
def C(con=connection(), commit=True):
    print("create")
    create(con)
    if commit:
        con.commit()
    return con

#option -Iで実行する
#テーブルにINSERT
def I(con=connection(), commit=True):
    print("insert")
    infiles = get_files()
    insert(con, infiles)
    con.commit()
    return con

#option -Dで実行する
#DELETE TABLE
def D(con=connection(), commit=True):
    print("delete")
    delete(con)
    con.commit()
    return con



if __name__ == '__main__':
    main()
