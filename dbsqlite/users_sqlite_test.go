package dbsqlite

import (
	"os"
	"testing"

	"natan/fingo/model"
	"natan/fingo/utils"
)

func setupDB(t *testing.T) {
	// Ensure a clean environment
	_ = os.Remove("fingo.db")

	// create database and schema
	if err := createDatabase(); err != nil {
		t.Fatalf("failed to create database for test setup: %v", err)
	}
}

func teardownDB() {
	_ = os.Remove("fingo.db")
}

func TestReturnAllUsers_Empty(t *testing.T) {
	setupDB(t)
	defer teardownDB()

	db, err := GetDatabaseConnection()
	if err != nil {
		t.Fatalf("GetDatabaseConnection() returned error: %v", err)
	}
	defer db.Close()

	users, err := ReturnAllUsers(db)
	if err != nil {
		t.Fatalf("ReturnAllUsers() returned error on empty DB: %v", err)
	}

	if len(users) != 0 {
		t.Fatalf("expected 0 users, got %d", len(users))
	}
}

func TestCreateUserAndReturnAllUsers(t *testing.T) {
	setupDB(t)
	defer teardownDB()

	db, err := GetDatabaseConnection()
	if err != nil {
		t.Fatalf("GetDatabaseConnection() returned error: %v", err)
	}
	defer db.Close()

	u := model.User{
		UserName:       "Alice",
		CurrentAmount:  utils.Money(1000), // R$10.00
		MonthlyInputs:  utils.Money(2000), // R$20.00
		MonthlyOutputs: utils.Money(500),  // R$5.00
	}

	_, err = CreateUser(u, db)
	if err != nil {
		t.Fatalf("CreateUser() returned error: %v", err)
	}

	users, err := ReturnAllUsers(db)
	if err != nil {
		t.Fatalf("ReturnAllUsers() returned error after insert: %v", err)
	}

	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}

	u2 := users[0]
	if u2.UserName != u.UserName {
		t.Errorf("username mismatch: expected %q, got %q", u.UserName, u2.UserName)
	}
	if u2.CurrentAmount != u.CurrentAmount {
		t.Errorf("current amount mismatch: expected %v, got %v", u.CurrentAmount, u2.CurrentAmount)
	}
	if u2.MonthlyInputs != u.MonthlyInputs {
		t.Errorf("monthly inputs mismatch: expected %v, got %v", u.MonthlyInputs, u2.MonthlyInputs)
	}
	if u2.MonthlyOutputs != u.MonthlyOutputs {
		t.Errorf("monthly outputs mismatch: expected %v, got %v", u.MonthlyOutputs, u2.MonthlyOutputs)
	}
}

func TestCreateUser_InvalidUserNameCausesError(t *testing.T) {
	setupDB(t)
	defer teardownDB()

	db, err := GetDatabaseConnection()
	if err != nil {
		t.Fatalf("GetDatabaseConnection() returned error: %v", err)
	}
	defer db.Close()

	// username contains an unescaped single quote which should break the SQL
	u := model.User{
		UserName:       "O'Neil",
		CurrentAmount:  utils.Money(0),
		MonthlyInputs:  utils.Money(0),
		MonthlyOutputs: utils.Money(0),
	}

	_, err = CreateUser(u, db)
	if err == nil {
		t.Fatalf("expected CreateUser() to return an error for invalid/unsafe username, but it returned nil")
	}
}
