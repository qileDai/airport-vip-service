package services

import (
	"airport-vip-service/src/models"
	"airport-vip-service/src/repositories"
	"airport-vip-service/src/schemas"
	"airport-vip-service/src/utils"
	"encoding/json"
	"fmt"
	"time"
)

type FlightScheduleService struct {
	repo      *repositories.FlightScheduleRepository
	auditRepo *repositories.AuditLogRepository
}

func NewFlightScheduleServiceWithRepos(
	repo *repositories.FlightScheduleRepository,
	auditRepo *repositories.AuditLogRepository,
) *FlightScheduleService {
	return &FlightScheduleService{
		repo:      repo,
		auditRepo: auditRepo,
	}
}

func (s *FlightScheduleService) Create(req *schemas.FlightScheduleCreateRequest) (*models.FlightSchedule, error) {
	schedule := &models.FlightSchedule{
		FlightNo:          req.FlightNo,
		DepartureAirport:  req.DepartureAirport,
		ArrivalAirport:    req.ArrivalAirport,
		ScheduledDepart:   req.ScheduledDepart,
		ScheduledArrive:   req.ScheduledArrive,
		ActualDepart:      req.ActualDepart,
		ActualArrive:      req.ActualArrive,
		FlightStatus:      req.FlightStatus,
		VipLoungeCapacity: req.VipLoungeCapacity,
		VipLoungeUsed:     req.VipLoungeUsed,
		Status:            req.Status,
		ResponsiblePerson: req.ResponsiblePerson,
		BatchNo:           req.BatchNo,
		Remarks:           req.Remarks,
	}

	if schedule.FlightStatus == "" {
		schedule.FlightStatus = models.FlightStatusScheduled
	}
	if schedule.Status == "" {
		schedule.Status = models.FlightRecordStatusActive
	}
	if schedule.VipLoungeCapacity == 0 {
		schedule.VipLoungeCapacity = 50
	}

	id, err := s.repo.Create(schedule)
	if err != nil {
		return nil, fmt.Errorf("failed to create flight schedule: %w", err)
	}

	schedule.ID = id

	s.logAudit(models.OperationCreate, models.EntityTypeFlight, id, req.FlightNo, "", schedule)

	return schedule, nil
}

func (s *FlightScheduleService) GetByID(id int64) (*models.FlightSchedule, error) {
	return s.repo.GetByID(id)
}

func (s *FlightScheduleService) List(page, perPage int, status string) (*schemas.FlightScheduleListResponse, error) {
	offset := (page - 1) * perPage
	schedules, total, err := s.repo.List(offset, perPage, status)
	if err != nil {
		return nil, fmt.Errorf("failed to list flight schedules: %w", err)
	}

	var data []schemas.FlightScheduleResponse
	for _, sch := range schedules {
		data = append(data, schemas.FlightScheduleResponse{
			ID:                sch.ID,
			FlightNo:          sch.FlightNo,
			DepartureAirport:  sch.DepartureAirport,
			ArrivalAirport:    sch.ArrivalAirport,
			ScheduledDepart:   sch.ScheduledDepart,
			ScheduledArrive:   sch.ScheduledArrive,
			ActualDepart:      sch.ActualDepart,
			ActualArrive:      sch.ActualArrive,
			FlightStatus:      sch.FlightStatus,
			VipLoungeCapacity: sch.VipLoungeCapacity,
			VipLoungeUsed:     sch.VipLoungeUsed,
			Status:            sch.Status,
			ResponsiblePerson: sch.ResponsiblePerson,
			BatchNo:           sch.BatchNo,
			Remarks:           sch.Remarks,
			SnapshotData:      sch.SnapshotData,
			CreatedAt:         sch.CreatedAt,
			UpdatedAt:         sch.UpdatedAt,
		})
	}

	return &schemas.FlightScheduleListResponse{
		Total:   total,
		Page:    page,
		PerPage: perPage,
		Data:    data,
	}, nil
}

