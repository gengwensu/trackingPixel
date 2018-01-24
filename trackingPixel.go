/*
trackingPixel: an OpenGraph data scrapper.
*/
package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"flag"

	"github.com/dyatlov/go-opengraph/opengraph"
)

//AppContext contains all global variables
type AppContext struct {
	Cache       map[string]int // cache for url
	FileHandler *os.File       // file handler
}

const base64GifPixel = "R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7" //1x1 transparant gif

func main() {
	outFile := flag.String("outfile", "scrap.txt", "output file for scrape results")
	flag.Parse()
	urlCache := map[string]int{}
	
	//Open file for append
	f, err := os.OpenFile(*outFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	globalVar := AppContext{
		FileHandler: f,
		Cache:       urlCache,
	}

	log.Fatal(http.ListenAndServe("localhost:8081", &globalVar))
}

func (ds *AppContext) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/tracking", "/tracking/": // for monitoring purpose
		if req.Method == "GET" {
			fmt.Fprint(w, "TrackingPixel service\n") // return signature of the service
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}

	case "/tracking/track.gif":
		if req.Method == "GET" {
			output, _ := base64.StdEncoding.DecodeString(base64GifPixel)
			w.Header().Set("Content-Type", "image/gif")
			w.Write(output) // return with a 1x1 transparent gif
			defer req.Body.Close()
			// now take care of the scrapping
			urlString := req.Header.Get("Referer") //get the Referer's url
			if _, ok := ds.Cache[urlString]; ok {  //url in cache
				ds.Cache[urlString]++
			} else { // not in cache
				ds.Cache[urlString]++                 //cache it
				go scrapOG(urlString, ds.FileHandler) // run the scrapper in a go routine
			}
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}

	default:
		w.WriteHeader(http.StatusNotFound) // 404
	}
}

func scrapOG(url string, f *os.File) {
	resp, _ := http.Get(url)
	bytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	og := opengraph.NewOpenGraph()
	err := og.ProcessHTML(strings.NewReader(string(bytes)))
	if err != nil {
		log.Fatal(err)
	} else {
		_, _ = f.WriteString(url + " : ")
		_, _ = f.WriteString(og.String())
		_, _ = f.WriteString(" ***** ")
	}
}
