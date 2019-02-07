package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func connect_param() (param string) {

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
			c.String(200, message)
		}
	} else { //GS1ならjanに変換する
		jan = gs1.ToJAN()
	}
	param := connect_param()
	db, err := sqlx.Connect("postgresql", param)
	//DB connection check
	if err != nil {
		message = "an ERROR occured in database connecting"
		c.String(200, message)
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
		c.String(200, message)
	}

	c.JSON(200, medis)

}

func handleY(c *gin.Context) {
	c.String("Hello World YYYYY")
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
