package api

import (
	"course-selection/model"
	"course-selection/service"
	"course-selection/tool"
	"fmt"
	"github.com/Chain-Zhang/pinyin"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func spiderMan(ctx *gin.Context) {
	client := http.Client{}
	req, err := http.NewRequest("GET", "http://jwzx.cqupt.edu.cn/kebiao/kb_stuList.php?jxb=A04212A1110022017", nil)
	if err != nil {
		fmt.Println("网页请求失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	//加入请求头
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cache-Control", "no-store")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36 Edg/100.0.1185.44")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Referer", "http://jwzx.cqupt.edu.cn/kebiao/kb_stu.php?xh=2021212196")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("cookie", "mLvnBZTNP4mtS=565PBR7jNqVZ5YjbZ63CBI9r1yzWMqDldJPIZW5JwfGupX9UaXLaXP4HRep2I7vGK8FU65eSLzC_9OLAIcYIRBq; PHPSESSID=ST-528345-o3-HgQ-F7EwRywjn2lQzAaYdFLQauthserver1; mLvnBZTNP4mtT=DxluYf8_dMOvA2UHBC4KvjYag55Tu7ZvUYC7Xzwpz_E5BrsZI1UtWKlI1epNh2d.3X97MCYW2JJNZOWJvd_sjIgLPPrcotjDRrEeWUJJqzkeiULmT5S.xpc5wyfYfea.LX9DVSqWdJyCFWFIgBVFBu.o6wOYNUK9SWJvIBa387kyWyW9yvFjxcwhKse7EWVSnJ8so1TeyfH3WgrEXD8y0Z1lkj7fUlYf1DsvfcJIDwiFGEAHSHgiPvt6I6RQjoh3BifeUyFFpvTNXgaDZPACpeoG6mbiAeKeSL2BIst7fYCDXke5GM1_CYgZ3NMUQALL2Zq9mqiIqSg2UndLE9fvjpkDl9gGucrgX0Z1aEXIP59")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("网络请求失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	defer resp.Body.Close()
	//解析网页
	docDetails, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("网页解析失败", err)
		tool.Failure(ctx, 500, "服务器错误")
		return
	}
	student := model.Student{
		RuleId: "1",
	}
	for j := 1; j <= 128; j++ {
		//#stuListTabs-current > table > tbody > tr:nth-child(1)  所有信息
		//#stuListTabs-current > table > tbody
		//#stuListTabs-current > table > tbody > tr:nth-child(1) > td:nth-child(2)  学号
		//#stuListTabs-current > table > tbody > tr:nth-child(1) > td:nth-child(3)  姓名
		//#stuListTabs-current > table > tbody > tr:nth-child(1) > td:nth-child(4)  性别
		//#stuListTabs-current > table > tbody > tr:nth-child(1) > td:nth-child(9)  年级
		//#stuListTabs-current > table > tbody > tr:nth-child(1) > td:nth-child(5)  班级
		//#stuListTabs-current > table > tbody > tr:nth-child(1) > td:nth-child(7)	专业
		//#stuListTabs-current > table > tbody > tr:nth-child(1) > td:nth-child(8)	院系
		k := strconv.Itoa(j)
		docDetails.Find("#stuListTabs-current > table > tbody > tr:nth-child(" + k + ")").
			Each(func(i int, s *goquery.Selection) {
				student.StudentId = s.Find("td:nth-child(2)").Text()
				student.StudentName = s.Find("td:nth-child(3)").Text()
				student.Password, err = pinyin.New(student.StudentName).Split("").Mode(pinyin.WithoutTone).Convert()
				if err != nil {
					fmt.Println("将汉字转化为拼音失败", err)
					tool.Failure(ctx, 500, "服务器错误")
					return
				}
				sex := s.Find("td:nth-child(4)").Text()
				student.Class = s.Find("td:nth-child(5)").Text()
				student.Major = s.Find("td:nth-child(7)").Text()
				student.Department = s.Find("td:nth-child(8)").Text()
				student.Grade = s.Find("td:nth-child(9)").Text()
				err = service.SpiderMan(student)
				if err != nil {
					fmt.Println("将学生信息导入MySQL失败", err)
					tool.Failure(ctx, 500, "服务器错误")
					return
				}
				err = service.HashSet("student", student.StudentId, "12345678900")
				if err != nil {
					fmt.Println("储存学生信息失败")
					tool.Failure(ctx, 500, "服务器错误")
					return
				}
				err = service.HashSet(student.StudentId, "studentName", student.Password)
				if err != nil {
					fmt.Println("储存学生姓名失败")
					tool.Failure(ctx, 500, "服务器错误")
					return
				}
				err = service.HashSet(student.StudentId, "password", student.Password)
				if err != nil {
					fmt.Println("储存学生密码失败")
					tool.Failure(ctx, 500, "服务器错误")
					return
				}
				err = service.HashSet(student.StudentId, "gender", sex)
				if err != nil {
					fmt.Println("储存学生性别失败")
					tool.Failure(ctx, 500, "服务器错误")
					return
				}
				err = service.HashSet(student.StudentId, "class", student.Class)
				if err != nil {
					fmt.Println("储存学生班级失败")
					tool.Failure(ctx, 500, "服务器错误")
					return
				}
				err = service.HashSet(student.StudentId, "department", student.Department)
				if err != nil {
					fmt.Println("储存学生院系失败")
					tool.Failure(ctx, 500, "服务器错误")
					return
				}
				err = service.HashSet(student.StudentId, "major", student.Major)
				if err != nil {
					fmt.Println("储存学生专业失败")
					tool.Failure(ctx, 500, "服务器错误")
					return
				}
				err = service.HashSet(student.StudentId, "roleLevel", "student")
				if err != nil {
					fmt.Println("储存学生权限等级错误", err)
					tool.Failure(ctx, 500, "服务器错误")
					return
				}
			})
	}
	tool.Success(ctx, http.StatusOK, "导入学生数据成功")
}
