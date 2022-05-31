package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const apiURL = "https://api.exchangerate.host/convert"

var fromCurr = flag.String("from", "USD", "Specify currecny convert from")
var toCurr = flag.String("to", "KZT", "Specify currecny convert to")
var amount = flag.Float64("amount", 1, "Specfy amount of currecny for conversion")
var date = flag.String("date", "", "Specfy date for historical data (format YYYY-MM-DD)")

func main() {
	flag.Parse()
	dt := time.Now().UTC()
	var err error
	if *date != "" {
		if dt, err = time.Parse("2006-01-02", *date); err != nil {
			log.Fatal(err)
		}
	}

	url := createAPIQuery(*fromCurr, *toCurr, *amount, dt)

	result, err := getResult(url)
	if err != nil {
		log.Fatal(err)
	}

	var resultTheDayBefore *Response
	var change float64
	if *amount == 1 {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			url = createAPIQuery(*fromCurr, *toCurr, *amount, dt.Add(-time.Hour*24))
			resultTheDayBefore, err = getResult(url)
			if err != nil {
				log.Fatal(err)
			}
			change = result.Result - resultTheDayBefore.Result
			if math.Abs(change) < 0.01 {
				change = 0
			}
			wg.Done()
		}()
		wg.Wait()
	}

	if result != nil && result.Success {
		fmt.Printf("[%s] %v %s = %.2f %s", result.Date, result.Query.Amount, result.Query.From, result.Result, result.Query.To)
	}

	if resultTheDayBefore != nil && resultTheDayBefore.Success {
		fmt.Printf(", change: %.2f\n", change)
	} else {
		fmt.Println()
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
