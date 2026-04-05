package test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"stock_backend/internal/model/request"
	"stock_backend/internal/model/response"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	registerPath = "/api/v1/auth/register"
	loginPath    = "/api/v1/auth/login"
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
	log.Println("Table is cleared")
}

func CreateTestUser(email string, password string) error {
	req := request.RegisterRequest{
		Email:    email,
		Password: password,
		Username: "test_username",
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	url := registerPath
	_, _, err := PerformRequest[*response.RegisterResponse](req, url, http.MethodPost, headers)
	return err
}

func CreateAdmin(email string, password string, username string) error {
	hashPw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	_, err = db.Exec("INSERT INTO users (id, email, password, verified, username, roleId) VALUES ($1, $2, $3, $4, $5, $6)",
		uuid.New(),
		email, string(hashPw), true, username, 2)
	if err != nil {
		return err
	}
	return nil
}

func GetUserToken(email string, password string) (string, error) {
	requestBody := request.LoginRequest{
		Email:    email,
		Password: password,
	}

	httpHeader := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	url := loginPath
	result, _, err := PerformRequest[*response.LoginResponse](requestBody, url, http.MethodPost, httpHeader)
	if err != nil {
		return "", err
	}

	return result.Token, nil
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
	if err != nil {
		return result, 0, err
	}

	defer func() {
		if err = res.Body.Close(); err != nil {
			log.Println("failed to close body")
		}
	}()

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
