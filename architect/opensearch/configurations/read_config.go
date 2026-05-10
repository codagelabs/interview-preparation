package configurations

import (
"encoding/json"
"fmt"
"io/ioutil"
"os"
)

func NewConfiguration(fileName string) *config {
	config := config{}
	fileName = fmt.Sprint(fileName, ".json")
	jsonFile, err := os.Open(fileName)
	if err != nil {
		panic("Error while reading application config. ")
	}

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &config)
	fmt.Println(config)
	return &config
}
