package services

import (
	"airport-vip-service/src/models"
	"airport-vip-service/src/repositories"
	"airport-vip-service/src/schemas"
	"fmt"
	"time"
)

type AuditService struct {
	repo *repositories.AuditLogRepository
}

func NewAuditServiceWithRepos(repo *repositories.AuditLogRepository) *AuditService {
	return &AuditService{repo: repo}
}

func (s *AuditService) GetByID(id int64) (*models.AuditLog, error) {
	return s.repo.GetByID(id)
}

func (s *AuditService) List(req *schemas.AuditLogQueryRequest) (*schemas.AuditLogListResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PerPage <= 0 {
		req.PerPage = 20
	}

	offset := (req.Page - 1) * req.PerPage

	logs, total, err := s.repo.List(offset, req.PerPage, req.EntityType, req.OperationType, req.Operator, req.StartTime, req.EndTime)
	if err != nil {
		return nil, fmt.Errorf("failed to list audit logs: %w", err)
	}

	var data []schemas.AuditLogResponse
	for _, l := range logs {
		data = append(data, schemas.AuditLogResponse{
			ID:                l.ID,
			LogNo:             l.LogNo,
			OperationType:     l.OperationType,
			EntityType:        l.EntityType,
			EntityID:          l.EntityID,
			EntityNo:          l.EntityNo,
			OldValue:          l.OldValue,
			NewValue:          l.NewValue,
			Operator:          l.Operator,
			OperatorRole:      l.OperatorRole,
			IPAddress:         l.IPAddress,
			UserAgent:         l.UserAgent,
			RequestID:         l.RequestID,
			Status:            l.Status,
			ResponsiblePerson: l.ResponsiblePerson,
			BatchNo:           l.BatchNo,
			Remarks:           l.Remarks,
			CreatedAt:         l.CreatedAt,
		})
	}

	return &schemas.AuditLogListResponse{
		Total:   total,
		Page:    req.Page,
		PerPage: req.PerPage,
		Data:    data,
	}, nil
}

func (s *AuditService) GetByEntity(entityType string, entityID int64) ([]schemas.AuditLogResponse, error) {
	logs, err := s.repo.GetByEntity(entityType, entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs by entity: %w", err)
	}

	var data []schemas.AuditLogResponse
	for _, l := range logs {
		data = append(data, schemas.AuditLogResponse{
			ID:                l.ID,
			LogNo:             l.LogNo,
			OperationType:     l.OperationType,
			EntityType:        l.EntityType,
			EntityID:          l.EntityID,
			EntityNo:          l.EntityNo,
			OldValue:          l.OldValue,
			NewValue:          l.NewValue,
			Operator:          l.Operator,
			OperatorRole:      l.OperatorRole,
			IPAddress:         l.IPAddress,
			UserAgent:         l.UserAgent,
			RequestID:         l.RequestID,
			Status:            l.Status,
			ResponsiblePerson: l.ResponsiblePerson,
			BatchNo:           l.BatchNo,
			Remarks:           l.Remarks,
			CreatedAt:         l.CreatedAt,
		})
	}

	return data, nil
}

func (s *AuditService) GetByDateRange(startDate, endDate time.Time) ([]schemas.AuditLogResponse, error) {
	logs, err := s.repo.GetByDateRange(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs by date range: %w", err)
	}

	var data []schemas.AuditLogResponse
	for _, l := range logs {
		data = append(data, schemas.AuditLogResponse{
			ID:                l.ID,
			LogNo:             l.LogNo,
			OperationType:     l.OperationType,
			EntityType:        l.EntityType,
			EntityID:          l.EntityID,
			EntityNo:          l.EntityNo,
			OldValue:          l.OldValue,
			NewValue:          l.NewValue,
			Operator:          l.Operator,
			OperatorRole:      l.OperatorRole,
			IPAddress:         l.IPAddress,
			UserAgent:         l.UserAgent,
			RequestID:         l.RequestID,
			Status:            l.Status,
			ResponsiblePerson: l.ResponsiblePerson,
			BatchNo:           l.BatchNo,
			Remarks:           l.Remarks,
			CreatedAt:         l.CreatedAt,
		})
	}

	return data, nil
}
