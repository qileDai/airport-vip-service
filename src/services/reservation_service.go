package services

import (
	"airport-vip-service/src/models"
	"airport-vip-service/src/repositories"
	"airport-vip-service/src/schemas"
	"airport-vip-service/src/utils"
	"fmt"
	"time"
)

type ReservationService struct {
	repo          *repositories.ReservationRecordRepository
	benefitRepo   *repositories.MemberBenefitRepository
	flightRepo    *repositories.FlightScheduleRepository
	companionRepo *repositories.CompanionRepository
	transitionRepo *repositories.StatusTransitionRepository
	auditRepo     *repositories.AuditLogRepository
}

func NewReservationServiceWithRepos(
	repo *repositories.ReservationRecordRepository,
	benefitRepo *repositories.MemberBenefitRepository,
	flightRepo *repositories.FlightScheduleRepository,
	companionRepo *repositories.CompanionRepository,
	transitionRepo *repositories.StatusTransitionRepository,
	auditRepo *repositories.AuditLogRepository,
) *ReservationService {
	return &ReservationService{
		repo:          repo,
		benefitRepo:   benefitRepo,
		flightRepo:    flightRepo,
		companionRepo: companionRepo,
		transitionRepo: transitionRepo,
		auditRepo:     auditRepo,
	}
}

func (s *ReservationService) Create(req *schemas.ReservationRecordCreateRequest) (*models.ReservationRecord, error) {
	record := &models.ReservationRecord{
		ReservationNo:     req.ReservationNo,
		MemberBenefitID:   req.MemberBenefitID,
		MemberName:        req.MemberName,
		FlightNo:          req.FlightNo,
		FlightScheduleID:  req.FlightScheduleID,
		VipLoungeName:     req.VipLoungeName,
		ReservationTime:   req.ReservationTime,
		GuestCount:        req.GuestCount,
		Status:            req.Status,
		ResponsiblePerson: req.ResponsiblePerson,
		BatchNo:           req.BatchNo,
		Remarks:           req.Remarks,
	}

	if record.Status == "" {
		record.Status = models.ReservationStatusDraft
	}

	id, err := s.repo.Create(record)
	if err != nil {
		return nil, fmt.Errorf("failed to create reservation record: %w", err)
	}

	record.ID = id

	s.logAudit(models.OperationCreate, models.EntityTypeReservation, id, req.ReservationNo, "", record)

	return record, nil
}

func (s *ReservationService) GetByID(id int64) (*models.ReservationRecord, error) {
	return s.repo.GetByID(id)
}

func (s *ReservationService) List(page, perPage int, status string, memberBenefitID int64) (*schemas.ReservationRecordListResponse, error) {
	offset := (page - 1) * perPage
	records, total, err := s.repo.List(offset, perPage, status, memberBenefitID)
	if err != nil {
		return nil, fmt.Errorf("failed to list reservation records: %w", err)
	}

	var data []schemas.ReservationRecordResponse
	for _, r := range records {
		data = append(data, schemas.ReservationRecordResponse{
			ID:                r.ID,
			ReservationNo:     r.ReservationNo,
			MemberBenefitID:   r.MemberBenefitID,
			MemberName:        r.MemberName,
			FlightNo:          r.FlightNo,
			FlightScheduleID:  r.FlightScheduleID,
			VipLoungeName:     r.VipLoungeName,
			ReservationTime:   r.ReservationTime,
			GuestCount:        r.GuestCount,
			Status:            r.Status,
			ResponsiblePerson: r.ResponsiblePerson,
			BatchNo:           r.BatchNo,
			Remarks:           r.Remarks,
			CreatedAt:         r.CreatedAt,
			UpdatedAt:         r.UpdatedAt,
		})
	}

	return &schemas.ReservationRecordListResponse{
		Total:   total,
		Page:    page,
		PerPage: perPage,
		Data:    data,
	}, nil
}

