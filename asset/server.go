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

func medis_handler(w http.ResponseWriter, r *http.Request) {
	barcode := r.FormValue("barcode")
	gs1 := GS1(barcode) //GS1として読んでみる
	var jan JAN
	//validate check digit
	if !gs1.CheckDigitOK() {
		jan = JAN(barcode) //ダメならJANとして読んでみる
		if !jan.CheckDigitOK() {
			fmt.Fprintf(w, "wrong barcode: checkdigit error")
			fmt.Fprintf(w, barcode)
			return
		}
	} else { //GS1ならjanに変換する
		jan = gs1.ToJAN()
	}
	param := connect_param()
	db, err := sqlx.Connect("postgresql", param)
	//DB connection check
	if err != nil {
		fmt.Fprintf(w, "an ERROR occured in database connecting")
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
		fmt.Fprintf(w, "an ERROR occured in executing SQL")
		fmt.Fprintf(w, fmt.Sprintf("SQL:%s", sql))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(medis)
	if err != nil {
		fmt.Fprintf(w, "an ERROR occured in marshaling JSON")
	}

	fmt.Fprintf(w, string(b))

	return
}
func y_handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World Y")
}

func hellow_handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World")
}

func main() {
	http.HandleFunc("/medis/", medis_handler)
	http.HandleFunc("/y/", y_handler)
	http.HandleFunc("/", hellow_handler)
	fmt.Println("listen: 8080")
	http.ListenAndServe(":8080", nil)
}
