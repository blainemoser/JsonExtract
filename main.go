package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	// "os"
	"errors"
	"log"
)

type jsonElem interface {
	extract(interface{}) interface{}
}

type mapSlice []map[string]interface{}

func (m mapSlice) extract(key interface{}) interface{} {

	i := key.(int)
	return m[i]

}

func Json_decode(data string) (interface{}, error) {
	var dat []map[string]interface{}
	err := json.Unmarshal([]byte(data), &dat)
	return dat, err
}

func getJSONPayload(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	return string(bytes)
}

func getType(i interface{}) (jsonElem, error) {
	switch v := i.(type) {
	// case int:

	// case slice of maps:
	case []map[string]interface{}:
		j := new(mapSlice)
		*j = v
		return *j, nil
	default:
		return nil, errors.New("No type found")
	}
}

func main() {
	here := getJSONPayload("https://restcountries.eu/rest/v2/name/usa")
	// fmt.Printf(here); os.Exit(3)
	decoded, err := Json_decode(here)
	if err != nil {
		panic(err.Error())
	}
	x, err := getType(decoded)
	if err != nil { // note extract this
		log.Fatal(err)
	}
	y := x.extract(0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(y)
}
