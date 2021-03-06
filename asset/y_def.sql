CREATE TABLE "y" (
    "変更区分"                                            integer                     ,
    "マスター種別"                                        character varying(1)        ,
    "医薬品コード"                                        integer                     ,
    "漢字有効桁数" 		                                  integer                     ,
    "漢字名称"                                            character varying(64)       ,
    "カナ有効桁数"                                        integer                     ,
    "カナ名称"                                            character varying(20)       ,
    "単位_コード"                                         integer                     ,
    "単位_漢字有効桁数"                                   integer                     ,
    "単位_漢字名称"                                       character varying(12)       ,
    "新_金額種別"                                         integer                     ,
    "新_金額"                                             float(32)                   ,
    "予備1"                                               character varying(10)                     ,
    "麻薬・毒薬・覚醒剤原料・向精神薬"                    integer                     ,
    "神経破壊剤"                                          integer                     ,
    "生物学的製剤"                                        integer                     ,
    "後発品"                                              integer                     ,
    "予備2"                                               character varying(10)                     ,
    "歯科特定薬剤"                                        integer                     ,
    "造影(補助)剤"                                        integer                     ,
    "注射容量"                                            float(32)                   ,
    "収載方式等識別"                                      integer                     ,
    "商品名等関連"                                        integer                     ,
    "旧_金額種別"                                         integer                     ,
    "旧_金額"                                             float(32)                   ,
    "漢字名称変更区分"                                    integer                     ,
    "カナ名称変更区分"                                    integer                     ,
    "剤形"                                                integer                     ,
    "予備3"                                               character varying(20)       ,
    "変更年月日"                                          character varying(8)        ,
    "廃止年月日"                                          character varying(8)        ,
    "薬価基準コード"                                      character varying(12)       ,
    "公表順序番号"                                        integer                     ,
    "経過措置年月日"                                      character varying(8)        ,
    "基本漢字名称"                                        character varying(200)      
);

CREATE INDEX index_y_name ON "y"(
    "漢字名称"
);

CREATE INDEX index_y_yj ON "y"(
    "薬価基準コード" 
);