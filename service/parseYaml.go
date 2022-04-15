package service

import (
	"course-selection/model"
	"github.com/storyicon/grbac"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// ParseRule 解析配置文件中的鉴权规则
func ParseRule() (grbac.Rules, error) {
	var rules []*grbac.Rule

	yamlFile, err := ioutil.ReadFile("ruleConfig.yaml")

	configMap := make(map[string]interface{})
	err = yaml.Unmarshal(yamlFile, configMap)
	if err != nil {
		return rules, err
	}

	roleInterface := configMap["role0"]
	roleValue, ok := roleInterface.(interface{})
	if ok {
		resourceMap, ok1 := roleValue.(map[interface{}]string)
		permissionMap, ok2 := roleValue.(map[interface{}][]string)
		if ok1 && ok2 {
			rules[0].ID = 0
			rules[0].Host = resourceMap["host"]
			rules[0].Path = resourceMap["path"]
			rules[0].Method = resourceMap["method"]
			rules[0].AuthorizedRoles = permissionMap["authorized_roles"]
			rules[0].ForbiddenRoles = permissionMap["forbidden_roles"]
			rules[0].AllowAnyone = false
		}
	}

	roleInterface = configMap["role1"]
	roleValue, ok = roleInterface.(interface{})
	if ok {
		resourceMap, ok1 := roleValue.(map[interface{}]string)
		permissionMap, ok2 := roleValue.(map[interface{}][]string)
		if ok1 && ok2 {
			rules[1].ID = 0
			rules[1].Host = resourceMap["host"]
			rules[1].Path = resourceMap["path"]
			rules[1].Method = resourceMap["method"]
			rules[1].AuthorizedRoles = permissionMap["authorized_roles"]
			rules[1].ForbiddenRoles = permissionMap["forbidden_roles"]
			rules[1].AllowAnyone = false
		}
	}

	roleInterface = configMap["role0"]
	roleValue, ok = roleInterface.(interface{})
	if ok {
		resourceMap, ok1 := roleValue.(map[interface{}]string)
		permissionMap, ok2 := roleValue.(map[interface{}][]string)
		if ok1 && ok2 {
			rules[2].ID = 0
			rules[2].Host = resourceMap["host"]
			rules[2].Path = resourceMap["path"]
			rules[2].Method = resourceMap["method"]
			rules[2].AuthorizedRoles = permissionMap["authorized_roles"]
			rules[2].ForbiddenRoles = permissionMap["forbidden_roles"]
			rules[2].AllowAnyone = false
		}
	}

	return rules, nil
}

// ParseSmsConfig 解析配置文件中的sms信息
func ParseSmsConfig(sms model.Message) (model.Message, error) {
	yamlFile, err := ioutil.ReadFile("smsConfig.yaml")
	if err != nil {
		return sms, err
	}
	configMap := make(map[string]interface{})
	err = yaml.Unmarshal(yamlFile, configMap)
	if err != nil {
		return sms, err
	}
	SMSInterface := configMap["SMS"]
	smsValue, ok := SMSInterface.(interface{})
	if ok {
		detailsValues, ok := smsValue.(map[interface{}]string)
		if ok {
			sms.SignId = detailsValues["secret_id"]
			sms.SecretKey = detailsValues["secret_key"]
			sms.AppId = detailsValues["app_id"]
			sms.AppKey = detailsValues["app_key"]
			sms.SignId = detailsValues["sign_id"]
			sms.TemplateId = detailsValues["template_id"]
			sms.Sign = detailsValues["sign"]
		}
	}
	return sms, nil
}
