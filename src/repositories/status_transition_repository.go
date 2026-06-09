package repositories

import (
	"airport-vip-service/src/models"
	"database/sql"
	"fmt"
	"time"
)

type StatusTransitionRepository struct {
	db *sql.DB
}

func NewStatusTransitionRepository(db *sql.DB) *StatusTransitionRepository {
	return &StatusTransitionRepository{db: db}
}

func (r *StatusTransitionRepository) Create(transition *models.StatusTransition) (int64, error) {
	query := `INSERT INTO status_transitions 
		(transition_no, entity_type, entity_id, from_status, to_status, action, reason, operator,
		status, responsible_person, batch_no, remarks, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := r.db.Exec(query,
		transition.TransitionNo, transition.EntityType, transition.EntityID,
		transition.FromStatus, transition.ToStatus, transition.Action, transition.Reason,
		transition.Operator, transition.Status, transition.ResponsiblePerson,
		transition.BatchNo, transition.Remarks, now, now)
	if err != nil {
		return 0, fmt.Errorf("failed to create status transition: %w", err)
	}

	return result.LastInsertId()
}

func (r *StatusTransitionRepository) GetByID(id int64) (*models.StatusTransition, error) {
	query := `SELECT id, transition_no, entity_type, entity_id, from_status, to_status, action, reason, operator,
		status, responsible_person, batch_no, remarks, created_at, updated_at 
		FROM status_transitions WHERE id = ?`

	transition := &models.StatusTransition{}
	err := r.db.QueryRow(query, id).Scan(
		&transition.ID, &transition.TransitionNo, &transition.EntityType, &transition.EntityID,
		&transition.FromStatus, &transition.ToStatus, &transition.Action, &transition.Reason,
		&transition.Operator, &transition.Status, &transition.ResponsiblePerson,
		&transition.BatchNo, &transition.Remarks, &transition.CreatedAt, &transition.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get status transition: %w", err)
	}

	return transition, nil
}

func (r *StatusTransitionRepository) List(offset, limit int, entityType string, entityID int64) ([]models.StatusTransition, int64, error) {
	var args []interface{}
	whereClause := "WHERE 1=1"

	if entityType != "" {
		whereClause += " AND entity_type = ?"
		args = append(args, entityType)
	}
	if entityID > 0 {
		whereClause += " AND entity_id = ?"
		args = append(args, entityID)
	}

	countQuery := "SELECT COUNT(*) FROM status_transitions " + whereClause
	var total int64
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count status transitions: %w", err)
	}

	query := `SELECT id, transition_no, entity_type, entity_id, from_status, to_status, action, reason, operator,
		status, responsible_person, batch_no, remarks, created_at, updated_at 
		FROM status_transitions ` + whereClause + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list status transitions: %w", err)
	}
	defer rows.Close()

	var transitions []models.StatusTransition
	for rows.Next() {
		var transition models.StatusTransition
		err := rows.Scan(
			&transition.ID, &transition.TransitionNo, &transition.EntityType, &transition.EntityID,
			&transition.FromStatus, &transition.ToStatus, &transition.Action, &transition.Reason,
			&transition.Operator, &transition.Status, &transition.ResponsiblePerson,
			&transition.BatchNo, &transition.Remarks, &transition.CreatedAt, &transition.UpdatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan status transition: %w", err)
		}
		transitions = append(transitions, transition)
	}

	return transitions, total, nil
}

func (r *StatusTransitionRepository) GetByEntity(entityType string, entityID int64) ([]models.StatusTransition, error) {
	query := `SELECT id, transition_no, entity_type, entity_id, from_status, to_status, action, reason, operator,
		status, responsible_person, batch_no, remarks, created_at, updated_at 
		FROM status_transitions WHERE entity_type = ? AND entity_id = ? ORDER BY created_at DESC`

	rows, err := r.db.Query(query, entityType, entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get status transitions by entity: %w", err)
	}
	defer rows.Close()

	var transitions []models.StatusTransition
	for rows.Next() {
		var transition models.StatusTransition
		err := rows.Scan(
			&transition.ID, &transition.TransitionNo, &transition.EntityType, &transition.EntityID,
			&transition.FromStatus, &transition.ToStatus, &transition.Action, &transition.Reason,
			&transition.Operator, &transition.Status, &transition.ResponsiblePerson,
			&transition.BatchNo, &transition.Remarks, &transition.CreatedAt, &transition.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan status transition: %w", err)
		}
		transitions = append(transitions, transition)
	}

	return transitions, nil
}

func (r *StatusTransitionRepository) GetLatestByEntity(entityType string, entityID int64) (*models.StatusTransition, error) {
	query := `SELECT id, transition_no, entity_type, entity_id, from_status, to_status, action, reason, operator,
		status, responsible_person, batch_no, remarks, created_at, updated_at 
		FROM status_transitions WHERE entity_type = ? AND entity_id = ? ORDER BY created_at DESC LIMIT 1`

	transition := &models.StatusTransition{}
	err := r.db.QueryRow(query, entityType, entityID).Scan(
		&transition.ID, &transition.TransitionNo, &transition.EntityType, &transition.EntityID,
		&transition.FromStatus, &transition.ToStatus, &transition.Action, &transition.Reason,
		&transition.Operator, &transition.Status, &transition.ResponsiblePerson,
		&transition.BatchNo, &transition.Remarks, &transition.CreatedAt, &transition.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get latest status transition: %w", err)
	}

	return transition, nil
}
