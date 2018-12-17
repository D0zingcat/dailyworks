package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Config struct {
	Users  []User `json:"users"`
	Port   string `json:"port"`
	Source string `json:"source"`
	Dir    string `json:"dir"`
	Cert   string `json:"cert"`
	Key    string `json:"key"`
	Safe   string `json:"safe"`
}

const (
	CONFIG_NAME = "config.json"
)

var dir string
var source string
var personal string
var config Config

func main() {
	// file, err := os.OpenFile(CONFIG_NAME, os.O_RDONLY, 0644)
	// if err != nil {
	// 	log.Println("fail to open file: ", CONFIG_NAME)
	// 	panic(err)
	// }
	bytesConfig, err := ioutil.ReadFile(CONFIG_NAME)
	if err != nil {
		log.Println("fail to open file: ", CONFIG_NAME)
		panic(err)
	}
	json.Unmarshal(bytesConfig, &config)
	cert := config.Cert
	key := config.Key
	port := config.Port
	dir = config.Dir
	source = config.Source
	personal = config.Safe
	http.HandleFunc("/", handle)
	http.ListenAndServe(":"+port, nil)
	http.ListenAndServeTLS(":"+port, cert, key, nil)
}

func handle(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request of %v from %v, Ip is %v\n", r.URL.Path, r.RemoteAddr, r.Header.Get("X-Real-IP"))
	r.RemoteAddr = strings.Replace(r.RemoteAddr, "[::1]", "localhost", -1)
	if !strings.HasPrefix(r.RemoteAddr, source) {
		http.Error(w, "Access Unauthorized!", http.StatusUnauthorized)
		return
	}
	log.Println(strings.HasPrefix(r.URL.Path, personal))
	log.Println(r.URL.Path, personal)
	if strings.HasPrefix(r.URL.Path, personal) {
		if !checkAuth(w, r) {
			w.Header().Set("WWW-Authenticate", `Basic realm="MY REALM"`)
			http.Error(w, "Access unauthorized! You are visiting personal space, should enter user and password first!", http.StatusUnauthorized)
			return
		}
	}
	http.FileServer(http.Dir((dir))).ServeHTTP(w, r)
}

func checkAuth(w http.ResponseWriter, r *http.Request) bool {
	usernamePair := strings.Split(r.URL.Path, "/")
	username := usernamePair[2]
	log.Println(usernamePair, len(usernamePair))
	if username == "" {
		return true
	}
	s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	log.Println(s)
	if len(s) != 2 {
		return false
	}
	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		log.Println("base64 decoding err! ", err)
		return false
	}
	pair := strings.SplitN(string(b), ":", 2)
	log.Println(pair, "pair")
	if len(pair) != 2 {
		log.Println("split pair err! ", pair)
		return false
	}
	if username != pair[0] {
		log.Println(username, pair[0])
		return false
	}
	users := config.Users
	for i := range users {
		if users[i].Username == username {
			hash := fmt.Sprintf("%x", md5.Sum([]byte(pair[1])))
			if hash != users[i].Password {
				log.Println(pair[1], hash)
				return false
			}
			return true
		}
	}
	return false
}
