package models

import "time"

type MemberBenefit struct {
	ID              int64     `json:"id"`
	BenefitNo       string    `json:"benefit_no"`
	MemberName      string    `json:"member_name"`
	MemberLevel     string    `json:"member_level"`
	RemainingQuota  int       `json:"remaining_quota"`
	TotalQuota      int       `json:"total_quota"`
	Status          string    `json:"status"`
	ResponsiblePerson string  `json:"responsible_person"`
	ValidFrom       time.Time `json:"valid_from"`
	ValidTo         time.Time `json:"valid_to"`
	BatchNo         string    `json:"batch_no"`
	Remarks         string    `json:"remarks"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

const (
	BenefitStatusActive    = "active"
	BenefitStatusExpired   = "expired"
	BenefitStatusSuspended = "suspended"
	BenefitStatusDraft     = "draft"
)

const (
	MemberLevelPlatinum = "platinum"
	MemberLevelGold     = "gold"
	MemberLevelSilver   = "silver"
	MemberLevelRegular  = "regular"
)
