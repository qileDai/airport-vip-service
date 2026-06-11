package schemas

import "time"

type MemberBenefitCreateRequest struct {
	BenefitNo         string    `json:"benefit_no"`
	MemberName        string    `json:"member_name" validate:"required"`
	MemberLevel       string    `json:"member_level" validate:"required,oneof=regular silver gold platinum"`
	RemainingQuota    int       `json:"remaining_quota" validate:"required,min=0"`
	TotalQuota        int       `json:"total_quota" validate:"required,min=1"`
	Status            string    `json:"status"`
	ValidFrom         time.Time `json:"valid_from"`
	ValidTo           time.Time `json:"valid_to"`
	ResponsiblePerson string    `json:"responsible_person"`
	BatchNo           string    `json:"batch_no"`
	Remarks           string    `json:"remarks"`
}

type MemberBenefitUpdateRequest struct {
	MemberName        string     `json:"member_name"`
	MemberLevel       string     `json:"member_level" validate:"omitempty,oneof=regular silver gold platinum"`
	RemainingQuota    *int       `json:"remaining_quota"`
	TotalQuota        *int       `json:"total_quota"`
	ValidFrom         *time.Time `json:"valid_from"`
	ValidTo           *time.Time `json:"valid_to"`
	Status            string     `json:"status" validate:"omitempty,oneof=active expired suspended"`
	ResponsiblePerson string     `json:"responsible_person"`
	BatchNo           string     `json:"batch_no"`
	Remarks           string     `json:"remarks"`
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
	Total   int64                   `json:"total"`
	Page    int                     `json:"page"`
	PerPage int                     `json:"per_page"`
	Data    []MemberBenefitResponse `json:"data"`
}

type ExpiryCheckResponse struct {
	ExpiredCount   int                      `json:"expired_count"`
	ExpiringSoon   int                      `json:"expiring_soon"`
	ExpiredEvents  []ExceptionEventResponse `json:"expired_events"`
	ExpiringEvents []ExceptionEventResponse `json:"expiring_events"`
}
