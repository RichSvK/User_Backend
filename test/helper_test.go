package test

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"stock_backend/internal/model/request"
	"stock_backend/internal/model/response"
	"strings"
)

func ClearTable(tableName string) {
	res, err := db.Exec(fmt.Sprintf("DELETE FROM %s", tableName))
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

func GetUserToken(email string, password string) (string, error){
	requestBody := request.LoginRequest{
		Email:    email,
		Password: password,
	}

	httpHeader := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	url := "/api/v1/users/login"
	result, _, err := PerformRequest[*response.LoginResponse](requestBody, url, http.MethodPost, httpHeader)
	if err != nil {
		return "", err
	}

	return result.Token, nil
}