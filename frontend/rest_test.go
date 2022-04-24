package frontend

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestRestPut(t *testing.T) {
	const key = "a-key"
	const value = "a-value"
	var restUrl = "http://localhost:" + restPort + "/v1/" + key

	client := http.Client{Timeout: 3 * time.Second}
	req, err := http.NewRequest("PUT", restUrl, strings.NewReader(value))
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 201 {
		t.Errorf("put returned wrong status code: %v", resp.StatusCode)
	}
}

func TestRestGet(t *testing.T) {
	const key = "a-key"
	const value = "a-value"
	var restUrl = "http://localhost:" + restPort + "/v1/" + key

	client := http.Client{Timeout: 3 * time.Second}
	req, err := http.NewRequest("PUT", restUrl, strings.NewReader(value))
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 201 {
		t.Errorf("put returned wrong status code: %v", resp.StatusCode)
	}

	resp, err = client.Get(restUrl)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("get returned wrong status code: %v", resp.StatusCode)
	}
	val, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	fmt.Println("string(val):", string(val))
	if string(val) != value {
		t.Errorf("val/value mismatch, val: %q, value: %q", val, value)
	}
}

func TestRestDelete(t *testing.T) {
	const key = "a-key"
	const value = "a-value"
	var restUrl = "http://localhost:" + restPort + "/v1/" + key

	client := http.Client{Timeout: 3 * time.Second}
	req, err := http.NewRequest("PUT", restUrl, strings.NewReader(value))
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 201 {
		t.Errorf("put returned wrong status code: %v", resp.StatusCode)
	}

	req, err = http.NewRequest("DELETE", restUrl, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("delete returned wrong status code: %v", resp.StatusCode)
	}

	resp, err = client.Get(restUrl)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 404 {
		t.Errorf("get returned wrong status code: %v", resp.StatusCode)
	}
}
