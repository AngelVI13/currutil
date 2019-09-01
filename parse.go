package main

import (
	"strings"
	"log"
	"net/http"
	"io/ioutil"
	"errors"
	"strconv"
	"fmt"
)

const (
        trendUp = iota
        trendDown
        trendEq
)

type currency struct {
        Id string
        Name string
        Rate float64
        Trend string
}

type currencies []currency

func getResponseString(url string) (string) {
        resp, err := http.Get(url)
        if err != nil {
                panic(err)
        }
        defer resp.Body.Close()

        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                panic(err)
        }
        return string(body)
}


const (
	forexTableStartString = "<div class=\"forextable\">"
	forexTableEndString = "</table>"
	tableBodyStartString = "<tbody>"
	tableBodyEndString = "</tbody>"
	rowSeparatorString = "</tr>"

	currencyIdStartString = "<td id=\""
	currencyIdEndString = "\" class=\"currency\""

	currencyNameOuterStartString = "<td class=\"alignLeft\">"
	currencyNameOuterEndString = "</td>"

	currencyNameInnerStartString = "\">"
	currencyNameInnerEndString = "</a>"

	currencyRateStartString = "span class=\"rate\">"
	currencyRateEndString = "</span>"

	currencyTrendStartString = "span class=\"trend "
	currencyTrendEndString = "\">"
)

const (
	removeMarkersFlag = iota
	keepMarkersFlag
)


func parseResponseString(response string) currencies {
	result := make(currencies, 0)

	body, err := extractForexTableBodyFromResponse(response)
	if err != nil {
		log.Fatal(err)
	}

	for _, textSlice := range(strings.Split(body, rowSeparatorString)) {
		// log.Print(slice)
		currencyItem, err := extractCurrencyInfo(textSlice)
		if err != nil {
			log.Print(err)
			continue
		}

		result = append(result, currencyItem)
	}
	
	// log.Print(result)
	return result
}

func extractForexTableBodyFromResponse(response string) (string, error) {
	result, err := getStringBetweenMarkers(response, forexTableStartString, forexTableEndString, keepMarkersFlag)
	if err != nil {  // return early -> here result is an empty string
		return result, err
	}

	result, err = getStringBetweenMarkers(result, tableBodyStartString, tableBodyEndString, keepMarkersFlag)
	if err != nil {  // return early -> here result is an empty string
		return result, err
	}
	return result, nil
}

func extractCurrencyInfo(text string) (currency, error) {
	id, err := getStringBetweenMarkers(text, currencyIdStartString, currencyIdEndString, removeMarkersFlag)
	if err != nil {
		return currency{}, errors.New(fmt.Sprintf("Couldn't extract currency ID from text %s", text))
	}

	var name string
	name, err = getStringBetweenMarkers(text, currencyNameOuterStartString, currencyNameOuterEndString, removeMarkersFlag)
	if err != nil{
		return currency{}, errors.New(fmt.Sprintf("Couldn't extract currency name from text %s", text))
	}
	name, err = getStringBetweenMarkers(name, currencyNameInnerStartString, currencyNameInnerEndString, removeMarkersFlag)
	if err != nil{
		return currency{}, errors.New(fmt.Sprintf("Couldn't extract currency name from text %s", text))
	}

	var rate string
	rate, err = getStringBetweenMarkers(text, currencyRateStartString, currencyRateEndString, removeMarkersFlag)
	if err != nil{
		return currency{}, errors.New(fmt.Sprintf("Couldn't extract currency rate from text %s", text))
	}
	rateF32, err := strconv.ParseFloat(rate, 32) 
	if err != nil {
		return currency{}, errors.New(fmt.Sprintf("Couldn't extract currency rate from text %s", text))	
	}

	var trend string
	trend, err = getStringBetweenMarkers(text, currencyTrendStartString, currencyTrendEndString, removeMarkersFlag)
	if err != nil{
		return currency{}, errors.New(fmt.Sprintf("Couldn't extract currency trend from text %s", text))
	}

	return currency{
		Id: id,
		Name: name,
		Rate: rateF32,
		Trend: trend,
	}, nil	
}


func getStringBetweenMarkers(text, start, end string, flag int) (string, error) {
	result := ""

	idx := strings.Index(text, start)
	if idx == -1 {
		return result, errors.New("Couldn't find start string")
	}

	// remove everything before start of table
	text = text[idx:]

	idx = strings.Index(text, end)
	if idx == -1 {
		return result, errors.New("Couldn't find end string")
	}

	// remove everything after end of table
	result = text[:idx + len(end)]

	if flag == removeMarkersFlag {  // if inclusive is false -> remove markers from result string 
		result = strings.Replace(result, start, "", -1)
		result = strings.Replace(result, end, "", -1)
	}
	return result, nil
}
