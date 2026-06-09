package models

import "time"

type ExceptionEvent struct {
	ID                  int64     `json:"id"`
	EventNo             string    `json:"event_no"`
	EventType           string    `json:"event_type"`
	EntityType          string    `json:"entity_type"`
	EntityID            int64     `json:"entity_id"`
	TriggerField        string    `json:"trigger_field"`
	ThresholdValue      string    `json:"threshold_value"`
	ActualValue         string    `json:"actual_value"`
	Severity            string    `json:"severity"`
	Handler             string    `json:"handler"`
	HandlingDeadline    *time.Time `json:"handling_deadline"`
	HandledAt           *time.Time `json:"handled_at"`
	HandlingResult      string    `json:"handling_result"`
	Status              string    `json:"status"`
	ResponsiblePerson   string    `json:"responsible_person"`
	BatchNo             string    `json:"batch_no"`
	Remarks             string    `json:"remarks"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

const (
	EventTypeBenefitExpired     = "benefit_expired"
	EventTypeQuotaExceeded      = "quota_exceeded"
	EventTypeCompanionViolation = "companion_violation"
	EventTypeWaitlistTimeout    = "waitlist_timeout"
	EventTypeVoucherInvalid     = "voucher_invalid"
	EventTypeFlightDelay        = "flight_delay"
)

const (
	ExceptionSeverityLow      = "low"
	ExceptionSeverityMedium   = "medium"
	ExceptionSeverityHigh     = "high"
	ExceptionSeverityCritical = "critical"
)

const (
	ExceptionStatusOpen      = "open"
	ExceptionStatusHandling  = "handling"
	ExceptionStatusResolved  = "resolved"
	ExceptionStatusClosed    = "closed"
	ExceptionStatusArchived  = "archived"
)
