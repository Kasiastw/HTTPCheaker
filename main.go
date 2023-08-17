package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	targetHost   = "https://example.com"
	numRequests  = 5
	intervalSecs = 10
	logFileName  = "log.txt"
)

type logger struct {
	file *os.File
}

func newLogger(fileName string) (*logger, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return &logger{file: file}, nil
}

func (l *logger) writeLog(message string) {
	_, err := l.file.WriteString(message)
	if err != nil {
		fmt.Println("Błąd podczas zapisu do pliku:", err)
	}
}

func (l *logger) close() {
	l.file.Close()
}

func checkStatusCode(resp *http.Response) int {
	return resp.StatusCode
}

func checkJSONContentType(resp *http.Response) bool {
	contentType := resp.Header.Get("Content-Type")
	return contentType == "application/json"
}

func validateJSONSyntax(resp *http.Response) bool {
	var responseJSON map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&responseJSON)
	return err == nil
}

func main() {
	logger, err := newLogger(logFileName)
	if err != nil {
		fmt.Println("Błąd podczas dostępu do pliku:", err)
		return
	}
	defer logger.close()

	httpClient := &http.Client{}
	for i := 0; i < numRequests; i++ {

		responseMessage := fmt.Sprintf("[%s] Próba: %d\n ", time.Now().Format(time.RFC3339), i+1)
		fmt.Printf(responseMessage)
		logger.writeLog(responseMessage)

		startTime := time.Now()
		resp, err := httpClient.Get(targetHost)
		duration := time.Since(startTime)
		if err != nil {
			duration = 0
			responseMessage := fmt.Sprintf("[%s] Błąd podczas wysyłania żądania: %s\n", time.Now().Format(time.RFC3339), err)
			fmt.Printf(responseMessage)
			logger.writeLog(responseMessage)
		}
		responseMessage = fmt.Sprintf("[%s] Żądanie %d wysłane\n", time.Now().Format(time.RFC3339), i+1)
		fmt.Printf(responseMessage)
		logger.writeLog(responseMessage)
		defer resp.Body.Close()

		responseMessage = fmt.Sprintf("[%s] Czas odpowiedzi: %v\n", time.Now().Format(time.RFC3339), duration)
		fmt.Printf(responseMessage)
		logger.writeLog(responseMessage)

		statusCode := checkStatusCode(resp)
		responseMessage = fmt.Sprintf("[%s] Kod odpowiedzi HTTP: %d\n", time.Now().Format(time.RFC3339), statusCode)
		fmt.Printf(responseMessage)
		logger.writeLog(responseMessage)

		if !checkJSONContentType(resp) {
			responseMessage = fmt.Sprintf("[%s] Odpowiedź nie jest typu JSON \n", time.Now().Format(time.RFC3339))
			fmt.Printf(responseMessage)
			logger.writeLog(responseMessage)
		}

		if !validateJSONSyntax(resp) {
			now := time.Now().Format(time.RFC3339)
			responseMessage = fmt.Sprintf("[%s] Błąd podczas walidacji JSON \n", now)
			fmt.Printf(responseMessage)
			logger.writeLog(responseMessage)
		}

		time.Sleep(intervalSecs * time.Second)
	}
}
