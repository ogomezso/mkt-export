package services

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"ogomez/mkt-export/pkg/client"
	"ogomez/mkt-export/pkg/config"
	"path/filepath"
	"strings"
	"time"
)

type RegistrationService struct {
	Config  config.Config
	RClient client.RestClient
}

type EventFiles struct {
	TopicName    string
	EventRegFile string
	SchemaFile   string
	ExampleFile  string
}

func NewRegistrationService(conf config.Config) *RegistrationService {
	restClient := client.New(conf.Marketplace.Mktplaceurl, conf.Marketplace.Credentials)

	return &RegistrationService{
		Config:  conf,
		RClient: *restClient,
	}
}

func (r *RegistrationService) Register() error {
	err := registerEvents(r)
	if err != nil {
		return err
	}

	err1 := registerSubscriptions(r)
	if err1 != nil {
		return err
	}

	return nil
}

func registerEvents(r *RegistrationService) error {
	eventsPath := r.Config.Input + "/events"
	filesByEvent, err := r.getFilePathsByEvent(eventsPath)
	if err != nil {
		return err
	}
	for _, files := range filesByEvent {
		if files.TopicName != "e" {
			log.Printf("Registering event for %s  topic on marketplace", files.TopicName)
			log.Println("-------------")
			log.Println(files.TopicName)
			log.Println("-------------")
			log.Println(files.EventRegFile)
			log.Println(files.SchemaFile)
			log.Println(files.ExampleFile)
			log.Println("-------------")
			reqBody, err := ioutil.ReadFile(files.EventRegFile)
			if err != nil {
				return err
			}
			multipartFiles := make(map[string]string)
			if files.SchemaFile != "" {
				multipartFiles["schema"] = files.SchemaFile
			}
			if files.ExampleFile != "" {
				multipartFiles["example"] = files.ExampleFile
			}
			headers := make(map[string]string)

			headers["Authorization"] = "Bearer " + r.RClient.Bearer
			headers["X-Clientid"] = r.Config.Marketplace.Appkey
			r.RClient.PostMultipart(r.Config.Marketplace.Mktplaceurl+"/v2/registry/aggregations/events", headers, reqBody, multipartFiles)
			time.Sleep(1 * time.Second)
		}
	}

	return nil
}

func (r *RegistrationService) getFilePathsByEvent(eventRootPath string) (map[string]*EventFiles, error) {

	filesByEvent := make(map[string]*EventFiles)
	err := filepath.Walk(eventRootPath,
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				log.Println(fmt.Sprint(err))
				return err
			}
			_, filename := filepath.Split(path)
			ext := filepath.Ext(filename)
			name := strings.Split(filename, ext)[0]
			topicName := strings.Split(name, "--")[0]
			if eventFile, ok := filesByEvent[topicName]; ok {
				switch {
				case strings.HasSuffix(name, "--schema"):
					eventFile.SchemaFile = eventRootPath + "/" + filename
          filesByEvent[topicName] = eventFile
          log.Printf("schema: filename: %s event file: %v filesByEvent: %v\n", filename, eventFile, filesByEvent[topicName])
				case strings.HasSuffix(name, "--example"):
					eventFile.ExampleFile = eventRootPath + "/" + filename
          filesByEvent[topicName] = eventFile
          log.Printf("example: filename: %s event file: %v filesByEvent: %v\n", filename, eventFile, filesByEvent[topicName])
        default:
          eventFile.EventRegFile = eventRootPath + "/" + filename
          filesByEvent[topicName] = eventFile
          log.Printf("event: filename: %s event file: %v filesByEvent: %v\n", filename, eventFile, filesByEvent[topicName])
				}
			} else {
				switch {
				case strings.HasSuffix(name, "--schema"):
					filesByEvent[topicName] = &EventFiles{
						TopicName:  topicName,
						SchemaFile: eventRootPath + "/" + filename,
					}
					log.Printf("file name - else schema: %s \n", filename)
				case strings.HasSuffix(name, "--example"):
					filesByEvent[topicName] = &EventFiles{
						TopicName:   topicName,
						ExampleFile: eventRootPath + "/" + filename,
					}
					log.Printf("topic name - else example: %s \n", topicName)

				default:
					filesByEvent[topicName] = &EventFiles{
						TopicName:    topicName,
						EventRegFile: eventRootPath + "/" + filename,
					}
				}
			}
			return nil
		})
	if err != nil {
		return nil, err
	}
  for _, ef := range filesByEvent {
    log.Printf("Event file: %v \n", ef)
  }
	return filesByEvent, nil
}

func registerSubscriptions(r *RegistrationService) error {
	subsPath := r.Config.Input + "/subscriptions"
	err := filepath.Walk(subsPath,
		func(path string, info fs.FileInfo, err error) error {
			_, filename := filepath.Split(path)
			if filename != "subscriptions" {
				reqBody, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				headers := make(map[string]string)
				headers["Authorization"] = "Bearer " + r.RClient.Bearer
				headers["X-ClientId"] = r.Config.Marketplace.Appkey
				headers["Content-type"] = "application/json"
				r.RClient.Post(r.Config.Marketplace.Mktplaceurl+"/v2/registry/subscriptions_events", reqBody, headers)
				time.Sleep(1 * time.Second)
			}
			return nil
		})
	if err != nil {
		return err
	}
	return nil
}
