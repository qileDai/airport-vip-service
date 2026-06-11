package schemas

type StatisticsRequest struct {
	StartDate string `query:"start_date"`
	EndDate   string `query:"end_date"`
	BatchNo   string `query:"batch_no"`
	Role      string `query:"role"`
}

type DateStatistics struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

type BatchStatistics struct {
	BatchNo string `json:"batch_no"`
	Count   int64  `json:"count"`
}

type RoleStatistics struct {
	Role  string `json:"role"`
	Count int64  `json:"count"`
}

type VerificationStatisticsResponse struct {
	TotalVerifications int64             `json:"total_verifications"`
	PassedCount        int64             `json:"passed_count"`
	FailedCount        int64             `json:"failed_count"`
	PendingCount       int64             `json:"pending_count"`
	PassRate           float64           `json:"pass_rate"`
	FailRate           float64           `json:"fail_rate"`
	ByDate             []DateStatistics  `json:"by_date"`
	ByBatch            []BatchStatistics `json:"by_batch"`
	ByRole             []RoleStatistics  `json:"by_role"`
}

type WaitlistStatisticsResponse struct {
	TotalEntries     int64             `json:"total_entries"`
	SeatedCount      int64             `json:"seated_count"`
	CancelledCount   int64             `json:"cancelled_count"`
	WaitingCount     int64             `json:"waiting_count"`
	AverageWaitMins  int64             `json:"average_wait_mins"`
	TransferRate     float64           `json:"transfer_rate"`
	CancellationRate float64           `json:"cancellation_rate"`
	ByDate           []DateStatistics  `json:"by_date"`
	ByBatch          []BatchStatistics `json:"by_batch"`
	ByRole           []RoleStatistics  `json:"by_role"`
}

type UsageStatisticsResponse struct {
	TotalReservations int64             `json:"total_reservations"`
	CompletedCount    int64             `json:"completed_count"`
	CancelledCount    int64             `json:"cancelled_count"`
	UniqueMembers     int64             `json:"unique_members"`
	AverageGuestCount float64           `json:"average_guest_count"`
	CompletionRate    float64           `json:"completion_rate"`
	CancellationRate  float64           `json:"cancellation_rate"`
	ByDate            []DateStatistics  `json:"by_date"`
	ByBatch           []BatchStatistics `json:"by_batch"`
	ByRole            []RoleStatistics  `json:"by_role"`
}
