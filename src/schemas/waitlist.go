package schemas

import "time"

type WaitlistEntryCreateRequest struct {
	MemberBenefitID   int64     `json:"member_benefit_id" validate:"required"`
	FlightScheduleID  int64     `json:"flight_schedule_id" validate:"required"`
	MemberName        string    `json:"member_name" validate:"required"`
	MemberLevel       string    `json:"member_level" validate:"oneof=regular silver gold platinum"`
	WaitingSince      time.Time `json:"waiting_since"`
	ResponsiblePerson string    `json:"responsible_person"`
	BatchNo           string    `json:"batch_no"`
	Remarks           string    `json:"remarks"`
}

type WaitlistEntryUpdateRequest struct {
	PriorityScore    *int   `json:"priority_score"`
	EstimatedWaitMins *int  `json:"estimated_wait_mins"`
	Status           string `json:"status"`
	Remarks          string `json:"remarks"`
}

type WaitlistEntryResponse struct {
	ID                int64     `json:"id"`
	WaitlistNo        string    `json:"waitlist_no"`
	MemberBenefitID   int64     `json:"member_benefit_id"`
	FlightScheduleID  int64     `json:"flight_schedule_id"`
	MemberName        string    `json:"member_name"`
	MemberLevel       string    `json:"member_level"`
	WaitingSince      time.Time `json:"waiting_since"`
	PriorityScore     int       `json:"priority_score"`
	EstimatedWaitMins int       `json:"estimated_wait_mins"`
	Status            string    `json:"status"`
	ResponsiblePerson string    `json:"responsible_person"`
	BatchNo           string    `json:"batch_no"`
	Remarks           string    `json:"remarks"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type WaitlistEntryListResponse struct {
	Total int                      `json:"total"`
	Page  int                      `json:"page"`
	Size  int                      `json:"size"`
	Data  []*WaitlistEntryResponse `json:"data"`
}

type PriorityCalculationRequest struct {
	MemberLevel    string `json:"member_level" validate:"required,oneof=regular silver gold platinum"`
	RemainingQuota int    `json:"remaining_quota" validate:"min=0"`
	WaitingMinutes int    `json:"waiting_minutes" validate:"min=0"`
}

type PriorityCalculationResponse struct {
	MemberLevel    string `json:"member_level"`
	RemainingQuota int    `json:"remaining_quota"`
	WaitingMinutes int    `json:"waiting_minutes"`
	PriorityScore  int    `json:"priority_score"`
	LevelScore     int    `json:"level_score"`
	QuotaScore     int    `json:"quota_score"`
	WaitScore      int    `json:"wait_score"`
}

type WaitlistArrangementRequest struct {
	FlightScheduleID int64  `json:"flight_schedule_id" validate:"required"`
	AvailableSeats   int    `json:"available_seats" validate:"required,min=1"`
	Operator         string `json:"operator"`
}

type WaitlistArrangementResponse struct {
	FlightScheduleID int64                     `json:"flight_schedule_id"`
	AvailableSeats   int                       `json:"available_seats"`
	ArrangedCount    int                       `json:"arranged_count"`
	ArrangedEntries  []*WaitlistEntryResponse  `json:"arranged_entries"`
}
