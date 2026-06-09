package services

import (
	"airport-vip-service/src/models"
	"airport-vip-service/src/repositories"
	"airport-vip-service/src/schemas"
	"airport-vip-service/src/utils"
	"fmt"
	"time"
)

type ExceptionService struct {
	repo      *repositories.ExceptionEventRepository
	auditRepo *repositories.AuditLogRepository
}

func NewExceptionServiceWithRepos(
	repo *repositories.ExceptionEventRepository,
	auditRepo *repositories.AuditLogRepository,
) *ExceptionService {
	return &ExceptionService{
		repo:      repo,
		auditRepo: auditRepo,
	}
}

func (s *ExceptionService) Create(req *schemas.ExceptionEventCreateRequest) (*models.ExceptionEvent, error) {
	event := &models.ExceptionEvent{
		EventNo:           req.EventNo,
		EventType:         req.EventType,
		EntityType:        req.EntityType,
		EntityID:          req.EntityID,
		TriggerField:      req.TriggerField,
		ThresholdValue:    req.ThresholdValue,
		ActualValue:       req.ActualValue,
		Severity:          req.Severity,
		Handler:           req.Handler,
		HandlingDeadline:  req.HandlingDeadline,
		HandledAt:         req.HandledAt,
		HandlingResult:    req.HandlingResult,
		Status:            req.Status,
		ResponsiblePerson: req.ResponsiblePerson,
		BatchNo:           req.BatchNo,
		Remarks:           req.Remarks,
	}

	if event.Status == "" {
		event.Status = models.ExceptionStatusOpen
	}
	if event.Severity == "" {
		event.Severity = models.ExceptionSeverityMedium
	}

	id, err := s.repo.Create(event)
	if err != nil {
		return nil, fmt.Errorf("failed to create exception event: %w", err)
	}

	event.ID = id

	s.logAudit(models.OperationCreate, "exception_event", id, req.EventNo, "", event)

	return event, nil
}

func (s *ExceptionService) GetByID(id int64) (*models.ExceptionEvent, error) {
	return s.repo.GetByID(id)
}

func (s *ExceptionService) GetByEventNo(eventNo string) (*models.ExceptionEvent, error) {
	return s.repo.GetByEventNo(eventNo)
}

func (s *ExceptionService) List(page, perPage int, eventType, status, severity string) (*schemas.ExceptionEventListResponse, error) {
	offset := (page - 1) * perPage
	events, total, err := s.repo.List(offset, perPage, eventType, status, severity)
	if err != nil {
		return nil, fmt.Errorf("failed to list exception events: %w", err)
	}

	var data []schemas.ExceptionEventResponse
	for _, e := range events {
		data = append(data, schemas.ExceptionEventResponse{
			ID:                e.ID,
			EventNo:           e.EventNo,
			EventType:         e.EventType,
			EntityType:        e.EntityType,
			EntityID:          e.EntityID,
			TriggerField:      e.TriggerField,
			ThresholdValue:    e.ThresholdValue,
			ActualValue:       e.ActualValue,
			Severity:          e.Severity,
			Handler:           e.Handler,
			HandlingDeadline:  e.HandlingDeadline,
			HandledAt:         e.HandledAt,
			HandlingResult:    e.HandlingResult,
			Status:            e.Status,
			ResponsiblePerson: e.ResponsiblePerson,
			BatchNo:           e.BatchNo,
			Remarks:           e.Remarks,
			CreatedAt:         e.CreatedAt,
			UpdatedAt:         e.UpdatedAt,
		})
	}

	return &schemas.ExceptionEventListResponse{
		Total:   total,
		Page:    page,
		PerPage: perPage,
		Data:    data,
	}, nil
}

func (s *ExceptionService) Handle(req *schemas.ExceptionHandleRequest) (*schemas.ExceptionEventResponse, error) {
	event, err := s.repo.GetByEventNo(req.EventNo)
	if err != nil {
		return nil, fmt.Errorf("failed to get exception event: %w", err)
	}
	if event == nil {
		return nil, fmt.Errorf("exception event not found")
	}

	oldStatus := event.Status

	event.Handler = req.Handler
	event.HandlingResult = req.HandlingResult
	event.Status = models.ExceptionStatusResolved

	if req.HandledAt != nil {
		event.HandledAt = req.HandledAt
	} else {
		now := time.Now()
		event.HandledAt = &now
	}

	if err := s.repo.Update(event.ID, event); err != nil {
		return nil, fmt.Errorf("failed to update exception event: %w", err)
	}

	s.logAudit(models.OperationUpdate, "exception_event", event.ID, event.EventNo, oldStatus, event.Status)

	return &schemas.ExceptionEventResponse{
		ID:                event.ID,
		EventNo:           event.EventNo,
		EventType:         event.EventType,
		EntityType:        event.EntityType,
		EntityID:          event.EntityID,
		TriggerField:      event.TriggerField,
		ThresholdValue:    event.ThresholdValue,
		ActualValue:       event.ActualValue,
		Severity:          event.Severity,
		Handler:           event.Handler,
		HandlingDeadline:  event.HandlingDeadline,
		HandledAt:         event.HandledAt,
		HandlingResult:    event.HandlingResult,
		Status:            event.Status,
		ResponsiblePerson: event.ResponsiblePerson,
		BatchNo:           event.BatchNo,
		Remarks:           event.Remarks,
		CreatedAt:         event.CreatedAt,
		UpdatedAt:         event.UpdatedAt,
	}, nil
}

func (s *ExceptionService) GetOpenEvents() ([]models.ExceptionEvent, error) {
	return s.repo.GetOpenEvents()
}

func (s *ExceptionService) logAudit(op, entityType string, entityID int64, entityNo string, oldValue, newValue interface{}) {
	oldValStr := ""
	newValStr := ""
	if oldValue != nil {
		oldValStr = fmt.Sprintf("%v", oldValue)
	}
	if newValue != nil {
		newValStr = fmt.Sprintf("%v", newValue)
	}

	audit := &models.AuditLog{
		LogNo:         utils.GenerateNo("AUD"),
		OperationType: op,
		EntityType:    entityType,
		EntityID:      entityID,
		EntityNo:      entityNo,
		OldValue:      oldValStr,
		NewValue:      newValStr,
		Operator:      "system",
		Status:        models.AuditStatusActive,
		CreatedAt:     time.Now(),
	}
	s.auditRepo.Create(audit)
}
