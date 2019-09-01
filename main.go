package main

import (
        "log"
        "os"
        "encoding/json"
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
                writeStringToFile(dataFile, resp)
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
        writeStringToFile(outputFile, string(jsonInfo))
}
