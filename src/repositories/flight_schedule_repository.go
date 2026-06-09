package repositories

import (
	"airport-vip-service/src/models"
	"database/sql"
	"fmt"
	"time"
)

type FlightScheduleRepository struct {
	db *sql.DB
}

func NewFlightScheduleRepository(db *sql.DB) *FlightScheduleRepository {
	return &FlightScheduleRepository{db: db}
}

func (r *FlightScheduleRepository) Create(schedule *models.FlightSchedule) (int64, error) {
	query := `INSERT INTO flight_schedules 
		(flight_no, departure_airport, arrival_airport, scheduled_depart, scheduled_arrive,
		actual_depart, actual_arrive, flight_status, vip_lounge_capacity, vip_lounge_used,
		status, responsible_person, batch_no, remarks, snapshot_data, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := r.db.Exec(query,
		schedule.FlightNo, schedule.DepartureAirport, schedule.ArrivalAirport,
		schedule.ScheduledDepart, schedule.ScheduledArrive,
		schedule.ActualDepart, schedule.ActualArrive, schedule.FlightStatus,
		schedule.VipLoungeCapacity, schedule.VipLoungeUsed,
		schedule.Status, schedule.ResponsiblePerson, schedule.BatchNo,
		schedule.Remarks, schedule.SnapshotData, now, now)
	if err != nil {
		return 0, fmt.Errorf("failed to create flight schedule: %w", err)
	}

	return result.LastInsertId()
}

func (r *FlightScheduleRepository) GetByID(id int64) (*models.FlightSchedule, error) {
	query := `SELECT id, flight_no, departure_airport, arrival_airport, scheduled_depart, scheduled_arrive,
		actual_depart, actual_arrive, flight_status, vip_lounge_capacity, vip_lounge_used,
		status, responsible_person, batch_no, remarks, snapshot_data, created_at, updated_at 
		FROM flight_schedules WHERE id = ?`

	schedule := &models.FlightSchedule{}
	err := r.db.QueryRow(query, id).Scan(
		&schedule.ID, &schedule.FlightNo, &schedule.DepartureAirport, &schedule.ArrivalAirport,
		&schedule.ScheduledDepart, &schedule.ScheduledArrive,
		&schedule.ActualDepart, &schedule.ActualArrive, &schedule.FlightStatus,
		&schedule.VipLoungeCapacity, &schedule.VipLoungeUsed,
		&schedule.Status, &schedule.ResponsiblePerson, &schedule.BatchNo,
		&schedule.Remarks, &schedule.SnapshotData, &schedule.CreatedAt, &schedule.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get flight schedule: %w", err)
	}

	return schedule, nil
}

func (r *FlightScheduleRepository) List(offset, limit int, status string) ([]models.FlightSchedule, int64, error) {
	var args []interface{}
	whereClause := "WHERE 1=1"

	if status != "" {
		whereClause += " AND status = ?"
		args = append(args, status)
	}

	countQuery := "SELECT COUNT(*) FROM flight_schedules " + whereClause
	var total int64
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count flight schedules: %w", err)
	}

	query := `SELECT id, flight_no, departure_airport, arrival_airport, scheduled_depart, scheduled_arrive,
		actual_depart, actual_arrive, flight_status, vip_lounge_capacity, vip_lounge_used,
		status, responsible_person, batch_no, remarks, snapshot_data, created_at, updated_at 
		FROM flight_schedules ` + whereClause + ` ORDER BY scheduled_depart DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list flight schedules: %w", err)
	}
	defer rows.Close()

	var schedules []models.FlightSchedule
	for rows.Next() {
		var schedule models.FlightSchedule
		err := rows.Scan(
			&schedule.ID, &schedule.FlightNo, &schedule.DepartureAirport, &schedule.ArrivalAirport,
			&schedule.ScheduledDepart, &schedule.ScheduledArrive,
			&schedule.ActualDepart, &schedule.ActualArrive, &schedule.FlightStatus,
			&schedule.VipLoungeCapacity, &schedule.VipLoungeUsed,
			&schedule.Status, &schedule.ResponsiblePerson, &schedule.BatchNo,
			&schedule.Remarks, &schedule.SnapshotData, &schedule.CreatedAt, &schedule.UpdatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan flight schedule: %w", err)
		}
		schedules = append(schedules, schedule)
	}

	return schedules, total, nil
}

