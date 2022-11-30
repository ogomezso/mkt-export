package services

import (
	"fmt"
	"log"
	"ogomez/mkt-export/pkg/config"
	"ogomez/mkt-export/pkg/export"
	"ogomez/mkt-export/pkg/model"
	"ogomez/mkt-export/pkg/util"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

type ExcelReader struct {
	Conf    config.Config
	jsonExp export.JsonExporter
}

func NewExcelReader(conf config.Config) *ExcelReader {
	return &ExcelReader{
		Conf: conf,
	}
}

func (e ExcelReader) ReadInput() error {
	f, _ := excelize.OpenFile(e.Conf.Input)

	sheet := f.GetSheetName(0)
	var topicReg model.TopicReg

	eventsOutPath := fmt.Sprintf("%s/events/", e.Conf.Output)
	util.BuildPath(eventsOutPath)

	subsOutPath := fmt.Sprintf("%s/subscriptions/", e.Conf.Output)
	util.BuildPath(subsOutPath)

	rows, err := f.Rows(sheet)
	if err != nil {
		log.Println("Error processing sheet rows")
		return err
	}
	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			log.Println("Error reading row colums")
			return err
		}

		if (row != nil) && (row[0] != "topicName") {

			topiccreationDate := time.Now().UTC().Format(time.RFC3339)
			topicFormatData := parseInt("topicFormatData", row[2])
			topicType := parseInt("topicType", row[3])
			topicStatus := parseInt("topicStatus", row[4])
			topicConfidentialityData := parseInt("topicConfidentialityData", row[5])
			topicPartitions := parseInt("topicPartitios", row[6])
			topicTTL := parseInt("topicTTL", row[7])
			topicPlatform := parseInt("topicPlatform", row[8])
			topicCategory := parseInt("topicCategory", row[10])

			topic := &model.Topic{
				TopicName:                row[0],
				TopicDescription:         row[1],
				TopicFormatData:          topicFormatData,
				TopicCreationDate:        topiccreationDate,
				TopicType:                topicType,
				TopicStatus:              topicStatus,
				TopicConfidentialityData: topicConfidentialityData,
				TopicPartitions:          topicPartitions,
				TopicTTL:                 topicTTL,
				TopicPlatform:            topicPlatform,
				TopicCDCSourceTable:      row[9],
				TopicCategory:            topicCategory,
			}
			app := &model.Application{
				Appkey: row[11],
			}
			event := &model.Event{
				EventName:                row[12],
				EventDescription:         row[13],
				EventSchemaCompatibility: 6,
				Topic:                    *topic,
			}

			topicReg = model.TopicReg{
				Application: *app,
				Event:       *event,
			}
		}
		if topicReg.Application.Appkey != "" {
			log.Printf("Exporting %s Event Subscription for %s appkey on topic %s", topicReg.Event.EventName, topicReg.Application.Appkey, topicReg.Event.Topic.TopicName)
			e.jsonExp.Export(topicReg, eventsOutPath+topicReg.Event.Topic.TopicName)
		}
		if len(row) > 14 {
			prods := strings.Split(row[14], "\n")
			for _, prod := range prods {
				if prod != "" {
					prodSub := &model.Subscription{
						AppKey:    prod,
						TopicName: topicReg.Event.Topic.TopicName,
						SubsType:  0,
					}
					finalPath := fmt.Sprintf("%s/producer-%s-%s", subsOutPath, prodSub.TopicName, prodSub.AppKey)
					if prodSub.AppKey != "" && prod != "subtProducers" {
						log.Printf("Exporting PRODUCER subscription for %s appkey on Topic %s", prodSub.AppKey, prodSub.TopicName)
						e.jsonExp.Export(prodSub, finalPath)
					}
				}
			}
		}

		if len(row) > 15 {
			cons := strings.Split(row[15], "\n")
			for _, con := range cons {
				if con != "" {
					conSub := &model.Subscription{
						AppKey:    con,
						TopicName: topicReg.Event.Topic.TopicName,
						SubsType:  1,
					}
					finalPath := fmt.Sprintf("%s/consumer-%s-%s", subsOutPath, conSub.TopicName, conSub.AppKey)
					if conSub.AppKey != "" && con != "subtConsumers" {
						log.Printf("Exporting CONSUMER subscription for %s appkey on Topic %s", conSub.AppKey, conSub.TopicName)
						e.jsonExp.Export(conSub, finalPath)
					}
				}
			}
		}
	}
	if err = rows.Close(); err != nil {
		log.Println("error closing excel file")
		return err
	}
	log.Printf("Event registration files exported to: %s", eventsOutPath)
	log.Printf("Event Subscription registration file exported to: %s", subsOutPath)
	return nil
}

func parseInt(field string, value string) int {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("Invalid not integer value for field: %s  ", field)
		panic(err)
	}

	return intValue

}
