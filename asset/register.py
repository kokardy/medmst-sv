#encoding:utf8

import psycopg2 as pg
import csv
import os, re
import os.path 
import 

SAVE_DIR = "/bootstrap/save"
CURRENT_DIR = "./"

def getConnection():
    return pg.connect(**param)

def get_files(path, hot="hot", y="y"):
    result = dict(hot=dict())
    for path, dirs, filenames in os.walk(SAVE_DIR):
        for filename in filenames:
            fullpath = os.path.join(path, filename)
            if os.path.basename(path) == hot:
                r = result[hot]
                if re.match("^MEDIS\d{8}.TXT$", filename):
                    r["main"] = fullpath
                elif re.match("^\d{8}.txt", filename):
                    r["extra"] = r.get("extra", []).append(fullpath)
                elif re.match("^h\d{8}del.txt", filename):
                    r["delete"] = fullpath
            if os.path.basename(path) == y:
                if filename == "y.csv":
                    result[y] = fullpath


def main():
    infiles = get_files(SAVE_DIR)
    print(infiles)


def create(con, tables=["hot", "y"]):
    if "hot" in tables:
        createHOT(con)
    if "y" in tables:
        createY(con)

def createHOT(con):
    filepath = os.path.join(CURRENT_DIR, "./medis_def.txt")
    f = codecs.open(filepath, "r", "utf8")
    lines = [line for line in f]
    sql = "\n".join(lines)
    cur = con.cursor()
    cur.execute(sql)

def createY(con):
    filepath = os.path.join(CURRENT_DIR, "./y_def.txt")
    f = codecs.open(filepath, "r", "utf8")
    lines = [line for line in f]
    sql = "\n".join(lines)
    cur = con.cursor()
    cur.execute(sql)



if __name__ == '__main__':
    main()
