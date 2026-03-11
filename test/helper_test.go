package test

import (
	"encoding/json"
	"io"
	"log"
	"net/http/httptest"
	"strings"
)

func ClearUser() {
	res, err := db.Exec("DELETE FROM users")
	if err != nil {
		log.Fatalf("Failed clear user data : %+v", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Failed clear user data : %+v", err)
	}

	if rows == 0 {
		log.Println("No Data Deleted")
		return
	}
	log.Println("Table users is cleared")
}

func PerformRequest[T any](requestBody any, url string, httpMethod string, httpHeader map[string]string) (T, int, error) {
	var result T
	requestJson, err := json.Marshal(requestBody)
	if err != nil {
		return result, 0, err
	}

	req := httptest.NewRequest(httpMethod, url, strings.NewReader(string(requestJson)))
	for key, val := range httpHeader {
		req.Header.Set(key, val)
	}

	res, err := app.Test(req)
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Println("Error in closing")
		}
	}()

	if err != nil {
		return result, res.StatusCode, err
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return result, res.StatusCode, err
	}

	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return result, res.StatusCode, err
	}

	return result, res.StatusCode, nil
}