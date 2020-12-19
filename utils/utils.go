package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fastjson"
)

func Read() {
	file, err := os.Open("LoginInfo.json")
	if err != nil {
		panic(fmt.Sprintf("Couldn't read config.json: %s", err.Error()))
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(fmt.Sprintf("Couldn't read file data: %s", err.Error()))
	}
	JsonData, err = parser.Parse(string(data))
	if err != nil {
		panic(fmt.Sprintf("Couldn't parse json data: %s", err.Error()))
	}
}

func SendRequest(method, url, ctype string) *fasthttp.Response {
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(method)
	req.Header.SetRequestURI(url)
	if len(ctype) != 0 {
		req.Header.SetContentType(ctype)
	}
	req.Header.Set("Authorization", string(JsonData.GetStringBytes("Token")))
	resp := fasthttp.AcquireResponse()

	err := client.DoTimeout(req, resp, 10*time.Second)
	if err != nil {
		fmt.Println(fmt.Sprintf("A fuckup happened: %s", err.Error()))
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
		return nil
	}
	return resp
}

var (
	JsonData *fastjson.Value
	parser   fastjson.Parser
	client   = &fasthttp.Client{
		MaxConnsPerHost:     2000,
		MaxIdleConnDuration: 30 * time.Second,
	}
)
