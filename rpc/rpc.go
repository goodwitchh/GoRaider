package rpc

import (
	"fmt"
	"time"

	"github.com/ananagame/rich-go/client"
)

func ChangeRPC(bans, channels, roles, proxies int, time time.Duration, serverName string) {
	err := client.Login("798025923173154887")
	if err != nil {
		fmt.Printf("Error with RPC[1]: %s", err.Error())
	}

	err = client.SetActivity(client.Activity{
		State:      fmt.Sprintf("in %s (In server: %s) on %d proxies", time, serverName, proxies),
		Details:    fmt.Sprintf("Got %d deleted channels | %d bans | %d deleted roles", channels, bans, roles),
		LargeImage: "large",
		LargeText:  "github.com/Not-Cyrus",
	})

	if err != nil {
		fmt.Printf("Error with RPC[2]: %s", err.Error())
	}
}
