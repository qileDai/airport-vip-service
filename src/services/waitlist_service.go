package services

import (
	"airport-vip-service/src/models"
	"airport-vip-service/src/repositories"
	"airport-vip-service/src/schemas"
	"airport-vip-service/src/utils"
	"fmt"
	"time"
)

type WaitlistService struct {
	repo        *repositories.WaitlistEntryRepository
	benefitRepo *repositories.MemberBenefitRepository
	flightRepo  *repositories.FlightScheduleRepository
	ruleRepo    *repositories.RuleConfigRepository
	auditRepo   *repositories.AuditLogRepository
}

func NewWaitlistServiceWithRepos(
	repo *repositories.WaitlistEntryRepository,
	benefitRepo *repositories.MemberBenefitRepository,
	flightRepo *repositories.FlightScheduleRepository,
	ruleRepo *repositories.RuleConfigRepository,
	auditRepo *repositories.AuditLogRepository,
) *WaitlistService {
	return &WaitlistService{
		repo:        repo,
		benefitRepo: benefitRepo,
		flightRepo:  flightRepo,
		ruleRepo:    ruleRepo,
		auditRepo:   auditRepo,
	}
}

func (s *WaitlistService) Create(req *schemas.WaitlistEntryCreateRequest) (*models.WaitlistEntry, error) {
	entry := &models.WaitlistEntry{
		WaitlistNo:        req.WaitlistNo,
		ReservationID:     req.ReservationID,
		MemberBenefitID:   req.MemberBenefitID,
		FlightScheduleID:  req.FlightScheduleID,
		MemberName:        req.MemberName,
		MemberLevel:       req.MemberLevel,
		WaitingSince:      req.WaitingSince,
		PriorityScore:     req.PriorityScore,
		EstimatedWaitMins: req.EstimatedWaitMins,
		Status:            req.Status,
		ResponsiblePerson: req.ResponsiblePerson,
		BatchNo:           req.BatchNo,
		Remarks:           req.Remarks,
	}

	if entry.Status == "" {
		entry.Status = models.WaitlistStatusWaiting
	}

	id, err := s.repo.Create(entry)
	if err != nil {
		return nil, fmt.Errorf("failed to create waitlist entry: %w", err)
	}

	entry.ID = id

	s.logAudit(models.OperationCreate, models.EntityTypeWaitlist, id, req.WaitlistNo, "", entry)

	return entry, nil
}

func (s *WaitlistService) GetByID(id int64) (*models.WaitlistEntry, error) {
	return s.repo.GetByID(id)
}

func (s *WaitlistService) List(page, perPage int, flightScheduleID int64, status string) (*schemas.WaitlistEntryListResponse, error) {
	offset := (page - 1) * perPage
	entries, total, err := s.repo.List(offset, perPage, flightScheduleID, status)
	if err != nil {
		return nil, fmt.Errorf("failed to list waitlist entries: %w", err)
	}

	var data []schemas.WaitlistEntryResponse
	for _, e := range entries {
		data = append(data, schemas.WaitlistEntryResponse{
			ID:                e.ID,
			WaitlistNo:        e.WaitlistNo,
			ReservationID:     e.ReservationID,
			MemberBenefitID:   e.MemberBenefitID,
			FlightScheduleID:  e.FlightScheduleID,
			MemberName:        e.MemberName,
			MemberLevel:       e.MemberLevel,
			WaitingSince:      e.WaitingSince,
			PriorityScore:     e.PriorityScore,
			EstimatedWaitMins: e.EstimatedWaitMins,
			Status:            e.Status,
			ResponsiblePerson: e.ResponsiblePerson,
			BatchNo:           e.BatchNo,
			Remarks:           e.Remarks,
			CreatedAt:         e.CreatedAt,
			UpdatedAt:         e.UpdatedAt,
		})
	}

	return &schemas.WaitlistEntryListResponse{
		Total:   total,
		Page:    page,
		PerPage: perPage,
		Data:    data,
	}, nil
}

func (s *WaitlistService) CalculatePriority(req *schemas.PriorityCalculationRequest) (*schemas.PriorityCalculationResponse, error) {
	levelWeight := s.getLevelWeight(req.MemberLevel)
	quotaWeight := s.getQuotaWeight(req.RemainingQuota)
	timeWeight := s.getTimeWeight(req.WaitingSince)
	flightWeight := s.getFlightWeight(req.FlightScheduleID)

	priorityScore := levelWeight + quotaWeight + timeWeight + flightWeight

	waitingEntries, _ := s.repo.GetWaitingByFlightScheduleID(req.FlightScheduleID)
	position := 1
	for _, e := range waitingEntries {
		if e.PriorityScore > priorityScore {
			position++
		}
	}

	estimatedWaitMins := position * 15

	return &schemas.PriorityCalculationResponse{
		PriorityScore:     priorityScore,
		LevelWeight:       levelWeight,
		QuotaWeight:       quotaWeight,
		TimeWeight:        timeWeight,
		FlightWeight:      flightWeight,
		EstimatedWaitMins: estimatedWaitMins,
		Position:          position,
		TotalWaiting:      len(waitingEntries),
	}, nil
}

