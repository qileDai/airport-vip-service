package models

import "time"

type AuditLog struct {
	ID              int64     `json:"id"`
	LogNo           string    `json:"log_no"`
	OperationType   string    `json:"operation_type"`
	EntityType      string    `json:"entity_type"`
	EntityID        int64     `json:"entity_id"`
	EntityNo        string    `json:"entity_no"`
	OldValue        string    `json:"old_value"`
	NewValue        string    `json:"new_value"`
	Operator        string    `json:"operator"`
	OperatorRole    string    `json:"operator_role"`
	IPAddress       string    `json:"ip_address"`
	UserAgent       string    `json:"user_agent"`
	RequestID       string    `json:"request_id"`
	Status          string    `json:"status"`
	ResponsiblePerson string    `json:"responsible_person"`
	BatchNo         string    `json:"batch_no"`
	Remarks         string    `json:"remarks"`
	CreatedAt       time.Time `json:"created_at"`
}

const (
	OperationCreate    = "create"
	OperationUpdate    = "update"
	OperationDelete    = "delete"
	OperationStatusChange = "status_change"
	OperationBatchImport = "batch_import"
	OperationVerify    = "verify"
	OperationArchive   = "archive"
)

const (
	AuditStatusActive   = "active"
	AuditStatusArchived = "archived"
)