func (s *ReservationService) Update(id int64, req *schemas.ReservationRecordUpdateRequest) (*models.ReservationRecord, error) {
	record, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get reservation record: %w", err)
	}
	if record == nil {
		return nil, fmt.Errorf("reservation record not found")
	}

	oldValue := *record

	if req.MemberName != "" {
		record.MemberName = req.MemberName
	}
	if req.FlightNo != "" {
		record.FlightNo = req.FlightNo
	}
	if req.FlightScheduleID != nil {
		record.FlightScheduleID = *req.FlightScheduleID
	}
	if req.VipLoungeName != "" {
		record.VipLoungeName = req.VipLoungeName
	}
	if req.ReservationTime != nil {
		record.ReservationTime = *req.ReservationTime
	}
	if req.GuestCount != nil {
		record.GuestCount = *req.GuestCount
	}
	if req.Status != "" {
		record.Status = req.Status
	}
	if req.ResponsiblePerson != "" {
		record.ResponsiblePerson = req.ResponsiblePerson
	}
	if req.BatchNo != "" {
		record.BatchNo = req.BatchNo
	}
	if req.Remarks != "" {
		record.Remarks = req.Remarks
	}

	if err := s.repo.Update(id, record); err != nil {
		return nil, fmt.Errorf("failed to update reservation record: %w", err)
	}

	s.logAudit(models.OperationUpdate, models.EntityTypeReservation, id, record.ReservationNo, oldValue, record)

	return record, nil
}

func (s *ReservationService) Delete(id int64) error {
	record, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get reservation record: %w", err)
	}
	if record == nil {
		return fmt.Errorf("reservation record not found")
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete reservation record: %w", err)
	}

	s.logAudit(models.OperationDelete, models.EntityTypeReservation, id, record.ReservationNo, record, "")

	return nil
}

func (s *ReservationService) BatchImportPreview(req *schemas.BatchImportRequest) (*schemas.BatchImportPreviewResponse, error) {
	var validRecords, invalidRecords int
	var errors []schemas.BatchImportError

	for i, record := range req.Records {
		if err := s.validateReservationRecord(&record); err != nil {
			invalidRecords++
			errors = append(errors, schemas.BatchImportError{
				RowIndex: i,
				Field:    "validation",
				Message:  err.Error(),
			})
		} else {
			validRecords++
		}
	}

	return &schemas.BatchImportPreviewResponse{
		TotalRecords:     len(req.Records),
		ValidRecords:     validRecords,
		InvalidRecords:   invalidRecords,
		ValidationErrors: errors,
	}, nil
}

func (s *ReservationService) BatchImport(req *schemas.BatchImportRequest) (*schemas.BatchImportResultResponse, error) {
	var successCount, failCount int
	var errors []schemas.BatchImportError
	var createdIDs []int64

	var validRecords []models.ReservationRecord
	for i, record := range req.Records {
		if err := s.validateReservationRecord(&record); err != nil {
			failCount++
			errors = append(errors, schemas.BatchImportError{
				RowIndex: i,
				Field:    "validation",
				Message:  err.Error(),
			})
			continue
		}

		validRecords = append(validRecords, models.ReservationRecord{
			ReservationNo:     record.ReservationNo,
			MemberBenefitID:   record.MemberBenefitID,
			MemberName:        record.MemberName,
			FlightNo:          record.FlightNo,
			FlightScheduleID:  record.FlightScheduleID,
			VipLoungeName:     record.VipLoungeName,
			ReservationTime:   record.ReservationTime,
			GuestCount:        record.GuestCount,
			Status:            record.Status,
			ResponsiblePerson: record.ResponsiblePerson,
			BatchNo:           record.BatchNo,
			Remarks:           record.Remarks,
		})
	}

	if len(validRecords) > 0 {
		ids, err := s.repo.BatchCreate(validRecords)
		if err != nil {
			return nil, fmt.Errorf("failed to batch create reservations: %w", err)
		}
		createdIDs = ids
		successCount = len(ids)
	}

	s.logAudit(models.OperationBatchImport, models.EntityTypeReservation, 0, "batch", "", fmt.Sprintf("created %d records", successCount))

	return &schemas.BatchImportResultResponse{
		SuccessCount: successCount,
		FailCount:    failCount,
		Errors:       errors,
		CreatedIDs:   createdIDs,
	}, nil
}

func (s *ReservationService) validateReservationRecord(req *schemas.ReservationRecordCreateRequest) error {
	benefit, err := s.benefitRepo.GetByID(req.MemberBenefitID)
	if err != nil {
		return fmt.Errorf("failed to verify member benefit: %w", err)
	}
	if benefit == nil {
		return fmt.Errorf("member benefit not found")
	}
	if benefit.Status != models.BenefitStatusActive {
		return fmt.Errorf("member benefit is not active")
	}
	if benefit.RemainingQuota < req.GuestCount {
		return fmt.Errorf("insufficient quota: have %d, need %d", benefit.RemainingQuota, req.GuestCount)
	}

	if req.FlightScheduleID > 0 {
		flight, err := s.flightRepo.GetByID(req.FlightScheduleID)
		if err != nil {
			return fmt.Errorf("failed to verify flight schedule: %w", err)
		}
		if flight == nil {
			return fmt.Errorf("flight schedule not found")
		}
	}

	return nil
}

