package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

const apiURL = "https://api.exchangerate.host/convert"

var fromCurr = flag.String("from", "USD", "Specify currecny convert from")
var toCurr = flag.String("to", "KZT", "Specify currecny convert to")
var amount = flag.Float64("amount", 1, "Specfy amount of currecny for conversion")
var date = flag.String("date", "", "Specfy date for historical data (format YYYY-MM-DD)")

func main() {
	flag.Parse()
	var dt time.Time
	if *date == "" {
		dt = time.Now()
	}
	url := createAPIQuery(*fromCurr, *toCurr, *amount, dt)

	result, err := getResult(url)
	if err != nil {
		log.Fatal(err)
	}

	if result.Success {
		fmt.Printf("[%s] %v %s = %.2f %s\n", result.Date, result.Query.Amount, result.Query.From, result.Result, result.Query.To)
	}
}

func createAPIQuery(from, to string, amount float64, date time.Time) string {
	endPoint, err := url.Parse(apiURL)
	if err != nil {
		log.Fatal(err)
	}
	values := url.Values{}
	values.Add("from", from)
	values.Add("to", to)
	values.Add("amount", fmt.Sprintf("%v", amount))
	values.Add("date", date.Format("2006-01-02"))

	endPoint.RawQuery = values.Encode()

	return endPoint.String()
}

func getResult(url string) (*Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if res.Body != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	result := Response{}

	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
