package http

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

// http.Get
func Test_httpGet(t *testing.T) {
	resp, err := http.Get("http://www.baidu.com")
	if err != nil { // 如果你得到的是重定向错误，那它俩的值都是 non-nil
		fmt.Println(err)
		return
	}
	defer close(resp)
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
	defer close(resp)
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

// netstat /n|find /c "80
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
	// 主动调用Response.Body.Close()将连接关闭,但并不说明,就能关闭连接的，若是没有读取过body那么会有泄漏
	defer close(resp)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	fmt.Println(string(body))

	time.Sleep(2 * time.Second)
}

// Go 的最新版本将读取并丢弃数据的任务交给了用户，如果你不处理，HTTP 连接可能会直接关闭而非重用
// 如果程序大量重用 HTTP 长连接，你可能要在处理响应的逻辑代码中加入
// 使用这种方式就没有大量TIME_WAIT连接不被回收了,完成TCP响应体读取流程
func close(resp *http.Response) error {
	io.Copy(ioutil.Discard, resp.Body) // 手动丢弃读取完毕的数据
	if resp != nil {
		resp.Body.Close()
	}
	return nil
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
