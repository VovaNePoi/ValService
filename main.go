package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

const (
	adress = "localhost:8080"
	market = "https://cbr.ru/scripts/XML_daily.asp"
)

func main() {
	s := newAPIService()
	s.StartServer(adress)
}

type ApiService struct {
}

func newAPIService() *ApiService {
	return &ApiService{}
}

// listen for get requests, starting server on different adresses, default = localhost:8080
func (s *ApiService) StartServer(adress string) error {
	http.HandleFunc("/currency/", s.HandleCurrencyService)
	log.Printf("Server start on http://localhost:8080/currency")
	return http.ListenAndServe(adress, nil)
}

type marketAPIparams struct {
	market   string
	date     string
	currency string
}

func (m *marketAPIparams) SetParams(params url.Values) {
	m.market = market
	m.currency = params.Get("currency")
	m.date = params.Get("date")
}

// type markertQueryParamsName struct {
// 	dateParamName string
// 	stockCode     string
// }

type Valute struct {
	XMLName   xml.Name `xml:"Valute"`
	ID        string   `xml:"ID,attr"`
	NumCode   string   `xml:"NumCode"`
	CharCode  string   `xml:"CharCode"`
	Nominal   string   `xml:"Nominal"`
	Name      string   `xml:"Name"`
	Value     string   `xml:"Value"`
	VunitRate string   `xml:"VunitRate"`
}

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Date    string   `xml:"Date,attr"`
	Name    string   `xml:"name,attr"`
	Valute  []Valute `xml:"Valute"`
}

// getting query params for market API request
// sending response to requesting service
func (s *ApiService) HandleCurrencyService(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	var queryParams marketAPIparams
	queryParams.SetParams(params)

	log.Printf("request params: market adress: %v, currency: %v, date: %v", queryParams.market, queryParams.currency, queryParams.date)
	GetCurrencyValue(queryParams)
}

// get data from market
// take query params
// if no market param, taking default
func GetCurrencyValue(params marketAPIparams) {
	var url string
	url = urlCreator(params)

	// creating browser user-agent, couse cbr dont allow to get data
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 YaBrowser/23.12.1.1067 Yowser/2.5 Safari/537.36")

	log.Printf("sending get request to %s", url)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("request to market error: %v", err)
	}
	log.Printf("status code of get request: %v", resp.StatusCode)
	log.Printf("Request done")

	defer resp.Body.Close()
	contentType := resp.Header.Get("Content-Type") // check if not utf-8 txt/xml
	log.Printf("content-type: %v", contentType)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("response reading error err: %v", err)
	}
	RestuctRespToJSON(body)
}

// taking querry params and return market API url
func urlCreator(params marketAPIparams) string {
	return fmt.Sprintf("%v?date_req=%v", params.market, params.date)
}

func RestuctRespToJSON(respBodyWin1251 []byte) {
	log.Printf("starting restruct response data")
	var valCurs ValCurs

	decoder := charmap.Windows1251.NewDecoder()
	respBodyUTF8, _, err := transform.Bytes(decoder, respBodyWin1251)
	if err != nil {
		log.Printf("transforming win1251 to UTf-8 err: %v", err)
	}

	respBodyUTF8Str := string(respBodyUTF8)
	respBodyUTF8Str = strings.Replace(respBodyUTF8Str, `encoding="windows-1251"`, `encoding="UTF-8"`, 1)
	respBodyUTF8 = []byte(respBodyUTF8Str)

	err = xml.Unmarshal(respBodyUTF8, &valCurs)
	if err != nil {
		log.Printf("xml data unmarshal err: %v", err)
	}

	fmt.Printf("%v", valCurs)

	log.Printf("data restructed to JSON, data: %v", valCurs)
}
