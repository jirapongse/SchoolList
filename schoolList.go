// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 241.

// Crawl2 crawls web links starting with the command-line arguments.
//
// This version uses a buffered channel as a counting semaphore
// to limit the number of concurrent calls to links.Extract.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

//!+sema
// tokens is a counting semaphore used to
// enforce a limit of 20 concurrent requests.
var tokens = make(chan struct{}, 20)
var block = make(chan struct{})

func crawl(url string) {
	//fmt.Println(url)
	tokens <- struct{}{} // acquire a token
	Extract(url)
	<-tokens // release the token

	block <- struct{}{}

}
func fileToLines(filename string) []string {
	contents, _ := ioutil.ReadFile(filename)
	return strings.Split((string(contents)), "\r\n")
}

func findString(filename string, text string) bool {
	contents, _ := ioutil.ReadFile(filename)
	strContents := string(contents)
	return strings.Contains(strContents, text)
}

func Extract(url string) {

	resp, err := http.Get(url)
	//	fmt.Printf("url: %s\n", url)
	if err != nil {
		log.Print(err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
	}

	b, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	//	str = strings.Replace(str, "\n","",-1)
	//	str = strings.Replace(str, "\r","",-1)

	if _, err := os.Stat("prev\\" + GetFileNameFromURL(url)); os.IsNotExist(err) {
		//	log.Print("New URL " + url)
		err := ioutil.WriteFile("prev\\"+GetFileNameFromURL(url), b, 0644)
		if err != nil {
			log.Print(err)

		}

	}
	err = ioutil.WriteFile("current\\"+GetFileNameFromURL(url), b, 0644)
	if err != nil {
		log.Print(err)

	}

	/*	f, err := os.Create(GetFileNameFromURL(url))
				if err != nil {
		    	log.Print(err)
		}
			_, err = f.Write(b)
					if err != nil {
		    	log.Print(err)
		}
			err = f.Close()
					if err != nil {
		    	log.Print(err)
		}*/

}

func GetFileNameFromURL(url string) string {
	str := strings.Replace(url, ":", "", -1)
	str = strings.Replace(str, "/", "", -1)
	str = strings.Replace(str, ".", "_", -1)
	return str
}
func main() {
	worklist := make(chan []string)

	var n int // number of pending sends to worklist

	// Start with the command-line arguments.

	// Start with the command-line arguments.

	go func() { worklist <- fileToLines("schoolList.txt") }()

	// Crawl the web concurrently.

	list := <-worklist
	for _, link := range list {
		//fmt.Printf("%s\n", link)
		n++
		go func(link string) {
			crawl(link)
			filename := GetFileNameFromURL(link)
			fi, err := os.Stat("prev\\" + filename)
			if err != nil {
				log.Print(err)
			} else {
				fi1, err := os.Stat("current\\" + filename)
				if err != nil {
					log.Print(err)
				} else {
					sizeChanged := fi1.Size() - fi.Size()
					if sizeChanged < 0 {
						sizeChanged *= -1
					}
					if sizeChanged >= 100 {
						//fmt.Printf("%s file size changed: %d\n", link, sizeChanged)
						if findString("current\\"+filename, "ครูอัตราจ้าง") {
							fmt.Printf("อัตราจ้าง %s file size changed: %d\n", link, sizeChanged)
						} else if findString("current\\"+filename, "พนักงานราชการ") {
							fmt.Printf("พนักงานราชการ %s file size changed: %d\n", link, sizeChanged)
						} else if findString("current\\"+filename, "รับสมัคร") {
							fmt.Printf("รับสมัคร %s file size changed: %d\n", link, sizeChanged)
						}

					}
				}
			}
		}(link)
	}

	for {
		<-block
		n--
		//fmt.Printf("n=%d\n", n)
		if n == 0 {
			break
		}
	}
	/*
		go func() { worklist <- fileToLines("schoolList.txt") }()

		list = <-worklist
		for _, link := range list {
			//fmt.Printf("%s\n", link)
			filename := GetFileNameFromURL(link)
			fi, err := os.Stat("prev\\" + filename)
			if err != nil {
				log.Print(err)
			} else {
				fi1, err := os.Stat("current\\" + filename)
				if err != nil {
					log.Print(err)
				} else {
					sizeChanged := fi1.Size() - fi.Size()
					if sizeChanged < 0 {
						sizeChanged *= -1
					}
					if sizeChanged >= 100 {
						//fmt.Printf("%s file size changed: %d\n", link, sizeChanged)
						if findString("current\\"+filename, "ครูอัตราจ้าง") {
							fmt.Printf("อัตราจ้าง %s file size changed: %d\n", link, sizeChanged)
						} else if findString("current\\"+filename, "พนักงานราชการ") {
							fmt.Printf("พนักงานราชการ %s file size changed: %d\n", link, sizeChanged)
						}

					}
				}
			}

		}
	*/
}

//!-
