package api

import (
	"course-selection/model"
	"course-selection/service"
	"course-selection/tool"
	"github.com/Chain-Zhang/pinyin"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func spiderMan(ctx *gin.Context) {
	client := http.Client{}
	req, err := http.NewRequest("GET", "http://jwzx.cqupt.edu.cn/kebiao/kb_stuList.php?jxb=A04212A1110022017", nil)
	tool.DealWithErr(ctx, err, "网页请求错误")
	//加入请求头
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cache-Control", "no-store")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36 Edg/100.0.1185.44")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Referer", "http://jwzx.cqupt.edu.cn/kebiao/kb_stu.php?xh=2021212196")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("cookie", "mLvnBZTNP4mtS=565PBR7jNqVZ5YjbZ63CBI9r1yzWMqDldJPIZW5JwfGupX9UaXLaXP4HRep2I7vGK8FU65eSLzC_9OLAIcYIRBq; PHPSESSID=ST-1068784-Q9AAJNzmPAxgNQ-zKsUuEjVzSiQauthserver1; mLvnBZTNP4mtT=kyZvtP4.CQl7aJLMPUHHmfF965j3JBZfBlZPwQDxm1dcVcYEICf_yaORLz9Ynuo5BFkXlYoO_fr8YO_qtBodkfIQxo7V2OlW.f9TOmofLrkSO36Ui2eWQK0SS5rWvHvbYPtt8Tj0WDS095nEz0R5G2Flb9dZ.8M8bDVpiAizVB1G8TWo7pSIC1vMrgL6n2rfqa7dPtxdEBNjcavcqM2y6O20K0UNpnehriVmh6ZC4a.iYS_7VFV_PRT_n8kCARI.TKACepeTCHSynaQSn9FTFuxDH3pYKEpHE669RJkfUnRbLgoUX2nJ6IgacpXOyHbcoviwCb4VsKk4KqAEq3MUarU2VsI9ld3LTkgsDdbJoQ9")
	resp, err := client.Do(req)
	tool.DealWithErr(ctx, err, "网络请求失败")
	defer resp.Body.Close()
	//解析网页
	docDetails, err := goquery.NewDocumentFromReader(resp.Body)
	tool.DealWithErr(ctx, err, "解析网页错误")
	student := model.Student{}
	for j := 1; j <= 128; j++ {
		//#stuListTabs-current > table > tbody > tr:nth-child(1)  所有信息
		//#stuListTabs-current > table > tbody > tr:nth-child(1) > td:nth-child(2)  学号
		//#stuListTabs-current > table > tbody > tr:nth-child(1) > td:nth-child(3)  姓名
		//#stuListTabs-current > table > tbody > tr:nth-child(1) > td:nth-child(4)  性别
		//#stuListTabs-current > table > tbody > tr:nth-child(1) > td:nth-child(7)  班级
		//#stuListTabs-current > table > tbody > tr:nth-child(1) > td:nth-child(8)	专业
		//#stuListTabs-current > table > tbody > tr:nth-child(1) > td:nth-child(9)	院系
		k := strconv.Itoa(j)
		docDetails.Find("#stuListTabs-current > table > tbody > tr:nth-child(" + k + ")").
			Each(func(i int, s *goquery.Selection) {
				student.StudentId = j
				student.StudentName = s.Find("td:nth-child(3)").Text()
				student.Password, err = pinyin.New(student.StudentName).Split("").Mode(pinyin.WithoutTone).Convert()
				tool.DealWithErr(ctx, err, "将汉字转化为拼音出错")
				sex := s.Find("td:nth-child(4)").Text()
				if sex == "女" {
					student.Gender = 2
				} else {
					student.Gender = 1
				}
				student.Class = s.Find("td:nth-child(7)").Text()
				student.Major = s.Find("td:nth-child(8)").Text()
				student.Department = s.Find("td:nth-child(9)").Text()
				err = service.SpiderMan(student)
				tool.DealWithErr(ctx, err, "将学生数据导入MySQL失败")
			})
	}
	tool.Success(ctx, http.StatusOK, "导入学生数据成功")
}
