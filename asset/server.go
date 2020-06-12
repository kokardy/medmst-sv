package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func connectParam() (param string) {
	param = fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"),
		os.Getenv("PG_DATABASE"),
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWORD"),
	)
	return
}

func handleBarcode(c *gin.Context) {
	var message string
	barcode := c.Param("barcode")
	gs1 := GS1(barcode) //GS1として読んでみる
	var jan JAN
	//validate check digit
	if !gs1.CheckDigitOK() {
		jan = JAN(barcode) //ダメならJANとして読んでみる
		if !jan.CheckDigitOK() {
			message = "wrong barcode: checkdigit error"
			c.String(500, message)
		}
	} else { //GS1ならjanに変換する
		jan = gs1.ToJAN()
	}
	param := connectParam()
	db, err := sqlx.Connect("postgres", param)
	//DB connection check
	if err != nil {
		message = "an ERROR occured in database connecting"
		c.String(500, message)
		return
	}

	sql := `SELECT 
				"販売名",
				"包装総量数",
				"包装総量単位" 
			FROM "medis"
			WHERE "ＪＡＮコード" = $1;`

	medis := Medis{}
	err = db.Get(&medis, sql, string(jan))
	//sql execute error
	if err != nil {
		message = "an ERROR occured in executing SQL: %s"
		message = fmt.Sprintf(message, sql)
		c.String(500, message)
		return
	}

	c.JSON(200, medis)

}

func handleAvailable(c *gin.Context) {
	var result []AvailableView
	var message string
	var err error
	queryString := c.DefaultQuery("query", "")
	sql := `SELECT
				*
			FROM "available_view"
			WHERE
				"販売名" like '%' || $1 || '%' OR
				"成分名" like $1 || '%' OR
				"個別医薬品コード" like $1 || '%' OR
				"HOT11" like $1 || '%' OR
				"製造会社" like $1 || '%' OR
				"販売会社" like $1 || '%'
			ORDER BY "個別医薬品コード", "販売会社", "製造会社"
			;
		`
	param := connectParam()
	db, err := sqlx.Connect("postgres", param)
	//DB connection check
	if err != nil {
		message = "an ERROR occured in database connecting: %s\n"
		message = fmt.Sprintf(message, err)
		message += fmt.Sprintf("queryString: %s", queryString)
		fmt.Printf("err: %s\n", message)
		c.String(500, message)
		return
	}
	defer db.Close()
	//Query
	err = db.Select(&result, sql, queryString)
	if err != nil {
		message = "an ERROR occured in executing SQL: %s \n%s\n"
		message = fmt.Sprintf(message, sql, err)
		message += fmt.Sprintf("queryString: %s", queryString)
		fmt.Printf("err: %s\n", message)
		c.String(500, message)
		return
	}

	c.Header("Cache-Control", "no-store")
	c.JSON(200, result)
}
func handleY(c *gin.Context) {
	var result []Y
	var message string
	var err error
	queryString := c.DefaultQuery("query", "")
	sql := `SELECT
				*
			FROM "y"
			WHERE
				"漢字名称" like '%' || $1 || '%' OR
				"カナ名称" like '%' || $1 || '%' OR
				"基本漢字名称" like '%' || $1 || '%'; 
		`
	param := connectParam()
	db, err := sqlx.Connect("postgres", param)
	//DB connection check
	if err != nil {
		message = "an ERROR occured in database connecting: %s\n"
		message = fmt.Sprintf(message, err)
		message += fmt.Sprintf("queryString: %s", queryString)
		fmt.Printf("err: %s\n", message)
		c.String(500, message)
		return
	}
	defer db.Close()
	//Query
	err = db.Select(&result, sql, queryString)
	if err != nil {
		message = "an ERROR occured in executing SQL: %s \n%s\n"
		message = fmt.Sprintf(message, sql, err)
		message += fmt.Sprintf("queryString: %s", queryString)
		fmt.Printf("err: %s\n", message)
		c.String(500, message)
		return
	}

	c.JSON(200, result)
}
func handleMedis(c *gin.Context) {
	var result []Medis
	var message string
	var err error
	queryString := c.DefaultQuery("query", "")
	sql := `SELECT
				*
			FROM "medis"
			WHERE
				"告示名称" like '%' || $1 || '%' OR
				"販売名" like '%' || $1 || '%' OR
				"レセプト電算処理システム医薬品名" like '%' || $1 || '%' OR
				"製造会社"  like '%' || $1 || '%' OR
				"販売会社"  like '%' || $1 || '%' 
				; 
		`
	param := connectParam()
	db, err := sqlx.Connect("postgres", param)
	//DB connection check
	if err != nil {
		message = "an ERROR occured in database connecting: %s\n"
		message = fmt.Sprintf(message, err)
		message += fmt.Sprintf("queryString: %s", queryString)
		fmt.Printf("err: %s\n", message)
		c.String(500, message)
		return
	}
	defer db.Close()
	//Query
	err = db.Select(&result, sql, queryString)
	if err != nil {
		message = "an ERROR occured in executing SQL: %s \n%s\n"
		message = fmt.Sprintf(message, sql, err)
		message += fmt.Sprintf("queryString: %s", queryString)
		fmt.Printf("err: %s\n", message)
		c.String(500, message)
		return
	}

	c.JSON(200, result)
}

