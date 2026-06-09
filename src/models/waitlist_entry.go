package models

import "time"

type WaitlistEntry struct {
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

const (
	WaitlistStatusWaiting   = "waiting"
	WaitlistStatusNotified  = "notified"
	WaitlistStatusSeated    = "seated"
	WaitlistStatusCancelled = "cancelled"
	WaitlistStatusExpired   = "expired"
)
