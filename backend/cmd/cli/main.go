package main

import (
 "fmt"
	"github.com/syncloud/platform/network"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Println("usage: ", os.Args[0], "")
		return
	}

	switch os.Args[1] {
	case "ipv6":
		ip, err := network.LocalIPv6()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(ip)
	case "ipv4":
		ip, err := network.LocalIPv4()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(ip)

	}

}
