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

type sliceString []interface{}

type mapString map[string]interface{}

func (m mapSlice) extract(key interface{}) interface{} {
	if i, ok := key.(int); ok {
		for j, _ := range m {
			if i == j {
				return m[i]
			}
		}
	}
	return nil
}

func (m sliceString) extract(key interface{}) interface{} {
	if i, ok := key.(int); ok {
		if i > len(m)-1 {
			return false
		}
		return m[i]
	}
	return nil
}

func (m mapString) extract(key interface{}) interface{} {
	if i, ok := key.(string); ok {
		for j, _ := range m {
			if i == j {
				return m[i]
			}
		}
	}
	return nil
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
	case map[string]interface{}:
		return mapString(v), nil
	case []map[string]interface{}:
		return mapSlice(v), nil
	case []interface{}:
		return sliceString(v), nil
	default:
		return nil, errors.New("No type found")
	}
}

func (j *JSONExtract) wrapperType() (bool, error) {
	j.RawJSON = strings.Trim(j.RawJSON, "\n")
	j.RawJSON = strings.Trim(j.RawJSON, " ")
	wrapper := string(j.RawJSON[0]) + string(j.RawJSON[len(j.RawJSON)-1])
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

func findInJSON(rawElem jsonElem, chain string) (interface{}, error) {
	properties := splitProperties(chain)
	var root interface{}
	for _, v := range properties {
		// Check whether v is an index reference (eg "[5]"):
		if isIndex(v) {
			root = rawElem.extract(getIndex(v))
		} else {
			root = rawElem.extract(v)
		}

		if root == nil {
			errorMsg := fmt.Sprintf("path `.../%s` not found in JSON", v)
			return nil, errors.New(errorMsg)
		}

		// Type check root and repeat the process if same it not string
		if _, ok := root.(string); ok {
			break
		}

		rawElem, _ = getElemType(root)
	}
	return root, nil
}

// JSONExtract creates an instance of the package for a raw JSON (stored as RawJSON)
type JSONExtract struct {
	RawJSON string
}

// Extract returns the value at the path specified or error
func (j *JSONExtract) Extract(path string) (interface{}, error) {

	wrapperType, err := j.wrapperType()
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

	result, err := findInJSON(unwrapped, path)
	if err != nil {
		return nil, err
	}
	return result, nil
}
