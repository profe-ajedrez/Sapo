package file

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// PrettyPrint  exports a true constant value to indicate to pretty print generated json
const PrettyPrint = true

// UglyPrint  exports a false constant value to indicate to not pretty print generated json
const UglyPrint = false

// JSONifyMap returns a jsonified []byte and the generated error from a map[string]string
func JSONifyMap(jsonMap map[string]string, prettyPrint bool) ([]byte, error) {
	if prettyPrint {
		return json.MarshalIndent(jsonMap, "", "  ")
	}
	return json.Marshal(jsonMap)
}

// ReadJSON   parses a json file as a GO map.
// Returns tuple composed of the map and the error if any, or nil.
//
// Use example:
// `myMap, err := json.parser.file.Read(pathToMyJsonFile)`
func ReadJSON(jsonFile string) (map[string]string, error) {
	file, err := os.Open(jsonFile)

	defer file.Close()

	if err == nil {
		bytes, _ := ioutil.ReadAll(file)
		var result map[string]string
		json.Unmarshal([]byte(bytes), &result)

		return result, nil
	}
	return nil, err
}

// WriteJSON  writes a map[string]string as a json file
func WriteJSON(jsonMap map[string]string, jsonFile string, prettyPrint bool) error {
	jasonDecoded, err := JSONifyMap(jsonMap, prettyPrint)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(jsonFile, jasonDecoded, 0644)
	return err
}
