package repositories

import (
	"airport-vip-service/src/models"
	"database/sql"
	"fmt"
	"time"
)

type ExceptionEventRepository struct {
	db *sql.DB
}

func NewExceptionEventRepository(db *sql.DB) *ExceptionEventRepository {
	return &ExceptionEventRepository{db: db}
}

func (r *ExceptionEventRepository) Create(event *models.ExceptionEvent) (int64, error) {
	query := `INSERT INTO exception_events 
		(event_no, event_type, entity_type, entity_id, trigger_field, threshold_value, actual_value,
		severity, handler, handling_deadline, handled_at, handling_result, status, responsible_person,
		batch_no, remarks, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := r.db.Exec(query,
		event.EventNo, event.EventType, event.EntityType, event.EntityID,
		event.TriggerField, event.ThresholdValue, event.ActualValue, event.Severity,
		event.Handler, event.HandlingDeadline, event.HandledAt, event.HandlingResult,
		event.Status, event.ResponsiblePerson, event.BatchNo, event.Remarks, now, now)
	if err != nil {
		return 0, fmt.Errorf("failed to create exception event: %w", err)
	}

	return result.LastInsertId()
}

func (r *ExceptionEventRepository) GetByID(id int64) (*models.ExceptionEvent, error) {
	query := `SELECT id, event_no, event_type, entity_type, entity_id, trigger_field, threshold_value, actual_value,
		severity, handler, handling_deadline, handled_at, handling_result, status, responsible_person,
		batch_no, remarks, created_at, updated_at 
		FROM exception_events WHERE id = ?`

	event := &models.ExceptionEvent{}
	err := r.db.QueryRow(query, id).Scan(
		&event.ID, &event.EventNo, &event.EventType, &event.EntityType, &event.EntityID,
		&event.TriggerField, &event.ThresholdValue, &event.ActualValue, &event.Severity,
		&event.Handler, &event.HandlingDeadline, &event.HandledAt, &event.HandlingResult,
		&event.Status, &event.ResponsiblePerson, &event.BatchNo, &event.Remarks,
		&event.CreatedAt, &event.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get exception event: %w", err)
	}

	return event, nil
}

func (r *ExceptionEventRepository) GetByEventNo(eventNo string) (*models.ExceptionEvent, error) {
	query := `SELECT id, event_no, event_type, entity_type, entity_id, trigger_field, threshold_value, actual_value,
		severity, handler, handling_deadline, handled_at, handling_result, status, responsible_person,
		batch_no, remarks, created_at, updated_at 
		FROM exception_events WHERE event_no = ?`

	event := &models.ExceptionEvent{}
	err := r.db.QueryRow(query, eventNo).Scan(
		&event.ID, &event.EventNo, &event.EventType, &event.EntityType, &event.EntityID,
		&event.TriggerField, &event.ThresholdValue, &event.ActualValue, &event.Severity,
		&event.Handler, &event.HandlingDeadline, &event.HandledAt, &event.HandlingResult,
		&event.Status, &event.ResponsiblePerson, &event.BatchNo, &event.Remarks,
		&event.CreatedAt, &event.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get exception event: %w", err)
	}

	return event, nil
}

func (r *ExceptionEventRepository) List(offset, limit int, eventType, status, severity string) ([]models.ExceptionEvent, int64, error) {
	var args []interface{}
	whereClause := "WHERE 1=1"

	if eventType != "" {
		whereClause += " AND event_type = ?"
		args = append(args, eventType)
	}
	if status != "" {
		whereClause += " AND status = ?"
		args = append(args, status)
	}
	if severity != "" {
		whereClause += " AND severity = ?"
		args = append(args, severity)
	}

	countQuery := "SELECT COUNT(*) FROM exception_events " + whereClause
	var total int64
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count exception events: %w", err)
	}

	query := `SELECT id, event_no, event_type, entity_type, entity_id, trigger_field, threshold_value, actual_value,
		severity, handler, handling_deadline, handled_at, handling_result, status, responsible_person,
		batch_no, remarks, created_at, updated_at 
		FROM exception_events ` + whereClause + ` ORDER BY severity DESC, created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list exception events: %w", err)
	}
	defer rows.Close()

	var events []models.ExceptionEvent
	for rows.Next() {
		var event models.ExceptionEvent
		err := rows.Scan(
			&event.ID, &event.EventNo, &event.EventType, &event.EntityType, &event.EntityID,
			&event.TriggerField, &event.ThresholdValue, &event.ActualValue, &event.Severity,
			&event.Handler, &event.HandlingDeadline, &event.HandledAt, &event.HandlingResult,
			&event.Status, &event.ResponsiblePerson, &event.BatchNo, &event.Remarks,
			&event.CreatedAt, &event.UpdatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan exception event: %w", err)
		}
		events = append(events, event)
	}

	return events, total, nil
}

