package repositories

import (
	"airport-vip-service/src/models"
	"database/sql"
	"fmt"
	"time"
)

type VerificationResultRepository struct {
	db *sql.DB
}

func NewVerificationResultRepository(db *sql.DB) *VerificationResultRepository {
	return &VerificationResultRepository{db: db}
}

func (r *VerificationResultRepository) Create(result *models.VerificationResult) (int64, error) {
	query := `INSERT INTO verification_results 
		(verification_no, reservation_id, member_benefit_id, flight_schedule_id, verification_type,
		result, failure_reason, verified_quota, verified_companions, verification_details, status,
		responsible_person, batch_no, remarks, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	res, err := r.db.Exec(query,
		result.VerificationNo, result.ReservationID, result.MemberBenefitID, result.FlightScheduleID,
		result.VerificationType, result.Result, result.FailureReason, result.VerifiedQuota,
		result.VerifiedCompanions, result.VerificationDetails, result.Status,
		result.ResponsiblePerson, result.BatchNo, result.Remarks, now, now)
	if err != nil {
		return 0, fmt.Errorf("failed to create verification result: %w", err)
	}

	return res.LastInsertId()
}

func (r *VerificationResultRepository) GetByID(id int64) (*models.VerificationResult, error) {
	query := `SELECT id, verification_no, reservation_id, member_benefit_id, flight_schedule_id, verification_type,
		result, failure_reason, verified_quota, verified_companions, verification_details, status,
		responsible_person, batch_no, remarks, created_at, updated_at 
		FROM verification_results WHERE id = ?`

	result := &models.VerificationResult{}
	var flightScheduleID sql.NullInt64
	err := r.db.QueryRow(query, id).Scan(
		&result.ID, &result.VerificationNo, &result.ReservationID, &result.MemberBenefitID,
		&flightScheduleID, &result.VerificationType, &result.Result, &result.FailureReason,
		&result.VerifiedQuota, &result.VerifiedCompanions, &result.VerificationDetails,
		&result.Status, &result.ResponsiblePerson, &result.BatchNo, &result.Remarks,
		&result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get verification result: %w", err)
	}
	if flightScheduleID.Valid {
		result.FlightScheduleID = flightScheduleID.Int64
	}

	return result, nil
}

func (r *VerificationResultRepository) List(offset, limit int, status, resultFilter string) ([]models.VerificationResult, int64, error) {
	var args []interface{}
	whereClause := "WHERE 1=1"

	if status != "" {
		whereClause += " AND status = ?"
		args = append(args, status)
	}
	if resultFilter != "" {
		whereClause += " AND result = ?"
		args = append(args, resultFilter)
	}

	countQuery := "SELECT COUNT(*) FROM verification_results " + whereClause
	var total int64
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count verification results: %w", err)
	}

	query := `SELECT id, verification_no, reservation_id, member_benefit_id, flight_schedule_id, verification_type,
		result, failure_reason, verified_quota, verified_companions, verification_details, status,
		responsible_person, batch_no, remarks, created_at, updated_at 
		FROM verification_results ` + whereClause + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list verification results: %w", err)
	}
	defer rows.Close()

	var results []models.VerificationResult
	for rows.Next() {
		var result models.VerificationResult
		var flightScheduleID sql.NullInt64
		err := rows.Scan(
			&result.ID, &result.VerificationNo, &result.ReservationID, &result.MemberBenefitID,
			&flightScheduleID, &result.VerificationType, &result.Result, &result.FailureReason,
			&result.VerifiedQuota, &result.VerifiedCompanions, &result.VerificationDetails,
			&result.Status, &result.ResponsiblePerson, &result.BatchNo, &result.Remarks,
			&result.CreatedAt, &result.UpdatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan verification result: %w", err)
		}
		if flightScheduleID.Valid {
			result.FlightScheduleID = flightScheduleID.Int64
		}
		results = append(results, result)
	}

	return results, total, nil
}

func (r *VerificationResultRepository) Update(id int64, result *models.VerificationResult) error {
	query := `UPDATE verification_results SET 
		result = ?, failure_reason = ?, verified_quota = ?, verified_companions = ?,
		verification_details = ?, status = ?, responsible_person = ?, batch_no = ?, remarks = ?, updated_at = ?
		WHERE id = ?`

	_, err := r.db.Exec(query,
		result.Result, result.FailureReason, result.VerifiedQuota, result.VerifiedCompanions,
		result.VerificationDetails, result.Status, result.ResponsiblePerson,
		result.BatchNo, result.Remarks, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update verification result: %w", err)
	}

	return nil
}

func (r *VerificationResultRepository) Delete(id int64) error {
	query := `DELETE FROM verification_results WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete verification result: %w", err)
	}
	return nil
}

