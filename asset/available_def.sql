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

CREATE TABLE "hot" (
    "HOT11" character varying(11) primary key,
    "status_no" integer REFERENCES status (no),
    "hot_comment" character varying(255)
);

CREATE TABLE "yj" (
    "yjcode" character varying(12) primary key,
    "status_no" integer REFERENCES status (no),
    "yj_comment" character varying(255) default '',
    "drug_code" character varying(10) default ''
);

CREATE TABLE "custom_yj"(
    "HOT9" character varying(9) primary key,
    "yjcode" character varying(12) 
);

CREATE FUNCTION "resolve_status" ("status_no1" integer, "status_no2" integer)
RETURNS varchar
    AS $$
select "name" FROM "status" WHERE "no" = COALESCE($1, 0) |  COALESCE($2, 0)
       $$ LANGUAGE SQL;


CREATE VIEW "available_view" as
	SELECT DISTINCT
		--"基準番号（ＨＯＴコード）",
		--"処方用番号（ＨＯＴ７）",
		COALESCE("ＪＡＮコード", '') as "JAN",
		"薬価基準収載医薬品コード",
		COALESCE(medis."個別医薬品コード", custom_yj."yjcode") as "個別医薬品コード",
		"告示名称",
		"販売名",
		"規格単位",
		"包装形態",
		"包装単位数",
		"包装単位単位",
		"製造会社",
		"販売会社",

        COALESCE(y."単位_漢字名称", '--') as "単位_漢字名称", 
        COALESCE(y."新_金額", -1) as "新_金額",


        SUBSTR("基準番号（ＨＯＴコード）", 1, 11) as "HOT11",
        COALESCE(yj.status_no, 0) as yj_status,
        COALESCE(yj_comment, '') as yj_comment,
        COALESCE(yj.drug_code, '') as drug_code,
        COALESCE(hot.status_no, 0) as hot_status,
        COALESCE(hot_comment, '') as hot_comment,
        COALESCE(yj.status_no, 0) | COALESCE(hot.status_no, 0) AS status_flag,
        resolve_status(yj."status_no", hot."status_no") as "採用状態",

        COALESCE(custom_yj."yjcode", '') as custom_yj
        

	FROM medis
	LEFT JOIN y
		ON y."薬価基準コード" = medis."薬価基準収載医薬品コード" 
    LEFT JOIN yj
        ON yj."yjcode" = medis."個別医薬品コード"
    LEFT JOIN hot
        ON hot."HOT11" = substr(medis."基準番号（ＨＯＴコード）", 1, 11)
    LEFT JOIN custom_yj
        ON SUBSTR("hot"."HOT11", 1 , 9) = custom_yj."HOT9"
    WHERE medis."更新区分" <> '4'; -- '4'は削除フラグ
