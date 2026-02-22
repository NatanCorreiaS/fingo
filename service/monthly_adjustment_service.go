package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"natan/fingo/dbsqlite"
)

// monthsBetween returns all year-month strings (format "YYYY-MM") from the month
// after `from` up to and including `to`. If `from` is empty, returns only `to`.
func monthsBetween(from, to string) ([]string, error) {
	if from == "" {
		return []string{to}, nil
	}

	fromTime, err := time.Parse("2006-01", from)
	if err != nil {
		return nil, fmt.Errorf("could not parse 'from' month %q: %w", from, err)
	}

	toTime, err := time.Parse("2006-01", to)
	if err != nil {
		return nil, fmt.Errorf("could not parse 'to' month %q: %w", to, err)
	}

	var months []string
	// Start from the month after `from`
	current := fromTime.AddDate(0, 1, 0)
	for !current.After(toTime) {
		months = append(months, current.Format("2006-01"))
		current = current.AddDate(0, 1, 0)
	}

	return months, nil
}

// currentYearMonth returns the current year-month string in "YYYY-MM" format.
func currentYearMonth() string {
	return time.Now().Format("2006-01")
}

// ProcessPendingAdjustments checks for any unprocessed months and applies
// the monthly adjustment for each one. On first run (no log entries), it
// records the current month without applying adjustments to establish a baseline.
func ProcessPendingAdjustments() error {
	db, err := dbsqlite.GetDatabaseConnection()
	if err != nil {
		return fmt.Errorf("could not get database connection: %w", err)
	}
	defer db.Close()

	ctx, cancel := dbsqlite.NewDBContext()
	defer cancel()

	lastProcessed, err := dbsqlite.GetLastProcessedMonth(ctx, db)
	if err != nil {
		return fmt.Errorf("could not get last processed month: %w", err)
	}

	now := currentYearMonth()

	// First run: no records exist yet. Record the current month as baseline
	// without applying any adjustment to avoid double-counting.
	if lastProcessed == "" {
		log.Printf("[MonthlyAdjustment] First run detected. Recording %s as baseline (no adjustment applied).", now)
		if err := dbsqlite.RecordMonthWithoutAdjustment(ctx, db, now); err != nil {
			return fmt.Errorf("could not record baseline month: %w", err)
		}
		return nil
	}

	// Already up to date
	if lastProcessed == now {
		log.Printf("[MonthlyAdjustment] Already up to date (last processed: %s).", lastProcessed)
		return nil
	}

	// Calculate all months that need processing
	pendingMonths, err := monthsBetween(lastProcessed, now)
	if err != nil {
		return fmt.Errorf("could not calculate pending months: %w", err)
	}

	if len(pendingMonths) == 0 {
		log.Printf("[MonthlyAdjustment] No pending months to process.")
		return nil
	}

	log.Printf("[MonthlyAdjustment] Processing %d pending month(s): %v", len(pendingMonths), pendingMonths)

	for _, month := range pendingMonths {
		// Each adjustment gets its own context to avoid timeout issues with many months
		adjCtx, adjCancel := dbsqlite.NewDBContext()

		if err := dbsqlite.ApplyMonthlyAdjustment(adjCtx, db, month); err != nil {
			adjCancel()
			return fmt.Errorf("could not apply adjustment for %s: %w", month, err)
		}

		log.Printf("[MonthlyAdjustment] Successfully applied adjustment for %s.", month)
		adjCancel()
	}

	log.Printf("[MonthlyAdjustment] All pending adjustments applied successfully.")
	return nil
}

// StartMonthlyAdjustmentScheduler starts a background goroutine that periodically
// checks whether a new month has started and applies the monthly adjustment.
// It checks every hour. The goroutine stops when the provided context is canceled.
func StartMonthlyAdjustmentScheduler(ctx context.Context) {
	go func() {
		// Run immediately on startup
		if err := ProcessPendingAdjustments(); err != nil {
			log.Printf("[MonthlyAdjustment] Error during startup processing: %v", err)
		}

		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		log.Println("[MonthlyAdjustment] Scheduler started. Checking every hour for pending adjustments.")

		for {
			select {
			case <-ctx.Done():
				log.Println("[MonthlyAdjustment] Scheduler stopped.")
				return
			case <-ticker.C:
				if err := ProcessPendingAdjustments(); err != nil {
					log.Printf("[MonthlyAdjustment] Error during periodic check: %v", err)
				}
			}
		}
	}()
}
