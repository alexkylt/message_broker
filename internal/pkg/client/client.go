package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	_ "net/url"
	"os"
	"strings"
)

type key struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func sendRequest(arguments string, host string, port int) error {

	arguments = strings.TrimSuffix(arguments, "\n")
	args := strings.Fields(arguments)
	fmt.Println("args - ", args)
	client := &http.Client{}

	switch strings.ToUpper(args[0]) {
	case "EXIT":
		os.Exit(0)
	case "GET":
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
	case "SET":
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
	case "DEL":
		url := fmt.Sprintf("http://%s:%d/?%s", host, port, args[1])
		req, err := http.NewRequest("DELETE", url, nil)
		fmt.Println(req)

		resp, err := client.Do(req)
		if err != nil {

			fmt.Println(err)
			return err
		}
		defer resp.Body.Close()
		io.Copy(os.Stdout, resp.Body)
	case "KEYS":
		url := fmt.Sprintf("http://%s:%d", host, port)
		str := fmt.Sprintf(`{"pattern": "%s"}`, args[1])
		var jsonStr = []byte(str)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

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
	flag.StringVar(&host, "host", "127.0.0.1", "specify host to use.  defaults to 127.0.0.1")
	flag.Parse()
	fmt.Println(port, host)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		inputStr, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		err = sendRequest(inputStr, host, port)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
