package main

import (
	"fmt"

	"github.com/sfi2k7/blueconfigclient"
)

func main() {
	client := blueconfigclient.NewClient("http://localhost:7891", "1234")
	// fmt.Println(client.CreatePath("/peoples/bill"))
	fmt.Println(client.GetNodes("/peoples"))
	// fmt.Println(client.SetValue("/peoples/bill/city/bristol"))
	fmt.Println(client.SetValue("/peoples/bill/state/RI"))
	fmt.Println(client.GetValue("/peoples/bill/city"))
	fmt.Println(client.GetValue("/peoples/bill/state"))
	// fmt.Println(client.GetValue("/peoples/faisal/email"))
	// fmt.Println(client.SetValue("peoples/faisal/email/sfi2k7@gmail.com"))
	// fmt.Println(client.GetValue("/peoples/faisal/email"))
	return
	// value, err := client.GetValue("/peoples/faisal/name")
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }
	// fmt.Println("Value:", value)

	// response, err := client.SimpleGet("/apps/smsdos/databases/url/value")
	// response, err := client.SimpleGet("/apps/smsdos/databases/password/value")
	response, err := client.SimpleGet("/apps/servers/nginx/port")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Response:", response)

	port, err := client.GetValueInt("/apps/servers/nginx/port", 4040)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("port:", port)

	// var nginx = struct {
	// 	Port string `json:"port"`
	// }{}

	// err = client.ParseValues("/apps/servers/nginx", &nginx)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }
	// fmt.Println("nginx:", nginx)
}
