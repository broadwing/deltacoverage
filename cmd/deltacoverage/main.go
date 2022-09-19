package main

import (
	"fmt"
	"log"

	"github.com/broadwing/deltacoverage"
)

func main() {
	_, err := deltacoverage.GenerateCoverProfiles("./")
	if err != nil {
		log.Fatal(err)
	}
	c, err := deltacoverage.NewCoverProfile("./")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(c)
}
