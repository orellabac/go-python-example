// xmltojson.go
package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
)

func xmlToMap(data []byte) (map[string]interface{}, error) {
	decoder := xml.NewDecoder(bytes.NewReader(data))
	var token xml.Token
	result := make(map[string]interface{})
	var current string

	for {
		t, err := decoder.Token()
		if err != nil {
			break
		}
		token = t
		switch se := token.(type) {
		case xml.StartElement:
			current = se.Name.Local
		case xml.CharData:
			if len(bytes.TrimSpace(se)) > 0 {
				result[current] = string(se)
			}
		}
	}
	return result, nil
}

//export XMLToJSON
func XMLToJSON(input *C.char) *C.char {
	xmlStr := C.GoString(input)
	mapped, err := xmlToMap([]byte(xmlStr))
	if err != nil {
		return C.CString(fmt.Sprintf(`{"error": "%s"}`, err.Error()))
	}

	jsonBytes, err := json.Marshal(mapped)
	if err != nil {
		return C.CString(fmt.Sprintf(`{"error": "marshal failed: %s"}`, err.Error()))
	}

	return C.CString(string(jsonBytes))
}

func main() {}
