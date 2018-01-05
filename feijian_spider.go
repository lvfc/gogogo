package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func Httpget(url string) (contents string, statuscode int) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		contents = ""
		statuscode = -100
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		statuscode = resp.StatusCode
		return
	}

	statuscode = resp.StatusCode
	contents = string(body)
	return
}

func SpiderSecond(url string) (title string, content string) {
	content, status := Httpget(url)
	if status != 200 {
		fmt.Println("url = ", url, " second get err!!! statuscode = ", status)
		return
	}
	myreg_title := regexp.MustCompile(`<h1>(.*?)</h1>`)
	myreg_content := regexp.MustCompile(`<div id="content">(?s:(.*?))</div>`)
	sp_title := myreg_title.FindAllStringSubmatch(content, -1)
	for _, value := range sp_title {
		title = value[1]
	}
	sp_content := myreg_content.FindAllStringSubmatch(content, -1)
	for _, value := range sp_content {
		content = value[1]
	}

	return

}

func WriteFile(title string, content string) {
	myfile, err := os.OpenFile("./feijian.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(" open file error")
		return
	}
	content = strings.Replace(content, "<br/><br/>", "\n", -1)
	content = strings.Replace(content, "&nbsp", "", -1)
	content = strings.Replace(content, ";", "", -1)
	defer myfile.Close()
	myfile.WriteString(title)
	myfile.WriteString("\n")
	myfile.WriteString(content)
	myfile.WriteString("\n")

}

func SpiderFirst(url string) {
	contents, status := Httpget(url)
	if status != 200 {
		fmt.Println(" url = ", url, " get err!!! statuscode = ", status)
		return
	}
	//	var flag string
	myreg := regexp.MustCompile(`<dd><a href="(.*?)"`)
	content_html := myreg.FindAllStringSubmatch(contents, -1)
	for _, target := range content_html {
		var u_url string
		u_url = "http://www.biquke.com/bq/3/3714/" + string(target[1])
		content_title, content_cont := SpiderSecond(u_url)
		WriteFile(content_title, content_cont)
		/*		fmt.Println("请按任意键继续爬取下一页，或者输入exit退出....")
				fmt.Scanf("%s\n", &flag)
				if flag == "exit" {
					break
				}*/
	}

}

func Dowork() {
	url := "http://www.biquke.com/bq/3/3714/"
	SpiderFirst(url)
}

func main() {
	Dowork()
}
