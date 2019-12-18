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
				"個別医薬品コード" like $1 || '%' OR
				"HOT11" like $1 || '%' OR
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

	//Binding
	err = c.Bind(&hot)
	fmt.Printf("Bind: %s", hot)

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
	//TODO:いつもINSERT成功している?
	//status commentが入らない。
	sql := `INSERT INTO "hot" 
				("HOT11", "status_no", "hot_comment")
				VALUES(:hot, :status, :comment)
			ON CONFLICT ("HOT11")
				DO UPDATE SET 
				"HOT11"=:hot, 
				"status_no"=:status, 
				"hot_comment"=:comment;
	`
	result, err := db.NamedExec(
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
		c.String(500, message)
		return
	}

	fmt.Println("result: " + result + "hot: " + hot.HOT)
	c.String(200, hot.String())
}

func putYJ(c *gin.Context) {
	var message string
	var yj YJStatus
	var err error
	err = c.Bind(&yj)
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
		message += fmt.Sprintf("data [%s]", yj)
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
				"yj"=:comment
				"drug_code"=:drug_code;
	`
	result, err := db.NamedExec(
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
		c.String(500, message)
		return
	}

	fmt.Println(result)
	c.String(200, yj.String())
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
		c.String(500, message)
		return
	}

	fmt.Println(result)
	c.String(200, cyj.String())
}

func redirectToPMDA(c *gin.Context) {
	var url string
	yjcode := c.Param("yjcode")
	redirectURL := os.Getenv("YJ_REDIRECT_URL") //http://localhost:8080/redirect/%s的なのが入る
	url = fmt.Sprintf(redirectURL, yjcode)
	c.Redirect(http.StatusTemporaryRedirect, url)
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
	r.GET("/redirect/pmda/:yjcode/", redirectToPMDA)

	//y
	r.GET("/json/y/", handleY)
	//medis
	r.GET("/json/medis/", handleMedis)
	//available
	r.GET("/json/available/", handleAvailable)
	r.PUT("/edit/hot/", putHOT)
	r.PUT("/edit/yj/", putYJ)
	r.PUT("/edit/cyj/", putCustomYJ)
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

	fmt.Println("listen: 8080")
	r.Run(":8080")
}
