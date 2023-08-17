package main

import (
	"encoding/json"
	"fmt"
	"httpchecker/logger"
	"net/http"
	"time"
)

const (
	targetHost   = "https://example.com"
	numRequests  = 5
	intervalSecs = 10
	logFileName  = "log.txt"
)

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
	logger, err := logger.NewLogger(logFileName)
	if err != nil {
		fmt.Println("Błąd podczas dostępu do pliku:", err)
		return
	}
	defer logger.Close()

	httpClient := &http.Client{}
	for i := 0; i < numRequests; i++ {

		responseMessage := fmt.Sprintf("[%s] Próba: %d\n ", time.Now().Format(time.RFC3339), i+1)
		logger.WriteLog(responseMessage)

		startTime := time.Now()
		resp, err := httpClient.Get(targetHost)
		duration := time.Since(startTime)
		if err != nil {
			duration = 0
			responseMessage := fmt.Sprintf("[%s] Błąd podczas wysyłania żądania: %s\n", time.Now().Format(time.RFC3339), err)
			logger.WriteLog(responseMessage)
		}
		responseMessage = fmt.Sprintf("[%s] Żądanie %d wysłane\n", time.Now().Format(time.RFC3339), i+1)
		logger.WriteLog(responseMessage)
		defer resp.Body.Close()

		responseMessage = fmt.Sprintf("[%s] Czas odpowiedzi: %v\n", time.Now().Format(time.RFC3339), duration)
		logger.WriteLog(responseMessage)

		statusCode := checkStatusCode(resp)
		responseMessage = fmt.Sprintf("[%s] Kod odpowiedzi HTTP: %d\n", time.Now().Format(time.RFC3339), statusCode)
		logger.WriteLog(responseMessage)

		if !checkJSONContentType(resp) {
			responseMessage = fmt.Sprintf("[%s] Odpowiedź nie jest typu JSON \n", time.Now().Format(time.RFC3339))
			logger.WriteLog(responseMessage)
		}

		if !validateJSONSyntax(resp) {
			now := time.Now().Format(time.RFC3339)
			responseMessage = fmt.Sprintf("[%s] Błąd podczas walidacji JSON \n", now)
			logger.WriteLog(responseMessage)
		}

		time.Sleep(intervalSecs * time.Second)
	}
}
