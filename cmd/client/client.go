package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"net/http"
	_ "net/url"
	"os"
)

type key struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type pattern struct {
	Pattern string `json:"pattern"`
}

func sendRequest(arguments string, host string, port int) error {

	arguments = strings.TrimSuffix(arguments, "\n")
	args := strings.Fields(arguments)
	fmt.Println("args - ", args, len(args))
	client := &http.Client{}
	if args[0] == "" {
		return nil
	}

	switch strings.ToUpper(args[0]) {
	case "EXIT":
		os.Exit(0)
	case "GET":
		if len(args) == 1 {
			fmt.Printf("%s\n", "You should specify a key value.")
		} else {
			url := fmt.Sprintf("http://%s:%d/?%s", host, port, args[1])

			req, err := http.NewRequest("GET", url, nil)
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println(err)
				return err
			}

			defer resp.Body.Close()
			// io.Copy(os.Stdout, resp.Body)
			contents, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("%s", err)
				os.Exit(1)
			}
			fmt.Printf("%s\n", string(contents))
		}
	case "SET":
		if len(args) == 2 {
			fmt.Printf("%s\n", "You should specify a key value along with a value.")
		} else {
			keymap := key{Name: args[1], Value: args[2]}
			buff := new(bytes.Buffer)
			json.NewEncoder(buff).Encode(keymap)
			url := fmt.Sprintf("http://%s:%d/", host, port)
			req, err := http.NewRequest("PUT", url, buff)

			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {

				fmt.Println(err)
				return err
			}
			defer resp.Body.Close()
			io.Copy(os.Stdout, resp.Body)
		}
	case "DEL":
		if len(args) == 1 {
			fmt.Printf("%s\n", "You should specify a key value.")
		} else {
			url := fmt.Sprintf("http://%s:%d/?%s", host, port, args[1])
			req, err := http.NewRequest("DELETE", url, nil)
			//fmt.Println(req)

			resp, err := client.Do(req)
			if err != nil {

				fmt.Println(err)
				return err
			}
			defer resp.Body.Close()
			io.Copy(os.Stdout, resp.Body)
		}
	case "KEYS":
		url := fmt.Sprintf("http://%s:%d", host, port)
		str := pattern{Pattern: args[1]}
		buff := new(bytes.Buffer)
		json.NewEncoder(buff).Encode(str)
		req, err := http.NewRequest("POST", url, buff)

		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return err
		}

		defer resp.Body.Close()
		// io.Copy(os.Stdout, resp.Body)
		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", string(contents))
	default:
		fmt.Println("TODO!!!!")
	}

	return nil
}

func main() {
	var port int
	var host string

	flag.IntVar(&port, "port", 9090, "specify port to use.  defaults to 9090.")
	flag.StringVar(&host, "host", "server", "specify host to use.  defaults to 'server' docker container")
	flag.Parse()
	fmt.Println(port, host)

	reader := bufio.NewReader(os.Stdin)

	for {

		fmt.Print("$ ")

		inputStr, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else if len(inputStr) < 5 {
			fmt.Fprintln(os.Stderr, "Please, specify the appropriate command and value for it: GET, SET or DEL.")
		} else {
			err = sendRequest(inputStr, host, port)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}
