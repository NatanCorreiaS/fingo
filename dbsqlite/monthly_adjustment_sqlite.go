package dbsqlite

import (
	"context"
	"database/sql"
	"fmt"
)

// GetLastProcessedMonth retrieves the most recent year_month from the monthly_adjustments_log.
// Returns an empty string if no adjustments have been recorded yet.
func GetLastProcessedMonth(ctx context.Context, db *sql.DB) (string, error) {
	const query = `SELECT year_month FROM monthly_adjustments_log ORDER BY year_month DESC LIMIT 1;`

	var yearMonth string
	err := db.QueryRowContext(ctx, query).Scan(&yearMonth)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("could not get last processed month: %w", err)
	}

	return yearMonth, nil
}

// IsMonthProcessed checks whether a specific year_month has already been processed.
func IsMonthProcessed(ctx context.Context, db *sql.DB, yearMonth string) (bool, error) {
	const query = `SELECT COUNT(*) FROM monthly_adjustments_log WHERE year_month = ?;`

	var count int
	err := db.QueryRowContext(ctx, query, yearMonth).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("could not check if month %s is processed: %w", yearMonth, err)
	}

	return count > 0, nil
}

// ApplyMonthlyAdjustment updates current_amount for all users by adding monthly_inputs
// and subtracting monthly_outputs, then records the year_month in the log.
// The entire operation runs inside a transaction to ensure atomicity.
func ApplyMonthlyAdjustment(ctx context.Context, db *sql.DB, yearMonth string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not begin transaction for monthly adjustment: %w", err)
	}
	defer tx.Rollback()

	// Update all users: current_amount = current_amount + monthly_inputs - monthly_outputs
	const updateStmt = `
		UPDATE users
		SET current_amount = current_amount + monthly_inputs - monthly_outputs;
	`

	_, err = tx.ExecContext(ctx, updateStmt)
	if err != nil {
		return fmt.Errorf("could not apply monthly adjustment to users: %w", err)
	}

	// Record that this month has been processed
	const insertStmt = `INSERT INTO monthly_adjustments_log(year_month) VALUES (?);`

	_, err = tx.ExecContext(ctx, insertStmt, yearMonth)
	if err != nil {
		return fmt.Errorf("could not record monthly adjustment for %s: %w", yearMonth, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit monthly adjustment transaction: %w", err)
	}

	return nil
}

// RecordMonthWithoutAdjustment records a year_month in the log without applying
// any adjustment. Used to mark the initial month when the system first starts.
func RecordMonthWithoutAdjustment(ctx context.Context, db *sql.DB, yearMonth string) error {
	const insertStmt = `INSERT OR IGNORE INTO monthly_adjustments_log(year_month) VALUES (?);`

	_, err := db.ExecContext(ctx, insertStmt, yearMonth)
	if err != nil {
		return fmt.Errorf("could not record initial month %s: %w", yearMonth, err)
	}

	return nil
}
