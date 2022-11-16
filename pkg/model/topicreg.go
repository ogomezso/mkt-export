package model

type TopicReg struct {
	Application Application `json:"application"`
	Event       Event       `json:"event"`
}

type Application struct {
	Appkey string `json:"appKey"`
}

type Event struct {
	EventName                string `json:"eventName"`
	EventSchemaId            int    `json:"eventSchemaId"`
	EventSchemaVersion       int    `json:"eventSchemaVersion"`
	EventSchemaCompatibility int    `json:"eventSchemaCompatibility"`
	EventDescription         string `json:"eventDescription"`
	Topic                    Topic  `json:"topic"`
}

type Topic struct {
	TopicName                string `json:"topicName"`
	TopicDescription         string `json:"topicDescription"`
	TopicFormatData          int    `json:"topicFormatData"`
	TopicCreationDate        string `json:"topicCreationDate"`
	TopicType                int    `json:"topicType"`
	TopicStatus              int    `json:"topicStatus"`
	TopicConfidentialityData int    `json:"topicConfidentialityData"`
	TopicPartitions          int    `json:"topicPartitions"`
	TopicTTL                 int    `json:"topicTTl"`
	TopicPlatform            int    `json:"topicPlatform"`
	TopicCDCsourceTable      string `json:"topicCDCsourceTable"`
	TopicCategory            int    `json:"topicCategory"`
}
