package services

import (
	"airport-vip-service/src/models"
	"airport-vip-service/src/repositories"
	"airport-vip-service/src/schemas"
	"airport-vip-service/src/utils"
	"fmt"
	"time"
)

type MemberBenefitService struct {
	repo          *repositories.MemberBenefitRepository
	exceptionRepo *repositories.ExceptionEventRepository
	auditRepo     *repositories.AuditLogRepository
}

func NewMemberBenefitService(db interface {
}) *MemberBenefitService {
	return &MemberBenefitService{}
}

func NewMemberBenefitServiceWithRepos(
	repo *repositories.MemberBenefitRepository,
	exceptionRepo *repositories.ExceptionEventRepository,
	auditRepo *repositories.AuditLogRepository,
) *MemberBenefitService {
	return &MemberBenefitService{
		repo:          repo,
		exceptionRepo: exceptionRepo,
		auditRepo:     auditRepo,
	}
}

func (s *MemberBenefitService) Create(req *schemas.MemberBenefitCreateRequest) (*models.MemberBenefit, error) {
	benefit := &models.MemberBenefit{
		BenefitNo:         req.BenefitNo,
		MemberName:        req.MemberName,
		MemberLevel:       req.MemberLevel,
		RemainingQuota:    req.RemainingQuota,
		TotalQuota:        req.TotalQuota,
		Status:            req.Status,
		ResponsiblePerson: req.ResponsiblePerson,
		ValidFrom:         req.ValidFrom,
		ValidTo:           req.ValidTo,
		BatchNo:           req.BatchNo,
		Remarks:           req.Remarks,
	}

	if benefit.Status == "" {
		benefit.Status = models.BenefitStatusActive
	}

	id, err := s.repo.Create(benefit)
	if err != nil {
		return nil, fmt.Errorf("failed to create member benefit: %w", err)
	}

	benefit.ID = id

	s.logAudit(models.OperationCreate, models.EntityTypeBenefit, id, req.BenefitNo, "", benefit)

	return benefit, nil
}

func (s *MemberBenefitService) GetByID(id int64) (*models.MemberBenefit, error) {
	return s.repo.GetByID(id)
}

func (s *MemberBenefitService) GetByBenefitNo(benefitNo string) (*models.MemberBenefit, error) {
	return s.repo.GetByBenefitNo(benefitNo)
}

func (s *MemberBenefitService) List(page, perPage int, status, memberLevel string) (*schemas.MemberBenefitListResponse, error) {
	offset := (page - 1) * perPage
	benefits, total, err := s.repo.List(offset, perPage, status, memberLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to list member benefits: %w", err)
	}

	var data []schemas.MemberBenefitResponse
	for _, b := range benefits {
		data = append(data, schemas.MemberBenefitResponse{
			ID:                b.ID,
			BenefitNo:         b.BenefitNo,
			MemberName:        b.MemberName,
			MemberLevel:       b.MemberLevel,
			RemainingQuota:    b.RemainingQuota,
			TotalQuota:        b.TotalQuota,
			Status:            b.Status,
			ResponsiblePerson: b.ResponsiblePerson,
			ValidFrom:         b.ValidFrom,
			ValidTo:           b.ValidTo,
			BatchNo:           b.BatchNo,
			Remarks:           b.Remarks,
			CreatedAt:         b.CreatedAt,
			UpdatedAt:         b.UpdatedAt,
		})
	}

	return &schemas.MemberBenefitListResponse{
		Total:   total,
		Page:    page,
		PerPage: perPage,
		Data:    data,
	}, nil
}