func putHOT(c *gin.Context) {
	var message string
	var hot HOTStatus
	var err error

	//Binding
	err = c.Bind(&hot)
	fmt.Printf("Bind: %s\n", hot)

	if err != nil {
		message = "an ERROR occured in binding form: %s\n"
		message = fmt.Sprintf(message, err)
		fmt.Printf("err: %s\n", message)
		c.String(500, message)
		return
	}
	param := connectParam()
	db, err := sqlx.Connect("postgres", param)
	//DB connection check
	if err != nil {
		message = "an ERROR occured in database connecting: %s\n"
		message = fmt.Sprintf(message, err)
		message += fmt.Sprintf("data [%s]", hot)
		fmt.Printf("err: %s\n", message)
		c.String(500, message)
		return
	}
	defer db.Close()

	sql := `INSERT INTO "hot" 
				("HOT11", "status_no", "hot_comment")
				VALUES(:hot, :status, :comment)
			ON CONFLICT ("HOT11")
				DO UPDATE SET 
				"HOT11"=:hot, 
				"status_no"=:status, 
				"hot_comment"=:comment;
	`
	_, err = db.NamedExec(
		sql,
		map[string]interface{}{
			"hot":     hot.HOT,
			"status":  hot.Status,
			"comment": hot.Comment,
		},
	)
	if err != nil {
		message = "an ERROR occured in executing SQL: %s \n%s\n"
		message = fmt.Sprintf(message, sql, err)
		message += fmt.Sprintf("data [%s]", hot)
		fmt.Printf("err: %s\n", message)
		c.String(500, message)
		return
	}

	fmt.Printf("result: OK hot: %s\n", hot)
	c.JSON(200, hot)
}

func putYJ(c *gin.Context) {
	var message string
	var yj YJStatus
	var err error
	//Binding
	err = c.Bind(&yj)
	fmt.Printf("Bind: %s\n", yj)
	if err != nil {
		message = "an ERROR occured in binding form: %s\n"
		message = fmt.Sprintf(message, err)
		fmt.Printf("err: %s\n", message)
		c.String(500, message)
		return
	}
	param := connectParam()
	db, err := sqlx.Connect("postgres", param)
	//DB connection check
	if err != nil {
		message = "an ERROR occured in database connecting: %s\n"
		message = fmt.Sprintf(message, err)
		message += fmt.Sprintf("data [%s]", yj)
		fmt.Printf("err: %s\n", message)
		c.String(500, message)
		return
	}
	defer db.Close()
	sql := `INSERT INTO "yj" 
					("yjcode", 
					"status_no",
					"yj_comment",
					"drug_code")
				VALUES(
					:yj,
					:status,
					:comment,
					:drug_code)
			ON CONFLICT ("yjcode")
				DO UPDATE SET 
				"yjcode"=:yj, 
				"status_no"=:status, 
				"yj_comment"=:comment,
				"drug_code"=:drug_code;
	`
	_, err = db.NamedExec(
		sql,
		map[string]interface{}{
			"yj":        yj.YJ,
			"status":    yj.Status,
			"comment":   yj.Comment,
			"drug_code": yj.DrugCode,
		},
	)
	if err != nil {
		message = "an ERROR occured in executing SQL: %s \n%s\n"
		message = fmt.Sprintf(message, sql, err)
		message += fmt.Sprintf("data [%s]", yj)
		fmt.Printf("err: %s\n", message)
		c.String(500, message)
		return
	}

	fmt.Printf("result: OK hot: %s\n", yj)
	c.JSON(200, yj)
}

func putCustomYJ(c *gin.Context) {
	var message string
	var cyj CustomYJ
	var err error
	err = c.Bind(&cyj)
	//Binding
	if err != nil {
		message = "an ERROR occured in binding form: %s\n"
		message = fmt.Sprintf(message, err)
		fmt.Printf("err: %s\n", message)
		c.String(500, message)
		return
	}
	param := connectParam()
	db, err := sqlx.Connect("postgres", param)
	//DB connection check
	if err != nil {
		message = "an ERROR occured in database connecting: %s\n"
		message = fmt.Sprintf(message, err)
		message += fmt.Sprintf("data [%s]", cyj)
		fmt.Printf("err: %s\n", message)
		c.String(500, message)
		return
	}
	defer db.Close()
	sql := `INSERT INTO "custom_yj" 
					("HOT9", 
					"yjcode")
				VALUES(
					:HOT9,
					:yjcode)
			ON CONFLICT ("HOT9")
				DO UPDATE SET 
				"HOT9"=:HOT9, 
				"yjcode"=:yjcode;
	`
	result, err := db.NamedExec(
		sql,
		map[string]interface{}{
			"HOT9":   cyj.HOT9,
			"yjcode": cyj.YJ,
		},
	)
	if err != nil {
		message = "an ERROR occured in executing SQL: %s \n%s\n"
		message = fmt.Sprintf(message, sql, err)
		message += fmt.Sprintf("data [%s]", cyj)
		fmt.Printf("err: %s\n", message)
		c.String(500, message)
		return
	}

	fmt.Println(result)
	c.String(200, cyj.String())
}

