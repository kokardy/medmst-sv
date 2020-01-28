#encoding:utf8

import psycopg2 as pg
import csv
import os, re
import sys
import os.path 
import codecs

ifile = "/asset/drug_code_list.txt"
backupfile = "/backup/backup.sql"

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

def sql_params():
    with codecs.open(ifile, mode="r", encoding="utf8") as f:
        params = [line.strip("\n").split("\t") for line in f]

    return params

def import_drugs():
    params = sql_params()

    for param in params:
        if len(param) != 3:
            print(param)
            raise Exception(u"line in drug_code_list.txt must have 3 fields")

    sql = """
INSERT INTO yj(
    yjcode,
    status_no,
    drug_code,
    yj_comment
) VALUES (
    %s,
    3,
    %s,
    %s
)
"""

    con = connection()
    cur = con.cursor()
    cur.executemany(sql, params)

    con.commit()


def restore_drugs():
    con = connection()
    cur = con.cursor()

    sql = """
    DROP TABLE yj;
    DROP TABLE hot;
    DROP TABLE custom_yj;
    """

    with codecs.open(backupfile, "r", encoding="utf8") as f:
        sql = sql +  "\n".join(f)

    try:
        cur.execute(sql)
    except:
        con.rollback()

    con.commit()

def main():
    if os.path.exists(backupfile):
        restore_drugs()
        print("restore drugs: {}".format(backupfile))
    else:
        print("restore file not found: {}".format(backupfile))

if __name__ == '__main__':
    main()
