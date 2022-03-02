package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type EnvVars struct {
	TwillioSid      string
	TwillioKey      string
	FromPhoneNumber string
	ToPhoneNumber   string
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var envVars EnvVars
	line := readLine(reader)
	err := json.Unmarshal([]byte(line), &envVars)
	if err != nil {
		log.Fatal(err)
	}

	serve(envVars)
}

func readLine(reader *bufio.Reader) string {
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	line = strings.TrimSpace(line)
	return line
}

func serve(env EnvVars) {
	countChan := make(chan chan int)
	go runCountServer(countChan)
	go runTimer(env)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			if strings.Contains(r.Header.Get("accept"), "text/html") {
				writeFile(w, "public/index.html", "text/html")
				return
			} else if r.URL.Path == "/main.js" {
				writeFile(w, "main.js", "application/javascript")
				return
			} else if r.URL.Path == "/favicon.ico" {
				writeFile(w, "public/favicon.svg", "image/svg+xml")
				return
			}
		} else if r.Method == "POST" {
			c := make(chan int)
			countChan <- c
			cnt := <-c
			fmt.Fprintf(w, "count: %d", cnt)
			return
		}
		w.WriteHeader(404)
	})

	port := ":2224"
	fmt.Println("Server running on port" + port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func writeFile(w http.ResponseWriter, path string, contentType string) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("content-type", contentType+"; charset=utf_8")
	w.Write(content)
}

func runCountServer(c chan chan int) {
	count := 0
	for ch := range c {
		count++
		ch <- count
		close(ch)
	}
}

func runTimer(env EnvVars) {
	for {
		loc := time.FixedZone("UTC-8", -8*60*60)
		now := time.Now().UnixMilli()
		date := time.Date(2022, time.February, 26, 17, 0, 0, 0, loc)
		for date.UnixMilli() <= now {
			date = time.Date(
				date.Year(),
				date.Month(),
				date.Day()+7,
				date.Hour(),
				date.Minute(),
				date.Second(),
				date.Nanosecond(),
				loc,
			)
		}
		fmt.Println("sleeping until date:", date)
		time.Sleep(time.Duration(date.UnixMilli()-now)*time.Millisecond + 10*time.Minute)
		sendText(env)
	}
}

func sendText(env EnvVars) {
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + env.TwillioSid + "/Messages.json"
	msgData := url.Values{}
	msgData.Set("To", env.ToPhoneNumber)
	msgData.Set("From", env.FromPhoneNumber)
	msgData.Set("Body", "hi from fire")
	msgDataReader := strings.NewReader(msgData.Encode())

	// Create HTTP request client
	client := &http.Client{}
	req, err1 := http.NewRequest("POST", urlStr, msgDataReader)
	if err1 != nil {
		log.Fatal(err1)
	}
	req.SetBasicAuth(env.TwillioSid, env.TwillioKey)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Make HTTP POST request and return message SID
	resp, err2 := client.Do(req)
	if err2 != nil {
		log.Fatal(err2)
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Got:", data)
	} else {
		fmt.Println(resp.Status)
	}

}
