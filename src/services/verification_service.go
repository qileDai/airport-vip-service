package services

import (
	"airport-vip-service/src/models"
	"airport-vip-service/src/repositories"
	"airport-vip-service/src/schemas"
	"airport-vip-service/src/utils"
	"fmt"
	"time"
)

type VerificationService struct {
	repo          *repositories.VerificationResultRepository
	benefitRepo   *repositories.MemberBenefitRepository
	reservationRepo *repositories.ReservationRecordRepository
	flightRepo    *repositories.FlightScheduleRepository
	companionRepo *repositories.CompanionRepository
	ruleRepo      *repositories.RuleConfigRepository
	auditRepo     *repositories.AuditLogRepository
}

func NewVerificationServiceWithRepos(
	repo *repositories.VerificationResultRepository,
	benefitRepo *repositories.MemberBenefitRepository,
	reservationRepo *repositories.ReservationRecordRepository,
	flightRepo *repositories.FlightScheduleRepository,
	companionRepo *repositories.CompanionRepository,
	ruleRepo *repositories.RuleConfigRepository,
	auditRepo *repositories.AuditLogRepository,
) *VerificationService {
	return &VerificationService{
		repo:          repo,
		benefitRepo:   benefitRepo,
		reservationRepo: reservationRepo,
		flightRepo:    flightRepo,
		companionRepo: companionRepo,
		ruleRepo:      ruleRepo,
		auditRepo:     auditRepo,
	}
}

func (s *VerificationService) Create(req *schemas.VerificationResultCreateRequest) (*models.VerificationResult, error) {
	result := &models.VerificationResult{
		VerificationNo:     req.VerificationNo,
		ReservationID:      req.ReservationID,
		MemberBenefitID:    req.MemberBenefitID,
		FlightScheduleID:   req.FlightScheduleID,
		VerificationType:   req.VerificationType,
		Result:             req.Result,
		FailureReason:      req.FailureReason,
		VerifiedQuota:      req.VerifiedQuota,
		VerifiedCompanions: req.VerifiedCompanions,
		VerificationDetails: req.VerificationDetails,
		Status:             req.Status,
		ResponsiblePerson:  req.ResponsiblePerson,
		BatchNo:            req.BatchNo,
		Remarks:            req.Remarks,
	}

	if result.Status == "" {
		result.Status = models.VerifyStatusDraft
	}
	if result.Result == "" {
		result.Result = models.VerifyResultPending
	}

	id, err := s.repo.Create(result)
	if err != nil {
		return nil, fmt.Errorf("failed to create verification result: %w", err)
	}

	result.ID = id

	s.logAudit(models.OperationCreate, models.EntityTypeVerify, id, req.VerificationNo, "", result)

	return result, nil
}

func (s *VerificationService) GetByID(id int64) (*models.VerificationResult, error) {
	return s.repo.GetByID(id)
}

func (s *VerificationService) List(page, perPage int, status, resultFilter string) (*schemas.VerificationResultListResponse, error) {
	offset := (page - 1) * perPage
	results, total, err := s.repo.List(offset, perPage, status, resultFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to list verification results: %w", err)
	}

	var data []schemas.VerificationResultResponse
	for _, r := range results {
		data = append(data, schemas.VerificationResultResponse{
			ID:                  r.ID,
			VerificationNo:      r.VerificationNo,
			ReservationID:       r.ReservationID,
			MemberBenefitID:     r.MemberBenefitID,
			FlightScheduleID:    r.FlightScheduleID,
			VerificationType:    r.VerificationType,
			Result:              r.Result,
			FailureReason:       r.FailureReason,
			VerifiedQuota:       r.VerifiedQuota,
			VerifiedCompanions:  r.VerifiedCompanions,
			VerificationDetails: r.VerificationDetails,
			Status:              r.Status,
			ResponsiblePerson:   r.ResponsiblePerson,
			BatchNo:             r.BatchNo,
			Remarks:             r.Remarks,
			CreatedAt:           r.CreatedAt,
			UpdatedAt:           r.UpdatedAt,
		})
	}

	return &schemas.VerificationResultListResponse{
		Total:   total,
		Page:    page,
		PerPage: perPage,
		Data:    data,
	}, nil
}

