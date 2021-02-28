package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var startTime time.Time

// Response contains the responses from restcountries
type Response struct {
	Country    string       `json:"name"`
	Currencies []Currency   `json:"currencies"`
	Border     []string     `json:"borders"`
	Exchange   ExchangeData `json:"exchangedata"`
}

// Currency contains the currency code responce from restcountries
type Currency struct {
	Code string `json:"code"`
}

// ExchangeData contains data for the individual countries in exchangeborder
type ExchangeData struct {
	Rates Rate   `json:"rates"`
	Name  string `json:"name"`
	Base  string `json:"base"`
	Date  string `json:"date"`
}

// Rate contains all the values for the different currencies
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

// Data contains the data for exchangeborder
type Data struct {
	Name     string `json:"name"`
	Currency string `json:"currency"`
	Rate     Rate   `json:"rate"`
	Base     string `json:"base"`
}

// Diagnostic contains the diagnostic data for the api's and uptime
type Diagnostic struct {
	ExchangeRateAPI int    `json:"exchangerateapi"`
	RestCountries   int    `json:"restcountries"`
	Version         string `json:"version"`
	Uptime          string `json:"uptime"`
}

// getBody returns the get request received from the api
func getBody(req string) []byte {
	respcurr, err := http.Get(req)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(respcurr.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return body
}

// Gets the exchange history of a provided country with provided start and end date
func exchangeHistory(w http.ResponseWriter, r *http.Request) {
	// Splits the URL into splices and then to individual strings
	search := strings.Split(r.URL.Path, "/")
	country := strings.Join(search[4:5], "")
	temp := strings.Split(strings.Join(search[5:], ""), "-")
	startdate := strings.Join(temp[:3], "-")
	enddate := strings.Join(temp[3:], "-")

	// Makes the get request to get the currency of said country and stores it in the struct
	var req = "https://restcountries.eu/rest/v2/name/" + country + "?fields=currencies"
	var currency []Response
	json.Unmarshal(getBody(req), &currency)

	// Makes a get request based on the currency code, startdate and enddate
	req = "https://api.exchangeratesapi.io/history?start_at=" + startdate + "&end_at=" + enddate + "&symbols=" + currency[0].Currencies[0].Code
	if currency[0].Currencies[0].Code == "EUR" {
		req += "&base=USD"
	}
	exrate := json.RawMessage(getBody(req))

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(exrate)
}

// Gets all the borders to a country and their rates
func exchangeBorder(w http.ResponseWriter, r *http.Request) {
	// Splits the URL into splices and then gets the country
	search := strings.Split(r.URL.Path, "/")
	country := strings.Join(search[4:], "")

	// Makes the get request and fill the struct
	var req = "https://restcountries.eu/rest/v2/name/" + country
	var responseObj []Response
	json.Unmarshal(getBody(req), &responseObj)

	var excBorder []Data
	var base string

	// Loops for each uniqe border and gets that country's data and exchangerate
	// based on the original country's currency
	for i := 0; i < len(responseObj[0].Border); i++ {
		Borders := responseObj[0].Border[i]

		// Gets the country bordering to the original
		req = "https://restcountries.eu/rest/v2/alpha?codes=" + Borders
		var responseObjCountry []Response
		json.Unmarshal(getBody(req), &responseObjCountry)

		// Gets the currency exchangerate for the country
		req = "https://api.exchangeratesapi.io/latest?symbols=" + responseObjCountry[0].Currencies[0].Code
		// if the countries around uses the same currency as the original
		// compare to some other currency
		if responseObjCountry[0].Currencies[0].Code == responseObj[0].Currencies[0].Code {
			if responseObjCountry[0].Currencies[0].Code == "EUR" {
				req += "&base=USD"
				base = "USD"
			} else {
				req += "&base=EUR"
				base = "EUR"
			}
		} else {
			req += "&base=" + responseObj[0].Currencies[0].Code
			base = responseObj[0].Currencies[0].Code
		}
		var responseObjExchange ExchangeData
		json.Unmarshal(getBody(req), &responseObjExchange)

		responseObjCountry[0].Exchange = responseObjExchange

		// Appends the data generated in the for loop to the Data array
		excBorder = append(excBorder, Data{Name: responseObjCountry[0].Country, Currency: responseObjCountry[0].Currencies[0].Code, Rate: responseObjCountry[0].Exchange.Rates, Base: base})
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(excBorder)
}

// Initializes the time of server start
func init() {
	startTime = time.Now()
}

// returns the duration the server has been running
func uptime() time.Duration {
	return time.Since(startTime)
}

// formats the durration output
func shortDur(d time.Duration) string {
	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = s[:len(s)-2]
	}
	if strings.HasSuffix(s, "h0m") {
		s = s[:len(s)-2]
	}
	return s
}

// Checks the api's and keeps track of uptime
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

	diagnostic := Diagnostic{ExchangeRateAPI: responseEx.StatusCode, RestCountries: responseCount.StatusCode, Version: "v1", Uptime: shortDur(uptime())}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(diagnostic)
}

// Contains all the links the user can take and how to use them
func homePage(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(rw, "This it the home page. Welcome!")
	fmt.Fprintf(rw, "\n/exchange/v1/exchangehistory/{country_name}/{begin_date-end_date}")
	fmt.Fprintf(rw, "\n/exchange/v1/exchangeborder/{country_name}")
	fmt.Fprintf(rw, "\n/exchange/v1/diag")
}

// Handles all of the endpoints
func handelRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/exchange/v1/exchangehistory/", exchangeHistory)
	http.HandleFunc("/exchange/v1/exchangeborder/", exchangeBorder)
	http.HandleFunc("/exchange/v1/diag/", diagnostics)

	log.Fatal(http.ListenAndServe(getport(), nil))
}

// Start point of execution
func main() {
	handelRequests()
}

// Get Port if it is set by environment or sets the port to 8080
func getport() string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return ":" + port
}
