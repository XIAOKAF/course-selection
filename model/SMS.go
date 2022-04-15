package model

type Message struct {
	SecretId   string `yaml:"secret_id"`
	SecretKey  string `yaml:"secret_key"`
	AppId      string `yaml:"app_id"`
	AppKey     string `yaml:"app_key"`
	SignId     string `yaml:"sign_id"`
	TemplateId string `yaml:"template_id"`
	Sign       string `yaml:"sign"`
}