func (s *VerificationService) VerifyEligibility(req *schemas.EligibilityVerificationRequest) (*schemas.EligibilityVerificationResponse, error) {
	response := &schemas.EligibilityVerificationResponse{
		IsEligible: true,
		Errors:     []string{},
		Warnings:   []string{},
	}

	benefit, err := s.benefitRepo.GetByID(req.MemberBenefitID)
	if err != nil {
		return nil, fmt.Errorf("failed to get member benefit: %w", err)
	}
	if benefit == nil {
		response.IsEligible = false
		response.Errors = append(response.Errors, "member benefit not found")
		response.BenefitValid = false
		return response, nil
	}

	response.BenefitValid = true
	if benefit.Status != models.BenefitStatusActive {
		response.IsEligible = false
		response.Errors = append(response.Errors, "member benefit is not active")
		response.BenefitValid = false
	}

	if benefit.ValidTo.Before(time.Now()) {
		response.IsEligible = false
		response.Errors = append(response.Errors, "member benefit has expired")
		response.BenefitValid = false
	}

	reservation, err := s.reservationRepo.GetByID(req.ReservationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reservation: %w", err)
	}
	if reservation == nil {
		response.IsEligible = false
		response.Errors = append(response.Errors, "reservation not found")
		return response, nil
	}

	if reservation.MemberBenefitID != req.MemberBenefitID {
		response.IsEligible = false
		response.Errors = append(response.Errors, "reservation does not belong to this member benefit")
	}

	flight, err := s.flightRepo.GetByID(req.FlightScheduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get flight schedule: %w", err)
	}
	if flight == nil {
		response.IsEligible = false
		response.Errors = append(response.Errors, "flight schedule not found")
		response.FlightValid = false
		return response, nil
	}

	response.FlightValid = true
	if flight.Status != models.FlightRecordStatusActive {
		response.Warnings = append(response.Warnings, "flight schedule is not active")
	}

	if flight.VipLoungeUsed >= flight.VipLoungeCapacity {
		response.IsEligible = false
		response.Errors = append(response.Errors, "VIP lounge is full")
		response.FlightValid = false
	}

	companions, err := s.companionRepo.GetByReservationID(req.ReservationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get companions: %w", err)
	}

	response.CompanionsValid = true
	maxCompanions := s.getMaxCompanions(benefit.MemberLevel)
	if len(companions) > maxCompanions {
		response.IsEligible = false
		response.Errors = append(response.Errors, fmt.Sprintf("exceeded maximum companions: allowed %d, have %d", maxCompanions, len(companions)))
		response.CompanionsValid = false
	}

	for _, c := range companions {
		if c.VerificationStatus == models.CompanionVerifyFailed {
			response.Warnings = append(response.Warnings, fmt.Sprintf("companion %s verification failed", c.CompanionName))
		}
	}

	response.QuotaValid = benefit.RemainingQuota >= reservation.GuestCount
	if !response.QuotaValid {
		response.IsEligible = false
		response.Errors = append(response.Errors, fmt.Sprintf("insufficient quota: have %d, need %d", benefit.RemainingQuota, reservation.GuestCount))
	}

	response.VerifiedQuota = reservation.GuestCount
	response.VerifiedCompanions = len(companions)

	verificationNo := utils.GenerateNo("VER")
	verification := &models.VerificationResult{
		VerificationNo:     verificationNo,
		ReservationID:      req.ReservationID,
		MemberBenefitID:    req.MemberBenefitID,
		FlightScheduleID:   req.FlightScheduleID,
		VerificationType:   models.VerifyTypeReservation,
		Result:             models.VerifyResultPassed,
		VerifiedQuota:      response.VerifiedQuota,
		VerifiedCompanions: response.VerifiedCompanions,
		Status:             models.VerifyStatusConfirmed,
	}

	if !response.IsEligible {
		verification.Result = models.VerifyResultFailed
		verification.FailureReason = response.Errors[0]
		verification.Status = models.VerifyStatusDraft
	}

	s.repo.Create(verification)
	response.VerificationNo = verificationNo

	s.logAudit(models.OperationVerify, models.EntityTypeVerify, verification.ID, verificationNo, "", verification)

	return response, nil
}

func (s *VerificationService) getMaxCompanions(memberLevel string) int {
	rules, err := s.ruleRepo.GetActiveByType(models.RuleTypeCompanionLimit)
	if err != nil || len(rules) == 0 {
		switch memberLevel {
		case models.MemberLevelPlatinum:
			return 4
		case models.MemberLevelGold:
			return 3
		case models.MemberLevelSilver:
			return 2
		default:
			return 1
		}
	}

	for _, rule := range rules {
		if rule.AppliesToLevel == memberLevel || rule.AppliesToLevel == "" {
			return int(rule.ThresholdValue)
		}
	}

	return 1
}

func (s *VerificationService) GetStatistics(startDate, endDate time.Time) (*schemas.VerificationStatisticsResponse, error) {
	stats, err := s.repo.GetStatistics(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get verification statistics: %w", err)
	}

	total := stats["total"]
	passed := stats["passed"]
	failed := stats["failed"]
	pending := stats["pending"]

	var passRate, failRate float64
	if total > 0 {
		passRate = float64(passed) / float64(total) * 100
		failRate = float64(failed) / float64(total) * 100
	}

	return &schemas.VerificationStatisticsResponse{
		TotalVerifications: total,
		PassedCount:        passed,
		FailedCount:        failed,
		PendingCount:       pending,
		PassRate:           passRate,
		FailRate:           failRate,
		ByDate:             []schemas.DateStatistics{},
		ByBatch:            []schemas.BatchStatistics{},
		ByRole:             []schemas.RoleStatistics{},
	}, nil
}

func (s *VerificationService) logAudit(op, entityType string, entityID int64, entityNo string, oldValue, newValue interface{}) {
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
