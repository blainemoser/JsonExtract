package jsonextract

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type jsonElem interface {
	extract(interface{}) interface{}
}

type mapSlice []map[string]interface{}

func (m mapSlice) extract(key interface{}) interface{} {
	if i, ok := key.(int); ok {
		for j, _ := range m {
			if i == j {
				return m[i]
			}
		}
	}
	return false
}

type sliceString []interface{}

func (m sliceString) extract(key interface{}) interface{} {
	if i, ok := key.(int); ok {
		if i > len(m)-1 {
			return false
		}
		return m[i]
	}
	return false
}

type mapString map[string]interface{}

func (m mapString) extract(key interface{}) interface{} {
	if i, ok := key.(string); ok {
		for j, _ := range m {
			if i == j {
				return m[i]
			}
		}
	}
	return false
}

func jsonDecode(data string, wrapperType bool) (interface{}, error) {
	if wrapperType {
		var dat map[string]interface{}
		err := json.Unmarshal([]byte(data), &dat)
		return dat, err
	}
	var dat []map[string]interface{}
	err := json.Unmarshal([]byte(data), &dat)
	return dat, err
}

func getElemType(i interface{}) (jsonElem, error) {
	switch v := i.(type) {
	// case map with string key:
	case map[string]interface{}:
		return mapString(v), nil
	// case slice of maps:
	case []map[string]interface{}:
		return mapSlice(v), nil
	// case of a json array
	case []interface{}:
		return sliceString(v), nil
	default:
		return nil, errors.New("No type found")
	}
}

func wrapperType(json string) (bool, error) {
	wrapper := string(json[0]) + string(json[len(json)-1])
	switch wrapper {
	case "{}":
		return true, nil
	case "[]":
		return false, nil
	default:
		return false, errors.New("Invalid JSON wrapper")
	}
}

func isIndex(prop string) bool {
	if len(prop) < 3 {
		return false
	}
	if string(prop[0])+string(prop[len(prop)-1]) == "[]" {
		return true
	}
	return false
}

func getIndex(prop string) int {
	ed := strings.Index(prop, "]")
	res, err := strconv.Atoi(prop[1:ed])
	if err != nil {
		panic(err.Error())
	}
	return res
}

func splitProperties(chain string) []string {
	splitter := regexp.MustCompile(`\/`)
	lines := splitter.Split(chain, -1)
	return lines
}

func checkRoot(root interface{}) bool {
	if root, ok := root.(bool); ok {
		if !root {
			return root
		}
	}
	return true
}

func findInJSON(rawElem jsonElem, chain string) (interface{}, error) {
	properties := splitProperties(chain)
	var root interface{}
	// var typeRoot interface{}
	for _, v := range properties {
		// Check whether v is an index reference:
		if isIndex(v) {
			root = rawElem.extract(getIndex(v))
		} else {
			root = rawElem.extract(v)
		}

		if !checkRoot(root) {
			errorMsg := fmt.Sprintf("Specified path not found in JSON at: `.../%s`", v)
			return nil, errors.New(errorMsg)
		}

		// Type check root and repeat the process if not string
		if _, ok := root.(string); ok {
			break
		}

		rawElem, _ = getElemType(root)
	}
	return root, nil
}

// JSONExtract is created with a raw json string. It implements one function: Extract which extracts the property of the provided path
type JSONExtract struct {
	RawJSON string
}

// Extract pulls data from a JSON according to the path specified
func (j *JSONExtract) Extract(chain string) (interface{}, error) {

	wrapperType, err := wrapperType(j.RawJSON)
	if err != nil {
		return nil, err
	}

	decoded, err := jsonDecode(j.RawJSON, wrapperType)
	if err != nil {
		return nil, err
	}

	unwrapped, err := getElemType(decoded)
	if err != nil {
		return nil, err
	}

	result, err := findInJSON(unwrapped, chain)
	if err != nil {
		return nil, err
	}
	return result, nil
}
