package file

import (
	"encoding/json"
	"io/ioutil"
)

// SaveJSON ...
func SaveJSON(filepath string, obj interface{}) error {
	data, err := json.MarshalIndent(obj, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath, data, 0644)
}

// LoadJSON ...
func LoadJSON(filepath string, obj interface{}) (error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, obj)
}
