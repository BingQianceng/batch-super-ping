package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
)

//将文件逐行读取至数组
func Getdomainlist() []string {
	file, err := os.Open("./domain.txt") // For read access.
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	data := make([]byte, 100)

	count, err := file.Read(data)
	if err != nil {
		log.Fatal(err)
	}

	domainlist := data[:count]
	return strings.Fields(string(domainlist))
}

// []string 去重
func RemoveDuplicate(list []string) []string {
	sort.Strings(list)
	i := 0
	var newlist = []string{""}
	for j := 0; j < len(list); j++ {
		if strings.Compare(newlist[i], list[j]) == -1 {
			newlist = append(newlist, list[j])
			i++
		}
	}
	return newlist
}

func Getnodelist() []string {
	client := &http.Client{}
	formValues := url.Values{}
	formValues.Set("node", "1,2,3,4,5,6")
	formValues.Set("host", "www.baidu.com")
	formValues.Set("csrfmiddlewaretoken", "QXNZxN49mxYTEkHnobxd4LhwZDR2uvWxpfxLUeO91myz7ErRCbiNizInk8vhsk8w")
	formDataStr := formValues.Encode()
	formDataBytes := []byte(formDataStr)
	formBytesReader := bytes.NewReader(formDataBytes)
	//fmt.Printf("formBytesReader: %s\n", formBytesReader)

	req, _ := http.NewRequest("POST", "http://www.wepcc.com", formBytesReader)
	req.Header.Add("Cookie", "csrftoken=AQTGFGiwH27Kp2nyMZPrWD2FTmxjHLQg98Ds272wmRHqSm720ZA1artweRbyFA2f")
	req.Header.Add("Origin", "http://www.wepcc.com")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
	req.Header.Add("Referer", "http://www.wepcc.com/")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	getnode := regexp.MustCompile(`.{8}-.{4}-.{4}-.{4}-.{12}`)
	nodelist := getnode.FindAllString(string(body), -1)
	return nodelist
}

func Getiplist(node string, domain string) (string, error) {
	client := &http.Client{}
	formValues := url.Values{}
	formValues.Set("node", node)
	formValues.Set("host", domain)
	formValues.Set("csrfmiddlewaretoken", "QXNZxN49mxYTEkHnobxd4LhwZDR2uvWxpfxLUeO91myz7ErRCbiNizInk8vhsk8w")
	formDataStr := formValues.Encode()
	formDataBytes := []byte(formDataStr)
	formBytesReader := bytes.NewReader(formDataBytes)
	//fmt.Printf("formBytesReader: %s\n", formBytesReader)

	req, _ := http.NewRequest("POST", "http://www.wepcc.com/check-ping.html", formBytesReader)
	req.Header.Add("Cookie", "csrftoken=AQTGFGiwH27Kp2nyMZPrWD2FTmxjHLQg98Ds272wmRHqSm720ZA1artweRbyFA2f")
	req.Header.Add("Origin", "http://www.wepcc.com")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
	req.Header.Add("Referer", "http://www.wepcc.com/")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	//fmt.Printf("string(body): %v\n", string(body))
	getip := regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`)
	ip := getip.FindAllString(string(body), -1)
	return ip[0], nil
}

func main() {
	domainlist := Getdomainlist()
	nodelist := Getnodelist()

	for _, domain := range domainlist {
		var ips []string

		for _, node := range nodelist {
			ip, err := Getiplist(node, domain)
			if err != nil {
				continue
			}
			ips = append(ips, ip)
		}

		newips := RemoveDuplicate(ips)
		fmt.Printf("%v: %v", domain, newips)
		if len(newips) > 2 {
			fmt.Println("    USED CDN")
		} else {
			fmt.Println("")
		}
	}

}
