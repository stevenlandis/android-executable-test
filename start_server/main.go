package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {

	logFile, err3 := os.Create("some.log")
	if err3 != nil {
		log.Fatal(err3)
	}

	cmd := exec.Command("./serve")
	stdin, err1 := cmd.StdinPipe()
	if err1 != nil {
		log.Fatal(err1)
	}
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	err2 := cmd.Start()
	if err2 != nil {
		log.Fatal(err2)
	}

	stdin.Write([]byte(os.Args[1] + "\n"))
	stdin.Close()
	os.Stdin.Close()
	os.Stdout.Close()
	os.Stderr.Close()
}
