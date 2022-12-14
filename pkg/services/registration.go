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
		log.Printf("Error registering events: %s", fmt.Sprint(err))
		return err
	}

	err1 := registerSubscriptions(r)
	if err1 != nil {
		log.Printf("Error Registering subscriptions: %s", fmt.Sprint(err1.Error()))
		return err1
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
		log.Printf("Registering event for %s  topic on marketplace", files.TopicName)
		log.Println("-------------")
		log.Println(files.TopicName)
		log.Println("-------------")
		log.Println(files.EventRegFile)
		log.Println(files.SchemaFile)
		log.Println(files.ExampleFile)
		log.Println("-------------")

		multipartFiles := make(map[string]string)
		requestBody, err := ioutil.ReadFile(files.EventRegFile)
		if err != nil {
			return err
		}
		if files.SchemaFile != "" {
			multipartFiles["schema"] = files.SchemaFile
		}
		if files.ExampleFile != "" {
			multipartFiles["example"] = files.ExampleFile
		}
		headers := make(map[string]string)

		headers["Authorization"] = "Bearer " + r.RClient.Bearer
		headers["X-ClientId"] = r.Config.Marketplace.Appkey
		r.RClient.PostMultipart(r.Config.Marketplace.Mktplaceurl+"/v2/registry/aggregations/events", headers, requestBody, multipartFiles)
		time.Sleep(1 * time.Second)
	}

	return nil
}

func (r *RegistrationService) getFilePathsByEvent(eventRootPath string) (map[string]EventFiles, error) {

	filesByEvent := make(map[string]EventFiles)
	err := filepath.Walk(eventRootPath,
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				log.Println(fmt.Sprint(err))
				return err
			}
			_, filename := filepath.Split(path)
			ext := filepath.Ext(filename)
			name := strings.Split(filename, ext)[0]
			topicName := strings.Split(name, "-")[0]
			if eventFile, ok := filesByEvent[topicName]; ok {
				switch {
				case strings.HasSuffix(name, "-schema"):
					eventFile.SchemaFile = eventRootPath + "/" + filename
				case strings.HasSuffix(name, "-example"):
					eventFile.ExampleFile = eventRootPath + "/" + filename
				}
			} else {
				switch {
				case strings.HasSuffix(name, "-schema"):
					filesByEvent[topicName] = EventFiles{
						TopicName:  topicName,
						SchemaFile: eventRootPath + "/" + filename,
					}
				case strings.HasSuffix(name, "-example"):
					filesByEvent[topicName] = EventFiles{
						TopicName:   topicName,
						ExampleFile: eventRootPath + "/" + filename,
					}
				default:
					filesByEvent[topicName] = EventFiles{
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
	return filesByEvent, nil
}

func registerSubscriptions(r *RegistrationService) error {
	subsPath := r.Config.Input + "/subscriptions"
	err := filepath.Walk(subsPath,
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				log.Println(fmt.Sprint(err.Error()))
				return err
			}
			_, filename := filepath.Split(path)
			if filename != "subscriptions" {
				log.Printf("Registrando Subscripcion para: %s", filename)
				log.Println("------------------------------")
				subfile := subsPath + "/" + filename
				reqBody, err := ioutil.ReadFile(subfile)
				if err != nil {
					log.Println(fmt.Sprint(err))
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
