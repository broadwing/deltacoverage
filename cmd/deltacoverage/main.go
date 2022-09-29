package main

import (
	"fmt"
	"log"

	"github.com/broadwing/deltacoverage"
)

func main() {
	c, err := deltacoverage.NewCoverProfile("./")
	if err != nil {
		log.Fatal(err)
	}
	err = c.Generate()
	if err != nil {
		log.Fatal(err)
	}
	err = c.Parse()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(c)
	err = c.Cleanup()
	if err != nil {
		log.Fatal(err)
	}
}
