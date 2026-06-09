package schemas

import "time"

type StatusTransitionCreateRequest struct {
	EntityType        string `json:"entity_type" validate:"required"`
	EntityID          int64  `json:"entity_id" validate:"required"`
	FromStatus        string `json:"from_status" validate:"required"`
	ToStatus          string `json:"to_status" validate:"required"`
	Action            string `json:"action" validate:"required"`
	Reason            string `json:"reason"`
	Operator          string `json:"operator" validate:"required"`
	ResponsiblePerson string `json:"responsible_person"`
	BatchNo           string `json:"batch_no"`
	Remarks           string `json:"remarks"`
}

type StatusTransitionResponse struct {
	ID                int64     `json:"id"`
	TransitionNo      string    `json:"transition_no"`
	EntityType        string    `json:"entity_type"`
	EntityID          int64     `json:"entity_id"`
	FromStatus        string    `json:"from_status"`
	ToStatus          string    `json:"to_status"`
	Action            string    `json:"action"`
	Reason            string    `json:"reason"`
	Operator          string    `json:"operator"`
	Status            string    `json:"status"`
	ResponsiblePerson string    `json:"responsible_person"`
	BatchNo           string    `json:"batch_no"`
	Remarks           string    `json:"remarks"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type StatusTransitionListResponse struct {
	Total int                        `json:"total"`
	Page  int                        `json:"page"`
	Size  int                        `json:"size"`
	Data  []*StatusTransitionResponse `json:"data"`
}
