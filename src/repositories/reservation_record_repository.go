package repositories

import (
	"airport-vip-service/src/models"
	"database/sql"
	"fmt"
	"time"
)

type ReservationRecordRepository struct {
	db *sql.DB
}

func NewReservationRecordRepository(db *sql.DB) *ReservationRecordRepository {
	return &ReservationRecordRepository{db: db}
}

func (r *ReservationRecordRepository) Create(record *models.ReservationRecord) (int64, error) {
	query := `INSERT INTO reservation_records 
		(reservation_no, member_benefit_id, member_name, flight_no, flight_schedule_id, vip_lounge_name, 
		reservation_time, guest_count, status, responsible_person, batch_no, remarks, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := r.db.Exec(query,
		record.ReservationNo, record.MemberBenefitID, record.MemberName, record.FlightNo,
		record.FlightScheduleID, record.VipLoungeName, record.ReservationTime,
		record.GuestCount, record.Status, record.ResponsiblePerson,
		record.BatchNo, record.Remarks, now, now)
	if err != nil {
		return 0, fmt.Errorf("failed to create reservation record: %w", err)
	}

	return result.LastInsertId()
}

func (r *ReservationRecordRepository) GetByID(id int64) (*models.ReservationRecord, error) {
	query := `SELECT id, reservation_no, member_benefit_id, member_name, flight_no, flight_schedule_id,
		vip_lounge_name, reservation_time, guest_count, status, responsible_person, batch_no, remarks, 
		created_at, updated_at 
		FROM reservation_records WHERE id = ?`

	record := &models.ReservationRecord{}
	var flightScheduleID sql.NullInt64
	err := r.db.QueryRow(query, id).Scan(
		&record.ID, &record.ReservationNo, &record.MemberBenefitID, &record.MemberName,
		&record.FlightNo, &flightScheduleID, &record.VipLoungeName, &record.ReservationTime,
		&record.GuestCount, &record.Status, &record.ResponsiblePerson,
		&record.BatchNo, &record.Remarks, &record.CreatedAt, &record.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get reservation record: %w", err)
	}

	if flightScheduleID.Valid {
		record.FlightScheduleID = flightScheduleID.Int64
	}

	return record, nil
}

func (r *ReservationRecordRepository) GetByReservationNo(reservationNo string) (*models.ReservationRecord, error) {
	query := `SELECT id, reservation_no, member_benefit_id, member_name, flight_no, flight_schedule_id,
		vip_lounge_name, reservation_time, guest_count, status, responsible_person, batch_no, remarks, 
		created_at, updated_at 
		FROM reservation_records WHERE reservation_no = ?`

	record := &models.ReservationRecord{}
	var flightScheduleID sql.NullInt64
	err := r.db.QueryRow(query, reservationNo).Scan(
		&record.ID, &record.ReservationNo, &record.MemberBenefitID, &record.MemberName,
		&record.FlightNo, &flightScheduleID, &record.VipLoungeName, &record.ReservationTime,
		&record.GuestCount, &record.Status, &record.ResponsiblePerson,
		&record.BatchNo, &record.Remarks, &record.CreatedAt, &record.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get reservation record: %w", err)
	}

	if flightScheduleID.Valid {
		record.FlightScheduleID = flightScheduleID.Int64
	}

	return record, nil
}

func (r *ReservationRecordRepository) List(offset, limit int, status string, memberBenefitID int64) ([]models.ReservationRecord, int64, error) {
	var args []interface{}
	whereClause := "WHERE 1=1"

	if status != "" {
		whereClause += " AND status = ?"
		args = append(args, status)
	}
	if memberBenefitID > 0 {
		whereClause += " AND member_benefit_id = ?"
		args = append(args, memberBenefitID)
	}

	countQuery := "SELECT COUNT(*) FROM reservation_records " + whereClause
	var total int64
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count reservation records: %w", err)
	}

	query := `SELECT id, reservation_no, member_benefit_id, member_name, flight_no, flight_schedule_id,
		vip_lounge_name, reservation_time, guest_count, status, responsible_person, batch_no, remarks, 
		created_at, updated_at 
		FROM reservation_records ` + whereClause + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list reservation records: %w", err)
	}
	defer rows.Close()

	var records []models.ReservationRecord
	for rows.Next() {
		var record models.ReservationRecord
		var flightScheduleID sql.NullInt64
		err := rows.Scan(
			&record.ID, &record.ReservationNo, &record.MemberBenefitID, &record.MemberName,
			&record.FlightNo, &flightScheduleID, &record.VipLoungeName, &record.ReservationTime,
			&record.GuestCount, &record.Status, &record.ResponsiblePerson,
			&record.BatchNo, &record.Remarks, &record.CreatedAt, &record.UpdatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan reservation record: %w", err)
		}
		if flightScheduleID.Valid {
			record.FlightScheduleID = flightScheduleID.Int64
		}
		records = append(records, record)
	}

	return records, total, nil
}

func (r *ReservationRecordRepository) Update(id int64, record *models.ReservationRecord) error {
	query := `UPDATE reservation_records SET 
		member_name = ?, flight_no = ?, flight_schedule_id = ?, vip_lounge_name = ?,
		reservation_time = ?, guest_count = ?, status = ?, responsible_person = ?,
		batch_no = ?, remarks = ?, updated_at = ?
		WHERE id = ?`

	_, err := r.db.Exec(query,
		record.MemberName, record.FlightNo, record.FlightScheduleID, record.VipLoungeName,
		record.ReservationTime, record.GuestCount, record.Status, record.ResponsiblePerson,
		record.BatchNo, record.Remarks, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update reservation record: %w", err)
	}

	return nil
}

func (r *ReservationRecordRepository) UpdateStatus(id int64, status string) error {
	query := `UPDATE reservation_records SET status = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update reservation record status: %w", err)
	}
	return nil
}

