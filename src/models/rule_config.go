package models

import "time"

type RuleConfig struct {
	ID                   int64     `json:"id"`
	RuleNo               string    `json:"rule_no"`
	RuleName             string    `json:"rule_name"`
	RuleType             string    `json:"rule_type"`
	RuleValue            string    `json:"rule_value"`
	ThresholdValue       float64   `json:"threshold_value"`
	AppliesToLevel       string    `json:"applies_to_level"`
	IsActive             bool      `json:"is_active"`
	EffectiveDate        time.Time `json:"effective_date"`
	ExpiryDate           *time.Time `json:"expiry_date"`
	Status               string    `json:"status"`
	ResponsiblePerson    string    `json:"responsible_person"`
	BatchNo              string    `json:"batch_no"`
	Remarks              string    `json:"remarks"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

const (
	RuleTypeQuotaLimit      = "quota_limit"
	RuleTypeCompanionLimit  = "companion_limit"
	RuleTypeWaitlistTimeout = "waitlist_timeout"
	RuleTypePriorityWeight  = "priority_weight"
	RuleTypeExpiryDays      = "expiry_days"
)

const (
	RuleStatusActive   = "active"
	RuleStatusInactive = "inactive"
	RuleStatusDraft    = "draft"
)