func (s *WaitlistService) getLevelWeight(level string) int {
	switch level {
	case models.MemberLevelPlatinum:
		return 40
	case models.MemberLevelGold:
		return 30
	case models.MemberLevelSilver:
		return 20
	case models.MemberLevelRegular:
		return 10
	default:
		return 5
	}
}

func (s *WaitlistService) getQuotaWeight(quota int) int {
	if quota >= 10 {
		return 15
	} else if quota >= 5 {
		return 10
	} else if quota >= 1 {
		return 5
	}
	return 0
}

func (s *WaitlistService) getTimeWeight(waitingSince time.Time) int {
	waitingMins := int(time.Since(waitingSince).Minutes())
	if waitingMins >= 120 {
		return 30
	} else if waitingMins >= 60 {
		return 20
	} else if waitingMins >= 30 {
		return 10
	}
	return 5
}

func (s *WaitlistService) getFlightWeight(flightScheduleID int64) int {
	flight, err := s.flightRepo.GetByID(flightScheduleID)
	if err != nil || flight == nil {
		return 0
	}

	timeUntilDepart := time.Until(flight.ScheduledDepart)
	if timeUntilDepart <= 30*time.Minute {
		return 25
	} else if timeUntilDepart <= 60*time.Minute {
		return 15
	} else if timeUntilDepart <= 120*time.Minute {
		return 10
	}
	return 5
}

func (s *WaitlistService) ArrangeWaitlist(req *schemas.WaitlistArrangementRequest) (*schemas.WaitlistArrangementResponse, error) {
	entries, err := s.repo.GetWaitingByFlightScheduleID(req.FlightScheduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get waiting entries: %w", err)
	}

	var arranged, notArranged []schemas.WaitlistEntryResponse
	arrangedCount := 0

	for _, entry := range entries {
		if arrangedCount < req.AvailableSeats {
			entry.Status = models.WaitlistStatusSeated
			s.repo.UpdateStatus(entry.ID, models.WaitlistStatusSeated)
			arrangedCount++

			arranged = append(arranged, schemas.WaitlistEntryResponse{
				ID:                entry.ID,
				WaitlistNo:        entry.WaitlistNo,
				ReservationID:     entry.ReservationID,
				MemberBenefitID:   entry.MemberBenefitID,
				FlightScheduleID:  entry.FlightScheduleID,
				MemberName:        entry.MemberName,
				MemberLevel:       entry.MemberLevel,
				WaitingSince:      entry.WaitingSince,
				PriorityScore:     entry.PriorityScore,
				EstimatedWaitMins: entry.EstimatedWaitMins,
				Status:            entry.Status,
				ResponsiblePerson: entry.ResponsiblePerson,
				BatchNo:           entry.BatchNo,
				Remarks:           entry.Remarks,
				CreatedAt:         entry.CreatedAt,
				UpdatedAt:         entry.UpdatedAt,
			})
		} else {
			notArranged = append(notArranged, schemas.WaitlistEntryResponse{
				ID:                entry.ID,
				WaitlistNo:        entry.WaitlistNo,
				ReservationID:     entry.ReservationID,
				MemberBenefitID:   entry.MemberBenefitID,
				FlightScheduleID:  entry.FlightScheduleID,
				MemberName:        entry.MemberName,
				MemberLevel:       entry.MemberLevel,
				WaitingSince:      entry.WaitingSince,
				PriorityScore:     entry.PriorityScore,
				EstimatedWaitMins: entry.EstimatedWaitMins,
				Status:            entry.Status,
				ResponsiblePerson: entry.ResponsiblePerson,
				BatchNo:           entry.BatchNo,
				Remarks:           entry.Remarks,
				CreatedAt:         entry.CreatedAt,
				UpdatedAt:         entry.UpdatedAt,
			})
		}
	}

	s.logAudit(models.OperationUpdate, models.EntityTypeWaitlist, req.FlightScheduleID, "arrangement", "", fmt.Sprintf("arranged %d entries", arrangedCount))

	return &schemas.WaitlistArrangementResponse{
		ArrangedCount: arrangedCount,
		Remaining:     len(entries) - arrangedCount,
		Arranged:      arranged,
		NotArranged:   notArranged,
	}, nil
}

func (s *WaitlistService) UpdateStatus(id int64, status string) error {
	return s.repo.UpdateStatus(id, status)
}

func (s *WaitlistService) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *WaitlistService) logAudit(op, entityType string, entityID int64, entityNo string, oldValue, newValue interface{}) {
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
