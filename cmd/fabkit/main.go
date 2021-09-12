package main

import (
	"github.com/czar0/fabkit-cli/internal/spinner"
	"github.com/czar0/fabkit-cli/pkg/cmd"
	"log"
)

func main() {
	spinner.Init()
	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
