package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const MAX = 5
const baseURL = "https://eluv.io/items/"

func MakeRequest(url string, ch chan<- string) []byte {
	client := &http.Client{}
	//start := time.Now()
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Y1JGMmR2RFpRc211MzdXR2dLNk1UY0w3WGpI")
	resp, _ := client.Do(req)
	//secs := time.Since(start).Seconds()
	body, _ := ioutil.ReadAll(resp.Body)
	//ch <- fmt.Sprintf("%.2f elapsed with response length: %d %s", secs, len(body), url)
	//ch <- fmt.Sprintf("%.2f, %d, %s, %d", secs, resp.StatusCode, err, len(body))
	return body
}

var responses map[string][]byte
var ids []string

func main() {
	if len(os.Args) != 2 {
		log.Fatal("You need provide a filename with all the ids, separated by newlines")
		os.Exit(-1)
	}
	//start := time.Now()
	ch := make(chan string)
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ids = append(ids, scanner.Text())
	}
	responses := make(map[string][]byte)
	guard := make(chan struct{}, MAX)

	for i, id := range ids {
		guard <- struct{}{}
		go func(n int, id string) {
			currURL := baseURL + id
			responses[id] = MakeRequest(currURL, ch)
			<-guard
		}(i, id)
	}
	for key, value := range responses {
		//can be changed to however you want the output stored
		fmt.Println("Key:", key, "Value:", string(value))
	}
}
