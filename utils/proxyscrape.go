package utils

import (
	"regexp"
	"sync"
	"time"
)

func GetProxies() (goodProxies []string) {
	resultChannel := make(chan ProxyResult)
	resultWaitGroup := new(sync.WaitGroup)

	resultWaitGroup.Add(1)

	allProxies := scrapeProxies()

	go func() {
		defer resultWaitGroup.Done()
		handleProxies(resultChannel, len(allProxies), &goodProxies)
	}()

	for _, proxy := range allProxies {
		go testProxy(proxy, resultChannel)
	}

	resultWaitGroup.Wait()

	return
}

func handleProxies(c chan ProxyResult, proxyAmount int, goodProxies *[]string) {
	for i := 0; i < proxyAmount; i++ {
		result := <-c

		switch {

		case !result.Sucess:
			continue

		case result.Sucess && result.StatusCode >= 200 && result.StatusCode < 400:
			*goodProxies = append(*goodProxies, result.ProxyIP)

		default:
			*goodProxies = append(*goodProxies, result.ProxyIP)

		}
	}

	close(c)
}

func scrapeProxies() (allProxies []string) {

	for _, url := range urls {
		response, _ := SendRequest("GET", url, "", "", nil)
		if response == "{}" {
			return
		}
		proxiesList := proxyRegex.FindAllString(string(response), -1)
		for _, proxy := range proxiesList {
			allProxies = append(allProxies, proxy)
		}
	}

	return
}

func testProxy(proxy string, c chan ProxyResult) {
	var (
		statusCode int
		sucess     bool = true
	)

	start := time.Now()

	_, statusCode = SendRequest("GET", "https://www.google.com", "", proxy, nil)
	if statusCode == 0 {
		sucess = false
		statusCode = -1
	}

	end := time.Now().Sub(start)

	c <- ProxyResult{ProxyIP: proxy, ProxySpeed: end, StatusCode: statusCode, Sucess: sucess}
}

type (
	ProxyResult struct {
		ProxyIP    string
		ProxySpeed time.Duration
		StatusCode int
		Sucess     bool
	}
)

var (
	proxyRegex = regexp.MustCompile(`(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?).){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?):([0-9]){1,4}`)
	urls       = []string{
		"https://api.proxyscrape.com/?request=getproxies&proxytype=http&timeout=10000&country=all&ssl=all&anonymity=all",
		"https://proxylistfree.net/All-free-proxy-list?page=1",
		"https://www.proxy-list.download/api/v1/get?type=https",
		"https://proxylist.icu/proxy/1",
		"https://free-proxy-list.net",
		"https://www.sslproxies.org",
		"https://socks-proxy.net",
		"https://us-proxy.org",
		"https://spys.one/",
	}
)
