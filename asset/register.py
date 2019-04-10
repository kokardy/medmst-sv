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
    result = dict(medis=[], y=[])
    for path, dirs, filenames in os.walk(save_dir):
        for filename in filenames:
            fullpath = os.path.join(path, filename)
            if os.path.basename(path) == hot:
                r = result[medis]
                if re.match("^MEDIS\d{8}.TXT$", filename):
                    r.append(fullpath)
                elif re.match("^\d{8}.txt", filename):
                    r.append(fullpath)
                elif re.match("^h\d{8}del.txt", filename):
                    r.append(fullpath)
            if os.path.basename(path) == y:
                if filename == "y.csv":
                    result[y] = [fullpath]

    print("result", result)
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
        os.path.join(ASSET_DIR, "available_def.sql"),
    ]

    for f in filepathlist:
        print("execute:", f)
        sql = _sql_from_file(f)
        cur = con.cursor()
        try:
            cur.execute(sql)
        except Exception as e:
            print(e)

#@param infiles={medis:[filelist], y:[filelist]}
#INSERT INTO TABLEを実行
def insert(con, infiles):
    infiles = get_files()
    
    insert_list = []
    for (table, skip) in [("medis", True), ("y", False)]:
        sql_template = os.path.join(ASSET_DIR, "{0}_insert.sql").format(table)
        insert_data =  infiles[table]
        insert_list.append([sql_template, insert_data, skip])
    for (sql_template, insert_data, skip) in insert_list:
        _insert(con, sql_template, insert_data, skip)

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
    ]
    cur = con.cursor()
    for sql in sqls:
        cur.execute(sql)
    con.commit()


#option -Cで実行する
#テーブル作成
def C():
    con = connection()
    create(con)
    con.commit()

#option -Iで実行する
#テーブルにINSERT
def I():
    con = connection()
    infiles = get_files()
    insert(con, infiles)
    con.commit()

#option -Dで実行する
#DELETE TABLE
def D():
    con = connection()
    delete(con)
    con.commit()


def main():
    infiles = get_files(SAVE_DIR)

    options = sys.argv[1].lstrip("-")

    
    if len(options) == 0:
        print("OPTION must be -[C][D][I]")
        print("C: create table")
        print("D: delete table data")
        print("I: insert data to table")
        return
    
    #OPTIONの分だけcreate delete insert関数登録
    exec_list = []
    if "C" in options:
        exec_list.append(C)
    if "D" in options:
        exec_list.append(D)
    if "I" in options:
        exec_list.append(I)

    #OPTIONにあったものだけ実行
    for func in exec_list:
        func()

if __name__ == '__main__':
    main()

#get_files
#@param save_dir 探すディレクトリ
#@return {medis:[filelist], y:[filelist]}
#ディレクトリからinsertするためのファイルのリスト探してくる
def get_files(save_dir=SAVE_DIR):
    hot = "hot"
    y = "y"
    medis = "medis"
    result = dict(medis=[], y=[])
    for path, dirs, filenames in os.walk(save_dir):
        for filename in filenames:
            fullpath = os.path.join(path, filename)
            if os.path.basename(path) == hot:
                r = result[medis]
                if re.match("^MEDIS\d{8}.TXT$", filename):
                    r.append(fullpath)
                elif re.match("^\d{8}.txt", filename):
                    r.append(fullpath)
                elif re.match("^h\d{8}del.txt", filename):
                    r.append(fullpath)
            if os.path.basename(path) == y:
                if filename == "y.csv":
                    result[y] = [fullpath]

    print("reuslt", result)
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
        os.path.join(ASSET_DIR, "available_def.sql"),
    ]

    for f in filepathlist:
        sql = _sql_from_file(f)
        cur = con.cursor()
        try:
            cur.execute(sql)
        except Exception as e:
            print(e)

#@param infiles={medis:[filelist], y:[filelist]}
#INSERT INTO TABLEを実行
def insert(con, infiles):
    infiles = get_files()
    
    insert_list = []
    for (table, skip) in [("medis", True), ("y", False)]:
        sql_template = os.path.join(ASSET_DIR, "{0}_insert.sql").format(table)
        insert_data =  infiles[table]
        insert_list.append([sql_template, insert_data, skip])
    for (sql_template, insert_data, skip) in insert_list:
        _insert(con, sql_template, insert_data, skip)

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
    ]
    cur = con.cursor()
    for sql in sqls:
        cur.execute(sql)
    con.commit()


#option -Cで実行する
#テーブル作成
def C():
    con = connection()
    print("create")
    create(con)
    con.commit()

#option -Iで実行する
#テーブルにINSERT
def I():
    con = connection()
    print("insert")
    infiles = get_files()
    insert(con, infiles)
    con.commit()

#option -Dで実行する
#DELETE TABLE
def D():
    con = connection()
    print("delete")
    delete(con)
    con.commit()


def main():
    infiles = get_files(SAVE_DIR)

    options = sys.argv[1].lstrip("-")

    
    if len(options) == 0:
        print("OPTION must be -[C][D][I]")
        print("C: create table")
        print("D: delete table data")
        print("I: insert data to table")
        return
    
    #OPTIONの分だけcreate delete insert関数登録
    exec_list = []
    if "C" in options:
        exec_list.append(C)
    if "D" in options:
        exec_list.append(D)
    if "I" in options:
        exec_list.append(I)

    #OPTIONにあったものだけ実行
    for func in exec_list:
        func()

if __name__ == '__main__':
    main()