//redirectToPMDA YJコードからredirect
func redirectToPMDA(c *gin.Context) {
	var url string
	yjcode := c.Param("yjcode")

	//http://localhost:8080/redirect/%s的なのが入る
	redirectURL := os.Getenv("YJ_REDIRECT_URL")

	url = fmt.Sprintf(redirectURL, yjcode)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

//redirectToPMDAWithDrugCode 薬品コードからredirect
func redirectToPMDAWithDrugCode(c *gin.Context) {
	var url string
	var message string
	drugCode := c.Param("drugCode")
	yjcode, err := drugcodeToYJ(drugCode)

	if err != nil {
		message = "%s\ncan not get YJ. 薬品コード:[%s]"
		message = fmt.Sprintf(message, err, drugCode)
		fmt.Printf("err: %s\n", message)
		c.String(404, message)
		return
	}

	//http://localhost:8080/redirect/%s的なのが入る
	redirectURL := os.Getenv("YJ_REDIRECT_URL")

	url = fmt.Sprintf(redirectURL, yjcode)
	//fmt.Printf("redirect url:%s\n", url)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

//drugcodeToYJ 薬日コードからYJコードに変換
func drugcodeToYJ(drugCode string) (yjcode string, err error) {
	param := connectParam()
	db, err := sqlx.Connect("postgres", param)
	//DB connection check
	var message string
	if err != nil {
		message = "an ERROR occured in database connecting: %s\n"
		message = fmt.Sprintf(message, err)
		message += fmt.Sprintf("drugcode [%s]", drugCode)
		err = fmt.Errorf("err: %s\n", message)
		return
	}
	defer db.Close()
	sql := `
SELECT
	"個別医薬品コード"
FROM available_view
WHERE "drug_code"=$1;
`

	result := struct {
		YJCode string `db:"個別医薬品コード"`
	}{}

	//resultに突っ込む
	err = db.Get(&result, sql, drugCode)
	if err != nil {
		message = "an ERROR occured in executing SQL: %s \n%s\n"
		message = fmt.Sprintf(message, sql, err)
		message += fmt.Sprintf("drugcode [%s]\n", drugCode)
		err = fmt.Errorf("err: %s\n", message)
		fmt.Println(err)
		return
	}

	yjcode = result.YJCode

	return

}

func main() {
	var user = os.Getenv("ADMIN_USER")
	var password = os.Getenv("ADMIN_PASSWORD")
	var accounts = gin.Accounts{
		user: password,
	}

	fmt.Printf("###################accounts: %s################\n", accounts)

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello world")
	})
	r.GET("/hoge", func(c *gin.Context) {
		c.String(200, "fuga")
	})

	//redirect to PMDA
	r.GET("/redirect/pmda/yj/:yjcode/", redirectToPMDA)
	r.GET("/redirect/pmda/drug_code/:drugCode/", redirectToPMDAWithDrugCode)

	//y
	r.GET("/json/y/", handleY)
	//medis
	r.GET("/json/medis/", handleMedis)
	//available
	r.GET("/json/available/", handleAvailable)
	r.PUT("/edit/hot/", putHOT)
	r.PUT("/edit/yj/", putYJ)
	r.PUT("/edit/custom_yj/", putCustomYJ)
	//barcode
	r.GET("/barcode/:barcode/", handleBarcode)

	//riot.js
	r.Static("/riot/", "/bootstrap/riot/")
	//fetch.js
	r.StaticFile("/fetch/fetch.umd.js", "/bootstrap/fetch/fetch.umd.js")

	//promise-polyfill.js
	r.StaticFile("/promise/polyfill.min.js", "/bootstrap/promise-polyfill/dist/polyfill.min.js")

	//static
	r.Static("/static/", "/asset/static/")
	//r.Static("/static_auth/", "/asset/static_auth/")
	auth_r := r.Group("/static_auth/", gin.BasicAuth(accounts))
	auth_r.Static("/", "/asset/static_auth/")

	//authorized := r.Group("/static_auth/", gin.BasicAuth(accounts))
	//authorized.Static("/", "/asset/static_auth/")

	var port string
	flag.StringVar(&port, "port", ":8080", "Port number default:':8080'")

	fmt.Printf("listen: %s \n", port)
	r.Run(port)
}
