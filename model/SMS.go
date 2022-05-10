package model

type Message struct {
	SecretId   string `json:"secret_id"`
	SecretKey  string `json:"secret_key"`
	AppId      string `json:"app_id"`
	AppKey     string `json:"app_key"`
	SignId     string `json:"sign_id"`
	TemplateId string `json:"template_id"`
	Sign       string `json:"sign"`
}
