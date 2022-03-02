package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
)

func main() {
	file, err := os.Open("/proc")
	if err != nil {
		log.Fatal(err)
	}
	files, err := file.Readdirnames(-1)
	if err != nil {
		log.Fatal(err)
	}
	for _, pid := range files {
		_, err := strconv.ParseInt(pid, 10, 0)
		if err != nil {
			continue
		}
		bts, err := ioutil.ReadFile(fmt.Sprintf("/proc/%s/cmdline", pid))
		if err != nil {
			log.Fatal(err)
		}
		if bytes.Equal(bts, []byte("./serve"+"\x00")) {
			cmd := exec.Command("kill", pid)
			err := cmd.Run()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
