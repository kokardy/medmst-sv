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

CREATE INDEX index_customyj_yj ON "custom_yj"(
    "yjcode" 
);

CREATE FUNCTION "resolve_status" ("status_no1" integer, "status_no2" integer)
RETURNS varchar
    AS $$
select "name" FROM "status" WHERE "no" = COALESCE($1, 0) |  COALESCE($2, 0)
       $$ LANGUAGE SQL;

CREATE VIEW "available_view" as
    SELECT 
    SUBSTR(m1."基準番号（ＨＯＴコード）", 1, 11)::varchar as "HOT11",
    -- m1."ＪＡＮコード" AS "JAN",
    COALESCE(
        (SELECT "ＪＡＮコード" FROM medis 
            WHERE SUBSTR("基準番号（ＨＯＴコード）", 1, 11) = SUBSTR(m1."基準番号（ＨＯＴコード）", 1, 11)
            LIMIT 1),
    '') as "JAN",
    m1."薬価基準収載医薬品コード"                  ,
    COALESCE(m1."個別医薬品コード", 
            custom_yj.yjcode) as "個別医薬品コード",
   -- m1."告示名称" ,
   -- (SELECT "告示名称" FROM medis
   --     WHERE SUBSTR("基準番号（ＨＯＴコード）", 1, 11) = "HOT11" LIMIT 1
   --     ) AS "告示名称",
    m1."販売名",
   -- (SELECT "販売名" FROM medis
   --     WHERE SUBSTR("基準番号（ＨＯＴコード）", 1, 11) = "HOT11" LIMIT 1
   --     ) AS "販売名",
    COALESCE(generic."成分名", '') as "成分名",
    COALESCE("後発情報", '') as "後発情報",
    "規格単位"                                     ,
    "包装形態"                                     ,
    "包装単位数"                                   ,
    "包装単位単位"                                 ,
    "製造会社"                                     ,
    "販売会社"                                     ,
    COALESCE(y1."単位_漢字名称", '--') as "単位_漢字名称",
    COALESCE(y1."新_金額", 0) as "新_金額",

    COALESCE(yj.status_no, 0) as yj_status,
    COALESCE(yj_comment, '') as yj_comment,
    COALESCE(yj.drug_code, '') as drug_code,
    COALESCE(hot.status_no, 0) as hot_status,
    COALESCE(hot_comment, '') as hot_comment,
    COALESCE(yj.status_no, 0) | COALESCE(hot.status_no, 0) AS status_flag,
    resolve_status(yj."status_no", hot."status_no") as "採用状態",
    COALESCE(custom_yj."yjcode", '') as custom_yj

	FROM medis m1
    LEFT JOIN yj
        ON m1."個別医薬品コード" = yj."yjcode"
    LEFT JOIN hot
        ON SUBSTR(m1."基準番号（ＨＯＴコード）", 1, 11) = hot."HOT11"
    LEFT JOIN custom_yj
        ON SUBSTR(hot."HOT11", 1 , 9) = custom_yj."HOT9"
	LEFT JOIN (SELECT DISTINCT "薬価基準コード", "単位_漢字名称", "新_金額" 
                    from y 
                    WHERE "薬価基準コード" IS NOT NULl
                        AND "薬価基準コード" <> '') as y1
		ON m1."薬価基準収載医薬品コード" = y1."薬価基準コード"
    LEFT JOIN generic
        ON m1."薬価基準収載医薬品コード" = generic."薬価基準収載医薬品コード"
    WHERE m1."更新区分" <> '2' AND m1."更新区分" <> '4' -- '2':中止,'4':削除
;