func (r *ReservationRecordRepository) Delete(id int64) error {
	query := `DELETE FROM reservation_records WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete reservation record: %w", err)
	}
	return nil
}

func (r *ReservationRecordRepository) BatchCreate(records []models.ReservationRecord) ([]int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	query := `INSERT INTO reservation_records 
		(reservation_no, member_benefit_id, member_name, flight_no, flight_schedule_id, vip_lounge_name, 
		reservation_time, guest_count, status, responsible_person, batch_no, remarks, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	var ids []int64
	now := time.Now()
	for _, record := range records {
		result, err := tx.Exec(query,
			record.ReservationNo, record.MemberBenefitID, record.MemberName, record.FlightNo,
			record.FlightScheduleID, record.VipLoungeName, record.ReservationTime,
			record.GuestCount, record.Status, record.ResponsiblePerson,
			record.BatchNo, record.Remarks, now, now)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create reservation record: %w", err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to get last insert id: %w", err)
		}
		ids = append(ids, id)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return ids, nil
}

func (r *ReservationRecordRepository) GetByDateRange(startDate, endDate time.Time) ([]models.ReservationRecord, error) {
	query := `SELECT id, reservation_no, member_benefit_id, member_name, flight_no, flight_schedule_id,
		vip_lounge_name, reservation_time, guest_count, status, responsible_person, batch_no, remarks, 
		created_at, updated_at 
		FROM reservation_records WHERE reservation_time BETWEEN ? AND ? ORDER BY reservation_time`

	rows, err := r.db.Query(query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get reservation records by date range: %w", err)
	}
	defer rows.Close()

	var records []models.ReservationRecord
	for rows.Next() {
		var record models.ReservationRecord
		var flightScheduleID sql.NullInt64
		err := rows.Scan(
			&record.ID, &record.ReservationNo, &record.MemberBenefitID, &record.MemberName,
			&record.FlightNo, &flightScheduleID, &record.VipLoungeName, &record.ReservationTime,
			&record.GuestCount, &record.Status, &record.ResponsiblePerson,
			&record.BatchNo, &record.Remarks, &record.CreatedAt, &record.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reservation record: %w", err)
		}
		if flightScheduleID.Valid {
			record.FlightScheduleID = flightScheduleID.Int64
		}
		records = append(records, record)
	}

	return records, nil
}

func (r *ReservationRecordRepository) GetByBatchNo(batchNo string) ([]models.ReservationRecord, error) {
	query := `SELECT id, reservation_no, member_benefit_id, member_name, flight_no, flight_schedule_id,
		vip_lounge_name, reservation_time, guest_count, status, responsible_person, batch_no, remarks, 
		created_at, updated_at 
		FROM reservation_records WHERE batch_no = ? ORDER BY created_at`

	rows, err := r.db.Query(query, batchNo)
	if err != nil {
		return nil, fmt.Errorf("failed to get reservation records by batch: %w", err)
	}
	defer rows.Close()

	var records []models.ReservationRecord
	for rows.Next() {
		var record models.ReservationRecord
		var flightScheduleID sql.NullInt64
		err := rows.Scan(
			&record.ID, &record.ReservationNo, &record.MemberBenefitID, &record.MemberName,
			&record.FlightNo, &flightScheduleID, &record.VipLoungeName, &record.ReservationTime,
			&record.GuestCount, &record.Status, &record.ResponsiblePerson,
			&record.BatchNo, &record.Remarks, &record.CreatedAt, &record.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reservation record: %w", err)
		}
		if flightScheduleID.Valid {
			record.FlightScheduleID = flightScheduleID.Int64
		}
		records = append(records, record)
	}

	return records, nil
}
