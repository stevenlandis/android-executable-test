package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
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
	statusChan := make(chan ProcStatus)
	fileCount := 0
	for _, pidStr := range files {
		pid64, err := strconv.ParseInt(pidStr, 10, 0)
		if err != nil {
			continue
		}
		pid := int(pid64)
		go readStatus(pid, statusChan)
		fileCount++
	}
	statuses := make([]ProcStatus, fileCount)
	for i := 0; i < fileCount; i++ {
		statuses[i] = <-statusChan
	}
	close(statusChan)
	sort.Sort(ProcStatusList(statuses))
	for _, status := range statuses {
		fmt.Printf("%d\t%s\t%d\n", status.Pid, status.Name, status.Threads)
	}
}

type ProcStatus struct {
	Pid     int
	Name    string
	Threads int
}
type ProcStatusList []ProcStatus

func (l ProcStatusList) Len() int {
	return len(l)
}
func (l ProcStatusList) Less(i int, j int) bool {
	return l[i].Pid < l[j].Pid
}
func (l ProcStatusList) Swap(i int, j int) {
	l[i], l[j] = l[j], l[i]
}

func readStatus(pid int, out chan ProcStatus) {
	bts, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/status", pid))
	if err != nil {
		log.Fatal(err)
	}
	status := string(bts)
	lines := strings.Split(status, "\n")
	info := make(map[string]string)
	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) < 2 {
			continue
		}
		info[strings.TrimSpace(parts[0])] = strings.TrimSpace(strings.Join(parts[1:], ":"))
	}
	threads, err := strconv.ParseInt(info["Threads"], 10, 0)
	if err != nil {
		log.Fatal(err)
	}
	out <- ProcStatus{
		Pid:     pid,
		Name:    info["Name"],
		Threads: int(threads),
	}
}
