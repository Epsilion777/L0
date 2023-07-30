package main

import (
	"L0/publisher"
	"log"
)

func main() {
	channelNats := "foo"
	log.Println("Start of publication")
	publisher.StartPublish(3, 15, channelNats)

}
