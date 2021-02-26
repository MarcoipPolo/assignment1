package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// A Response struct to restcountries
type Response struct {
	Country    string       `json:"name"`
	Currencies []Currency   `json:"currencies"`
	Border     []string     `json:"borders"`
	Exchange   ExchangeData `json:exchangedata`
}

type Currency struct {
	Code string `json:"code"`
}

type ExchangeData struct {
	Rates Rate   `json:"rates"`
	Name  string `json:"name"`
	Base  string `json:"base"`
	Date  string `json:"date"`
}

type Rate struct {
	CAD float64 `json:"CAD,omitempty"`
	HKD float64 `json:"HKD,omitempty"`
	ISK float64 `json:"ISK,omitempty"`
	PHP float64 `json:"PHP,omitempty"`
	DKK float64 `json:"DKK,omitempty"`
	HUF float64 `json:"HUF,omitempty"`
	CZK float64 `json:"CZK,omitempty"`
	AUD float64 `json:"AUD,omitempty"`
	RON float64 `json:"RON,omitempty"`
	SEK float64 `json:"SEK,omitempty"`
	IDR float64 `json:"IDR,omitempty"`
	INR float64 `json:"INR,omitempty"`
	BRL float64 `json:"BRL,omitempty"`
	RUB float64 `json:"RUB,omitempty"`
	HRK float64 `json:"HRK,omitempty"`
	JPY float64 `json:"JPY,omitempty"`
	THB float64 `json:"THB,omitempty"`
	CHF float64 `json:"CHF,omitempty"`
	SGD float64 `json:"SGD,omitempty"`
	PLN float64 `json:"PLN,omitempty"`
	BGN float64 `json:"BGN,omitempty"`
	TRY float64 `json:"TRY,omitempty"`
	CNY float64 `json:"CNY,omitempty"`
	NOK float64 `json:"NOK,omitempty"`
	NZD float64 `json:"NZD,omitempty"`
	ZAR float64 `json:"ZAR,omitempty"`
	USD float64 `json:"USD,omitempty"`
	MXN float64 `json:"MXN,omitempty"`
	ILS float64 `json:"ILS,omitempty"`
	GBP float64 `json:"GBP,omitempty"`
	KRW float64 `json:"KRW,omitempty"`
	MYR float64 `json:"MYR,omitempty"`
	EUR float64 `json:"EUR,omitempty"`
}

type Final struct {
	Rate []Data `json:"rates"`
}

type Data struct {
	Name     string `json:"name"`
	Currency string `json:"currency"`
	Rate     Rate   `json:"rate"`
	Base     string `json:"base"`
}

type Diagnostic struct {
	ExchangeRateAPI int    `json:"exchangerateapi"`
	RestCountries   int    `json:"restcountries"`
	Version         string `json:"version"`
	Uptime          string `json:"uptime"`
}

//// Gets the exchange history of a provided country with provided start and end date
func exchangeHistory(w http.ResponseWriter, r *http.Request) {
	search := strings.Split(r.URL.Path, "/")
	country := strings.Join(search[4:5], "")
	temp := strings.Split(strings.Join(search[5:], ""), "-")
	startdate := strings.Join(temp[:3], "-")
	enddate := strings.Join(temp[3:], "-")

	var req = "https://restcountries.eu/rest/v2/name/" + country + "?fields=currencies"

	respcurr, err := http.Get(req)
	if err != nil {
		log.Fatalln(err)
	}

	bodycurr, err := ioutil.ReadAll(respcurr.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var curr []Response
	json.Unmarshal(bodycurr, &curr)

	req = "https://api.exchangeratesapi.io/history?start_at=" + startdate + "&end_at=" + enddate + "&symbols=" + curr[0].Currencies[0].Code
	if curr[0].Currencies[0].Code == "EUR" {
		req += "&base=USD"
	}

	resp, err := http.Get(req)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	exrate := json.RawMessage(body)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(exrate)
}

func exchangeBorder(w http.ResponseWriter, r *http.Request) {
	search := strings.Split(r.URL.Path, "/")
	country := strings.Join(search[4:], "")

	url := "https://restcountries.eu/rest/v2/name/" + country

	response, err := http.Get(url)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject []Response
	json.Unmarshal(responseData, &responseObject)

	var data []Data
	var base string

	for i := 0; i < len(responseObject[0].Border); i++ {
		Borders := responseObject[0].Border[i]

		url1 := "https://restcountries.eu/rest/v2/alpha?codes=" + Borders

		responseCountry, err := http.Get(url1)
		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}

		responseDataCountry, err := ioutil.ReadAll(responseCountry.Body)
		if err != nil {
			log.Fatal(err)
		}

		var responseObjectCountry []Response
		json.Unmarshal(responseDataCountry, &responseObjectCountry)

		url2 := "https://api.exchangeratesapi.io/latest?symbols=" + responseObjectCountry[0].Currencies[0].Code
		if responseObjectCountry[0].Currencies[0].Code == responseObject[0].Currencies[0].Code {
			if responseObjectCountry[0].Currencies[0].Code == "EUR" {
				url2 += "&base=USD"
				base = "USD"
			} else {
				url2 += "&base=EUR"
				base = "EUR"
			}
		} else {
			url2 += "&base=" + responseObject[0].Currencies[0].Code
			base = responseObject[0].Currencies[0].Code
		}

		responseExchange, err := http.Get(url2)
		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}

		responseDataExchange, err := ioutil.ReadAll(responseExchange.Body)
		if err != nil {
			log.Fatal(err)
		}

		var responseObjectExchange ExchangeData
		json.Unmarshal(responseDataExchange, &responseObjectExchange)

		responseObjectCountry[0].Exchange = responseObjectExchange

		data = append(data, Data{Name: responseObjectCountry[0].Country, Currency: responseObjectCountry[0].Currencies[0].Code, Rate: responseObjectCountry[0].Exchange.Rates, Base: base})

	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(data)
}

func diagnostics(w http.ResponseWriter, r *http.Request) {

	responseEx, err := http.Get("https://api.exchangeratesapi.io")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseCount, err := http.Get("https://restcountries.eu")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	diagnostic := Diagnostic{ExchangeRateAPI: responseEx.StatusCode, RestCountries: responseCount.StatusCode, Version: "v1"}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(diagnostic)
}

func homePage(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(rw, "This it the home page. Welcome!")
}

func handelRequests() {
	/// We have two endpoints, for the main root, like localhost:4747, it runs homepage function and for localhost:4747/articles it executes AllArticles function
	http.HandleFunc("/", homePage)
	http.HandleFunc("/exchange/v1/exchangehistory/", exchangeHistory)
	http.HandleFunc("/exchange/v1/exchangeborder/", exchangeBorder)
	http.HandleFunc("/exchange/v1/diag/", diagnostics)

	log.Fatal(http.ListenAndServe(getport(), nil))
}

func main() {
	handelRequests()
}

//// Get Port if it is set by environment, else use a defined one like "8080"
func getport() string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return ":" + port
}
