package repositories

import (
	"airport-vip-service/src/models"
	"database/sql"
	"fmt"
	"time"
)

type MemberBenefitRepository struct {
	db *sql.DB
}

func NewMemberBenefitRepository(db *sql.DB) *MemberBenefitRepository {
	return &MemberBenefitRepository{db: db}
}

func (r *MemberBenefitRepository) Create(benefit *models.MemberBenefit) (int64, error) {
	query := `INSERT INTO member_benefits 
		(benefit_no, member_name, member_level, remaining_quota, total_quota, status, responsible_person, valid_from, valid_to, batch_no, remarks, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := r.db.Exec(query,
		benefit.BenefitNo, benefit.MemberName, benefit.MemberLevel,
		benefit.RemainingQuota, benefit.TotalQuota, benefit.Status,
		benefit.ResponsiblePerson, benefit.ValidFrom, benefit.ValidTo,
		benefit.BatchNo, benefit.Remarks, now, now)
	if err != nil {
		return 0, fmt.Errorf("failed to create member benefit: %w", err)
	}

	return result.LastInsertId()
}

func (r *MemberBenefitRepository) GetByID(id int64) (*models.MemberBenefit, error) {
	query := `SELECT id, benefit_no, member_name, member_level, remaining_quota, total_quota, status, 
		responsible_person, valid_from, valid_to, batch_no, remarks, created_at, updated_at 
		FROM member_benefits WHERE id = ?`

	benefit := &models.MemberBenefit{}
	err := r.db.QueryRow(query, id).Scan(
		&benefit.ID, &benefit.BenefitNo, &benefit.MemberName, &benefit.MemberLevel,
		&benefit.RemainingQuota, &benefit.TotalQuota, &benefit.Status,
		&benefit.ResponsiblePerson, &benefit.ValidFrom, &benefit.ValidTo,
		&benefit.BatchNo, &benefit.Remarks, &benefit.CreatedAt, &benefit.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get member benefit: %w", err)
	}

	return benefit, nil
}

func (r *MemberBenefitRepository) GetByBenefitNo(benefitNo string) (*models.MemberBenefit, error) {
	query := `SELECT id, benefit_no, member_name, member_level, remaining_quota, total_quota, status, 
		responsible_person, valid_from, valid_to, batch_no, remarks, created_at, updated_at 
		FROM member_benefits WHERE benefit_no = ?`

	benefit := &models.MemberBenefit{}
	err := r.db.QueryRow(query, benefitNo).Scan(
		&benefit.ID, &benefit.BenefitNo, &benefit.MemberName, &benefit.MemberLevel,
		&benefit.RemainingQuota, &benefit.TotalQuota, &benefit.Status,
		&benefit.ResponsiblePerson, &benefit.ValidFrom, &benefit.ValidTo,
		&benefit.BatchNo, &benefit.Remarks, &benefit.CreatedAt, &benefit.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get member benefit: %w", err)
	}

	return benefit, nil
}

func (r *MemberBenefitRepository) List(offset, limit int, status, memberLevel string) ([]models.MemberBenefit, int64, error) {
	var args []interface{}
	whereClause := "WHERE 1=1"

	if status != "" {
		whereClause += " AND status = ?"
		args = append(args, status)
	}
	if memberLevel != "" {
		whereClause += " AND member_level = ?"
		args = append(args, memberLevel)
	}

	countQuery := "SELECT COUNT(*) FROM member_benefits " + whereClause
	var total int64
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count member benefits: %w", err)
	}

	query := `SELECT id, benefit_no, member_name, member_level, remaining_quota, total_quota, status, 
		responsible_person, valid_from, valid_to, batch_no, remarks, created_at, updated_at 
		FROM member_benefits ` + whereClause + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list member benefits: %w", err)
	}
	defer rows.Close()

	var benefits []models.MemberBenefit
	for rows.Next() {
		var benefit models.MemberBenefit
		err := rows.Scan(
			&benefit.ID, &benefit.BenefitNo, &benefit.MemberName, &benefit.MemberLevel,
			&benefit.RemainingQuota, &benefit.TotalQuota, &benefit.Status,
			&benefit.ResponsiblePerson, &benefit.ValidFrom, &benefit.ValidTo,
			&benefit.BatchNo, &benefit.Remarks, &benefit.CreatedAt, &benefit.UpdatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan member benefit: %w", err)
		}
		benefits = append(benefits, benefit)
	}

	return benefits, total, nil
}

func (r *MemberBenefitRepository) Update(id int64, benefit *models.MemberBenefit) error {
	query := `UPDATE member_benefits SET 
		member_name = ?, member_level = ?, remaining_quota = ?, total_quota = ?,
		status = ?, responsible_person = ?, valid_from = ?, valid_to = ?,
		batch_no = ?, remarks = ?, updated_at = ?
		WHERE id = ?`

	_, err := r.db.Exec(query,
		benefit.MemberName, benefit.MemberLevel, benefit.RemainingQuota, benefit.TotalQuota,
		benefit.Status, benefit.ResponsiblePerson, benefit.ValidFrom, benefit.ValidTo,
		benefit.BatchNo, benefit.Remarks, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update member benefit: %w", err)
	}

	return nil
}

func (r *MemberBenefitRepository) UpdateStatus(id int64, status string) error {
	query := `UPDATE member_benefits SET status = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update member benefit status: %w", err)
	}
	return nil
}

