package schemas

import "time"

type ReservationRecordCreateRequest struct {
	ReservationNo     string    `json:"reservation_no"`
	MemberBenefitID   int64     `json:"member_benefit_id" validate:"required"`
	MemberName        string    `json:"member_name" validate:"required"`
	FlightNo          string    `json:"flight_no" validate:"required"`
	FlightScheduleID  int64     `json:"flight_schedule_id"`
	VipLoungeName     string    `json:"vip_lounge_name"`
	ReservationTime   time.Time `json:"reservation_time"`
	GuestCount        int       `json:"guest_count" validate:"min=0"`
	Status            string    `json:"status"`
	ResponsiblePerson string    `json:"responsible_person"`
	BatchNo           string    `json:"batch_no"`
	Remarks           string    `json:"remarks"`
}

type ReservationRecordUpdateRequest struct {
	MemberName        string     `json:"member_name"`
	FlightNo          string     `json:"flight_no"`
	FlightScheduleID  *int64     `json:"flight_schedule_id"`
	VipLoungeName     string     `json:"vip_lounge_name"`
	ReservationTime   *time.Time `json:"reservation_time"`
	GuestCount        *int       `json:"guest_count"`
	Status            string     `json:"status"`
	ResponsiblePerson string     `json:"responsible_person"`
	BatchNo           string     `json:"batch_no"`
	Remarks           string     `json:"remarks"`
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
	Total   int64                       `json:"total"`
	Page    int                         `json:"page"`
	PerPage int                         `json:"per_page"`
	Data    []ReservationRecordResponse `json:"data"`
}

type BatchImportRequest struct {
	Records      []ReservationRecordCreateRequest  `json:"records"`
	Reservations []*ReservationRecordCreateRequest `json:"reservations" validate:"required,dive"`
	BatchNo      string                            `json:"batch_no" validate:"required"`
	Operator     string                            `json:"operator"`
}

type BatchImportError struct {
	RowIndex int    `json:"row_index"`
	Field    string `json:"field"`
	Message  string `json:"message"`
}

type BatchImportPreviewResponse struct {
	TotalRecords     int                `json:"total_records"`
	ValidRecords     int                `json:"valid_records"`
	InvalidRecords   int                `json:"invalid_records"`
	ValidationErrors []BatchImportError `json:"validation_errors"`
}

type BatchImportResultResponse struct {
	SuccessCount int                `json:"success_count"`
	FailCount    int                `json:"fail_count"`
	Errors       []BatchImportError `json:"errors"`
	CreatedIDs   []int64            `json:"created_ids"`
}

type StatusChangeRequest struct {
	Action     string `json:"action" validate:"required,oneof=submit approve reject supplement cancel complete archive"`
	ToStatus   string `json:"to_status"`
	Reason     string `json:"reason"`
	Remarks    string `json:"remarks"`
	Operator   string `json:"operator" validate:"required"`
	EntityType string `json:"-"`
	EntityID   int64  `json:"-"`
}

type StatusChangeResponse struct {
	Success        bool                     `json:"success"`
	TransitionNo   string                   `json:"transition_no,omitempty"`
	PreviousStatus string                   `json:"previous_status"`
	NewStatus      string                   `json:"new_status,omitempty"`
	Errors         []string                 `json:"errors,omitempty"`
	Transition     StatusTransitionResponse `json:"transition"`
}

type AllowedTransitionsResponse struct {
	CurrentStatus  string   `json:"current_status"`
	AllowedNext    []string `json:"allowed_next"`
	RequiresReason []string `json:"requires_reason"`
	Transitions    []string `json:"transitions,omitempty"`
}
