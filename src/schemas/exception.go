package schemas

import "time"

type ExceptionEventCreateRequest struct {
	EventNo           string     `json:"event_no"`
	EventType         string     `json:"event_type" validate:"required"`
	EntityType        string     `json:"entity_type" validate:"required"`
	EntityID          int64      `json:"entity_id" validate:"required"`
	TriggerField      string     `json:"trigger_field" validate:"required"`
	ThresholdValue    string     `json:"threshold_value"`
	ActualValue       string     `json:"actual_value"`
	Severity          string     `json:"severity" validate:"oneof=low medium high critical"`
	Handler           string     `json:"handler"`
	HandlingDeadline  *time.Time `json:"handling_deadline"`
	HandledAt         *time.Time `json:"handled_at"`
	HandlingResult    string     `json:"handling_result"`
	Status            string     `json:"status"`
	ResponsiblePerson string     `json:"responsible_person"`
	BatchNo           string     `json:"batch_no"`
	Remarks           string     `json:"remarks"`
}

type ExceptionEventUpdateRequest struct {
	Handler           string     `json:"handler"`
	HandlingDeadline  *time.Time `json:"handling_deadline"`
	HandledAt         *time.Time `json:"handled_at"`
	HandlingResult    string     `json:"handling_result"`
	Status            string     `json:"status"`
	ResponsiblePerson string     `json:"responsible_person"`
	Remarks           string     `json:"remarks"`
}

type ExceptionEventResponse struct {
	ID                int64      `json:"id"`
	EventNo           string     `json:"event_no"`
	EventType         string     `json:"event_type"`
	EntityType        string     `json:"entity_type"`
	EntityID          int64      `json:"entity_id"`
	TriggerField      string     `json:"trigger_field"`
	ThresholdValue    string     `json:"threshold_value"`
	ActualValue       string     `json:"actual_value"`
	Severity          string     `json:"severity"`
	Handler           string     `json:"handler"`
	HandlingDeadline  *time.Time `json:"handling_deadline"`
	HandledAt         *time.Time `json:"handled_at"`
	HandlingResult    string     `json:"handling_result"`
	Status            string     `json:"status"`
	ResponsiblePerson string     `json:"responsible_person"`
	BatchNo           string     `json:"batch_no"`
	Remarks           string     `json:"remarks"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type ExceptionEventListResponse struct {
	Total   int64                    `json:"total"`
	Page    int                      `json:"page"`
	PerPage int                      `json:"per_page"`
	Data    []ExceptionEventResponse `json:"data"`
}

type ExceptionHandleRequest struct {
	EventID        int64      `json:"event_id" validate:"required"`
	EventNo        string     `json:"event_no"`
	Handler        string     `json:"handler" validate:"required"`
	HandledAt      *time.Time `json:"handled_at"`
	HandlingResult string     `json:"handling_result" validate:"required"`
	Remarks        string     `json:"remarks"`
}
