package jsonextract

import (
	"reflect"
	"testing"
)

func checkType(prop interface{}, extracted string, expect reflect.Kind, t *testing.T) {
	i := reflect.ValueOf(prop).Kind()

	// Run type-assertion. We expect string.
	if i != expect {
		t.Errorf("Error in General Test - incorrect underlying type: got %s, want %s", i, expect)
	}
}

// General test
func TestGeneral(t *testing.T) {

	json := &JSONExtract{RawJSON: `[{"blainemoser": "here", "list": ["1", "3", "34"]}]`}

	extracted := "[0]/list/[2]"
	prop, err := json.Extract(extracted)
	if err != nil {
		t.Errorf("Error in JSON extract; got %s, want 34", err.Error())
	}

	checkType(prop, extracted, reflect.String, t)

	// Here we want to retrieve a slice of strings
	extracted = "[0]/list"
	prop, err = json.Extract(extracted)
	if err != nil {
		t.Errorf("Error in JSON extract; got error %s, want [\"1\", \"3\", \"34\"]", err.Error())
	}

	checkType(prop, extracted, reflect.Slice, t)

	if prop, ok := prop.([]interface{}); ok {
		for _, v := range prop {
			checkType(v, extracted, reflect.String, t)
		}
	}

	json = &JSONExtract{RawJSON: `{
		"prop1": "prop1",
		"prop2": {
			"prop2.1": {
				"prop2.1.1": "prop2.1.1"
			},
			"prop2.2": [
				"prop2.2.1",
				"prop2.2.2"
			]
		}
	}`}

	prop, err = json.Extract("prop2")
	checkType(prop, "prop1", reflect.Map, t)

	prop, err = json.Extract("prop2/prop2.1")
	checkType(prop, "prop2/prop2.1", reflect.Map, t)

	prop, err = json.Extract("prop2/prop2.1/prop2.1.1")
	checkType(prop, "prop2/prop2.1/prop2.1.1", reflect.String, t)

	prop, _ = prop.(string)

	if prop != "prop2.1.1" {
		t.Errorf("unexpected property extracted at `prop2/prop2.1/prop2.1.1`; want: %s, got %s", "prop2.1.1", prop)
	}

	prop, err = json.Extract("prop2/prop2.2/[1]")
	checkType(prop, "prop2/prop2.1/[1]", reflect.String, t)

	prop, _ = prop.(string)

	if prop != "prop2.2.2" {
		t.Errorf("unexpected property extracted at `prop2/prop2.2/[1]`; want: %s, got %s", "prop2.2.2", prop)
	}

}

func checkError(expect string, err error, t *testing.T) {
	if err.Error() != expect {
		t.Errorf("incorrect error thrown; got: %s, want: %s", err.Error(), expect)
	}
}

// The purpose of this test is to ensure correct error-reporting
func TestError(t *testing.T) {
	json := &JSONExtract{RawJSON: `{"prop1": "test", "prop2": "test2"`}
	prop, err := json.Extract("prop1")

	if prop != nil {
		t.Errorf("Error not thrown as expected; got: no error, want: error")
	}

	checkError("Invalid JSON wrapper", err, t)

	json = &JSONExtract{RawJSON: `{"prop1": "test", "prop2": "test2"}`}
	prop, err = json.Extract("prop3")

	if prop != nil {
		t.Errorf("Error not thrown as expected after requesting extract of non-existant property; got: no error, want: error")
	}

	checkError("path `.../prop3` not found in JSON", err, t)
}