func (r *FlightScheduleRepository) Update(id int64, schedule *models.FlightSchedule) error {
	query := `UPDATE flight_schedules SET 
		departure_airport = ?, arrival_airport = ?, scheduled_depart = ?, scheduled_arrive = ?,
		actual_depart = ?, actual_arrive = ?, flight_status = ?, vip_lounge_capacity = ?,
		vip_lounge_used = ?, status = ?, responsible_person = ?, batch_no = ?, remarks = ?,
		snapshot_data = ?, updated_at = ?
		WHERE id = ?`

	_, err := r.db.Exec(query,
		schedule.DepartureAirport, schedule.ArrivalAirport, schedule.ScheduledDepart, schedule.ScheduledArrive,
		schedule.ActualDepart, schedule.ActualArrive, schedule.FlightStatus, schedule.VipLoungeCapacity,
		schedule.VipLoungeUsed, schedule.Status, schedule.ResponsiblePerson, schedule.BatchNo,
		schedule.Remarks, schedule.SnapshotData, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update flight schedule: %w", err)
	}

	return nil
}

func (r *FlightScheduleRepository) UpdateStatus(id int64, status string) error {
	query := `UPDATE flight_schedules SET status = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update flight schedule status: %w", err)
	}
	return nil
}

func (r *FlightScheduleRepository) Delete(id int64) error {
	query := `DELETE FROM flight_schedules WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete flight schedule: %w", err)
	}
	return nil
}

func (r *FlightScheduleRepository) Archive(id int64, snapshotData string) error {
	query := `UPDATE flight_schedules SET status = 'archived', snapshot_data = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, snapshotData, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to archive flight schedule: %w", err)
	}
	return nil
}

func (r *FlightScheduleRepository) BatchArchive(ids []int64, createSnapshot bool) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	var count int
	for _, id := range ids {
		var query string
		if createSnapshot {
			query = `UPDATE flight_schedules SET status = 'archived', snapshot_data = json_extract(json_object('flight_no', flight_no, 'departure', departure_airport, 'arrival', arrival_airport, 'scheduled_depart', scheduled_depart, 'flight_status', flight_status), '$'), updated_at = ? WHERE id = ?`
		} else {
			query = `UPDATE flight_schedules SET status = 'archived', updated_at = ? WHERE id = ?`
		}
		result, err := tx.Exec(query, time.Now(), id)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to archive flight schedule: %w", err)
		}
		rowsAffected, _ := result.RowsAffected()
		count += int(rowsAffected)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return count, nil
}

func (r *FlightScheduleRepository) GetByFlightNo(flightNo string) ([]models.FlightSchedule, error) {
	query := `SELECT id, flight_no, departure_airport, arrival_airport, scheduled_depart, scheduled_arrive,
		actual_depart, actual_arrive, flight_status, vip_lounge_capacity, vip_lounge_used,
		status, responsible_person, batch_no, remarks, snapshot_data, created_at, updated_at 
		FROM flight_schedules WHERE flight_no = ? ORDER BY scheduled_depart`

	rows, err := r.db.Query(query, flightNo)
	if err != nil {
		return nil, fmt.Errorf("failed to get flight schedules by flight no: %w", err)
	}
	defer rows.Close()

	var schedules []models.FlightSchedule
	for rows.Next() {
		var schedule models.FlightSchedule
		err := rows.Scan(
			&schedule.ID, &schedule.FlightNo, &schedule.DepartureAirport, &schedule.ArrivalAirport,
			&schedule.ScheduledDepart, &schedule.ScheduledArrive,
			&schedule.ActualDepart, &schedule.ActualArrive, &schedule.FlightStatus,
			&schedule.VipLoungeCapacity, &schedule.VipLoungeUsed,
			&schedule.Status, &schedule.ResponsiblePerson, &schedule.BatchNo,
			&schedule.Remarks, &schedule.SnapshotData, &schedule.CreatedAt, &schedule.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan flight schedule: %w", err)
		}
		schedules = append(schedules, schedule)
	}

	return schedules, nil
}

func (r *FlightScheduleRepository) UpdateLoungeUsage(id int64, used int) error {
	query := `UPDATE flight_schedules SET vip_lounge_used = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, used, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update lounge usage: %w", err)
	}
	return nil
}
