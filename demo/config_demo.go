package main

import (
	"log"
	"zyj.com/golang-study/config"
)

func main() {
	err := config.Init()
	if err != nil {
		log.Printf("Failed to init config: %v", err)
		return
	}
}
