package models

import "time"

type UsageVoucher struct {
	ID                 int64     `json:"id"`
	VoucherNo          string    `json:"voucher_no"`
	ReservationID      int64     `json:"reservation_id"`
	MemberBenefitID    int64     `json:"member_benefit_id"`
	VoucherType        string    `json:"voucher_type"`
	QRCode             string    `json:"qr_code"`
	ValidFrom          time.Time `json:"valid_from"`
	ValidTo            time.Time `json:"valid_to"`
	UsedAt             *time.Time `json:"used_at"`
	UsedLocation       string    `json:"used_location"`
	VerificationStatus string    `json:"verification_status"`
	Status             string    `json:"status"`
	ResponsiblePerson  string    `json:"responsible_person"`
	BatchNo            string    `json:"batch_no"`
	Remarks            string    `json:"remarks"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

const (
	VoucherTypeSingle = "single_use"
	VoucherTypeMulti  = "multi_use"
	VoucherTypeGuest  = "guest_pass"
)

const (
	VoucherStatusActive   = "active"
	VoucherStatusUsed     = "used"
	VoucherStatusExpired  = "expired"
	VoucherStatusCancelled = "cancelled"
)

const (
	VoucherVerifyValid   = "valid"
	VoucherVerifyInvalid = "invalid"
	VoucherVerifyUsed    = "used"
)
