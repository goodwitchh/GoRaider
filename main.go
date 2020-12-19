package main

import (
	"fmt"

	"github.com/Not-Cyrus/GoRaider/api"
	"github.com/Not-Cyrus/GoRaider/utils"
)

func main() {
	utils.Read()
	response := utils.SendRequest("GET", "https://discordapp.com/api/v7/users/@me", "application/json")
	if response == nil {
		return
	}
	if response.StatusCode() != 200 {
		panic("Your bot token is incorrect")
	}
	checkGuild := utils.SendRequest("GET", fmt.Sprintf("https://discord.com/api/v6/guilds/%s", utils.JsonData.GetStringBytes("GuildID")), "application/json")
	if checkGuild.StatusCode() != 200 {
		panic("The bot can't access that guild.")
	}

	api.Nuke()
}
