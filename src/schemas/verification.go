package schemas

import "time"

type VerificationResultCreateRequest struct {
	VerificationNo     string `json:"verification_no"`
	ReservationID      int64  `json:"reservation_id" validate:"required"`
	MemberBenefitID    int64  `json:"member_benefit_id" validate:"required"`
	FlightScheduleID   int64  `json:"flight_schedule_id"`
	VerificationType   string `json:"verification_type" validate:"required,oneof=reservation checkin companion"`
	Result             string `json:"result" validate:"required,oneof=passed failed pending"`
	FailureReason      string `json:"failure_reason"`
	VerifiedQuota      int    `json:"verified_quota"`
	VerifiedCompanions int    `json:"verified_companions"`
	VerificationDetails string `json:"verification_details"`
	Status             string `json:"status"`
	ResponsiblePerson  string `json:"responsible_person"`
	BatchNo            string `json:"batch_no"`
	Remarks            string `json:"remarks"`
}

type VerificationResultUpdateRequest struct {
	Result              string `json:"result" validate:"oneof=passed failed pending"`
	FailureReason       string `json:"failure_reason"`
	VerifiedQuota       *int   `json:"verified_quota"`
	VerifiedCompanions  *int   `json:"verified_companions"`
	VerificationDetails string `json:"verification_details"`
	Status              string `json:"status"`
	Remarks             string `json:"remarks"`
}

type VerificationResultResponse struct {
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

type VerificationResultListResponse struct {
	Total   int64                        `json:"total"`
	Page    int                          `json:"page"`
	PerPage int                          `json:"per_page"`
	Data    []VerificationResultResponse `json:"data"`
}

type EligibilityVerificationRequest struct {
	ReservationID   int64  `json:"reservation_id" validate:"required"`
	MemberBenefitID int64  `json:"member_benefit_id" validate:"required"`
	FlightScheduleID int64 `json:"flight_schedule_id" validate:"required"`
	GuestCount      int    `json:"guest_count" validate:"min=0"`
	Operator        string `json:"operator"`
}

type EligibilityVerificationResponse struct {
	IsEligible          bool                      `json:"is_eligible"`
	VerificationNo      string                    `json:"verification_no,omitempty"`
	BenefitValid        bool                      `json:"benefit_valid"`
	FlightValid         bool                      `json:"flight_valid"`
	CompanionsValid     bool                      `json:"companions_valid"`
	QuotaValid          bool                      `json:"quota_valid"`
	VerifiedQuota       int                       `json:"verified_quota"`
	VerifiedCompanions  int                       `json:"verified_companions"`
	Errors              []string                  `json:"errors"`
	Warnings            []string                  `json:"warnings"`
	MemberBenefit       *MemberBenefitResponse    `json:"member_benefit,omitempty"`
	FlightSchedule      *FlightScheduleResponse   `json:"flight_schedule,omitempty"`
	VerificationResult  *VerificationResultResponse `json:"verification_result,omitempty"`
	ValidationErrors    []string                  `json:"validation_errors,omitempty"`
}
