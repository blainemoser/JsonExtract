# JsonExtract
## A package for Go for extracting a JSON's properties
### Details
The function `jsonextract.JSONExtract.Extract` decodes a raw JSON object and returns the property at the specified path (if same exists).

To use this function, create an instance of the struct `JSONExtract`, with the property `RawJSON` containing the raw JSON object:
```json := &jsonextract.JSONExtract{RawJSON: `{"path": {"to": {"property": "value"}}}`}```

Retrieve a property:
`property, err := json.Extract("path/to/property")`

The path to the property must be separated by forward-slashes (`/`).

For list indexes, indicate the element position by enclosing the index in square brackets, e.g. `"path/to/list/[0]"`.

The user must perform type-checking (since `jsonextract.JSONExtract.Extract` returns the type `interface{}`) for any one of the following underlying types:
- `string`
- `[]interface{}`
- `[]map[string]interface{}`
- `map[string]interface{}`

### Usage Example
```
package main

import (
	"fmt"

	"github.com/blainemoser/jsonextract"
)

func main() {

  	// Example of a raw JSON object
	jsonText := `{
		"prop1": "prop1_text",
		"prop2": {
			"prop2.1": "prop2.1_text",
			"prop2.2": [
				"item1",
				"item2"
			]
		},
		"list": [
			"1", 
			"2", 
			"3"
		]
	}`
	json := &jsonextract.JSONExtract{RawJSON: jsonText}
	prop, err := json.Extract("prop1")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(prop) // prints "prop1_text"

	prop, err = json.Extract("prop2/prop2.1")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(prop) // prints "prop2.1_text"

	prop, err = json.Extract("list/[2]")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(prop) // prints "3"

	prop, err = json.Extract("prop2/prop2.2/[0]")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(prop) // prints "item1"

	prop, err = json.Extract("prop2/prop2.2/[4]")
	if err != nil {
		fmt.Println(err.Error()) // prints "path `.../[4]` not found in JSON"
	}
	fmt.Println(prop) // prints <nil>

	prop, err = json.Extract("prop2")

	// Type assertion to find string properties nested in "prop2"
	if prop, ok := prop.(map[string]interface{}); ok {
		for _, value := range prop {
			if value, ok := value.(string); ok {
				fmt.Printf("String value found: %s\n", value)
			}
		}
	}
}

```

