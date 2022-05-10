package service

import (
	"course-selection/model"
	"encoding/json"
	"io/ioutil"
	"os"
)

// ParseSmsConfig 解析配置文件中的sms信息
func ParseSmsConfig() (model.Message, error) {
	var s model.Message
	file, err := os.Open("config/smsConfig.json")
	if err != nil {
		return s, err
	}
	fileByte, err := ioutil.ReadAll(file)
	if err != nil {
		return s, err
	}
	err = json.Unmarshal(fileByte, &s)
	if err != nil {
		return s, err
	}
	return s, nil
}

// ParseBucket 解析储存桶配置文件
func ParseBucket() (model.Bucket, error) {
	bucket := model.Bucket{}
	file, err := os.Open("config/bucketConfig.json")
	if err != nil {
		return bucket, err
	}
	fileByte, err := ioutil.ReadAll(file)
	if err != nil {
		return bucket, err
	}
	err = json.Unmarshal(fileByte, &bucket)
	if err != nil {
		return bucket, err
	}
	return bucket, nil
}
