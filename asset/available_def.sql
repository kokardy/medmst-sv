CREATE TABLE "status" (
    "id" primary key,
    "name" character varying(10)
);
INSERT INTO "status" VALUES(
    "id" = 0,
    "name" = "不採用"
);
INSERT INTO "status" VALUES(
    "id" = 1,
    "name" = "採用薬"
);
INSERT INTO "status" VALUES(
    "id" = 2,
    "name" = "院外専用"
);
INSERT INTO "status" VALUES(
    "id" = 3,
    "name" = "院内専用"
);

CREATE TABLE "yj" (
    "yjcode" character varying(12) primary key,
    "status_id" integer REFERENCIES status(id)
);

CREATE TABLE "hot" (
    "hot11" character varying(11) primary key,
    "status_id" integer REFERENCIES status(id)
);
