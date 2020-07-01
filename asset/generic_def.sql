CREATE TABLE "generic" (
    "薬価基準収載医薬品コード"              character varying(255) PRIMARY KEY,
    "成分名"                                character varying(255) ,
    "品名"                                  character varying(255) ,
    "後発情報"                              character varying(4)   ,
    "収載年月日"                            character varying(8)   ,
    "経過措置による使用期限"                character varying(20)  ,
    "備考"                                  character varying(255)
);

CREATE INDEX "index_generic_gen" ON "generic"(
    "後発情報"
);
CREATE INDEX "index_generic_gname" ON "generic"(
    "成分名"
);
