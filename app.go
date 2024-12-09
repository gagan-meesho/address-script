package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

func makeRequest(url string, payload []byte, headers map[string]string, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	fmt.Printf("Response: %s\n", string(body))
}

func main() {
	url := "http://address-search.int.meesho.int/api/2.0/addresses/414091697?context=cart&cart_identifier=default&source=SUPERSTORE"
	payload := []byte(`{
    "id":414091697,
    "alternative_mobile": null,
    "pin": "560109",
    "city": "नमस्ते",
    "name": "नमस्ते",
    "address_line_1": "lkjasdf",
    "mobile": "9909890299",
    "address_line_2": "नमस्ते",
    "state": "नमस्ते",
    "landmark": "नमस्ते",
    "country_id": 1,
    "user_id": 340284318
	}`)

	headers := map[string]string{
		"Content-Type":            "application/json",
		"Authentication":          "<get this from vault>", //note: not commiting. get this from vault
		"app-user-id":             "340284318",
		"MEESHO-ISO-COUNTRY-CODE": "IN",
	}

	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go makeRequest(url, payload, headers, &wg)
	}

	wg.Wait()
	fmt.Println("All requests completed")
}
