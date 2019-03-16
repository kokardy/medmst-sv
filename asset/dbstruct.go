package main

import "fmt"

//Medis is HOTコードのマスタ
type Medis struct {
	HOT            string `db:"基準番号（ＨＯＴコード）" json:"基準番号（ＨＯＴコード）"`
	HOT7           string `db:"処方用番号（ＨＯＴ７）" json:"処方用番号（ＨＯＴ７）"`
	Company        string `db:"会社識別用番号" json:"会社識別用番号"`
	ChozaiNo       string `db:"調剤用番号" json:"調剤用番号"`
	DeliNo         string `db:"物流用番号" json:"物流用番号"`
	JANCode        string `db:"ＪＡＮコード" json:"ＪＡＮコード"`
	PriceCode      string `db:"薬価基準収載医薬品コード" json:"薬価基準収載医薬品コード"`
	YJCode         string `db:"個別医薬品コード" json:"個別医薬品コード"`
	ReceCode1      string `db:"レセプト電算処理システムコード（１）" json:"レセプト電算処理システムコード（１）"`
	ReceCode2      string `db:"レセプト電算処理システムコード（２）" json:"レセプト電算処理システムコード（２）"`
	PublicName     string `db:"告示名称" json:"告示名称"`
	ConsName       string `db:"販売名" json:"販売名"`
	ReceName       string `db:"レセプト電算処理システム医薬品名" json:"レセプト電算処理システム医薬品名"`
	Unit           string `db:"規格単位" json:"規格単位"`
	CoverType      string `db:"包装形態" json:"包装形態"`
	CoverNum       string `db:"包装単位数" json:"包装単位数"`
	CoverUnit      string `db:"包装単位単位" json:"包装単位単位"`
	CoverTotal     string `db:"包装総量数" json:"包装総量数"`
	CoverTotalUnit string `db:"包装総量単位" json:"包装総量単位"`
	Kubun          string `db:"区分" json:"区分"`
	ManuCompany    string `db:"製造会社" json:"製造会社"`
	ConsCampany    string `db:"販売会社" json:"販売会社"`
	UpdateType     string `db:"更新区分" json:"更新区分"`
	UpdateDate     string `db:"更新年月日" json:"更新年月日"`
}

//Y is 薬価のマスタ
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

type AvailableView struct {
	//medis
	HOT            string `db:"基準番号（ＨＯＴコード）" json:"基準番号（ＨＯＴコード）"`
	HOT7           string `db:"処方用番号（ＨＯＴ７）" json:"処方用番号（ＨＯＴ７）"`
	JANCode        string `db:"ＪＡＮコード" json:"ＪＡＮコード"`
	PriceCode      string `db:"薬価基準収載医薬品コード" json:"薬価基準収載医薬品コード"`
	YJCode         string `db:"個別医薬品コード" json:"個別医薬品コード"`
	PublicName     string `db:"告示名称" json:"告示名称"`
	ConsName       string `db:"販売名" json:"販売名"`
	Unit           string `db:"規格単位" json:"規格単位"`
	CoverType      string `db:"包装形態" json:"包装形態"`
	CoverNum       string `db:"包装単位数" json:"包装単位数"`
	CoverUnit      string `db:"包装単位単位" json:"包装単位単位"`
	CoverTotal     string `db:"包装総量数" json:"包装総量数"`
	CoverTotalUnit string `db:"包装総量単位" json:"包装総量単位"`
	ManuCompany    string `db:"製造会社" json:"製造会社"`
	ConsCampany    string `db:"販売会社" json:"販売会社"`
	//y
	//Name     string  `db:"漢字名称" json:"漢字名称"`
	UnitName string  `db:"単位_漢字名称" json:"単位_漢字名称"`
	Price    float32 `db:"新_金額" json:"新_金額"`
	//custom
	HOT11      string `db:"HOT11" json:"HOT11"`
	YJStatus   int    `db:"yj_status" json:"yj_status"`
	HOTStatus  int    `db:"hot_status" json:"hot_status"`
	StatusFlag int    `db:"status_flag" json:"status_flag"`
	Status     string `db:"採用状態" json:"採用状態"`
}

type HOTStatus struct {
	HOT       string `db:"HOT11" json:"HOT"`
	Status_no int    `db:"status_no" json:"status_no"`
}

func (hs HOTStatus) String() string {
	return fmt.Sprintf("HOT:%s status:%d", hs.HOT, hs.Status_no)
}

type YJStatus struct {
	YJ        string `db:"yjcode"`
	Status_no int    `db:status_no`
}
