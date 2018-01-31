package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"bytes"
	"github.com/gorilla/mux"
)

type Url struct {
	NormalURL string `json:"normalurl"`
	ShortURL string `json:"shorturl"`
}

type DefaultMsg struct {
	MSG string `json:"msg"`
}

var urls []Url

func createMessage(msg string) DefaultMsg {
	var message DefaultMsg
	message.MSG = msg
	return message
}

func defaultMessage(w http.ResponseWriter, r *http.Request){
	mymsg := createMessage("Please enter the url")
	json.NewEncoder(w).Encode(mymsg)
}

func getRandomNum(spectrum int) string {
	return strconv.Itoa(rand.Intn(spectrum))
}

func redirectToUrl(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	for _, item := range urls {
		if item.ShortURL == params["url"] {
			var buffer bytes.Buffer
			buffer.WriteString("http://")
			buffer.WriteString(item.NormalURL)
			http.Redirect(w, r, buffer.String(), 302)
			return
		} else if item.NormalURL == params["url"]{
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	mymsg := createMessage("There is no website like this in our dataset")
	json.NewEncoder(w).Encode(mymsg)
}

func saveUrl(w http.ResponseWriter, r *http.Request){
	
	params := mux.Vars(r)

	for _, item := range urls {
		if item.NormalURL == params["url"] {
			json.NewEncoder(w).Encode(item)
			return
		} else if item.ShortURL == params["url"]{
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	var newUrl Url
	newUrl.NormalURL = params["url"]	
	newUrl.ShortURL = getRandomNum(100000)
	urls = append(urls, newUrl)
	json.NewEncoder(w).Encode(newUrl)
}

func main(){

	r := mux.NewRouter()

	urls = append(urls, Url{NormalURL: "www.github.com", ShortURL: "12345"})
	urls = append(urls, Url{NormalURL: "www.emberjs.com", ShortURL: "56789"})
	urls = append(urls, Url{NormalURL: "www.reactjs.org", ShortURL: "98765"})

	r.HandleFunc("/", defaultMessage)
	r.HandleFunc("/{url}", redirectToUrl).Methods("GET")
	r.HandleFunc("/new/{url}", saveUrl).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", r))
}