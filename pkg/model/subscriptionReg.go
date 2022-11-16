package model

type Subscription struct {
	AppKey        string `json:"appKey"`
	TopicName     string `json:"topicName"`
	SubsType      int    `json:"subtType"`
	ConsumerGroup string `json:"subtConsumerGroup"`
}
