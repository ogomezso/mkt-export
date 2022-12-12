package services

import "ogomez/mkt-export/pkg/config"

type Service interface {
	Export()
}

type ExportHandler struct {
	Services map[string]Service
}

const (
	EVENT_SERVICE   = "events"
)

func NewExportHandler(conf config.Config) *ExportHandler {
	services := make(map[string]Service)
	services[EVENT_SERVICE] = NewEventsService(conf) 
	return &ExportHandler{
		Services: services,
	}
}

func (exp *ExportHandler) BuildExport() {
	done := make(chan bool, len(exp.Services))
	for _, v := range exp.Services {
		go func(s Service) {
			s.Export()
			done <- true
		}(v)
	}
	for i := 0; i < len(exp.Services); i++ {
		<-done
	}
	close(done)
}
