package tests

import (
	"airport-vip-service/src/db"
	"airport-vip-service/src/seed"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/labstack/echo/v4"
)

var testDB *sql.DB

func setupTestDB(t *testing.T) *sql.DB {
	if testDB == nil {
		var err error
		testDB, err = sql.Open("sqlite3", ":memory:")
		if err != nil {
			t.Fatalf("Failed to open test database: %v", err)
		}
		db.InitTables(testDB)
		seed.RunSeed(testDB)
	}
	return testDB
}

func TestHealthEndpoint(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "ok",
			"version": "1.0.0",
		})
	}

	if err := handler(c); err != nil {
		t.Errorf("Health endpoint failed: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestMemberBenefitRepository(t *testing.T) {
	database := setupTestDB(t)

	repo := db.NewMemberBenefitRepository(database)

	benefits, err := repo.List(1, 10, "", "")
	if err != nil {
		t.Errorf("Failed to list member benefits: %v", err)
	}

	if len(benefits) == 0 {
		t.Error("Expected at least one member benefit from seed data")
	}

	benefit, err := repo.GetByID(1)
	if err != nil {
		t.Errorf("Failed to get member benefit by ID: %v", err)
	}

	if benefit == nil {
		t.Error("Expected member benefit with ID 1 to exist")
	}
}

func TestReservationRepository(t *testing.T) {
	database := setupTestDB(t)

	repo := db.NewReservationRecordRepository(database)

	reservations, err := repo.List(1, 10, "", 0)
	if err != nil {
		t.Errorf("Failed to list reservations: %v", err)
	}

	if len(reservations) == 0 {
		t.Error("Expected at least one reservation from seed data")
	}
}

func TestFlightScheduleRepository(t *testing.T) {
	database := setupTestDB(t)

	repo := db.NewFlightScheduleRepository(database)

	schedules, err := repo.List(1, 10, "")
	if err != nil {
		t.Errorf("Failed to list flight schedules: %v", err)
	}

	if len(schedules) == 0 {
		t.Error("Expected at least one flight schedule from seed data")
	}
}

func TestWaitlistRepository(t *testing.T) {
	database := setupTestDB(t)

	repo := db.NewWaitlistEntryRepository(database)

	entries, err := repo.List(1, 10, 0, "")
	if err != nil {
		t.Errorf("Failed to list waitlist entries: %v", err)
	}

	if len(entries) == 0 {
		t.Error("Expected at least one waitlist entry from seed data")
	}
}

func TestVerificationRepository(t *testing.T) {
	database := setupTestDB(t)

	repo := db.NewVerificationResultRepository(database)

	results, err := repo.List(1, 10, "", "")
	if err != nil {
		t.Errorf("Failed to list verification results: %v", err)
	}

	if len(results) == 0 {
		t.Error("Expected at least one verification result from seed data")
	}
}

func TestExceptionRepository(t *testing.T) {
	database := setupTestDB(t)

	repo := db.NewExceptionEventRepository(database)

	events, err := repo.List(1, 10, "", "", "")
	if err != nil {
		t.Errorf("Failed to list exception events: %v", err)
	}

	if len(events) == 0 {
		t.Error("Expected at least one exception event from seed data")
	}
}

func TestAuditLogRepository(t *testing.T) {
	database := setupTestDB(t)

	repo := db.NewAuditLogRepository(database)

	logs, err := repo.List(1, 10, "", "", "")
	if err != nil {
		t.Errorf("Failed to list audit logs: %v", err)
	}

	if len(logs) == 0 {
		t.Error("Expected at least one audit log from seed data")
	}
}

func TestRuleConfigRepository(t *testing.T) {
	database := setupTestDB(t)

	repo := db.NewRuleConfigRepository(database)

	rules, err := repo.List(1, 10, "", "")
	if err != nil {
		t.Errorf("Failed to list rule configs: %v", err)
	}

	if len(rules) == 0 {
		t.Error("Expected at least one rule config from seed data")
	}
}
