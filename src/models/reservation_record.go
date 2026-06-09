package models

import "time"

type ReservationRecord struct {
	ID                int64     `json:"id"`
	ReservationNo     string    `json:"reservation_no"`
	MemberBenefitID   int64     `json:"member_benefit_id"`
	MemberName        string    `json:"member_name"`
	FlightNo          string    `json:"flight_no"`
	FlightScheduleID  int64     `json:"flight_schedule_id"`
	VipLoungeName     string    `json:"vip_lounge_name"`
	ReservationTime   time.Time `json:"reservation_time"`
	GuestCount        int       `json:"guest_count"`
	Status            string    `json:"status"`
	ResponsiblePerson string    `json:"responsible_person"`
	BatchNo           string    `json:"batch_no"`
	Remarks           string    `json:"remarks"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

const (
	ReservationStatusDraft      = "draft"
	ReservationStatusPending    = "pending_review"
	ReservationStatusConfirmed  = "confirmed"
	ReservationStatusCancelled  = "cancelled"
	ReservationStatusCompleted  = "completed"
	ReservationStatusRejected   = "rejected"
	ReservationStatusSupplement = "pending_supplement"
	ReservationStatusArchived   = "archived"
)
