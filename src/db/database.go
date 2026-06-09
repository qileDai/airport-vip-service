package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

func InitDB(dataSourceName string) (*sql.DB, error) {
	dir := filepath.Dir(dataSourceName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if err = createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	tables := []string{
		createMemberBenefitTable(),
		createReservationRecordTable(),
		createFlightScheduleTable(),
		createCompanionTable(),
		createUsageVoucherTable(),
		createWaitlistEntryTable(),
		createVerificationResultTable(),
		createStatusTransitionTable(),
		createRuleConfigTable(),
		createExceptionEventTable(),
		createAuditLogTable(),
	}

	for _, table := range tables {
		if _, err := db.Exec(table); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	return nil
}

func createMemberBenefitTable() string {
	return `
	CREATE TABLE IF NOT EXISTS member_benefits (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		benefit_no TEXT NOT NULL UNIQUE,
		member_name TEXT NOT NULL,
		member_level TEXT NOT NULL,
		remaining_quota INTEGER NOT NULL DEFAULT 0,
		total_quota INTEGER NOT NULL DEFAULT 0,
		status TEXT NOT NULL DEFAULT 'active',
		responsible_person TEXT,
		valid_from DATETIME NOT NULL,
		valid_to DATETIME NOT NULL,
		batch_no TEXT,
		remarks TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_member_benefits_benefit_no ON member_benefits(benefit_no);
	CREATE INDEX IF NOT EXISTS idx_member_benefits_status ON member_benefits(status);
	CREATE INDEX IF NOT EXISTS idx_member_benefits_member_level ON member_benefits(member_level);
	`
}

func createReservationRecordTable() string {
	return `
	CREATE TABLE IF NOT EXISTS reservation_records (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		reservation_no TEXT NOT NULL UNIQUE,
		member_benefit_id INTEGER NOT NULL,
		member_name TEXT NOT NULL,
		flight_no TEXT NOT NULL,
		flight_schedule_id INTEGER,
		vip_lounge_name TEXT NOT NULL,
		reservation_time DATETIME NOT NULL,
		guest_count INTEGER NOT NULL DEFAULT 1,
		status TEXT NOT NULL DEFAULT 'draft',
		responsible_person TEXT,
		batch_no TEXT,
		remarks TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (member_benefit_id) REFERENCES member_benefits(id),
		FOREIGN KEY (flight_schedule_id) REFERENCES flight_schedules(id)
	);
	CREATE INDEX IF NOT EXISTS idx_reservation_records_reservation_no ON reservation_records(reservation_no);
	CREATE INDEX IF NOT EXISTS idx_reservation_records_status ON reservation_records(status);
	CREATE INDEX IF NOT EXISTS idx_reservation_records_member_benefit_id ON reservation_records(member_benefit_id);
	`
}

func createFlightScheduleTable() string {
	return `
	CREATE TABLE IF NOT EXISTS flight_schedules (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		flight_no TEXT NOT NULL,
		departure_airport TEXT NOT NULL,
		arrival_airport TEXT NOT NULL,
		scheduled_depart DATETIME NOT NULL,
		scheduled_arrive DATETIME NOT NULL,
		actual_depart DATETIME,
		actual_arrive DATETIME,
		flight_status TEXT NOT NULL DEFAULT 'scheduled',
		vip_lounge_capacity INTEGER NOT NULL DEFAULT 50,
		vip_lounge_used INTEGER NOT NULL DEFAULT 0,
		status TEXT NOT NULL DEFAULT 'active',
		responsible_person TEXT,
		batch_no TEXT,
		remarks TEXT,
		snapshot_data TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_flight_schedules_flight_no ON flight_schedules(flight_no);
	CREATE INDEX IF NOT EXISTS idx_flight_schedules_status ON flight_schedules(status);
	CREATE INDEX IF NOT EXISTS idx_flight_schedules_scheduled_depart ON flight_schedules(scheduled_depart);
	`
}

func createCompanionTable() string {
	return `
	CREATE TABLE IF NOT EXISTS companions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		companion_no TEXT NOT NULL UNIQUE,
		reservation_id INTEGER NOT NULL,
		companion_name TEXT NOT NULL,
		companion_id_card TEXT,
		relation_type TEXT NOT NULL,
		is_vip_eligible INTEGER NOT NULL DEFAULT 0,
		verification_status TEXT NOT NULL DEFAULT 'pending',
		status TEXT NOT NULL DEFAULT 'active',
		responsible_person TEXT,
		batch_no TEXT,
		remarks TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (reservation_id) REFERENCES reservation_records(id)
	);
	CREATE INDEX IF NOT EXISTS idx_companions_companion_no ON companions(companion_no);
	CREATE INDEX IF NOT EXISTS idx_companions_reservation_id ON companions(reservation_id);
	CREATE INDEX IF NOT EXISTS idx_companions_status ON companions(status);
	`
}

func createUsageVoucherTable() string {
	return `
	CREATE TABLE IF NOT EXISTS usage_vouchers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		voucher_no TEXT NOT NULL UNIQUE,
		reservation_id INTEGER NOT NULL,
		member_benefit_id INTEGER NOT NULL,
		voucher_type TEXT NOT NULL,
		qr_code TEXT,
		valid_from DATETIME NOT NULL,
		valid_to DATETIME NOT NULL,
		used_at DATETIME,
		used_location TEXT,
		verification_status TEXT NOT NULL DEFAULT 'valid',
		status TEXT NOT NULL DEFAULT 'active',
		responsible_person TEXT,
		batch_no TEXT,
		remarks TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (reservation_id) REFERENCES reservation_records(id),
		FOREIGN KEY (member_benefit_id) REFERENCES member_benefits(id)
	);
	CREATE INDEX IF NOT EXISTS idx_usage_vouchers_voucher_no ON usage_vouchers(voucher_no);
	CREATE INDEX IF NOT EXISTS idx_usage_vouchers_reservation_id ON usage_vouchers(reservation_id);
	CREATE INDEX IF NOT EXISTS idx_usage_vouchers_status ON usage_vouchers(status);
	`
}

func createWaitlistEntryTable() string {
	return `
	CREATE TABLE IF NOT EXISTS waitlist_entries (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		waitlist_no TEXT NOT NULL UNIQUE,
		reservation_id INTEGER,
		member_benefit_id INTEGER NOT NULL,
		flight_schedule_id INTEGER NOT NULL,
		member_name TEXT NOT NULL,
		member_level TEXT NOT NULL,
		waiting_since DATETIME NOT NULL,
		priority_score INTEGER NOT NULL DEFAULT 0,
		estimated_wait_mins INTEGER NOT NULL DEFAULT 0,
		status TEXT NOT NULL DEFAULT 'waiting',
		responsible_person TEXT,
		batch_no TEXT,
		remarks TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (reservation_id) REFERENCES reservation_records(id),
		FOREIGN KEY (member_benefit_id) REFERENCES member_benefits(id),
		FOREIGN KEY (flight_schedule_id) REFERENCES flight_schedules(id)
	);
	CREATE INDEX IF NOT EXISTS idx_waitlist_entries_waitlist_no ON waitlist_entries(waitlist_no);
	CREATE INDEX IF NOT EXISTS idx_waitlist_entries_status ON waitlist_entries(status);
	CREATE INDEX IF NOT EXISTS idx_waitlist_entries_priority_score ON waitlist_entries(priority_score);
	`
}

func createVerificationResultTable() string {
	return `
	CREATE TABLE IF NOT EXISTS verification_results (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		verification_no TEXT NOT NULL UNIQUE,
		reservation_id INTEGER NOT NULL,
		member_benefit_id INTEGER NOT NULL,
		flight_schedule_id INTEGER,
		verification_type TEXT NOT NULL,
		result TEXT NOT NULL DEFAULT 'pending',
		failure_reason TEXT,
		verified_quota INTEGER NOT NULL DEFAULT 0,
		verified_companions INTEGER NOT NULL DEFAULT 0,
		verification_details TEXT,
		status TEXT NOT NULL DEFAULT 'draft',
		responsible_person TEXT,
		batch_no TEXT,
		remarks TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (reservation_id) REFERENCES reservation_records(id),
		FOREIGN KEY (member_benefit_id) REFERENCES member_benefits(id),
		FOREIGN KEY (flight_schedule_id) REFERENCES flight_schedules(id)
	);
	CREATE INDEX IF NOT EXISTS idx_verification_results_verification_no ON verification_results(verification_no);
	CREATE INDEX IF NOT EXISTS idx_verification_results_status ON verification_results(status);
	CREATE INDEX IF NOT EXISTS idx_verification_results_result ON verification_results(result);
	`
}

func createStatusTransitionTable() string {
	return `
	CREATE TABLE IF NOT EXISTS status_transitions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		transition_no TEXT NOT NULL UNIQUE,
		entity_type TEXT NOT NULL,
		entity_id INTEGER NOT NULL,
		from_status TEXT NOT NULL,
		to_status TEXT NOT NULL,
		action TEXT NOT NULL,
		reason TEXT,
		operator TEXT NOT NULL,
		status TEXT NOT NULL DEFAULT 'active',
		responsible_person TEXT,
		batch_no TEXT,
		remarks TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_status_transitions_transition_no ON status_transitions(transition_no);
	CREATE INDEX IF NOT EXISTS idx_status_transitions_entity_type ON status_transitions(entity_type);
	CREATE INDEX IF NOT EXISTS idx_status_transitions_entity_id ON status_transitions(entity_id);
	`
}

func createRuleConfigTable() string {
	return `
	CREATE TABLE IF NOT EXISTS rule_configs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		rule_no TEXT NOT NULL UNIQUE,
		rule_name TEXT NOT NULL,
		rule_type TEXT NOT NULL,
		rule_value TEXT NOT NULL,
		threshold_value REAL NOT NULL DEFAULT 0,
		applies_to_level TEXT,
		is_active INTEGER NOT NULL DEFAULT 1,
		effective_date DATETIME NOT NULL,
		expiry_date DATETIME,
		status TEXT NOT NULL DEFAULT 'active',
		responsible_person TEXT,
		batch_no TEXT,
		remarks TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_rule_configs_rule_no ON rule_configs(rule_no);
	CREATE INDEX IF NOT EXISTS idx_rule_configs_rule_type ON rule_configs(rule_type);
	CREATE INDEX IF NOT EXISTS idx_rule_configs_status ON rule_configs(status);
	`
}

func createExceptionEventTable() string {
	return `
	CREATE TABLE IF NOT EXISTS exception_events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		event_no TEXT NOT NULL UNIQUE,
		event_type TEXT NOT NULL,
		entity_type TEXT NOT NULL,
		entity_id INTEGER NOT NULL,
		trigger_field TEXT NOT NULL,
		threshold_value TEXT,
		actual_value TEXT,
		severity TEXT NOT NULL DEFAULT 'medium',
		handler TEXT,
		handling_deadline DATETIME,
		handled_at DATETIME,
		handling_result TEXT,
		status TEXT NOT NULL DEFAULT 'open',
		responsible_person TEXT,
		batch_no TEXT,
		remarks TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_exception_events_event_no ON exception_events(event_no);
	CREATE INDEX IF NOT EXISTS idx_exception_events_event_type ON exception_events(event_type);
	CREATE INDEX IF NOT EXISTS idx_exception_events_status ON exception_events(status);
	CREATE INDEX IF NOT EXISTS idx_exception_events_severity ON exception_events(severity);
	`
}

func createAuditLogTable() string {
	return `
	CREATE TABLE IF NOT EXISTS audit_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		log_no TEXT NOT NULL UNIQUE,
		operation_type TEXT NOT NULL,
		entity_type TEXT NOT NULL,
		entity_id INTEGER NOT NULL,
		entity_no TEXT NOT NULL,
		old_value TEXT,
		new_value TEXT,
		operator TEXT NOT NULL,
		operator_role TEXT,
		ip_address TEXT,
		user_agent TEXT,
		request_id TEXT,
		status TEXT NOT NULL DEFAULT 'active',
		responsible_person TEXT,
		batch_no TEXT,
		remarks TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_audit_logs_log_no ON audit_logs(log_no);
	CREATE INDEX IF NOT EXISTS idx_audit_logs_operation_type ON audit_logs(operation_type);
	CREATE INDEX IF NOT EXISTS idx_audit_logs_entity_type ON audit_logs(entity_type);
	CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at);
	`
}
