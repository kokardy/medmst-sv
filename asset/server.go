package main

import (
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
				"告示名称" like '%' || $1 || '%' OR
				"薬価基準収載医薬品コード" like '%' || $1 || '%' OR
				"個別医薬品コード" like '%' || $1 || '%' OR
				"HOT11" like '%' || $1 || '%' OR
				"製造会社" like '%' || $1 || '%' OR
				"販売会社" like '%' || $1 || '%'; 
		`
	param := connectParam()
	db, err := sqlx.Connect("postgres", param)
	//DB connection check
	if err != nil {
		message = "an ERROR occured in database connecting: %s\n"
		message = fmt.Sprintf(message, err)
		message += fmt.Sprintf("queryString: %s", queryString)
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
		c.String(500, message)
		return
	}

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
		c.String(500, message)
		return
	}

	c.JSON(200, result)
}

func putHOT(c *gin.Context) {
	var message string
	var hot HOTStatus
	var err error
	err = c.Bind(&hot)
	//Binding
	if err != nil {
		message = "an ERROR occured in binding form: %s\n"
		message = fmt.Sprintf(message, err)
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
		c.String(500, message)
		return
	}
	defer db.Close()
	sql := `INSERT INTO "hot" ("HOT11","status_no")
		VALUES(:hot,:status)
		ON CONFLICT ("HOT11")
			DO UPDATE SET "HOT11"=:hot, "status_no"=:status;
	`
	result, err := db.NamedExec(
		sql,
		map[string]interface{}{
			"hot":    hot.HOT,
			"status": hot.Status,
		},
	)
	if err != nil {
		message = "an ERROR occured in executing SQL: %s \n%s\n"
		message = fmt.Sprintf(message, sql, err)
		message += fmt.Sprintf("data [%s]", hot)
		c.String(500, message)
		return
	}

	fmt.Println(result)
	c.String(200, hot.String())
}

func putYJ(c *gin.Context) {
}

func redirectToPMDA(c *gin.Context) {
	var url string
	yjcode := c.Param("yjcode")
	domain := os.Getenv("YJ_REDIRECTER")
	url = fmt.Sprintf("//%s/%s/", domain, yjcode)
	c.Redirect(http.StatusPermanentRedirect, url)
}

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello world")
	})
	r.GET("/hoge", func(c *gin.Context) {
		c.String(200, "fuga")
	})

	//redirect to PMDA
	r.GET("/redirect/pmda/:yjcode/", redirectToPMDA)

	//y
	r.GET("/json/y/", handleY)
	//medis
	r.GET("/json/medis/", handleMedis)
	//available
	r.GET("/json/available/", handleAvailable)
	r.PUT("/edit/hot/", putHOT)
	r.PUT("/edit/yj/", putYJ)
	//barcode
	r.GET("/barcode/:barcode/", handleBarcode)

	//riot.js
	r.StaticFile("/riot/riot+compiler.min.js", "/bootstrap/riot/riot+compiler.min.js")
	//fetch.js
	r.StaticFile("/fetch/fetch.umd.js", "/bootstrap/fetch/fetch.umd.js")

	//promise-polyfill.js
	r.StaticFile("/promise/polyfill.min.js", "/bootstrap/promise-polyfill/dist/polyfill.min.js")

	//sttic
	r.Static("/static/", "/asset/static/")

	fmt.Println("listen: 8080")
	r.Run(":8080")
}
