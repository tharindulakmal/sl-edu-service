package models

type SubTopic struct {
	SubTopicID   int    `json:"subTopicId"`
	TopicID      int    `json:"topicId"`
	SubTopicName string `json:"subTopicName"`
}

type Topic struct {
	TopicID       int        `json:"topicId"`
	TopicName     string     `json:"topicName"`
	MainTopicName string     `json:"mainTopicName"`
	SubTopicList  []SubTopic `json:"subTopicList"`
}

type SmartNote struct {
	SubTopicName    string `json:"subTopicName"`
	ImageDefUrl     string `json:"imageDefUrl"`
	Definition      string `json:"definition"`
	Theory          string `json:"theory"`
	ImageTheoryUrl  string `json:"imageTheoryUrl"`
	Example         string `json:"example"`
	ImageExampleUrl string `json:"imageExampleUrl"`
}

type TopicsResponse struct {
	Topics           []Topic   `json:"topics"`
	DefaultSmartNote SmartNote `json:"defaultSmartNote"`
}
