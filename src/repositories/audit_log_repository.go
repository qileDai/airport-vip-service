package repositories

import (
	"airport-vip-service/src/models"
	"database/sql"
	"fmt"
	"time"
)

type AuditLogRepository struct {
	db *sql.DB
}

func NewAuditLogRepository(db *sql.DB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

func (r *AuditLogRepository) Create(log *models.AuditLog) (int64, error) {
	query := `INSERT INTO audit_logs 
		(log_no, operation_type, entity_type, entity_id, entity_no, old_value, new_value,
		operator, operator_role, ip_address, user_agent, request_id, status, responsible_person,
		batch_no, remarks, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := r.db.Exec(query,
		log.LogNo, log.OperationType, log.EntityType, log.EntityID, log.EntityNo,
		log.OldValue, log.NewValue, log.Operator, log.OperatorRole,
		log.IPAddress, log.UserAgent, log.RequestID, log.Status,
		log.ResponsiblePerson, log.BatchNo, log.Remarks, now)
	if err != nil {
		return 0, fmt.Errorf("failed to create audit log: %w", err)
	}

	return result.LastInsertId()
}

func (r *AuditLogRepository) GetByID(id int64) (*models.AuditLog, error) {
	query := `SELECT id, log_no, operation_type, entity_type, entity_id, entity_no, old_value, new_value,
		operator, operator_role, ip_address, user_agent, request_id, status, responsible_person,
		batch_no, remarks, created_at 
		FROM audit_logs WHERE id = ?`

	log := &models.AuditLog{}
	err := r.db.QueryRow(query, id).Scan(
		&log.ID, &log.LogNo, &log.OperationType, &log.EntityType, &log.EntityID, &log.EntityNo,
		&log.OldValue, &log.NewValue, &log.Operator, &log.OperatorRole,
		&log.IPAddress, &log.UserAgent, &log.RequestID, &log.Status,
		&log.ResponsiblePerson, &log.BatchNo, &log.Remarks, &log.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get audit log: %w", err)
	}

	return log, nil
}

func (r *AuditLogRepository) List(offset, limit int, entityType, operationType, operator string, startTime, endTime *time.Time) ([]models.AuditLog, int64, error) {
	var args []interface{}
	whereClause := "WHERE 1=1"

	if entityType != "" {
		whereClause += " AND entity_type = ?"
		args = append(args, entityType)
	}
	if operationType != "" {
		whereClause += " AND operation_type = ?"
		args = append(args, operationType)
	}
	if operator != "" {
		whereClause += " AND operator = ?"
		args = append(args, operator)
	}
	if startTime != nil {
		whereClause += " AND created_at >= ?"
		args = append(args, startTime)
	}
	if endTime != nil {
		whereClause += " AND created_at <= ?"
		args = append(args, endTime)
	}

	countQuery := "SELECT COUNT(*) FROM audit_logs " + whereClause
	var total int64
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count audit logs: %w", err)
	}

	query := `SELECT id, log_no, operation_type, entity_type, entity_id, entity_no, old_value, new_value,
		operator, operator_role, ip_address, user_agent, request_id, status, responsible_person,
		batch_no, remarks, created_at 
		FROM audit_logs ` + whereClause + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list audit logs: %w", err)
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var log models.AuditLog
		err := rows.Scan(
			&log.ID, &log.LogNo, &log.OperationType, &log.EntityType, &log.EntityID, &log.EntityNo,
			&log.OldValue, &log.NewValue, &log.Operator, &log.OperatorRole,
			&log.IPAddress, &log.UserAgent, &log.RequestID, &log.Status,
			&log.ResponsiblePerson, &log.BatchNo, &log.Remarks, &log.CreatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan audit log: %w", err)
		}
		logs = append(logs, log)
	}

	return logs, total, nil
}

func (r *AuditLogRepository) GetByEntity(entityType string, entityID int64) ([]models.AuditLog, error) {
	query := `SELECT id, log_no, operation_type, entity_type, entity_id, entity_no, old_value, new_value,
		operator, operator_role, ip_address, user_agent, request_id, status, responsible_person,
		batch_no, remarks, created_at 
		FROM audit_logs WHERE entity_type = ? AND entity_id = ? ORDER BY created_at DESC`

	rows, err := r.db.Query(query, entityType, entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs by entity: %w", err)
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var log models.AuditLog
		err := rows.Scan(
			&log.ID, &log.LogNo, &log.OperationType, &log.EntityType, &log.EntityID, &log.EntityNo,
			&log.OldValue, &log.NewValue, &log.Operator, &log.OperatorRole,
			&log.IPAddress, &log.UserAgent, &log.RequestID, &log.Status,
			&log.ResponsiblePerson, &log.BatchNo, &log.Remarks, &log.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan audit log: %w", err)
		}
		logs = append(logs, log)
	}

	return logs, nil
}

func (r *AuditLogRepository) GetByDateRange(startDate, endDate time.Time) ([]models.AuditLog, error) {
	query := `SELECT id, log_no, operation_type, entity_type, entity_id, entity_no, old_value, new_value,
		operator, operator_role, ip_address, user_agent, request_id, status, responsible_person,
		batch_no, remarks, created_at 
		FROM audit_logs WHERE created_at BETWEEN ? AND ? ORDER BY created_at DESC`

	rows, err := r.db.Query(query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs by date range: %w", err)
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var log models.AuditLog
		err := rows.Scan(
			&log.ID, &log.LogNo, &log.OperationType, &log.EntityType, &log.EntityID, &log.EntityNo,
			&log.OldValue, &log.NewValue, &log.Operator, &log.OperatorRole,
			&log.IPAddress, &log.UserAgent, &log.RequestID, &log.Status,
			&log.ResponsiblePerson, &log.BatchNo, &log.Remarks, &log.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan audit log: %w", err)
		}
		logs = append(logs, log)
	}

	return logs, nil
}
