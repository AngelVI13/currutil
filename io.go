package main

import (
	"os"
	"io/ioutil"
)

func writeStringToFile(filename string, text string) {
        file, err := os.Create(filename)
        if err != nil {
                panic(err)
        }
        defer file.Close()

        _, err = file.WriteString(text)
        if err != nil {
                panic(err)
        }
        file.Sync()  // flush writes to stable storage
}

func readStringFromFile(filename string) (text string) {
        data, err := ioutil.ReadFile(filename)
        if err != nil {
                panic(err)
        }
        text = string(data)
        return
}
