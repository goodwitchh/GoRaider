package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
	"github.com/valyala/fastjson"
)

func GetData(url string, keys ...string) (*fastjson.Value, []string) {
	members, _ := SendRequest("GET", fmt.Sprintf(url, GuildID), "application/json", "", nil)

	jsonData, err := parser.Parse(members)
	if err != nil {
		panic(fmt.Sprintf("Couldn't parse member JSON: %s", err.Error()))
	}
	madeArray := make([]string, len(jsonData.GetArray()))

	for _, jsonValue := range jsonData.GetArray() {
		madeArray = append(madeArray, string(jsonValue.GetStringBytes(keys...)))
	}

	return jsonData, madeArray
}

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

	GuildID = string(JsonData.GetStringBytes("GuildID"))
}

func ReadData(data string) []string {

	jsonData, err := parser.Parse(data)
	if err != nil {
		return nil
	}

	array := jsonData.GetArray()
	var tempSlice = make([]string, len(array))

	for _, memberID := range jsonData.GetArray() {
		tempSlice = append(tempSlice, string(memberID.GetStringBytes()))
	}

	return tempSlice
}

func SendRequest(method, URL, ctype, proxy string, body []byte) (string, int) {
	client := &fasthttp.Client{}

	if len(proxy) > 0 {
		client = &fasthttp.Client{MaxConnsPerHost: 10000, MaxConnDuration: 10, Dial: fasthttpproxy.FasthttpHTTPDialerTimeout(proxy, time.Second*10)}
	}

	var (
		req *fasthttp.Request
		res *fasthttp.Response
	)

	req = fasthttp.AcquireRequest()

	defer fasthttp.ReleaseRequest(req)

	if body != nil {
		req.SetBody(body)
	}

	req.Header.SetMethod(method)
	req.Header.SetRequestURI(URL)
	req.Header.Set("Accept", "*/*")

	if len(ctype) != 0 {
		req.Header.Set("Content-Type", ctype)
	}

	if strings.HasPrefix(URL, "https://discordapp.com") || strings.HasPrefix(URL, "https://discord.com") {
		switch JsonData.GetBool("Bot") {
		case true:
			req.Header.Set("Authorization", "Bot "+string(JsonData.GetStringBytes("Token")))
		default:
			req.Header.Set("Authorization", string(JsonData.GetStringBytes("Token")))
		}
	}

	res = fasthttp.AcquireResponse()
	res.SetConnectionClose()

	err := client.DoTimeout(req, res, 10*time.Second)

	if err != nil {
		fasthttp.ReleaseResponse(res)
		return "{}", 0
	}

	return string(res.Body()), res.StatusCode()
}

var (
	GuildID  string
	JsonData *fastjson.Value
	parser   fastjson.Parser
	res      = fasthttp.AcquireResponse()
)
