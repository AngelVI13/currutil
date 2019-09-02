package main

import (
	"errors"
	"io/ioutil"
	"os"
)

func writeStringToFile(filename string, text string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(text)
	if err != nil {
		return err
	}
	file.Sync() // flush writes to stable storage
	return nil
}

func readStringFromFile(filename string) (text string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	text = string(data)
	return
}