func (s *MemberBenefitService) Update(id int64, req *schemas.MemberBenefitUpdateRequest) (*models.MemberBenefit, error) {
	benefit, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get member benefit: %w", err)
	}
	if benefit == nil {
		return nil, fmt.Errorf("member benefit not found")
	}

	oldValue := *benefit

	if req.MemberName != "" {
		benefit.MemberName = req.MemberName
	}
	if req.MemberLevel != "" {
		benefit.MemberLevel = req.MemberLevel
	}
	if req.RemainingQuota != nil {
		benefit.RemainingQuota = *req.RemainingQuota
	}
	if req.TotalQuota != nil {
		benefit.TotalQuota = *req.TotalQuota
	}
	if req.Status != "" {
		benefit.Status = req.Status
	}
	if req.ResponsiblePerson != "" {
		benefit.ResponsiblePerson = req.ResponsiblePerson
	}
	if req.ValidFrom != nil {
		benefit.ValidFrom = *req.ValidFrom
	}
	if req.ValidTo != nil {
		benefit.ValidTo = *req.ValidTo
	}
	if req.BatchNo != "" {
		benefit.BatchNo = req.BatchNo
	}
	if req.Remarks != "" {
		benefit.Remarks = req.Remarks
	}

	if err := s.repo.Update(id, benefit); err != nil {
		return nil, fmt.Errorf("failed to update member benefit: %w", err)
	}

	s.logAudit(models.OperationUpdate, models.EntityTypeBenefit, id, benefit.BenefitNo, oldValue, benefit)

	return benefit, nil
}

func (s *MemberBenefitService) Delete(id int64) error {
	benefit, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get member benefit: %w", err)
	}
	if benefit == nil {
		return fmt.Errorf("member benefit not found")
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete member benefit: %w", err)
	}

	s.logAudit(models.OperationDelete, models.EntityTypeBenefit, id, benefit.BenefitNo, benefit, "")

	return nil
}

func (s *MemberBenefitService) CheckExpiry() (*schemas.ExpiryCheckResponse, error) {
	expired, err := s.repo.GetExpired()
	if err != nil {
		return nil, fmt.Errorf("failed to get expired benefits: %w", err)
	}

	expiringSoon, err := s.repo.GetExpiringSoon(7)
	if err != nil {
		return nil, fmt.Errorf("failed to get expiring benefits: %w", err)
	}

	var expiredEvents []schemas.ExceptionEventResponse
	for _, b := range expired {
		event := &models.ExceptionEvent{
			EventNo:           utils.GenerateNo("EXP"),
			EventType:         models.EventTypeBenefitExpired,
			EntityType:        models.EntityTypeBenefit,
			EntityID:          b.ID,
			TriggerField:      "valid_to",
			ThresholdValue:    time.Now().Format("2006-01-02"),
			ActualValue:       b.ValidTo.Format("2006-01-02"),
			Severity:          models.ExceptionSeverityHigh,
			Status:            models.ExceptionStatusOpen,
			ResponsiblePerson: b.ResponsiblePerson,
		}
		if _, err := s.exceptionRepo.Create(event); err == nil {
			s.repo.UpdateStatus(b.ID, models.BenefitStatusExpired)
			expiredEvents = append(expiredEvents, schemas.ExceptionEventResponse{
				ID:                event.ID,
				EventNo:           event.EventNo,
				EventType:         event.EventType,
				EntityType:        event.EntityType,
				EntityID:          event.EntityID,
				TriggerField:      event.TriggerField,
				ThresholdValue:    event.ThresholdValue,
				ActualValue:       event.ActualValue,
				Severity:          event.Severity,
				Status:            event.Status,
				ResponsiblePerson: event.ResponsiblePerson,
				CreatedAt:         event.CreatedAt,
				UpdatedAt:         event.UpdatedAt,
			})
		}
	}

	var expiringEvents []schemas.ExceptionEventResponse
	for _, b := range expiringSoon {
		expiringEvents = append(expiringEvents, schemas.ExceptionEventResponse{
			EventNo:     utils.GenerateNo("EXP"),
			EventType:   "benefit_expiring_soon",
			EntityType:  models.EntityTypeBenefit,
			EntityID:    b.ID,
			TriggerField: "valid_to",
			ActualValue: b.ValidTo.Format("2006-01-02"),
			Severity:    models.ExceptionSeverityMedium,
			Status:      models.ExceptionStatusOpen,
		})
	}

	return &schemas.ExpiryCheckResponse{
		ExpiredCount:   len(expiredEvents),
		ExpiringSoon:   len(expiringSoon),
		ExpiredEvents:  expiredEvents,
		ExpiringEvents: expiringEvents,
	}, nil
}

func (s *MemberBenefitService) logAudit(op, entityType string, entityID int64, entityNo string, oldValue, newValue interface{}) {
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
	}
	s.auditRepo.Create(audit)
}
