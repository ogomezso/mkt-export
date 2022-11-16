package export

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type JsonExporter struct{}

func (e JsonExporter) Export(res interface{}, outputPath string) error {
	file, errJson := json.MarshalIndent(res, "", " ")
	if errJson != nil {
    log.Println("Error parsing Json")
    return errJson
	}
	err := ioutil.WriteFile(outputPath+".json", file, 0644)
	if err != nil {
    log.Println("Error writing json file")
		return err
	}
	return nil
}
