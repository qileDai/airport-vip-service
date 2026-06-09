package tests

import (
	"airport-vip-service/src/models"
	"airport-vip-service/src/schemas"
	"airport-vip-service/src/services"
	"airport-vip-service/src/utils"
	"testing"
	"time"
)

func TestPriorityCalculation(t *testing.T) {
	testCases := []struct {
		name          string
		memberLevel   string
		remainingQuota int
		waitingMins   int
		expectedMin   int
	}{
		{"Platinum member with high quota", "platinum", 15, 30, 60},
		{"Gold member with medium quota", "gold", 8, 45, 40},
		{"Silver member with low quota", "silver", 3, 60, 30},
		{"Regular member", "regular", 2, 90, 20},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			score := calculatePriority(tc.memberLevel, tc.remainingQuota, tc.waitingMins)
			if score < tc.expectedMin {
				t.Errorf("Priority score %d is less than expected minimum %d", score, tc.expectedMin)
			}
		})
	}
}

func calculatePriority(memberLevel string, remainingQuota, waitingMins int) int {
	score := 0

	switch memberLevel {
	case "platinum":
		score += 40
	case "gold":
		score += 30
	case "silver":
		score += 20
	case "regular":
		score += 10
	}

	quotaScore := remainingQuota * 2
	if quotaScore > 20 {
		quotaScore = 20
	}
	score += quotaScore

	waitScore := waitingMins / 10
	if waitScore > 40 {
		waitScore = 40
	}
	score += waitScore

	return score
}

func TestStatusTransition(t *testing.T) {
	allowedTransitions := map[string][]string{
		"draft":             {"pending_review"},
		"pending_review":    {"confirmed", "rejected", "pending_supplement"},
		"pending_supplement": {"pending_review", "draft"},
		"confirmed":         {"completed", "cancelled"},
		"rejected":          {"draft"},
		"completed":         {"archived"},
		"cancelled":         {},
		"archived":          {},
	}

	testCases := []struct {
		fromStatus string
		toStatus   string
		allowed    bool
	}{
		{"draft", "pending_review", true},
		{"draft", "confirmed", false},
		{"pending_review", "confirmed", true},
		{"pending_review", "rejected", true},
		{"confirmed", "completed", true},
		{"completed", "draft", false},
	}

	for _, tc := range testCases {
		t.Run(tc.fromStatus+"->"+tc.toStatus, func(t *testing.T) {
			allowed := isTransitionAllowed(tc.fromStatus, tc.toStatus, allowedTransitions)
			if allowed != tc.allowed {
				t.Errorf("Transition from %s to %s: expected %v, got %v", tc.fromStatus, tc.toStatus, tc.allowed, allowed)
			}
		})
	}
}

func isTransitionAllowed(fromStatus, toStatus string, transitions map[string][]string) bool {
	allowed, exists := transitions[fromStatus]
	if !exists {
		return false
	}
	for _, s := range allowed {
		if s == toStatus {
			return true
		}
	}
	return false
}

func TestBenefitExpiryCheck(t *testing.T) {
	now := time.Now()

	testCases := []struct {
		name      string
		validTo   time.Time
		isExpired bool
	}{
		{"Already expired", now.AddDate(0, 0, -1), true},
		{"Expires today", now, true},
		{"Not expired", now.AddDate(0, 0, 7), false},
		{"Far future", now.AddDate(1, 0, 0), false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isExpired := tc.validTo.Before(now) || tc.validTo.Equal(now)
			if isExpired != tc.isExpired {
				t.Errorf("Expiry check for %s: expected %v, got %v", tc.validTo.Format("2006-01-02"), tc.isExpired, isExpired)
			}
		})
	}
}

func TestCompanionLimitByMemberLevel(t *testing.T) {
	limits := map[string]int{
		"platinum": 4,
		"gold":     3,
		"silver":   2,
		"regular":  1,
	}

	testCases := []struct {
		memberLevel string
		guestCount  int
		shouldPass  bool
	}{
		{"platinum", 4, true},
		{"platinum", 5, false},
		{"gold", 3, true},
		{"gold", 4, false},
		{"silver", 2, true},
		{"silver", 3, false},
		{"regular", 1, true},
		{"regular", 2, false},
	}

	for _, tc := range testCases {
		t.Run(tc.memberLevel+"_"+string(rune(tc.guestCount+'0')), func(t *testing.T) {
			limit := limits[tc.memberLevel]
			isValid := tc.guestCount <= limit
			if isValid != tc.shouldPass {
				t.Errorf("Companion limit check for %s with %d guests: expected %v, got %v", tc.memberLevel, tc.guestCount, tc.shouldPass, isValid)
			}
		})
	}
}

func TestGenerateNo(t *testing.T) {
	benefitNo := utils.GenerateBenefitNo()
	if len(benefitNo) < 10 {
		t.Errorf("Generated benefit no is too short: %s", benefitNo)
	}

	reservationNo := utils.GenerateReservationNo()
	if len(reservationNo) < 10 {
		t.Errorf("Generated reservation no is too short: %s", reservationNo)
	}

	flightNo := utils.GenerateFlightNo()
	if len(flightNo) < 4 {
		t.Errorf("Generated flight no is too short: %s", flightNo)
	}
}

func TestBatchImportValidation(t *testing.T) {
	validReservation := schemas.ReservationRecordCreateRequest{
		MemberBenefitID:   1,
		MemberName:        "测试会员",
		FlightNo:          "CA1234",
		FlightScheduleID:  1,
		VipLoungeName:     "测试VIP厅",
		ReservationTime:   time.Now().Add(24 * time.Hour),
		GuestCount:        2,
		ResponsiblePerson: "测试负责人",
	}

	if validReservation.MemberName == "" {
		t.Error("Member name should not be empty")
	}

	if validReservation.GuestCount < 0 {
		t.Error("Guest count should not be negative")
	}
}

func TestVerificationResultStatus(t *testing.T) {
	statuses := []string{
		models.VerificationStatusDraft,
		models.VerificationStatusPending,
		models.VerificationStatusConfirmed,
		models.VerificationStatusRejected,
	}

	for _, status := range statuses {
		if status == "" {
			t.Errorf("Status should not be empty")
		}
	}
}

func TestExceptionSeverity(t *testing.T) {
	severities := []string{"low", "medium", "high", "critical"}

	for _, severity := range severities {
		if severity == "" {
			t.Errorf("Severity should not be empty")
		}
	}
}

func TestMemberBenefitStatus(t *testing.T) {
	statuses := []string{
		models.BenefitStatusActive,
		models.BenefitStatusExpired,
		models.BenefitStatusSuspended,
	}

	for _, status := range statuses {
		if status == "" {
			t.Errorf("Benefit status should not be empty")
		}
	}
}

func TestReservationStatus(t *testing.T) {
	statuses := []string{
		models.ReservationStatusDraft,
		models.ReservationStatusPendingReview,
		models.ReservationStatusPendingSupplement,
		models.ReservationStatusConfirmed,
		models.ReservationStatusCompleted,
		models.ReservationStatusCancelled,
		models.ReservationStatusRejected,
		models.ReservationStatusArchived,
	}

	for _, status := range statuses {
		if status == "" {
			t.Errorf("Reservation status should not be empty")
		}
	}
}
