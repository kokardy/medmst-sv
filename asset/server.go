package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"net/http"
)

func connect_param() (param string) {

	return
}

func medis_handler(w http.ResponseWriter, r *http.Request) {
	barcode := r.FormValue("barcode")
	gs1 := GS1(barcode)
	if !gs1.CheckDigitOK() {
		fmt.Fprintf(w, "wrong barcode: checkdigit error")
		fmt.Fprintf(w, barcode)
	}
	jan := gs1.ToJAN()
	param := connect_param()
	db, err := sqlx.Connect("postgresql", param)
	if err != nil {
		fmt.Fprintf(w, "an ERROR occured in database connecting")
	}

	sql := `SELECT "薬品名" FROM "medis" WHERE "JANコード";`

	medis := Medis{}
	db.Get(&medis, sql, string(jan))

	fmt.Fprintf(w, "Hello, World")
	fmt.Fprintf(w, barcode)
}
func y_handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World")
}

func main() {
	http.HandleFunc("/medis", medis_handler)
	http.HandleFunc("/y", y_handler)
	http.ListenAndServe(":8080", nil)
}
