package main

import (
	"log"

	"github.com/apaichon/thesis-itf/itf/internal/process"
)

func main() {
	log.Println("Process Manager Started")
	pm := process.NewProcessManager()
	pm.Start()
}
