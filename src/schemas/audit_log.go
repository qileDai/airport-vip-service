package schemas

import "time"

type AuditLogCreateRequest struct {
	OperationType     string `json:"operation_type" validate:"required"`
	EntityType        string `json:"entity_type" validate:"required"`
	EntityID          int64  `json:"entity_id"`
	EntityNo          string `json:"entity_no"`
	OldValue          string `json:"old_value"`
	NewValue          string `json:"new_value"`
	Operator          string `json:"operator" validate:"required"`
	OperatorRole      string `json:"operator_role"`
	IPAddress         string `json:"ip_address"`
	UserAgent         string `json:"user_agent"`
	RequestID         string `json:"request_id"`
	ResponsiblePerson string `json:"responsible_person"`
	BatchNo           string `json:"batch_no"`
	Remarks           string `json:"remarks"`
}

type AuditLogResponse struct {
	ID                int64     `json:"id"`
	LogNo             string    `json:"log_no"`
	OperationType     string    `json:"operation_type"`
	EntityType        string    `json:"entity_type"`
	EntityID          int64     `json:"entity_id"`
	EntityNo          string    `json:"entity_no"`
	OldValue          string    `json:"old_value"`
	NewValue          string    `json:"new_value"`
	Operator          string    `json:"operator"`
	OperatorRole      string    `json:"operator_role"`
	IPAddress         string    `json:"ip_address"`
	UserAgent         string    `json:"user_agent"`
	RequestID         string    `json:"request_id"`
	Status            string    `json:"status"`
	ResponsiblePerson string    `json:"responsible_person"`
	BatchNo           string    `json:"batch_no"`
	Remarks           string    `json:"remarks"`
	CreatedAt         time.Time `json:"created_at"`
}

type AuditLogListResponse struct {
	Total   int64              `json:"total"`
	Page    int                `json:"page"`
	PerPage int                `json:"per_page"`
	Data    []AuditLogResponse `json:"data"`
}

type AuditLogQueryRequest struct {
	Page          int        `json:"page"`
	PerPage       int        `json:"per_page"`
	EntityType    string     `json:"entity_type"`
	EntityID      int64      `json:"entity_id"`
	OperationType string     `json:"operation_type"`
	Operator      string     `json:"operator"`
	StartDate     string     `json:"start_date"`
	EndDate       string     `json:"end_date"`
	StartTime     *time.Time `json:"start_time"`
	EndTime       *time.Time `json:"end_time"`
}