func (s *FlightScheduleService) Update(id int64, req *schemas.FlightScheduleUpdateRequest) (*models.FlightSchedule, error) {
	schedule, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get flight schedule: %w", err)
	}
	if schedule == nil {
		return nil, fmt.Errorf("flight schedule not found")
	}

	oldValue := *schedule

	if req.DepartureAirport != "" {
		schedule.DepartureAirport = req.DepartureAirport
	}
	if req.ArrivalAirport != "" {
		schedule.ArrivalAirport = req.ArrivalAirport
	}
	if req.ScheduledDepart != nil {
		schedule.ScheduledDepart = *req.ScheduledDepart
	}
	if req.ScheduledArrive != nil {
		schedule.ScheduledArrive = *req.ScheduledArrive
	}
	if req.ActualDepart != nil {
		schedule.ActualDepart = req.ActualDepart
	}
	if req.ActualArrive != nil {
		schedule.ActualArrive = req.ActualArrive
	}
	if req.FlightStatus != "" {
		schedule.FlightStatus = req.FlightStatus
	}
	if req.VipLoungeCapacity != nil {
		schedule.VipLoungeCapacity = *req.VipLoungeCapacity
	}
	if req.VipLoungeUsed != nil {
		schedule.VipLoungeUsed = *req.VipLoungeUsed
	}
	if req.Status != "" {
		schedule.Status = req.Status
	}
	if req.ResponsiblePerson != "" {
		schedule.ResponsiblePerson = req.ResponsiblePerson
	}
	if req.BatchNo != "" {
		schedule.BatchNo = req.BatchNo
	}
	if req.Remarks != "" {
		schedule.Remarks = req.Remarks
	}

	if err := s.repo.Update(id, schedule); err != nil {
		return nil, fmt.Errorf("failed to update flight schedule: %w", err)
	}

	s.logAudit(models.OperationUpdate, models.EntityTypeFlight, id, schedule.FlightNo, oldValue, schedule)

	return schedule, nil
}

func (s *FlightScheduleService) Delete(id int64) error {
	schedule, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get flight schedule: %w", err)
	}
	if schedule == nil {
		return fmt.Errorf("flight schedule not found")
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete flight schedule: %w", err)
	}

	s.logAudit(models.OperationDelete, models.EntityTypeFlight, id, schedule.FlightNo, schedule, "")

	return nil
}

func (s *FlightScheduleService) Archive(id int64) (*schemas.FlightSnapshotResponse, error) {
	schedule, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get flight schedule: %w", err)
	}
	if schedule == nil {
		return nil, fmt.Errorf("flight schedule not found")
	}

	snapshot := map[string]interface{}{
		"flight_no":          schedule.FlightNo,
		"departure_airport":  schedule.DepartureAirport,
		"arrival_airport":    schedule.ArrivalAirport,
		"scheduled_depart":   schedule.ScheduledDepart,
		"scheduled_arrive":   schedule.ScheduledArrive,
		"actual_depart":      schedule.ActualDepart,
		"actual_arrive":      schedule.ActualArrive,
		"flight_status":      schedule.FlightStatus,
		"vip_lounge_capacity": schedule.VipLoungeCapacity,
		"vip_lounge_used":    schedule.VipLoungeUsed,
		"archived_at":        time.Now(),
	}

	snapshotData, _ := json.Marshal(snapshot)

	if err := s.repo.Archive(id, string(snapshotData)); err != nil {
		return nil, fmt.Errorf("failed to archive flight schedule: %w", err)
	}

	s.logAudit(models.OperationArchive, models.EntityTypeFlight, id, schedule.FlightNo, schedule.Status, models.FlightRecordStatusArchived)

	return &schemas.FlightSnapshotResponse{
		FlightScheduleID: id,
		SnapshotData:     string(snapshotData),
		ArchivedAt:       time.Now(),
	}, nil
}

func (s *FlightScheduleService) BatchArchive(req *schemas.FlightArchiveRequest) (int, error) {
	count, err := s.repo.BatchArchive(req.FlightScheduleIDs, req.CreateSnapshot)
	if err != nil {
		return 0, fmt.Errorf("failed to batch archive flight schedules: %w", err)
	}

	s.logAudit(models.OperationArchive, models.EntityTypeFlight, 0, "batch", "", fmt.Sprintf("archived %d records", count))

	return count, nil
}

func (s *FlightScheduleService) logAudit(op, entityType string, entityID int64, entityNo string, oldValue, newValue interface{}) {
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
