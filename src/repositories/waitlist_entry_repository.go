package repositories

import (
	"airport-vip-service/src/models"
	"database/sql"
	"fmt"
	"time"
)

type WaitlistEntryRepository struct {
	db *sql.DB
}

func NewWaitlistEntryRepository(db *sql.DB) *WaitlistEntryRepository {
	return &WaitlistEntryRepository{db: db}
}

func (r *WaitlistEntryRepository) Create(entry *models.WaitlistEntry) (int64, error) {
	query := `INSERT INTO waitlist_entries 
		(waitlist_no, reservation_id, member_benefit_id, flight_schedule_id, member_name, member_level,
		waiting_since, priority_score, estimated_wait_mins, status, responsible_person, batch_no, remarks, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := r.db.Exec(query,
		entry.WaitlistNo, entry.ReservationID, entry.MemberBenefitID, entry.FlightScheduleID,
		entry.MemberName, entry.MemberLevel, entry.WaitingSince, entry.PriorityScore,
		entry.EstimatedWaitMins, entry.Status, entry.ResponsiblePerson,
		entry.BatchNo, entry.Remarks, now, now)
	if err != nil {
		return 0, fmt.Errorf("failed to create waitlist entry: %w", err)
	}

	return result.LastInsertId()
}

func (r *WaitlistEntryRepository) GetByID(id int64) (*models.WaitlistEntry, error) {
	query := `SELECT id, waitlist_no, reservation_id, member_benefit_id, flight_schedule_id, member_name, member_level,
		waiting_since, priority_score, estimated_wait_mins, status, responsible_person, batch_no, remarks, created_at, updated_at 
		FROM waitlist_entries WHERE id = ?`

	entry := &models.WaitlistEntry{}
	var reservationID sql.NullInt64
	err := r.db.QueryRow(query, id).Scan(
		&entry.ID, &entry.WaitlistNo, &reservationID, &entry.MemberBenefitID, &entry.FlightScheduleID,
		&entry.MemberName, &entry.MemberLevel, &entry.WaitingSince, &entry.PriorityScore,
		&entry.EstimatedWaitMins, &entry.Status, &entry.ResponsiblePerson,
		&entry.BatchNo, &entry.Remarks, &entry.CreatedAt, &entry.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get waitlist entry: %w", err)
	}
	if reservationID.Valid {
		entry.ReservationID = reservationID.Int64
	}

	return entry, nil
}

func (r *WaitlistEntryRepository) List(offset, limit int, flightScheduleID int64, status string) ([]models.WaitlistEntry, int64, error) {
	var args []interface{}
	whereClause := "WHERE 1=1"

	if flightScheduleID > 0 {
		whereClause += " AND flight_schedule_id = ?"
		args = append(args, flightScheduleID)
	}
	if status != "" {
		whereClause += " AND status = ?"
		args = append(args, status)
	}

	countQuery := "SELECT COUNT(*) FROM waitlist_entries " + whereClause
	var total int64
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count waitlist entries: %w", err)
	}

	query := `SELECT id, waitlist_no, reservation_id, member_benefit_id, flight_schedule_id, member_name, member_level,
		waiting_since, priority_score, estimated_wait_mins, status, responsible_person, batch_no, remarks, created_at, updated_at 
		FROM waitlist_entries ` + whereClause + ` ORDER BY priority_score DESC, waiting_since ASC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list waitlist entries: %w", err)
	}
	defer rows.Close()

	var entries []models.WaitlistEntry
	for rows.Next() {
		var entry models.WaitlistEntry
		var reservationID sql.NullInt64
		err := rows.Scan(
			&entry.ID, &entry.WaitlistNo, &reservationID, &entry.MemberBenefitID, &entry.FlightScheduleID,
			&entry.MemberName, &entry.MemberLevel, &entry.WaitingSince, &entry.PriorityScore,
			&entry.EstimatedWaitMins, &entry.Status, &entry.ResponsiblePerson,
			&entry.BatchNo, &entry.Remarks, &entry.CreatedAt, &entry.UpdatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan waitlist entry: %w", err)
		}
		if reservationID.Valid {
			entry.ReservationID = reservationID.Int64
		}
		entries = append(entries, entry)
	}

	return entries, total, nil
}

func (r *WaitlistEntryRepository) Update(id int64, entry *models.WaitlistEntry) error {
	query := `UPDATE waitlist_entries SET 
		priority_score = ?, estimated_wait_mins = ?, status = ?, responsible_person = ?,
		batch_no = ?, remarks = ?, updated_at = ?
		WHERE id = ?`

	_, err := r.db.Exec(query,
		entry.PriorityScore, entry.EstimatedWaitMins, entry.Status, entry.ResponsiblePerson,
		entry.BatchNo, entry.Remarks, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update waitlist entry: %w", err)
	}

	return nil
}

func (r *WaitlistEntryRepository) UpdateStatus(id int64, status string) error {
	query := `UPDATE waitlist_entries SET status = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update waitlist entry status: %w", err)
	}
	return nil
}

func (r *WaitlistEntryRepository) Delete(id int64) error {
	query := `DELETE FROM waitlist_entries WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete waitlist entry: %w", err)
	}
	return nil
}

func (r *WaitlistEntryRepository) GetWaitingByFlightScheduleID(flightScheduleID int64) ([]models.WaitlistEntry, error) {
	query := `SELECT id, waitlist_no, reservation_id, member_benefit_id, flight_schedule_id, member_name, member_level,
		waiting_since, priority_score, estimated_wait_mins, status, responsible_person, batch_no, remarks, created_at, updated_at 
		FROM waitlist_entries WHERE flight_schedule_id = ? AND status = 'waiting' 
		ORDER BY priority_score DESC, waiting_since ASC`

	rows, err := r.db.Query(query, flightScheduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get waiting entries: %w", err)
	}
	defer rows.Close()

	var entries []models.WaitlistEntry
	for rows.Next() {
		var entry models.WaitlistEntry
		var reservationID sql.NullInt64
		err := rows.Scan(
			&entry.ID, &entry.WaitlistNo, &reservationID, &entry.MemberBenefitID, &entry.FlightScheduleID,
			&entry.MemberName, &entry.MemberLevel, &entry.WaitingSince, &entry.PriorityScore,
			&entry.EstimatedWaitMins, &entry.Status, &entry.ResponsiblePerson,
			&entry.BatchNo, &entry.Remarks, &entry.CreatedAt, &entry.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan waitlist entry: %w", err)
		}
		if reservationID.Valid {
			entry.ReservationID = reservationID.Int64
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func (r *WaitlistEntryRepository) UpdatePriorityScore(id int64, score int) error {
	query := `UPDATE waitlist_entries SET priority_score = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, score, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update priority score: %w", err)
	}
	return nil
}

func (r *WaitlistEntryRepository) CountByStatus(status string) (int64, error) {
	query := `SELECT COUNT(*) FROM waitlist_entries WHERE status = ?`
	var count int64
	err := r.db.QueryRow(query, status).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count waitlist entries: %w", err)
	}
	return count, nil
}