func (s *ReservationService) ChangeStatus(req *schemas.StatusChangeRequest) (*schemas.StatusChangeResponse, error) {
	record, err := s.repo.GetByID(req.EntityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reservation record: %w", err)
	}
	if record == nil {
		return nil, fmt.Errorf("reservation record not found")
	}

	allowedNext, requiresReason := s.getAllowedTransitions(record.Status)
	allowed := false
	for _, s := range allowedNext {
		if s == req.ToStatus {
			allowed = true
			break
		}
	}

	if !allowed {
		return &schemas.StatusChangeResponse{
			Success:       false,
			PreviousStatus: record.Status,
			Errors:        []string{fmt.Sprintf("cannot transition from %s to %s", record.Status, req.ToStatus)},
		}, nil
	}

	needsReason := false
	for _, s := range requiresReason {
		if s == req.ToStatus {
			needsReason = true
			break
		}
	}

	if needsReason && req.Reason == "" {
		return &schemas.StatusChangeResponse{
			Success:       false,
			PreviousStatus: record.Status,
			Errors:        []string{"reason is required for this status change"},
		}, nil
	}

	previousStatus := record.Status
	record.Status = req.ToStatus

	if err := s.repo.Update(req.EntityID, record); err != nil {
		return nil, fmt.Errorf("failed to update reservation status: %w", err)
	}

	transition := &models.StatusTransition{
		TransitionNo:      utils.GenerateNo("TRN"),
		EntityType:        models.EntityTypeReservation,
		EntityID:          req.EntityID,
		FromStatus:        previousStatus,
		ToStatus:          req.ToStatus,
		Action:            "status_change",
		Reason:            req.Reason,
		Operator:          req.Operator,
		Status:            models.TransitionStatusActive,
		ResponsiblePerson: record.ResponsiblePerson,
	}
	s.transitionRepo.Create(transition)

	s.logAudit(models.OperationStatusChange, models.EntityTypeReservation, req.EntityID, record.ReservationNo, previousStatus, req.ToStatus)

	return &schemas.StatusChangeResponse{
		Success:        true,
		TransitionNo:   transition.TransitionNo,
		PreviousStatus: previousStatus,
		NewStatus:      req.ToStatus,
		Transition: schemas.StatusTransitionResponse{
			ID:           transition.ID,
			TransitionNo: transition.TransitionNo,
			EntityType:   transition.EntityType,
			EntityID:     transition.EntityID,
			FromStatus:   transition.FromStatus,
			ToStatus:     transition.ToStatus,
			Action:       transition.Action,
			Reason:       transition.Reason,
			Operator:     transition.Operator,
			Status:       transition.Status,
			CreatedAt:    transition.CreatedAt,
			UpdatedAt:    transition.UpdatedAt,
		},
	}, nil
}

func (s *ReservationService) GetAllowedTransitions(currentStatus string) *schemas.AllowedTransitionsResponse {
	allowed, requiresReason := s.getAllowedTransitions(currentStatus)
	return &schemas.AllowedTransitionsResponse{
		CurrentStatus: currentStatus,
		AllowedNext:   allowed,
		RequiresReason: requiresReason,
	}
}

func (s *ReservationService) getAllowedTransitions(currentStatus string) ([]string, []string) {
	var allowed []string
	var requiresReason []string

	switch currentStatus {
	case models.ReservationStatusDraft:
		allowed = []string{models.ReservationStatusPending, models.ReservationStatusCancelled}
	case models.ReservationStatusPending:
		allowed = []string{models.ReservationStatusConfirmed, models.ReservationStatusSupplement, models.ReservationStatusRejected}
		requiresReason = []string{models.ReservationStatusRejected}
	case models.ReservationStatusSupplement:
		allowed = []string{models.ReservationStatusPending, models.ReservationStatusCancelled}
	case models.ReservationStatusConfirmed:
		allowed = []string{models.ReservationStatusCompleted, models.ReservationStatusCancelled}
	case models.ReservationStatusCompleted:
		allowed = []string{models.ReservationStatusArchived}
	case models.ReservationStatusRejected:
		allowed = []string{models.ReservationStatusArchived}
		requiresReason = []string{models.ReservationStatusArchived}
	case models.ReservationStatusCancelled:
		allowed = []string{models.ReservationStatusArchived}
	}

	return allowed, requiresReason
}

func (s *ReservationService) logAudit(op, entityType string, entityID int64, entityNo string, oldValue, newValue interface{}) {
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
