package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const (
	url = "https://raw.githubusercontent.com/pcm-dpc/COVID-19/master/dati-province/dpc-covid19-ita-province.csv"
)

type covidData struct {
	Date           string `json:"data"`
	City           string `json:"city"`
	DailyIncrement int    `json:"daily"`
	NumeroCasi     string `json:"casi"`
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":5900", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	rows, err := fetchURL(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	t, err := template.ParseFiles("templates/index.html")
	t.Execute(w, rows)

}

func fetchURL(url string) (m map[string][]covidData, err error) {
	m = make(map[string][]covidData)
	res, err := http.Get(url)
	if err != nil {
		return m, fmt.Errorf("error fetching url: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return m, fmt.Errorf("response error. StatusCode: %d", res.StatusCode)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return m, fmt.Errorf("error reading body: %v", err)

	}

	csvFile := csv.NewReader(bytes.NewBuffer(body))
	records, err := csvFile.ReadAll()

	if err != nil {
		return m, fmt.Errorf("error converting csv: %v", err)
	}
	prev := 0
	for i, r := range records {
		prov := r[6]
		if prov != "" {
			if len(m[prov]) > 0 {
				if i > 0 {
					prev = conv(m[prov][len(m[prov])-1].NumeroCasi)
				}
			}
			m[prov] = append(m[prov], covidData{
				Date:           strings.Split(r[0], " ")[0],
				City:           r[5],
				NumeroCasi:     r[len(r)-1],
				DailyIncrement: conv(r[len(r)-1]) - prev,
			})
		}
	}
	return m, nil
}

func conv(s string) (i int) {
	i, _ = strconv.Atoi(s)
	return i
}
