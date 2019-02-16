package main

type Medis struct {
	Name  string
	Total string
	Unit  string
}

type Y struct {
	Change         int     `db:"変更区分" json:"変更区分"`
	Master         string  `db:"マスター種別" json:"マスター種別"`
	MedCode        int     `db:"医薬品コード" json:"医薬品コード"`
	KanjiDigit     int     `db:"漢字有効桁数" json:"漢字有効桁数"`
	Name           string  `db:"漢字名称" json:"漢字名称"`
	KanaDigit      int     `db:"カナ有効桁数" json:"カナ有効桁数"`
	KanaName       string  `db:"カナ名称" json:"カナ名称" `
	UnitCode       int     `db:"単位_コード" json:"単位_コード"`
	UnitDigit      int     `db:"単位_漢字有効桁数" json:"単位_漢字有効桁数"`
	UnitName       string  `db:"単位_漢字名称" json:"単位_漢字名称"`
	PriceType      int     `db:"新_金額種別" json:"新_金額種別"`
	Price          float32 `db:"新_金額" json:"新_金額"`
	Yobi1          string  `db:"予備1" json:"予備1"`
	Law            int     `db:"麻薬・毒薬・覚醒剤原料・向精神薬" json:"麻薬・毒薬・覚醒剤原料・向精神薬"`
	Neuro          int     `db:"神経破壊剤" json:"神経破壊剤"`
	Bio            int     `db:"生物学的製剤" json:"生物学的製剤"`
	Generic        int     `db:"後発品" json:"後発品"`
	Yobi2          string  `db:"予備2" json:"予備2"`
	Dental         int     `db:"歯科特定薬剤" json:"歯科特定薬剤"`
	Radio          int     `db:"造影(補助)剤" json:"造影(補助)剤"`
	InjAmount      float32 `db:"注射容量" json:"注射容量"`
	ListType       int     `db:"収載方式等識別" json:"収載方式等識別"`
	AboutName      int     `db:"商品名等関連" json:"商品名等関連"`
	PriceTypeOld   int     `db:"旧_金額種別" json:"旧_金額種別"`
	PriceOld       float32 `db:"旧_金額" json:"旧_金額"`
	NameChange     int     `db:"漢字名称変更区分" json:"漢字名称変更区分"`
	KanaNameChange int     `db:"カナ名称変更区分" json:"カナ名称変更区分"`
	MedType        int     `db:"剤形" json:"剤形"`
	Yobi3          string  `db:"予備3" json:"予備3"`
	ChangeDate     string  `db:"変更年月日" json:"変更年月日"`
	FinishDate     string  `db:"廃止年月日" json:"廃止年月日"`
	PriceCode      string  `db:"薬価基準コード" json:"薬価基準コード"`
	Publish        int     `db:"公表順序番号" json:"公表順序番号"`
	ExpDate        string  `db:"経過措置年月日" json:"経過措置年月日"`
	BaseName       string  `db:"基本漢字名称" json:"基本漢字名称"`
}
