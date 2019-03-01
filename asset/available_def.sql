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
    '院内専用'
);
INSERT INTO "status" VALUES(
    2,
    '院外専用'
);
INSERT INTO "status" VALUES(
    3,
    '採用'
);

CREATE TABLE "yj" (
    "yjcode" character varying(12) primary key,
    "status_no" integer REFERENCES status (no)
);

CREATE TABLE "hot" (
    "HOT11" character varying(11) primary key,
    "status_no" integer REFERENCES status (no)
);

CREATE FUNCTION "resolve_status" ("status_no1" integer, "status_no2" integer)
RETURNS varchar
    AS $$
select "name" FROM "status" WHERE "no" = CASE 
            WHEN $1 IS NULL THEN $2 
            WHEN $2 IS NULL THEN $1
        ELSE $1 | $2
        END
       $$ LANGUAGE SQL;


CREATE VIEW "available_view" as
	SELECT DISTINCT
		"基準番号（ＨＯＴコード）",
		"処方用番号（ＨＯＴ７）",
		"ＪＡＮコード",
		"薬価基準収載医薬品コード",
		"個別医薬品コード",
		"告示名称",
		"販売名",
		"規格単位",
		"包装形態",
		"包装単位数",
		"包装単位単位",
		"包装総量数",
		"包装総量単位",
		"製造会社",
		"販売会社",

        "漢字名称",
        "単位_漢字名称",
        "新_金額",

		"HOT11",
        "yj"."status_no" as "yj_status",
        "hot"."status_no" as "hot_status",

        CASE 
            WHEN "yj"."status_no" IS NULL THEN "hot"."status_no" 
            WHEN "hot"."status_no" IS NULL THEN "yj"."status_no" 
        ELSE "yj"."status_no" | "hot"."status_no" 
        END AS "status_flag",

        resolve_status("yj"."status_no", "hot"."status_no") as "採用状態"
                

	FROM "medis"
	LEFT JOIN "y"
		ON "y"."薬価基準コード" = "medis"."薬価基準収載医薬品コード" 
    LEFT JOIN "yj"
        ON "yj"."yjcode" = "medis"."薬価基準収載医薬品コード"
    LEFT JOIN "hot"
        ON "hot"."HOT11" = substr("medis"."基準番号（ＨＯＴコード）", 1, 11);
