package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const dataFile = "data.txt"
const outputFile = "currencies.json"

func main() {
	resp := ""
	if _, err := os.Stat(dataFile); !os.IsNotExist(err) {
		// file exists
		log.Print("Data file exists. Reading from file")
		resp = readStringFromFile(dataFile)
	} else if os.IsNotExist(err) {
		log.Print("Obtaining currency exchange rates")

		resp = getResponseString("https://www.ecb.europa.eu/stats/policy_and_exchange_rates/euro_reference_exchange_rates/html/index.en.html")
		if err = writeStringToFile(dataFile, resp); err != nil {
			log.Fatal(fmt.Sprintf("Failed to write currency info to file. %s", err))
		}
	} else {
		panic(err)
	}

	currencySlice := parseResponseString(resp)

	log.Printf("Extracting currency finished. Extracted %d currencies", len(currencySlice))
	jsonInfo, err := json.Marshal(currencySlice)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Writing currency info to file %s", outputFile)
	if err = writeStringToFile(outputFile, string(jsonInfo)); err != nil {
		log.Fatal(fmt.Sprintf("Failed to write currency info to file. %s", err))
	}
}
