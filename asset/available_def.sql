CREATE TABLE "status" (
    "no" integer primary key,
    "name" character varying(10)
);

INSERT INTO "status" VALUES(
    0,
    '不採用'
);
INSERT INTO "status" VALUES(
    1,
    '採用薬'
);
INSERT INTO "status" VALUES(
    2,
    '院外専用'
);
INSERT INTO "status" VALUES(
    4,
    '院内専用'
);

CREATE TABLE "yj" (
    "yjcode" character varying(12) primary key,
    "status_no" integer REFERENCES status (no)
);

CREATE TABLE "hot" (
    "hot11" character varying(11) primary key,
    "status_no" integer REFERENCES status (no)
);
