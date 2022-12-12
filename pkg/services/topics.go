package services

import (
	"log"
	"ogomez/mkt-export/pkg/config"
)

type TopicSubtService struct {
	Conf        config.Config
	ExcelReader ExcelReader
}

func NewEventsService(conf config.Config) *TopicSubtService {

	excelReader := NewExcelReader(conf)

	return &TopicSubtService{
		Conf:        conf,
		ExcelReader: *excelReader,
	}
}

func (e *TopicSubtService) Export() {
	log.Printf("Exporting Events from: %s ", e.Conf.Input)

	err := e.ExcelReader.ReadInput()

	if err != nil {
		log.Println("Error processing Excel File")
	}  
}
