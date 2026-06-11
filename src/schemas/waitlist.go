package schemas

import "time"

type WaitlistEntryCreateRequest struct {
	WaitlistNo        string    `json:"waitlist_no"`
	ReservationID     int64     `json:"reservation_id"`
	MemberBenefitID   int64     `json:"member_benefit_id" validate:"required"`
	FlightScheduleID  int64     `json:"flight_schedule_id" validate:"required"`
	MemberName        string    `json:"member_name" validate:"required"`
	MemberLevel       string    `json:"member_level" validate:"oneof=regular silver gold platinum"`
	WaitingSince      time.Time `json:"waiting_since"`
	PriorityScore     int       `json:"priority_score"`
	EstimatedWaitMins int       `json:"estimated_wait_mins"`
	Status            string    `json:"status"`
	ResponsiblePerson string    `json:"responsible_person"`
	BatchNo           string    `json:"batch_no"`
	Remarks           string    `json:"remarks"`
}

type WaitlistEntryUpdateRequest struct {
	PriorityScore      *int   `json:"priority_score"`
	EstimatedWaitMins  *int   `json:"estimated_wait_mins"`
	Status             string `json:"status"`
	Remarks            string `json:"remarks"`
}

type WaitlistEntryResponse struct {
	ID                int64     `json:"id"`
	WaitlistNo        string    `json:"waitlist_no"`
	ReservationID     int64     `json:"reservation_id"`
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
	Total   int64                   `json:"total"`
	Page    int                     `json:"page"`
	PerPage int                     `json:"per_page"`
	Data    []WaitlistEntryResponse `json:"data"`
}

type PriorityCalculationRequest struct {
	MemberLevel      string    `json:"member_level" validate:"required,oneof=regular silver gold platinum"`
	RemainingQuota   int       `json:"remaining_quota" validate:"min=0"`
	WaitingMinutes   int       `json:"waiting_minutes" validate:"min=0"`
	WaitingSince     time.Time `json:"waiting_since"`
	FlightScheduleID int64     `json:"flight_schedule_id"`
}

type PriorityCalculationResponse struct {
	MemberLevel       string `json:"member_level,omitempty"`
	RemainingQuota    int    `json:"remaining_quota,omitempty"`
	WaitingMinutes    int    `json:"waiting_minutes,omitempty"`
	PriorityScore     int    `json:"priority_score"`
	LevelScore        int    `json:"level_score,omitempty"`
	QuotaScore        int    `json:"quota_score,omitempty"`
	WaitScore         int    `json:"wait_score,omitempty"`
	LevelWeight       int    `json:"level_weight"`
	QuotaWeight       int    `json:"quota_weight"`
	TimeWeight        int    `json:"time_weight"`
	FlightWeight      int    `json:"flight_weight"`
	EstimatedWaitMins int    `json:"estimated_wait_mins"`
	Position          int    `json:"position"`
	TotalWaiting      int    `json:"total_waiting"`
}

type WaitlistArrangementRequest struct {
	FlightScheduleID int64  `json:"flight_schedule_id" validate:"required"`
	AvailableSeats   int    `json:"available_seats" validate:"required,min=1"`
	Operator         string `json:"operator"`
}

type WaitlistArrangementResponse struct {
	FlightScheduleID int64                   `json:"flight_schedule_id,omitempty"`
	AvailableSeats   int                     `json:"available_seats,omitempty"`
	ArrangedCount    int                     `json:"arranged_count"`
	Remaining        int                     `json:"remaining"`
	Arranged         []WaitlistEntryResponse `json:"arranged"`
	NotArranged      []WaitlistEntryResponse `json:"not_arranged"`
	ArrangedEntries  []*WaitlistEntryResponse `json:"arranged_entries,omitempty"`
}
