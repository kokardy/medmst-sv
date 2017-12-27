#encoding:utf8

import psycopg2 as pg
import csv
import os, re
import os.path 

if "MEDMST_SAVE" in os.environ:
    SAVE_DIR = os.environ["MEDMST_SAVE"]
else:
    SAVE_DIR = "/bootstrap/save"

PARAM = dict(
            host = os.environ["MEDMST_PG_HOST"],
            port = os.environ["MEDMST_PG_PORT"],
            user = os.environ["MEDMST_PG_USER"],
            password = os.environ["MEDMST_PG_PASSWORD"],
)

def connection():
    return _getConnection(PARAM)

def _connection(param):
    return pg.connect(**param)

def get_files(save_dir=SAVE_DIR):
    hot = "hot"
    y = "y"
    result = dict(hot=[], y=[])
    for path, dirs, filenames in os.walk(save_dir):
        for filename in filenames:
            fullpath = os.path.join(path, filename)
            if os.path.basename(path) == hot:
                r = result[hot]
                if re.match("^MEDIS\d{8}.TXT$", filename):
                    r.append(fullpath)
                elif re.match("^\d{8}.txt", filename):
                    r.append(fullpath)
                elif re.match("^h\d{8}del.txt", filename):
                    r.append(fullpath)
            if os.path.basename(path) == y:
                if filename == "y.csv":
                    result[y] = [fullpath]

    return result

def create(con):
    filepath = os.path.join(SAVE_DIR, "medis_def.txt")
    _sql_from_file(filepath)
    cur = con.cursor()
    cur.execute(sql)
    filepath = os.path.join(SAVE_DIR, "y_def.txt")
    _sql_from_file(filepath)
    cur = con.cursor()
    cur.execute(sql)

def _sql_from_file(filepath):
    with codecs.open(filepath, "r", "utf8") as f:
        lines = [line for line in f]
        sql = "\n".join(lines)
    return sql

def insert(con, infiles):
    infiles = get_files()
    
    insert_list = [(os.path.join(SAVE_DIR, "{0}_insert.txt").format(table),
        infiles[table], skip)
            for (table, skip) in [("medis", True), ("y", False)]]
    for (sql_file, insert_data, skip) in insert_list:
        _insert(con, sql_file, insert_data, skip)

def _insert(con, sql_file, insert_files, line1skip):
    sql = _sql_from_file(sql_file)
    for insert_file in insert_files:
        with codecs.open(insert_file, "r", "utf-8") as f:
            r = csv.reader(f)
            if line1skip:
                r.next()
            cur = con.cursor()
            cur.executemany(sql, r)


def main():
    infiles = get_files(SAVE_DIR)
    print(infiles)
    connection()

if __name__ == '__main__':
    main()
