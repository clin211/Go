package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	url := "https://jsonplaceholder.typicode.com/posts"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Cookie", "__cfduid=d974efa095d4dc9676d383482e35e77d51618042392")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
