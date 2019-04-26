#encoding:utf8

import psycopg2 as pg
import csv
import os, re
import sys
import os.path 
import codecs

ifile = "/asset/drug_code_list.txt"

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
        params = [line.split("\t") for line in f]

    return params

def import_drugs():
    params = sql_params()

    sql = """
INSERT INTO yj(
    yjcode,
    status_no,
    drug_code    
    yj_comment,
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


def main():
    import_drugs()
    

if __name__ == '__main__':
    main()
