package main

import (
	"choose-your-adventure/story"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main()  {
	port := flag.Int("port", 3000, "the port to start the Choose Your Own Adventure web application")
	fileName := flag.String("json", "gopher.json", "a JSON file with the Choose Your Own Adventure Story")
	flag.Parse()

	fmt.Printf("Using the story in %s\n", *fileName)

	file, err := os.Open(*fileName)

	if err != nil {
		panic(err)
	}

	stry, err := story.JsonStory(file)
	if err != nil {
		panic(err)
	}

	h := story.NewHandler(stry)
	fmt.Printf("Starting the server on port: %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
}