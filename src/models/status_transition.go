package models

import "time"

type StatusTransition struct {
	ID                int64     `json:"id"`
	TransitionNo      string    `json:"transition_no"`
	EntityType        string    `json:"entity_type"`
	EntityID          int64     `json:"entity_id"`
	FromStatus        string    `json:"from_status"`
	ToStatus          string    `json:"to_status"`
	Action            string    `json:"action"`
	Reason            string    `json:"reason"`
	Operator          string    `json:"operator"`
	Status            string    `json:"status"`
	ResponsiblePerson string    `json:"responsible_person"`
	BatchNo           string    `json:"batch_no"`
	Remarks           string    `json:"remarks"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

const (
	EntityTypeBenefit     = "member_benefit"
	EntityTypeReservation = "reservation"
	EntityTypeFlight      = "flight_schedule"
	EntityTypeCompanion   = "companion"
	EntityTypeVoucher     = "usage_voucher"
	EntityTypeWaitlist    = "waitlist"
	EntityTypeVerify      = "verification"
)

const (
	TransitionStatusActive   = "active"
	TransitionStatusArchived = "archived"
)