func (r *VerificationResultRepository) GetStatistics(startDate, endDate time.Time) (map[string]int64, error) {
	query := `SELECT 
		COUNT(*) as total,
		SUM(CASE WHEN result = 'passed' THEN 1 ELSE 0 END) as passed,
		SUM(CASE WHEN result = 'failed' THEN 1 ELSE 0 END) as failed,
		SUM(CASE WHEN result = 'pending' THEN 1 ELSE 0 END) as pending
		FROM verification_results WHERE created_at BETWEEN ? AND ?`

	stats := make(map[string]int64)
	var total, passed, failed, pending int64
	err := r.db.QueryRow(query, startDate, endDate).Scan(&total, &passed, &failed, &pending)
	if err != nil {
		return nil, fmt.Errorf("failed to get verification statistics: %w", err)
	}

	stats["total"] = total
	stats["passed"] = passed
	stats["failed"] = failed
	stats["pending"] = pending

	return stats, nil
}

func (r *VerificationResultRepository) GetByDateRange(startDate, endDate time.Time) ([]models.VerificationResult, error) {
	query := `SELECT id, verification_no, reservation_id, member_benefit_id, flight_schedule_id, verification_type,
		result, failure_reason, verified_quota, verified_companions, verification_details, status,
		responsible_person, batch_no, remarks, created_at, updated_at 
		FROM verification_results WHERE created_at BETWEEN ? AND ? ORDER BY created_at`

	rows, err := r.db.Query(query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get verification results by date range: %w", err)
	}
	defer rows.Close()

	var results []models.VerificationResult
	for rows.Next() {
		var result models.VerificationResult
		var flightScheduleID sql.NullInt64
		err := rows.Scan(
			&result.ID, &result.VerificationNo, &result.ReservationID, &result.MemberBenefitID,
			&flightScheduleID, &result.VerificationType, &result.Result, &result.FailureReason,
			&result.VerifiedQuota, &result.VerifiedCompanions, &result.VerificationDetails,
			&result.Status, &result.ResponsiblePerson, &result.BatchNo, &result.Remarks,
			&result.CreatedAt, &result.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan verification result: %w", err)
		}
		if flightScheduleID.Valid {
			result.FlightScheduleID = flightScheduleID.Int64
		}
		results = append(results, result)
	}

	return results, nil
}

func (r *VerificationResultRepository) GetByBatchNo(batchNo string) ([]models.VerificationResult, error) {
	query := `SELECT id, verification_no, reservation_id, member_benefit_id, flight_schedule_id, verification_type,
		result, failure_reason, verified_quota, verified_companions, verification_details, status,
		responsible_person, batch_no, remarks, created_at, updated_at 
		FROM verification_results WHERE batch_no = ? ORDER BY created_at`

	rows, err := r.db.Query(query, batchNo)
	if err != nil {
		return nil, fmt.Errorf("failed to get verification results by batch: %w", err)
	}
	defer rows.Close()

	var results []models.VerificationResult
	for rows.Next() {
		var result models.VerificationResult
		var flightScheduleID sql.NullInt64
		err := rows.Scan(
			&result.ID, &result.VerificationNo, &result.ReservationID, &result.MemberBenefitID,
			&flightScheduleID, &result.VerificationType, &result.Result, &result.FailureReason,
			&result.VerifiedQuota, &result.VerifiedCompanions, &result.VerificationDetails,
			&result.Status, &result.ResponsiblePerson, &result.BatchNo, &result.Remarks,
			&result.CreatedAt, &result.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan verification result: %w", err)
		}
		if flightScheduleID.Valid {
			result.FlightScheduleID = flightScheduleID.Int64
		}
		results = append(results, result)
	}

	return results, nil
}
