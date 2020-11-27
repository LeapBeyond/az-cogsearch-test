// Package http wraps up some stuff for doing raw API calls
package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Post attempts to post a payload to a url with a set of headers, and returns the status code of the request.
// note that status codes outside the 200-399 range are treated as errors.
func Post(url string, payload interface{}, headers map[string]string) (int, error) {
	return doHttpCall(url, "POST", payload, headers)
}

// Put attempts to post a payload to a url with a set of headers, and returns the status code of the request.
// note that status codes outside the 200-399 range are treated as errors.
func Put(url string, payload interface{}, headers map[string]string) (int, error) {
	return doHttpCall(url, "PUT", payload, headers)
}

// Get tries to get from the supplied url, with the supplied headers, and return the body
func Get(url string, headers map[string]string) ([]byte, error) {
	// construct the request
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, fmt.Errorf("Failed to consruct the post request: %v", err)
	}
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	// perform the  request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return []byte{}, fmt.Errorf("Post request failed: %v", err)
	}
	if response.StatusCode < 200 || response.StatusCode > 399 {
		return []byte{}, fmt.Errorf("Received an unexpected response %d", response.StatusCode)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("Failed to retrieve body from response: %v", err)
	}

	return body, nil
}

func doHttpCall(url, callType string, payload interface{}, headers map[string]string) (int, error) {
	// marshalls the object to a []byte
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return 400, fmt.Errorf("Failed to marshall the payload to JSON: %v", err)
	}

	// construct the request
	request, err := http.NewRequest(callType, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return 400, fmt.Errorf("Failed to consruct the post request: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	// perform the request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return 400, fmt.Errorf("%s request failed: %v", callType, err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 399 {
		body, _ := ioutil.ReadAll(response.Body)
		return response.StatusCode, fmt.Errorf("Received an unexpected response %d: %s", response.StatusCode, body)
	}

	return response.StatusCode, nil
}
