package schemas

import "time"

type RuleConfigCreateRequest struct {
	RuleName          string    `json:"rule_name" validate:"required"`
	RuleType          string    `json:"rule_type" validate:"required"`
	RuleValue         string    `json:"rule_value" validate:"required"`
	ThresholdValue    float64   `json:"threshold_value"`
	AppliesToLevel    string    `json:"applies_to_level"`
	IsActive          bool      `json:"is_active"`
	EffectiveDate     time.Time `json:"effective_date"`
	ResponsiblePerson string    `json:"responsible_person"`
	BatchNo           string    `json:"batch_no"`
	Remarks           string    `json:"remarks"`
}

type RuleConfigUpdateRequest struct {
	RuleName          string    `json:"rule_name"`
	RuleValue         string    `json:"rule_value"`
	ThresholdValue    *float64  `json:"threshold_value"`
	AppliesToLevel    string    `json:"applies_to_level"`
	IsActive          *bool     `json:"is_active"`
	EffectiveDate     time.Time `json:"effective_date"`
	Status            string    `json:"status"`
	ResponsiblePerson string    `json:"responsible_person"`
	Remarks           string    `json:"remarks"`
}

type RuleConfigResponse struct {
	ID                int64     `json:"id"`
	RuleNo            string    `json:"rule_no"`
	RuleName          string    `json:"rule_name"`
	RuleType          string    `json:"rule_type"`
	RuleValue         string    `json:"rule_value"`
	ThresholdValue    float64   `json:"threshold_value"`
	AppliesToLevel    string    `json:"applies_to_level"`
	IsActive          bool      `json:"is_active"`
	EffectiveDate     time.Time `json:"effective_date"`
	Status            string    `json:"status"`
	ResponsiblePerson string    `json:"responsible_person"`
	BatchNo           string    `json:"batch_no"`
	Remarks           string    `json:"remarks"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type RuleConfigListResponse struct {
	Total int                   `json:"total"`
	Page  int                   `json:"page"`
	Size  int                   `json:"size"`
	Data  []*RuleConfigResponse `json:"data"`
}
