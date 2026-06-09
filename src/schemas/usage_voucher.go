package schemas

import "time"

type UsageVoucherCreateRequest struct {
	ReservationID     int64     `json:"reservation_id" validate:"required"`
	MemberBenefitID   int64     `json:"member_benefit_id" validate:"required"`
	VoucherType       string    `json:"voucher_type" validate:"oneof=single_use multi_use guest_pass"`
	ValidFrom         time.Time `json:"valid_from"`
	ValidTo           time.Time `json:"valid_to"`
	ResponsiblePerson string    `json:"responsible_person"`
	BatchNo           string    `json:"batch_no"`
	Remarks           string    `json:"remarks"`
}

type UsageVoucherUpdateRequest struct {
	VerificationStatus string    `json:"verification_status"`
	Status             string    `json:"status"`
	UsedLocation       string    `json:"used_location"`
	ResponsiblePerson  string    `json:"responsible_person"`
	Remarks            string    `json:"remarks"`
}

type UsageVoucherResponse struct {
	ID                int64     `json:"id"`
	VoucherNo         string    `json:"voucher_no"`
	ReservationID     int64     `json:"reservation_id"`
	MemberBenefitID   int64     `json:"member_benefit_id"`
	VoucherType       string    `json:"voucher_type"`
	QrCode            string    `json:"qr_code"`
	ValidFrom         time.Time `json:"valid_from"`
	ValidTo           time.Time `json:"valid_to"`
	VerificationStatus string   `json:"verification_status"`
	Status            string    `json:"status"`
	UsedLocation      string    `json:"used_location"`
	ResponsiblePerson string    `json:"responsible_person"`
	BatchNo           string    `json:"batch_no"`
	Remarks           string    `json:"remarks"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type UsageVoucherListResponse struct {
	Total int                    `json:"total"`
	Page  int                    `json:"page"`
	Size  int                    `json:"size"`
	Data  []*UsageVoucherResponse `json:"data"`
}
