package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestCheckStatusCode(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusOK,
	}

	statusCode := checkStatusCode(resp)
	if statusCode != http.StatusOK {
		t.Errorf("Oczekiwany status code %d, ale mam %d", http.StatusOK, statusCode)
	}
}

func TestCheckJSONContentType(t *testing.T) {
	resp := &http.Response{
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
	}

	if !checkJSONContentType(resp) {
		t.Error("Oczekiwany JSON content type, ale mam false")
	}
}

func TestValidateJSONSyntax(t *testing.T) {
	jsonResponse := `{"key": "value"}`
	resp := &http.Response{
		Body:       ioutil.NopCloser(strings.NewReader(jsonResponse)),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		StatusCode: http.StatusOK,
	}

	if !validateJSONSyntax(resp) {
		t.Error("Oczekiwana walidacja JSON syntax, ale mam false")
	}

	invalidJSONResponse := `{"key": "value"`
	invalidResp := &http.Response{
		Body:       ioutil.NopCloser(strings.NewReader(invalidJSONResponse)),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		StatusCode: http.StatusOK,
	}

	if validateJSONSyntax(invalidResp) {
		t.Error("Oczekiwana walidacja JSON syntax, ale mam true")
	}
}
