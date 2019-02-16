package main

import (
	"fmt"
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

func handleY(c *gin.Context) {
	var yList []Y
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
		c.String(500, message)
		return
	}
	//Query
	err = db.Select(&yList, sql, queryString)
	if err != nil {
		message = "an ERROR occured in executing SQL: %s \n%s\n"
		message = fmt.Sprintf(message, sql, err)
		c.String(500, message)
		return
	}

	c.JSON(200, yList)
}

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello world")
	})
	r.GET("/hoge", func(c *gin.Context) {
		c.String(200, "fuga")
	})

	//y
	r.GET("/y/", handleY)
	//medis
	r.GET("/medis/", handleY)
	//barcode
	r.GET("/barcode/:barcode/", handleBarcode)

	fmt.Println("listen: 8080")
	r.Run(":8080")
}
