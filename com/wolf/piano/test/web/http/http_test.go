package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

// http.Get
func Test_httpGet(t *testing.T) {
	resp, err := http.Get("http://www.baidu.com")
	if err != nil {
		fmt.Println(err)
		return
	}
	//defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("11111" + string(body))

	time.Sleep(2 * time.Second)
}

func httpPost() {
	resp, err := http.Post("http://www.baidu.com",
		"application/x-www-form-urlencode",
		strings.NewReader("name=abc")) // Content-Type post请求必须设置
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

// http.Client
func Test_httpDo(t *testing.T) {
	client := &http.Client{}
	req, err := http.NewRequest("POST",
		"http://www.baidu.com",
		strings.NewReader("name=abc"))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	fmt.Println(string(body))

	time.Sleep(2 * time.Second)
}

// http.Client
func Test_httpGET(t *testing.T) {
	resp, err := http.Get("http://xxx/addrs?sv=210&p=qqq&lan=java&label=session.huitian")
	if err != nil {
		fmt.Println(err)
		return
	}
	//defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

	time.Sleep(2 * time.Second)
}
