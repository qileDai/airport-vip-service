package schemas

import "time"

type CompanionCreateRequest struct {
	ReservationID     int64  `json:"reservation_id" validate:"required"`
	CompanionName     string `json:"companion_name" validate:"required"`
	CompanionIdCard   string `json:"companion_id_card"`
	RelationType      string `json:"relation_type" validate:"oneof=spouse child parent colleague friend other"`
	IsVipEligible     bool   `json:"is_vip_eligible"`
	ResponsiblePerson string `json:"responsible_person"`
	BatchNo           string `json:"batch_no"`
	Remarks           string `json:"remarks"`
}

type CompanionUpdateRequest struct {
	CompanionName     string `json:"companion_name"`
	CompanionIdCard   string `json:"companion_id_card"`
	RelationType      string `json:"relation_type" validate:"oneof=spouse child parent colleague friend other"`
	IsVipEligible     *bool  `json:"is_vip_eligible"`
	VerificationStatus string `json:"verification_status"`
	Status            string `json:"status"`
	ResponsiblePerson string `json:"responsible_person"`
	Remarks           string `json:"remarks"`
}

type CompanionResponse struct {
	ID                int64     `json:"id"`
	CompanionNo       string    `json:"companion_no"`
	ReservationID     int64     `json:"reservation_id"`
	CompanionName     string    `json:"companion_name"`
	CompanionIdCard   string    `json:"companion_id_card"`
	RelationType      string    `json:"relation_type"`
	IsVipEligible     bool      `json:"is_vip_eligible"`
	VerificationStatus string   `json:"verification_status"`
	Status            string    `json:"status"`
	ResponsiblePerson string    `json:"responsible_person"`
	BatchNo           string    `json:"batch_no"`
	Remarks           string    `json:"remarks"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type CompanionListResponse struct {
	Total int                  `json:"total"`
	Page  int                  `json:"page"`
	Size  int                  `json:"size"`
	Data  []*CompanionResponse `json:"data"`
}
