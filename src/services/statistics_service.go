package services

import (
	"airport-vip-service/src/models"
	"airport-vip-service/src/repositories"
	"airport-vip-service/src/schemas"
	"fmt"
	"time"
)

type StatisticsService struct {
	verificationRepo *repositories.VerificationResultRepository
	waitlistRepo     *repositories.WaitlistEntryRepository
	reservationRepo  *repositories.ReservationRecordRepository
}

func NewStatisticsServiceWithRepos(
	verificationRepo *repositories.VerificationResultRepository,
	waitlistRepo *repositories.WaitlistEntryRepository,
	reservationRepo *repositories.ReservationRecordRepository,
) *StatisticsService {
	return &StatisticsService{
		verificationRepo: verificationRepo,
		waitlistRepo:     waitlistRepo,
		reservationRepo:  reservationRepo,
	}
}

func (s *StatisticsService) GetVerificationStatistics(req *schemas.StatisticsRequest) (*schemas.VerificationStatisticsResponse, error) {
	start, end := parseDateRange(req.StartDate, req.EndDate)
	stats, err := s.verificationRepo.GetStatistics(start, end)
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

func (s *StatisticsService) GetWaitlistStatistics(req *schemas.StatisticsRequest) (*schemas.WaitlistStatisticsResponse, error) {
	entries, _, err := s.waitlistRepo.List(0, 1000000, 0, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get waitlist statistics: %w", err)
	}

	var seatedCount, cancelledCount, waitingCount int64
	var totalWaitMins int64
	for _, entry := range entries {
		switch entry.Status {
		case models.WaitlistStatusSeated:
			seatedCount++
		case models.WaitlistStatusCancelled:
			cancelledCount++
		case models.WaitlistStatusWaiting:
			waitingCount++
		}
		totalWaitMins += int64(entry.EstimatedWaitMins)
	}

	totalEntries := int64(len(entries))
	var averageWaitMins int64
	var transferRate, cancellationRate float64
	if totalEntries > 0 {
		averageWaitMins = totalWaitMins / totalEntries
		transferRate = float64(seatedCount) / float64(totalEntries) * 100
		cancellationRate = float64(cancelledCount) / float64(totalEntries) * 100
	}

	return &schemas.WaitlistStatisticsResponse{
		TotalEntries:     totalEntries,
		SeatedCount:      seatedCount,
		CancelledCount:   cancelledCount,
		WaitingCount:     waitingCount,
		AverageWaitMins:  averageWaitMins,
		TransferRate:     transferRate,
		CancellationRate: cancellationRate,
		ByDate:           []schemas.DateStatistics{},
		ByBatch:          []schemas.BatchStatistics{},
		ByRole:           []schemas.RoleStatistics{},
	}, nil
}

func (s *StatisticsService) GetUsageStatistics(req *schemas.StatisticsRequest) (*schemas.UsageStatisticsResponse, error) {
	start, end := parseDateRange(req.StartDate, req.EndDate)
	records, err := s.reservationRepo.GetByDateRange(start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get usage statistics: %w", err)
	}

	var completedCount, cancelledCount int64
	var totalGuests int64
	members := make(map[string]struct{})
	for _, record := range records {
		switch record.Status {
		case models.ReservationStatusCompleted:
			completedCount++
		case models.ReservationStatusCancelled:
			cancelledCount++
		}
		totalGuests += int64(record.GuestCount)
		if record.MemberName != "" {
			members[record.MemberName] = struct{}{}
		}
	}

	totalReservations := int64(len(records))
	var averageGuestCount, completionRate, cancellationRate float64
	if totalReservations > 0 {
		averageGuestCount = float64(totalGuests) / float64(totalReservations)
		completionRate = float64(completedCount) / float64(totalReservations) * 100
		cancellationRate = float64(cancelledCount) / float64(totalReservations) * 100
	}

	return &schemas.UsageStatisticsResponse{
		TotalReservations: totalReservations,
		CompletedCount:    completedCount,
		CancelledCount:    cancelledCount,
		UniqueMembers:     int64(len(members)),
		AverageGuestCount: averageGuestCount,
		CompletionRate:    completionRate,
		CancellationRate:  cancellationRate,
		ByDate:            []schemas.DateStatistics{},
		ByBatch:           []schemas.BatchStatistics{},
		ByRole:            []schemas.RoleStatistics{},
	}, nil
}

func parseDateRange(startDate, endDate string) (time.Time, time.Time) {
	var start, end time.Time
	var err error

	if startDate != "" {
		start, err = time.Parse("2006-01-02", startDate)
		if err != nil {
			start = time.Now().AddDate(0, -1, 0)
		}
	} else {
		start = time.Now().AddDate(0, -1, 0)
	}

	if endDate != "" {
		end, err = time.Parse("2006-01-02", endDate)
		if err != nil {
			end = time.Now()
		}
	} else {
		end = time.Now()
	}

	return start, end
}
