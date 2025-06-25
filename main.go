package main

import (
	"fmt"
	"os"

	"thiennguyen.dev/welab-healthcare-app/cmd"
)

func main() {
	app := cmd.New()
	if err := app.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
