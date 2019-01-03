package server

import (
	_ "database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	_ "strings"
	_ "time"

	"github.com/alexkylt/message_broker/internal/pkg/storage"
	_ "github.com/lib/pq"
)

type serverCfg struct {
	host    string
	port    int
	storage storage.StorageInterface
}

type key struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

var keymap key

func InitServer(port int, storage storage.StorageInterface) *serverCfg {
	server := &serverCfg{port: port, storage: storage}
	http.HandleFunc("/", server.generalHandler)
	return server
}

func (s *serverCfg) Run() error {
	url := fmt.Sprintf("%s:%d", s.host, s.port)

	log.Fatal(http.ListenAndServe(url, nil))
	return nil
}

func (s *serverCfg) generalHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("generalHandler")

	switch r.Method {
	case "GET":

		if err := s.getHandler(w, r); err != nil {
			// w.httpCode(404)

			log.Print("ALARM!")
		}
	case "PUT":
		fmt.Println("PUT")
		if r.Body == nil {

			http.Error(w, "Please send a request body", 400)
			return
		}
		if err := s.setHandler(w, r); err != nil {
			// w.httpCode(404)

			log.Print("ALARM!")
		}
	case "DELETE":
		fmt.Println("DEL")
		if err := s.delHandler(w, r); err != nil {
			// w.httpCode(404)

			log.Print("ALARM!")
		}
	case "POST":
		fmt.Println("POST")
		if r.Body == nil {

			http.Error(w, "Please send a request body", 400)
			return
		}
		if err := s.postHandler(w, r); err != nil {
			// w.httpCode(404)

			log.Print("ALARM!")
		}
	}
}

func (s *serverCfg) getHandler(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("GET params were:", r.URL.Query())
	fmt.Println("getHandler", r.Form)
	fmt.Println("getHandler", r.ParseForm)
	if err := r.ParseForm(); err != nil {
		fmt.Println("error - ", err)
		log.Print(err)
		return err
	}
	for k, v := range r.Form {
		fmt.Println("Form - ", k, v)
		// value, _ := s.storage.Get(v[0])
		value, err := s.storage.Get(k)
		if err != nil {
			fmt.Fprintf(w, value)
		} else {
			tuple_value := fmt.Sprintf("(%s, %s)", k, value)
			fmt.Fprintf(w, tuple_value)
		}

	}
	return nil
}

func (s *serverCfg) setHandler(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("SET params were:", r.URL.Query())
	err := json.NewDecoder(r.Body).Decode(&keymap)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return err
	}
	// TODO think about var
	s.storage.Set(keymap.Name, keymap.Value)
	// fmt.Println(k.Name, k.Value, value, s.storage, act)
	return nil
}

func (s *serverCfg) delHandler(w http.ResponseWriter, r *http.Request) error {

	if err := r.ParseForm(); err != nil {
		fmt.Println("error - ", err)
		log.Print(err)
		return err
	}
	for k, _ := range r.Form {
		fmt.Println("Before delete - ", k, s.storage)
		value := s.storage.Delete(k)
		//fmt.Fprintf(w, value)
		fmt.Println("After delete - ", k, s.storage, value)
	}
	return nil
}

func (s *serverCfg) postHandler(w http.ResponseWriter, r *http.Request) error {

	fmt.Println("POST params were:", r.URL.Query())
	err := json.NewDecoder(r.Body).Decode(&keymap)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return err
	}
	// TODO think about var
	s.storage.Keys(keymap.Name)
	// fmt.Println(k.Name, k.Value, value, s.storage, act)
	return nil
}
