package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/alexkylt/message_broker/internal/pkg/storage"
)

// Config ... - structure for host, port and storage
type Config struct {
	host    string
	port    int
	storage storage.StrgInterface
}

type key struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type pattern struct {
	Pattern string `json:"pattern"`
}

var keymap key

var patternmap pattern

// InitServer : initialize server with the appropriate storage
func InitServer(port int, storage storage.StrgInterface) *Config {
	server := &Config{port: port, storage: storage}
	http.HandleFunc("/", server.generalHandler)
	return server
}

// Run ... - run the Server
func (s *Config) Run() error {
	url := fmt.Sprintf("%s:%d", s.host, s.port)

	log.Fatal(http.ListenAndServe(url, nil))
	return nil
}

func (s *Config) generalHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("generalHandler")

	switch r.Method {
	case "GET":

		if err := s.getHandler(w, r); err != nil {

			log.Print("ERROR:")
		}
	case "PUT":
		fmt.Println("PUT")
		if r.Body == nil {

			http.Error(w, "Please send a request body", 400)
			return
		}
		if err := s.setHandler(w, r); err != nil {

			log.Print("ERROR:")
		}
	case "DELETE":
		fmt.Println("DEL")
		if err := s.delHandler(w, r); err != nil {

			log.Print("ERROR:")
		}
	case "POST":
		fmt.Println("POST")
		if r.Body == nil {

			http.Error(w, "Please send a request body", 400)
			return
		}
		if err := s.postHandler(w, r); err != nil {

			log.Print("ERROR:", err)
		}
	}
}

func (s *Config) getHandler(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		fmt.Println("error - ", err)
		log.Print(err)
		return err
	}
	for k, v := range r.Form {
		fmt.Println("Form - ", k, v)
		value, err := s.storage.Get(k)
		if err != nil {
			fmt.Fprintf(w, value)
		} else {
			tupleValue := fmt.Sprintf("(%s, %s)", k, value)
			fmt.Fprintf(w, tupleValue)
		}

	}
	return nil
}

func (s *Config) setHandler(w http.ResponseWriter, r *http.Request) error {

	err := json.NewDecoder(r.Body).Decode(&keymap)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return err
	}
	// TODO think about var
	s.storage.Set(keymap.Name, keymap.Value)
	return nil
}

func (s *Config) delHandler(w http.ResponseWriter, r *http.Request) error {

	if err := r.ParseForm(); err != nil {
		fmt.Println("error - ", err)
		log.Print(err)
		return err
	}
	for k := range r.Form {
		s.storage.Delete(k)
	}
	return nil
}

func (s *Config) postHandler(w http.ResponseWriter, r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(&patternmap)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return err
	}
	// TODO think about var
	keys, err := s.storage.Keys(patternmap.Pattern)
	fmt.Fprintf(w, strings.Join(keys, ","))
	return nil
}
