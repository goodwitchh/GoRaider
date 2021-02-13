package main

import (
	"fmt"

	"github.com/Not-Cyrus/GoRaider/api"
	"github.com/Not-Cyrus/GoRaider/utils"
	"github.com/valyala/fastjson"
)

func main() {
	var (
		res  string
		stat int
	)

	utils.Read()
	_, stat = utils.SendRequest("GET", "https://discordapp.com/api/v7/users/@me", "application/json", "", nil)
	if stat != 200 {
		panic("Your bot token is incorrect")
	}
	res, stat = utils.SendRequest("GET", fmt.Sprintf("https://discord.com/api/v6/guilds/%s", utils.JsonData.GetStringBytes("GuildID")), "application/json", "", nil)
	if stat != 200 {
		panic("The bot can't access that guild.")
	}
	parsed, err := parser.Parse(res)
	if err != nil {
		panic("Couldn't read guild name")
	}
	api.Nuke(string(parsed.GetStringBytes("name")))
}

var (
	parser fastjson.Parser
)
