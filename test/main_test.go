package test

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	ClearTable("users")
	if err := CreateTestUser(email, password); err != nil {
		panic(err)
	}

	tok, err := GetUserToken(email, password)
	if err != nil {
		panic(err)
	}
	token = tok

	if err := CreateAdmin("admin@gmail.com", password, "admin"); err != nil {
		panic(err)
	}

	adminToken, err = GetUserToken("admin@gmail.com", password)
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}
