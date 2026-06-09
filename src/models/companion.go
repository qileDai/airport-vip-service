package models

import "time"

type Companion struct {
	ID                int64     `json:"id"`
	CompanionNo       string    `json:"companion_no"`
	ReservationID     int64     `json:"reservation_id"`
	CompanionName     string    `json:"companion_name"`
	CompanionIDCard   string    `json:"companion_id_card"`
	RelationType      string    `json:"relation_type"`
	IsVipEligible     bool      `json:"is_vip_eligible"`
	VerificationStatus string   `json:"verification_status"`
	Status            string    `json:"status"`
	ResponsiblePerson string    `json:"responsible_person"`
	BatchNo           string    `json:"batch_no"`
	Remarks           string    `json:"remarks"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

const (
	RelationSpouse    = "spouse"
	RelationChild     = "child"
	RelationParent    = "parent"
	RelationColleague = "colleague"
	RelationFriend    = "friend"
	RelationOther     = "other"
)

const (
	CompanionStatusActive   = "active"
	CompanionStatusInactive = "inactive"
	CompanionStatusPending  = "pending"
)

const (
	CompanionVerifyPassed  = "passed"
	CompanionVerifyFailed  = "failed"
	CompanionVerifyPending = "pending"
)
