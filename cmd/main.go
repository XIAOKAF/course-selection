package main

import (
	"course-selection/api"
	"course-selection/dao"
)

func main() {
	dao.InitDB()
	api.InitEngine()
}
