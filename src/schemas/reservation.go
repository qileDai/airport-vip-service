package schemas

import "time"

type ReservationRecordCreateRequest struct {
	MemberBenefitID   int64     `json:"member_benefit_id" validate:"required"`
	MemberName        string    `json:"member_name" validate:"required"`
	FlightNo          string    `json:"flight_no" validate:"required"`
	FlightScheduleID  int64     `json:"flight_schedule_id"`
	VipLoungeName     string    `json:"vip_lounge_name"`
	ReservationTime   time.Time `json:"reservation_time"`
	GuestCount        int       `json:"guest_count" validate:"min=0"`
	ResponsiblePerson string    `json:"responsible_person"`
	BatchNo           string    `json:"batch_no"`
	Remarks           string    `json:"remarks"`
}

type ReservationRecordUpdateRequest struct {
	MemberName        string    `json:"member_name"`
	FlightNo          string    `json:"flight_no"`
	FlightScheduleID  *int64    `json:"flight_schedule_id"`
	VipLoungeName     string    `json:"vip_lounge_name"`
	ReservationTime   time.Time `json:"reservation_time"`
	GuestCount        *int      `json:"guest_count"`
	Status            string    `json:"status"`
	ResponsiblePerson string    `json:"responsible_person"`
	Remarks           string    `json:"remarks"`
}

type ReservationRecordResponse struct {
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
	RejectionReason   string    `json:"rejection_reason"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type ReservationRecordListResponse struct {
	Total int                         `json:"total"`
	Page  int                         `json:"page"`
	Size  int                         `json:"size"`
	Data  []*ReservationRecordResponse `json:"data"`
}

type BatchImportRequest struct {
	Reservations      []*ReservationRecordCreateRequest `json:"reservations" validate:"required,dive"`
	BatchNo           string                            `json:"batch_no" validate:"required"`
	Operator          string                            `json:"operator"`
}

type BatchImportPreviewResponse struct {
	TotalCount      int                          `json:"total_count"`
	ValidCount      int                          `json:"valid_count"`
	InvalidCount    int                          `json:"invalid_count"`
	ValidationErrors []BatchImportValidationError `json:"validation_errors"`
}

type BatchImportValidationError struct {
	RowIndex int    `json:"row_index"`
	Field    string `json:"field"`
	Message  string `json:"message"`
}

type BatchImportResponse struct {
	TotalCount   int                           `json:"total_count"`
	SuccessCount int                           `json:"success_count"`
	FailedCount  int                           `json:"failed_count"`
	Errors       []BatchImportValidationError  `json:"errors"`
}

type StatusChangeRequest struct {
	Action           string `json:"action" validate:"required,oneof=submit approve reject supplement cancel complete archive"`
	Reason           string `json:"reason"`
	Remarks          string `json:"remarks"`
	Operator         string `json:"operator" validate:"required"`
	EntityType       string `json:"-"`
	EntityID         int64  `json:"-"`
}

type AllowedTransitionsResponse struct {
	CurrentStatus string   `json:"current_status"`
	Transitions   []string `json:"transitions"`
}
