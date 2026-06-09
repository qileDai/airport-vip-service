package schemas

import "time"

type MemberBenefitCreateRequest struct {
	MemberName        string    `json:"member_name" validate:"required"`
	MemberLevel       string    `json:"member_level" validate:"required,oneof=regular silver gold platinum"`
	RemainingQuota    int       `json:"remaining_quota" validate:"required,min=0"`
	TotalQuota        int       `json:"total_quota" validate:"required,min=1"`
	ValidFrom         time.Time `json:"valid_from"`
	ValidTo           time.Time `json:"valid_to"`
	ResponsiblePerson string    `json:"responsible_person"`
	BatchNo           string    `json:"batch_no"`
	Remarks           string    `json:"remarks"`
}

type MemberBenefitUpdateRequest struct {
	MemberName        string    `json:"member_name"`
	MemberLevel       string    `json:"member_level" validate:"omitempty,oneof=regular silver gold platinum"`
	RemainingQuota    *int      `json:"remaining_quota"`
	TotalQuota        *int      `json:"total_quota"`
	ValidFrom         time.Time `json:"valid_from"`
	ValidTo           time.Time `json:"valid_to"`
	Status            string    `json:"status" validate:"omitempty,oneof=active expired suspended"`
	ResponsiblePerson string    `json:"responsible_person"`
	Remarks           string    `json:"remarks"`
}

type MemberBenefitResponse struct {
	ID                int64     `json:"id"`
	BenefitNo         string    `json:"benefit_no"`
	MemberName        string    `json:"member_name"`
	MemberLevel       string    `json:"member_level"`
	RemainingQuota    int       `json:"remaining_quota"`
	TotalQuota        int       `json:"total_quota"`
	Status            string    `json:"status"`
	ResponsiblePerson string    `json:"responsible_person"`
	ValidFrom         time.Time `json:"valid_from"`
	ValidTo           time.Time `json:"valid_to"`
	BatchNo           string    `json:"batch_no"`
	Remarks           string    `json:"remarks"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type MemberBenefitListResponse struct {
	Total int                      `json:"total"`
	Page  int                      `json:"page"`
	Size  int                      `json:"size"`
	Data  []*MemberBenefitResponse `json:"data"`
}

type ExpiryCheckResponse struct {
	TotalChecked   int                      `json:"total_checked"`
	ExpiredCount   int                      `json:"expired_count"`
	ExpiringCount  int                      `json:"expiring_count"`
	ExpiredBenefits []*MemberBenefitResponse `json:"expired_benefits"`
	ExpiringBenefits []*MemberBenefitResponse `json:"expiring_benefits"`
}
