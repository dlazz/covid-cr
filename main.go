package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var url = "https://raw.githubusercontent.com/pcm-dpc/COVID-19/master/dati-province/dpc-covid19-ita-province.csv"

const (
	timeLayout = "2006-01-02 15:04:05"
)

type covidCR struct {
	Date           string `json:"data"`
	City           string `json:"city"`
	DailyIncrement int    `json:"daily`
	NumeroCasi     string `json:"casi"`
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":5900", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	rows, err := getValues(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	t, err := template.ParseFiles("templates/index.html")
	t.Execute(w, rows)

}
func getValues(url string) (rows []covidCR, err error) {
	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return rows, fmt.Errorf("response error. StatusCode: %d", res.StatusCode)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return rows, fmt.Errorf("error reading body: %v", err)

	}

	csvFile := csv.NewReader(bytes.NewBuffer(body))
	record, err := csvFile.ReadAll()

	if err != nil {
		return rows, fmt.Errorf("error convertin csv: %v", err)
	}

	rows = []covidCR{}
	pvDay := 0
	for _, item := range record {
		if item[6] == "CR" {
			nCasi := item[len(item)-1]
			nc, _ := strconv.Atoi(nCasi)
			row := covidCR{}
			row.Date = strings.Split(item[0], " ")[0]
			row.City = item[5]
			row.DailyIncrement = nc - pvDay
			row.NumeroCasi = nCasi

			rows = append(rows, row)
			pvDay = nc
		}
	}

	return rows, nil
}
