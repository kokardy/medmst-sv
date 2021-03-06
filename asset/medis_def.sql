CREATE TABLE "medis" (
    "基準番号（ＨＯＴコード）"              character varying(255) PRIMARY KEY,
    "処方用番号（ＨＯＴ７）"                character varying(255) ,
    "会社識別用番号"                        character varying(255) ,
    "調剤用番号"                            character varying(255) ,
    "物流用番号"                            character varying(255) ,
    "ＪＡＮコード"                          character varying(255) ,
    "薬価基準収載医薬品コード"              character varying(255) ,
    "個別医薬品コード"                      character varying(255) ,
    "レセプト電算処理システムコード（１）"  character varying(255) ,
    "レセプト電算処理システムコード（２）"  character varying(255) ,
    "告示名称"                              character varying(255) ,
    "販売名"                                character varying(255) ,
    "レセプト電算処理システム医薬品名"      character varying(255) ,
    "規格単位"                              character varying(255) ,
    "包装形態"                              character varying(255) ,
    "包装単位数"                            character varying(255) ,
    "包装単位単位"                          character varying(255) ,
    "包装総量数"                            character varying(255) ,
    "包装総量単位"                          character varying(255) ,
    "区分"                                  character varying(255) ,
    "製造会社"                              character varying(255) ,
    "販売会社"                              character varying(255) ,
    "更新区分"                              character varying(255) ,
    "更新年月日"                            character varying(255)
);

CREATE INDEX "index_medis_jan" ON "medis"(
    "ＪＡＮコード"
);
CREATE INDEX "index_medis_y" ON "medis"(
    "薬価基準収載医薬品コード"
);
CREATE INDEX "index_medis_yj" ON "medis"(
    "個別医薬品コード"
);
CREATE INDEX "index_medis_hot" ON "medis"(
    "基準番号（ＨＯＴコード）"
);
CREATE INDEX "index_medis_hot11" ON "medis"(
    SUBSTR("基準番号（ＨＯＴコード）", 1, 11)
);
CREATE INDEX "index_medis_pname" ON "medis"(
    "製造会社"
);
CREATE INDEX "index_medis_cname" ON "medis"(
    "販売会社"
);
CREATE INDEX "index_medis_update" ON "medis"(
    "更新区分"
);