package dbsqlite

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("could not open in-memory database: %v", err)
	}

	_, err = db.Exec(`
		PRAGMA foreign_keys = ON;
		CREATE TABLE users(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_name TEXT NOT NULL,
			current_amount REAL NOT NULL,
			monthly_inputs REAL NOT NULL,
			monthly_outputs REAL NOT NULL
		);
		CREATE TABLE monthly_adjustments_log(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			year_month TEXT NOT NULL UNIQUE,
			applied_at TEXT DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		t.Fatalf("could not create tables: %v", err)
	}

	return db
}

func insertTestUser(t *testing.T, db *sql.DB, name string, currentAmount, monthlyInputs, monthlyOutputs int64) int64 {
	t.Helper()

	res, err := db.Exec(
		`INSERT INTO users(user_name, current_amount, monthly_inputs, monthly_outputs) VALUES (?, ?, ?, ?)`,
		name, currentAmount, monthlyInputs, monthlyOutputs,
	)
	if err != nil {
		t.Fatalf("could not insert test user %q: %v", name, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		t.Fatalf("could not get last insert id: %v", err)
	}
	return id
}

func getUserCurrentAmount(t *testing.T, db *sql.DB, id int64) int64 {
	t.Helper()

	var amount int64
	err := db.QueryRow(`SELECT current_amount FROM users WHERE id = ?`, id).Scan(&amount)
	if err != nil {
		t.Fatalf("could not get current_amount for user %d: %v", id, err)
	}
	return amount
}

func TestGetLastProcessedMonth_Empty(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	last, err := GetLastProcessedMonth(ctx, db)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if last != "" {
		t.Errorf("expected empty string for no records, got %q", last)
	}
}

func TestRecordMonthWithoutAdjustment(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Insert a user to verify no adjustment is made
	id := insertTestUser(t, db, "Alice", 10000, 5000, 3000)

	err := RecordMonthWithoutAdjustment(ctx, db, "2025-07")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify the month was recorded
	last, err := GetLastProcessedMonth(ctx, db)
	if err != nil {
		t.Fatalf("unexpected error getting last processed month: %v", err)
	}
	if last != "2025-07" {
		t.Errorf("expected last processed month %q, got %q", "2025-07", last)
	}

	// Verify user's current_amount was NOT changed
	amount := getUserCurrentAmount(t, db, id)
	if amount != 10000 {
		t.Errorf("expected current_amount to remain 10000, got %d", amount)
	}
}

func TestRecordMonthWithoutAdjustment_Idempotent(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := RecordMonthWithoutAdjustment(ctx, db, "2025-07")
	if err != nil {
		t.Fatalf("unexpected error on first insert: %v", err)
	}

	// Second call should not error (INSERT OR IGNORE)
	err = RecordMonthWithoutAdjustment(ctx, db, "2025-07")
	if err != nil {
		t.Fatalf("unexpected error on duplicate insert: %v", err)
	}
}

func TestIsMonthProcessed(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Should be false initially
	processed, err := IsMonthProcessed(ctx, db, "2025-07")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if processed {
		t.Error("expected month to not be processed initially")
	}

	// Record the month
	err = RecordMonthWithoutAdjustment(ctx, db, "2025-07")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Now it should be true
	processed, err = IsMonthProcessed(ctx, db, "2025-07")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !processed {
		t.Error("expected month to be processed after recording")
	}

	// A different month should still be false
	processed, err = IsMonthProcessed(ctx, db, "2025-08")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if processed {
		t.Error("expected different month to not be processed")
	}
}

func TestApplyMonthlyAdjustment_SingleUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// User with current_amount=10000, monthly_inputs=5000, monthly_outputs=3000
	// Expected after adjustment: 10000 + 5000 - 3000 = 12000
	id := insertTestUser(t, db, "Alice", 10000, 5000, 3000)

	err := ApplyMonthlyAdjustment(ctx, db, "2025-07")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	amount := getUserCurrentAmount(t, db, id)
	if amount != 12000 {
		t.Errorf("expected current_amount 12000, got %d", amount)
	}

	// Verify the month was logged
	processed, err := IsMonthProcessed(ctx, db, "2025-07")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !processed {
		t.Error("expected month 2025-07 to be recorded as processed")
	}
}

func TestApplyMonthlyAdjustment_MultipleUsers(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Alice: 10000 + 5000 - 3000 = 12000
	aliceID := insertTestUser(t, db, "Alice", 10000, 5000, 3000)
	// Bob: 20000 + 8000 - 6000 = 22000
	bobID := insertTestUser(t, db, "Bob", 20000, 8000, 6000)
	// Carol: 0 + 1000 - 500 = 500
	carolID := insertTestUser(t, db, "Carol", 0, 1000, 500)

	err := ApplyMonthlyAdjustment(ctx, db, "2025-07")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := []struct {
		name     string
		id       int64
		expected int64
	}{
		{"Alice", aliceID, 12000},
		{"Bob", bobID, 22000},
		{"Carol", carolID, 500},
	}

	for _, tc := range tests {
		amount := getUserCurrentAmount(t, db, tc.id)
		if amount != tc.expected {
			t.Errorf("%s: expected current_amount %d, got %d", tc.name, tc.expected, amount)
		}
	}
}

func TestApplyMonthlyAdjustment_NegativeBalance(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// User where outputs > inputs: 1000 + 2000 - 5000 = -2000
	id := insertTestUser(t, db, "Dave", 1000, 2000, 5000)

	err := ApplyMonthlyAdjustment(ctx, db, "2025-07")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	amount := getUserCurrentAmount(t, db, id)
	if amount != -2000 {
		t.Errorf("expected current_amount -2000, got %d", amount)
	}
}

func TestApplyMonthlyAdjustment_DuplicateMonthFails(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	insertTestUser(t, db, "Alice", 10000, 5000, 3000)

	// First application should succeed
	err := ApplyMonthlyAdjustment(ctx, db, "2025-07")
	if err != nil {
		t.Fatalf("unexpected error on first apply: %v", err)
	}

	// Second application of the same month should fail (UNIQUE constraint)
	err = ApplyMonthlyAdjustment(ctx, db, "2025-07")
	if err == nil {
		t.Fatal("expected error when applying the same month twice, got nil")
	}

	// Verify the amount was only adjusted once (10000 + 5000 - 3000 = 12000)
	amount := getUserCurrentAmount(t, db, 1)
	if amount != 12000 {
		t.Errorf("expected current_amount 12000 (single adjustment), got %d", amount)
	}
}

func TestApplyMonthlyAdjustment_MultiMonthCatchup(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Alice: start at 10000, monthly net = +5000 - 3000 = +2000
	id := insertTestUser(t, db, "Alice", 10000, 5000, 3000)

	// Simulate processing 3 missed months
	months := []string{"2025-05", "2025-06", "2025-07"}
	for _, month := range months {
		err := ApplyMonthlyAdjustment(ctx, db, month)
		if err != nil {
			t.Fatalf("unexpected error applying month %s: %v", month, err)
		}
	}

	// Expected: 10000 + (2000 * 3) = 16000
	amount := getUserCurrentAmount(t, db, id)
	if amount != 16000 {
		t.Errorf("expected current_amount 16000 after 3 months, got %d", amount)
	}

	// Verify all months were logged
	for _, month := range months {
		processed, err := IsMonthProcessed(ctx, db, month)
		if err != nil {
			t.Fatalf("unexpected error checking month %s: %v", month, err)
		}
		if !processed {
			t.Errorf("expected month %s to be recorded as processed", month)
		}
	}
}

func TestGetLastProcessedMonth_ReturnsLatest(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	insertTestUser(t, db, "Alice", 10000, 5000, 3000)

	// Apply months out of order
	months := []string{"2025-03", "2025-05", "2025-04"}
	for _, month := range months {
		err := ApplyMonthlyAdjustment(ctx, db, month)
		if err != nil {
			t.Fatalf("unexpected error applying month %s: %v", month, err)
		}
	}

	last, err := GetLastProcessedMonth(ctx, db)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should return the latest month by string ordering (which works for YYYY-MM format)
	if last != "2025-05" {
		t.Errorf("expected last processed month %q, got %q", "2025-05", last)
	}
}

func TestApplyMonthlyAdjustment_ZeroInputsAndOutputs(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// User with zero monthly inputs and outputs: amount should not change
	id := insertTestUser(t, db, "Eve", 5000, 0, 0)

	err := ApplyMonthlyAdjustment(ctx, db, "2025-07")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	amount := getUserCurrentAmount(t, db, id)
	if amount != 5000 {
		t.Errorf("expected current_amount 5000 (unchanged), got %d", amount)
	}
}

func TestApplyMonthlyAdjustment_NoUsers(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// No users in the database — should succeed without error
	err := ApplyMonthlyAdjustment(ctx, db, "2025-07")
	if err != nil {
		t.Fatalf("unexpected error when no users exist: %v", err)
	}

	// Month should still be recorded
	processed, err := IsMonthProcessed(ctx, db, "2025-07")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !processed {
		t.Error("expected month to be recorded even with no users")
	}
}

func TestApplyMonthlyAdjustment_TransactionAtomicity(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	id := insertTestUser(t, db, "Alice", 10000, 5000, 3000)

	// First apply succeeds
	err := ApplyMonthlyAdjustment(ctx, db, "2025-07")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Attempting the same month again should fail, and the user amount
	// should remain at the value after the first (successful) adjustment
	_ = ApplyMonthlyAdjustment(ctx, db, "2025-07")

	amount := getUserCurrentAmount(t, db, id)
	if amount != 12000 {
		t.Errorf("expected current_amount to remain 12000 after failed duplicate, got %d", amount)
	}
}
