package main

import (
	"fmt"

	"github.com/minmax1996/aoimdb/client"
)

func main() {
	client, err := client.NewTcpClient(":1593")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	client.AuthWithUserPassPair("admin", "pass")

	if resp, err := client.Get("2.avatar:111222"); err != nil {
		fmt.Println("no avatar:(")
		if _, err := client.Set("2.avatar:111222", "https://example.com"); err != nil {
			fmt.Println("no set")
		}
	} else {
		fmt.Println("userAvatar: " + resp.Value.(string))
	}

	if resp, err := client.Get("2.avatar:111222"); err != nil {
		fmt.Println("no avatar:(")
	} else {
		fmt.Println("now userAvatar: " + resp.Value.(string))
	}

}
