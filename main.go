package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"fmt"
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

func defaultMessage(w http.ResponseWriter, r *http.Request){
	var msg DefaultMsg
	msg.MSG = "Please enter the url"
	json.NewEncoder(w).Encode(msg)
}

func redirectToUrl(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	fmt.Println("Params ", params)
	fmt.Println(urls)
	for _, item := range urls {
		if item.ShortURL == params["url"] {
			var buffer bytes.Buffer
			buffer.WriteString("http://")
			buffer.WriteString(item.NormalURL)
			http.Redirect(w, r, buffer.String(), 302) //to do: seperate function
			return
		} else if item.NormalURL == params["url"]{
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	var mymsg DefaultMsg
	mymsg.MSG = "Redirect to Url and from redirectToUrl function"
	json.NewEncoder(w).Encode(mymsg)
}


func saveUrl(w http.ResponseWriter, r *http.Request){
	
	params := mux.Vars(r)
	w.Header().Set("Content_Type", "application/json")

	fmt.Println("saveUrl", params["url"])
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
	newUrl.ShortURL = strconv.Itoa(rand.Intn(1000000))
	urls = append(urls, newUrl)
	json.NewEncoder(w).Encode(newUrl)
}

func main(){

	r := mux.NewRouter()

	urls = append(urls, Url{NormalURL: "www.youtube.com", ShortURL: "90"})
	urls = append(urls, Url{NormalURL: "www.github.com", ShortURL: "44"})
	urls = append(urls, Url{NormalURL: "www.emberjs.com", ShortURL: "123"})

	r.HandleFunc("/", defaultMessage)
	r.HandleFunc("/{url}", redirectToUrl).Methods("GET")
	r.HandleFunc("/new/{url}", saveUrl).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", r))
}