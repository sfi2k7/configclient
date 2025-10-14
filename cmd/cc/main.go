package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfi2k7/configclient"
)

func main() {
	remoteUrl := flag.String("url", "https://blue-config-app.bluebbb.org", "")
	remoteTokne := flag.String("token", "1234", "")
	flag.Parse()

	args := os.Args

	if len(args) == 1 {
		fmt.Println("Blue Config - Version 0.1")
		return
	}

	var method string
	if len(args) > 1 {
		method = args[1]
		if method != "get" && method != "set" {
			fmt.Println("invalid method:", method)
			return
		}
	}

	if len(args) == 2 {
		fmt.Println("must provide path for method:", method)
		return
	}

	p := args[2]

	c := configclient.NewClient(*remoteUrl, *remoteTokne)
	fmt.Println("path:", p)
	response, err := c.SimpleGet(p)
	if err != nil {
		fmt.Println("error:", err.Error())
		return
	}

	if response.Error != "" {
		fmt.Println("error:", response.Error)
		return
	}

	fmt.Println(response.Result)
}
