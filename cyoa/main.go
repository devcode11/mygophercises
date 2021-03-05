package main

import (
	"os"
	"log"
	"main/story"
	"net/http"
)


func main() {
	f, err := os.Open("gopher.json")
	fatalErr(err)
	stor, err := story.FromJson(f)
	fatalErr(err)
	err = f.Close()
	fatalErr(err)
	handler := story.StoryHandler{stor}
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func fatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
