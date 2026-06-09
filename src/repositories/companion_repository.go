package repositories

import (
	"airport-vip-service/src/models"
	"database/sql"
	"fmt"
	"time"
)

type CompanionRepository struct {
	db *sql.DB
}

func NewCompanionRepository(db *sql.DB) *CompanionRepository {
	return &CompanionRepository{db: db}
}

func (r *CompanionRepository) Create(companion *models.Companion) (int64, error) {
	query := `INSERT INTO companions 
		(companion_no, reservation_id, companion_name, companion_id_card, relation_type,
		is_vip_eligible, verification_status, status, responsible_person, batch_no, remarks, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := r.db.Exec(query,
		companion.CompanionNo, companion.ReservationID, companion.CompanionName,
		companion.CompanionIDCard, companion.RelationType, companion.IsVipEligible,
		companion.VerificationStatus, companion.Status, companion.ResponsiblePerson,
		companion.BatchNo, companion.Remarks, now, now)
	if err != nil {
		return 0, fmt.Errorf("failed to create companion: %w", err)
	}

	return result.LastInsertId()
}

func (r *CompanionRepository) GetByID(id int64) (*models.Companion, error) {
	query := `SELECT id, companion_no, reservation_id, companion_name, companion_id_card, relation_type,
		is_vip_eligible, verification_status, status, responsible_person, batch_no, remarks, created_at, updated_at 
		FROM companions WHERE id = ?`

	companion := &models.Companion{}
	var isVipEligible int
	err := r.db.QueryRow(query, id).Scan(
		&companion.ID, &companion.CompanionNo, &companion.ReservationID, &companion.CompanionName,
		&companion.CompanionIDCard, &companion.RelationType, &isVipEligible,
		&companion.VerificationStatus, &companion.Status, &companion.ResponsiblePerson,
		&companion.BatchNo, &companion.Remarks, &companion.CreatedAt, &companion.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get companion: %w", err)
	}
	companion.IsVipEligible = isVipEligible == 1

	return companion, nil
}

func (r *CompanionRepository) GetByReservationID(reservationID int64) ([]models.Companion, error) {
	query := `SELECT id, companion_no, reservation_id, companion_name, companion_id_card, relation_type,
		is_vip_eligible, verification_status, status, responsible_person, batch_no, remarks, created_at, updated_at 
		FROM companions WHERE reservation_id = ? ORDER BY created_at`

	rows, err := r.db.Query(query, reservationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get companions by reservation: %w", err)
	}
	defer rows.Close()

	var companions []models.Companion
	for rows.Next() {
		var companion models.Companion
		var isVipEligible int
		err := rows.Scan(
			&companion.ID, &companion.CompanionNo, &companion.ReservationID, &companion.CompanionName,
			&companion.CompanionIDCard, &companion.RelationType, &isVipEligible,
			&companion.VerificationStatus, &companion.Status, &companion.ResponsiblePerson,
			&companion.BatchNo, &companion.Remarks, &companion.CreatedAt, &companion.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan companion: %w", err)
		}
		companion.IsVipEligible = isVipEligible == 1
		companions = append(companions, companion)
	}

	return companions, nil
}

func (r *CompanionRepository) List(offset, limit int, reservationID int64, status string) ([]models.Companion, int64, error) {
	var args []interface{}
	whereClause := "WHERE 1=1"

	if reservationID > 0 {
		whereClause += " AND reservation_id = ?"
		args = append(args, reservationID)
	}
	if status != "" {
		whereClause += " AND status = ?"
		args = append(args, status)
	}

	countQuery := "SELECT COUNT(*) FROM companions " + whereClause
	var total int64
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count companions: %w", err)
	}

	query := `SELECT id, companion_no, reservation_id, companion_name, companion_id_card, relation_type,
		is_vip_eligible, verification_status, status, responsible_person, batch_no, remarks, created_at, updated_at 
		FROM companions ` + whereClause + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list companions: %w", err)
	}
	defer rows.Close()

	var companions []models.Companion
	for rows.Next() {
		var companion models.Companion
		var isVipEligible int
		err := rows.Scan(
			&companion.ID, &companion.CompanionNo, &companion.ReservationID, &companion.CompanionName,
			&companion.CompanionIDCard, &companion.RelationType, &isVipEligible,
			&companion.VerificationStatus, &companion.Status, &companion.ResponsiblePerson,
			&companion.BatchNo, &companion.Remarks, &companion.CreatedAt, &companion.UpdatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan companion: %w", err)
		}
		companion.IsVipEligible = isVipEligible == 1
		companions = append(companions, companion)
	}

	return companions, total, nil
}

func (r *CompanionRepository) Update(id int64, companion *models.Companion) error {
	query := `UPDATE companions SET 
		companion_name = ?, companion_id_card = ?, relation_type = ?, is_vip_eligible = ?,
		verification_status = ?, status = ?, responsible_person = ?, batch_no = ?, remarks = ?, updated_at = ?
		WHERE id = ?`

	isVipEligible := 0
	if companion.IsVipEligible {
		isVipEligible = 1
	}

	_, err := r.db.Exec(query,
		companion.CompanionName, companion.CompanionIDCard, companion.RelationType, isVipEligible,
		companion.VerificationStatus, companion.Status, companion.ResponsiblePerson,
		companion.BatchNo, companion.Remarks, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update companion: %w", err)
	}

	return nil
}

func (r *CompanionRepository) Delete(id int64) error {
	query := `DELETE FROM companions WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete companion: %w", err)
	}
	return nil
}

func (r *CompanionRepository) CountByReservationID(reservationID int64) (int, error) {
	query := `SELECT COUNT(*) FROM companions WHERE reservation_id = ? AND status = 'active'`
	var count int
	err := r.db.QueryRow(query, reservationID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count companions: %w", err)
	}
	return count, nil
}
