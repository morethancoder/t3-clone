package main

import (
	"log"

	"github.com/pocketbase/pocketbase"
)


func main() {
	  app := pocketbase.New()


	  err := app.Start()
	  if err != nil {
		  log.Fatal(err)
	  }
}
