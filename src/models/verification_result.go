package models

import "time"

type VerificationResult struct {
	ID                  int64     `json:"id"`
	VerificationNo      string    `json:"verification_no"`
	ReservationID       int64     `json:"reservation_id"`
	MemberBenefitID     int64     `json:"member_benefit_id"`
	FlightScheduleID    int64     `json:"flight_schedule_id"`
	VerificationType    string    `json:"verification_type"`
	Result              string    `json:"result"`
	FailureReason       string    `json:"failure_reason"`
	VerifiedQuota       int       `json:"verified_quota"`
	VerifiedCompanions  int       `json:"verified_companions"`
	VerificationDetails string    `json:"verification_details"`
	Status              string    `json:"status"`
	ResponsiblePerson   string    `json:"responsible_person"`
	BatchNo             string    `json:"batch_no"`
	Remarks             string    `json:"remarks"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

const (
	VerifyTypeReservation = "reservation"
	VerifyTypeCheckin     = "checkin"
	VerifyTypeCompanion   = "companion"
	VerifyTypeVoucher     = "voucher"
)

const (
	VerifyResultPassed = "passed"
	VerifyResultFailed = "failed"
	VerifyResultPending = "pending"
)

const (
	VerifyStatusDraft     = "draft"
	VerifyStatusPending   = "pending_review"
	VerifyStatusConfirmed = "confirmed"
	VerifyStatusArchived  = "archived"
	VerifyStatusRejected  = "rejected"
)