func (r *MemberBenefitRepository) Delete(id int64) error {
	query := `DELETE FROM member_benefits WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete member benefit: %w", err)
	}
	return nil
}

func (r *MemberBenefitRepository) GetExpiringSoon(days int) ([]models.MemberBenefit, error) {
	query := `SELECT id, benefit_no, member_name, member_level, remaining_quota, total_quota, status, 
		responsible_person, valid_from, valid_to, batch_no, remarks, created_at, updated_at 
		FROM member_benefits WHERE status = 'active' AND valid_to BETWEEN datetime('now') AND datetime('now', '+' || ? || ' days')`

	rows, err := r.db.Query(query, days)
	if err != nil {
		return nil, fmt.Errorf("failed to get expiring benefits: %w", err)
	}
	defer rows.Close()

	var benefits []models.MemberBenefit
	for rows.Next() {
		var benefit models.MemberBenefit
		err := rows.Scan(
			&benefit.ID, &benefit.BenefitNo, &benefit.MemberName, &benefit.MemberLevel,
			&benefit.RemainingQuota, &benefit.TotalQuota, &benefit.Status,
			&benefit.ResponsiblePerson, &benefit.ValidFrom, &benefit.ValidTo,
			&benefit.BatchNo, &benefit.Remarks, &benefit.CreatedAt, &benefit.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan member benefit: %w", err)
		}
		benefits = append(benefits, benefit)
	}

	return benefits, nil
}

func (r *MemberBenefitRepository) GetExpired() ([]models.MemberBenefit, error) {
	query := `SELECT id, benefit_no, member_name, member_level, remaining_quota, total_quota, status, 
		responsible_person, valid_from, valid_to, batch_no, remarks, created_at, updated_at 
		FROM member_benefits WHERE status = 'active' AND valid_to < datetime('now')`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get expired benefits: %w", err)
	}
	defer rows.Close()

	var benefits []models.MemberBenefit
	for rows.Next() {
		var benefit models.MemberBenefit
		err := rows.Scan(
			&benefit.ID, &benefit.BenefitNo, &benefit.MemberName, &benefit.MemberLevel,
			&benefit.RemainingQuota, &benefit.TotalQuota, &benefit.Status,
			&benefit.ResponsiblePerson, &benefit.ValidFrom, &benefit.ValidTo,
			&benefit.BatchNo, &benefit.Remarks, &benefit.CreatedAt, &benefit.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan member benefit: %w", err)
		}
		benefits = append(benefits, benefit)
	}

	return benefits, nil
}

func (r *MemberBenefitRepository) UpdateQuota(id int64, delta int) error {
	query := `UPDATE member_benefits SET remaining_quota = remaining_quota + ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, delta, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update quota: %w", err)
	}
	return nil
}