func (r *ExceptionEventRepository) Update(id int64, event *models.ExceptionEvent) error {
	query := `UPDATE exception_events SET 
		handler = ?, handling_deadline = ?, handled_at = ?, handling_result = ?, status = ?,
		responsible_person = ?, batch_no = ?, remarks = ?, updated_at = ?
		WHERE id = ?`

	_, err := r.db.Exec(query,
		event.Handler, event.HandlingDeadline, event.HandledAt, event.HandlingResult,
		event.Status, event.ResponsiblePerson, event.BatchNo, event.Remarks, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update exception event: %w", err)
	}

	return nil
}

func (r *ExceptionEventRepository) UpdateStatus(id int64, status string) error {
	query := `UPDATE exception_events SET status = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update exception event status: %w", err)
	}
	return nil
}

func (r *ExceptionEventRepository) Delete(id int64) error {
	query := `DELETE FROM exception_events WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete exception event: %w", err)
	}
	return nil
}

func (r *ExceptionEventRepository) GetByEntityTypeAndID(entityType string, entityID int64) ([]models.ExceptionEvent, error) {
	query := `SELECT id, event_no, event_type, entity_type, entity_id, trigger_field, threshold_value, actual_value,
		severity, handler, handling_deadline, handled_at, handling_result, status, responsible_person,
		batch_no, remarks, created_at, updated_at 
		FROM exception_events WHERE entity_type = ? AND entity_id = ? ORDER BY created_at DESC`

	rows, err := r.db.Query(query, entityType, entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get exception events by entity: %w", err)
	}
	defer rows.Close()

	var events []models.ExceptionEvent
	for rows.Next() {
		var event models.ExceptionEvent
		err := rows.Scan(
			&event.ID, &event.EventNo, &event.EventType, &event.EntityType, &event.EntityID,
			&event.TriggerField, &event.ThresholdValue, &event.ActualValue, &event.Severity,
			&event.Handler, &event.HandlingDeadline, &event.HandledAt, &event.HandlingResult,
			&event.Status, &event.ResponsiblePerson, &event.BatchNo, &event.Remarks,
			&event.CreatedAt, &event.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan exception event: %w", err)
		}
		events = append(events, event)
	}

	return events, nil
}

func (r *ExceptionEventRepository) GetOpenEvents() ([]models.ExceptionEvent, error) {
	query := `SELECT id, event_no, event_type, entity_type, entity_id, trigger_field, threshold_value, actual_value,
		severity, handler, handling_deadline, handled_at, handling_result, status, responsible_person,
		batch_no, remarks, created_at, updated_at 
		FROM exception_events WHERE status IN ('open', 'handling') ORDER BY severity DESC, created_at ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get open exception events: %w", err)
	}
	defer rows.Close()

	var events []models.ExceptionEvent
	for rows.Next() {
		var event models.ExceptionEvent
		err := rows.Scan(
			&event.ID, &event.EventNo, &event.EventType, &event.EntityType, &event.EntityID,
			&event.TriggerField, &event.ThresholdValue, &event.ActualValue, &event.Severity,
			&event.Handler, &event.HandlingDeadline, &event.HandledAt, &event.HandlingResult,
			&event.Status, &event.ResponsiblePerson, &event.BatchNo, &event.Remarks,
			&event.CreatedAt, &event.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan exception event: %w", err)
		}
		events = append(events, event)
	}

	return events, nil
}
