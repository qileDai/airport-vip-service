package repositories

import (
	"airport-vip-service/src/models"
	"database/sql"
	"fmt"
	"time"
)

type RuleConfigRepository struct {
	db *sql.DB
}

func NewRuleConfigRepository(db *sql.DB) *RuleConfigRepository {
	return &RuleConfigRepository{db: db}
}

func (r *RuleConfigRepository) Create(rule *models.RuleConfig) (int64, error) {
	query := `INSERT INTO rule_configs 
		(rule_no, rule_name, rule_type, rule_value, threshold_value, applies_to_level, is_active,
		effective_date, expiry_date, status, responsible_person, batch_no, remarks, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	isActive := 0
	if rule.IsActive {
		isActive = 1
	}
	result, err := r.db.Exec(query,
		rule.RuleNo, rule.RuleName, rule.RuleType, rule.RuleValue, rule.ThresholdValue,
		rule.AppliesToLevel, isActive, rule.EffectiveDate, rule.ExpiryDate,
		rule.Status, rule.ResponsiblePerson, rule.BatchNo, rule.Remarks, now, now)
	if err != nil {
		return 0, fmt.Errorf("failed to create rule config: %w", err)
	}

	return result.LastInsertId()
}

func (r *RuleConfigRepository) GetByID(id int64) (*models.RuleConfig, error) {
	query := `SELECT id, rule_no, rule_name, rule_type, rule_value, threshold_value, applies_to_level, is_active,
		effective_date, expiry_date, status, responsible_person, batch_no, remarks, created_at, updated_at 
		FROM rule_configs WHERE id = ?`

	rule := &models.RuleConfig{}
	var isActive int
	var expiryDate sql.NullTime
	err := r.db.QueryRow(query, id).Scan(
		&rule.ID, &rule.RuleNo, &rule.RuleName, &rule.RuleType, &rule.RuleValue,
		&rule.ThresholdValue, &rule.AppliesToLevel, &isActive,
		&rule.EffectiveDate, &expiryDate, &rule.Status, &rule.ResponsiblePerson,
		&rule.BatchNo, &rule.Remarks, &rule.CreatedAt, &rule.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get rule config: %w", err)
	}
	rule.IsActive = isActive == 1
	if expiryDate.Valid {
		rule.ExpiryDate = &expiryDate.Time
	}

	return rule, nil
}

func (r *RuleConfigRepository) GetByRuleNo(ruleNo string) (*models.RuleConfig, error) {
	query := `SELECT id, rule_no, rule_name, rule_type, rule_value, threshold_value, applies_to_level, is_active,
		effective_date, expiry_date, status, responsible_person, batch_no, remarks, created_at, updated_at 
		FROM rule_configs WHERE rule_no = ?`

	rule := &models.RuleConfig{}
	var isActive int
	var expiryDate sql.NullTime
	err := r.db.QueryRow(query, ruleNo).Scan(
		&rule.ID, &rule.RuleNo, &rule.RuleName, &rule.RuleType, &rule.RuleValue,
		&rule.ThresholdValue, &rule.AppliesToLevel, &isActive,
		&rule.EffectiveDate, &expiryDate, &rule.Status, &rule.ResponsiblePerson,
		&rule.BatchNo, &rule.Remarks, &rule.CreatedAt, &rule.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get rule config: %w", err)
	}
	rule.IsActive = isActive == 1
	if expiryDate.Valid {
		rule.ExpiryDate = &expiryDate.Time
	}

	return rule, nil
}

func (r *RuleConfigRepository) List(offset, limit int, ruleType, status string) ([]models.RuleConfig, int64, error) {
	var args []interface{}
	whereClause := "WHERE 1=1"

	if ruleType != "" {
		whereClause += " AND rule_type = ?"
		args = append(args, ruleType)
	}
	if status != "" {
		whereClause += " AND status = ?"
		args = append(args, status)
	}

	countQuery := "SELECT COUNT(*) FROM rule_configs " + whereClause
	var total int64
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count rule configs: %w", err)
	}

	query := `SELECT id, rule_no, rule_name, rule_type, rule_value, threshold_value, applies_to_level, is_active,
		effective_date, expiry_date, status, responsible_person, batch_no, remarks, created_at, updated_at 
		FROM rule_configs ` + whereClause + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list rule configs: %w", err)
	}
	defer rows.Close()

	var rules []models.RuleConfig
	for rows.Next() {
		var rule models.RuleConfig
		var isActive int
		var expiryDate sql.NullTime
		err := rows.Scan(
			&rule.ID, &rule.RuleNo, &rule.RuleName, &rule.RuleType, &rule.RuleValue,
			&rule.ThresholdValue, &rule.AppliesToLevel, &isActive,
			&rule.EffectiveDate, &expiryDate, &rule.Status, &rule.ResponsiblePerson,
			&rule.BatchNo, &rule.Remarks, &rule.CreatedAt, &rule.UpdatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan rule config: %w", err)
		}
		rule.IsActive = isActive == 1
		if expiryDate.Valid {
			rule.ExpiryDate = &expiryDate.Time
		}
		rules = append(rules, rule)
	}

	return rules, total, nil
}

func (r *RuleConfigRepository) Update(id int64, rule *models.RuleConfig) error {
	query := `UPDATE rule_configs SET 
		rule_name = ?, rule_value = ?, threshold_value = ?, applies_to_level = ?, is_active = ?,
		effective_date = ?, expiry_date = ?, status = ?, responsible_person = ?, batch_no = ?, remarks = ?, updated_at = ?
		WHERE id = ?`

	isActive := 0
	if rule.IsActive {
		isActive = 1
	}

	_, err := r.db.Exec(query,
		rule.RuleName, rule.RuleValue, rule.ThresholdValue, rule.AppliesToLevel, isActive,
		rule.EffectiveDate, rule.ExpiryDate, rule.Status, rule.ResponsiblePerson,
		rule.BatchNo, rule.Remarks, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update rule config: %w", err)
	}

	return nil
}

func (r *RuleConfigRepository) Delete(id int64) error {
	query := `DELETE FROM rule_configs WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete rule config: %w", err)
	}
	return nil
}

func (r *RuleConfigRepository) GetActiveByType(ruleType string) ([]models.RuleConfig, error) {
	query := `SELECT id, rule_no, rule_name, rule_type, rule_value, threshold_value, applies_to_level, is_active,
		effective_date, expiry_date, status, responsible_person, batch_no, remarks, created_at, updated_at 
		FROM rule_configs WHERE rule_type = ? AND is_active = 1 AND status = 'active' 
		AND effective_date <= datetime('now') AND (expiry_date IS NULL OR expiry_date > datetime('now'))`

	rows, err := r.db.Query(query, ruleType)
	if err != nil {
		return nil, fmt.Errorf("failed to get active rules by type: %w", err)
	}
	defer rows.Close()

	var rules []models.RuleConfig
	for rows.Next() {
		var rule models.RuleConfig
		var isActive int
		var expiryDate sql.NullTime
		err := rows.Scan(
			&rule.ID, &rule.RuleNo, &rule.RuleName, &rule.RuleType, &rule.RuleValue,
			&rule.ThresholdValue, &rule.AppliesToLevel, &isActive,
			&rule.EffectiveDate, &expiryDate, &rule.Status, &rule.ResponsiblePerson,
			&rule.BatchNo, &rule.Remarks, &rule.CreatedAt, &rule.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rule config: %w", err)
		}
		rule.IsActive = isActive == 1
		if expiryDate.Valid {
			rule.ExpiryDate = &expiryDate.Time
		}
		rules = append(rules, rule)
	}

	return rules, nil
}
