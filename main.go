package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const apiKey = "9CZs1BoRo1OjDycZhyNfRwwLhYyAYPMP"
const apiURL = "https://api.apilayer.com/fixer/convert"

var fromCurr = flag.String("from", "USD", "Specify currecny convert from")
var toCurr = flag.String("to", "KZT", "Specify currecny convert to")
var amount = flag.Float64("amount", 1, "Specfy amount of currecny for conversion")

func main() {

	url := createAPIQuery()

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("apikey", apiKey)

	if err != nil {
		fmt.Println(err)
	}
	res, err := client.Do(req)
	if res.Body != nil {
		defer res.Body.Close()
	}
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	result := struct {
		Success bool
		Result  float32 `json:"result"`
		Query   struct {
			From   string  `json:"from"`
			To     string  `json:"to"`
			Amount float32 `json:"amount"`
		}
	}{}

	if err = json.Unmarshal(body, &result); err != nil {
		log.Fatal(err)
	}

	if result.Success {
		fmt.Printf("%v %s = %.2f %s\n", result.Query.Amount, result.Query.From, result.Result, result.Query.To)
	}
}

func createAPIQuery() string {
	endPoint, err := url.Parse(apiURL)
	if err != nil {
		log.Fatal(err)
	}

	flag.Parse()

	values := url.Values{}
	values.Add("from", *fromCurr)
	values.Add("to", *toCurr)
	values.Add("amount", fmt.Sprintf("%v", *amount))

	endPoint.RawQuery = values.Encode()

	return endPoint.String()
}
