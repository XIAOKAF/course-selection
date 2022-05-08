package service

import (
	"course-selection/model"
	feat "encoding/json"
	"github.com/storyicon/grbac"
	"io/ioutil"
	"os"
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
			rules[1].ID = 1
			rules[1].Host = resourceMap["host"]
			rules[1].Path = resourceMap["path"]
			rules[1].Method = resourceMap["method"]
			rules[1].AuthorizedRoles = permissionMap["authorized_roles"]
			rules[1].ForbiddenRoles = permissionMap["forbidden_roles"]
			rules[1].AllowAnyone = false
		}
	}

	roleInterface = configMap["role2"]
	roleValue, ok = roleInterface.(interface{})
	if ok {
		resourceMap, ok1 := roleValue.(map[interface{}]string)
		permissionMap, ok2 := roleValue.(map[interface{}][]string)
		if ok1 && ok2 {
			rules[2].ID = 2
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
	file, err := os.Open("config/smsConfig.json")
	if err != nil {
		return sms, err
	}
	fileByte, err := ioutil.ReadAll(file)
	if err != nil {
		return sms, err
	}
	err = json.Unmarshal(fileByte, &sms)
	if err != nil {
		return sms, err
	}
	return sms, nil
}
