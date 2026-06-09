package schemas

type StatisticsRequest struct {
	StartDate string `query:"start_date"`
	EndDate   string `query:"end_date"`
	BatchNo   string `query:"batch_no"`
	Role      string `query:"role"`
}

type DateStatistics struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type BatchStatistics struct {
	BatchNo string `json:"batch_no"`
	Count   int    `json:"count"`
}

type RoleStatistics struct {
	Role  string `json:"role"`
	Count int    `json:"count"`
}

type VerificationStatisticsResponse struct {
	TotalVerifications int                `json:"total_verifications"`
	PassedCount        int                `json:"passed_count"`
	FailedCount        int                `json:"failed_count"`
	PendingCount       int                `json:"pending_count"`
	PassRate           float64            `json:"pass_rate"`
	FailRate           float64            `json:"fail_rate"`
	ByDate             []DateStatistics   `json:"by_date"`
	ByBatch            []BatchStatistics  `json:"by_batch"`
	ByRole             []RoleStatistics   `json:"by_role"`
}

type WaitlistStatisticsResponse struct {
	TotalEntries      int               `json:"total_entries"`
	SeatedCount       int               `json:"seated_count"`
	CancelledCount    int               `json:"cancelled_count"`
	WaitingCount      int               `json:"waiting_count"`
	AverageWaitMins   int               `json:"average_wait_mins"`
	TransferRate      float64           `json:"transfer_rate"`
	CancellationRate  float64           `json:"cancellation_rate"`
	ByDate            []DateStatistics  `json:"by_date"`
	ByBatch           []BatchStatistics `json:"by_batch"`
	ByRole            []RoleStatistics  `json:"by_role"`
}

type UsageStatisticsResponse struct {
	TotalReservations  int               `json:"total_reservations"`
	CompletedCount     int               `json:"completed_count"`
	CancelledCount     int               `json:"cancelled_count"`
	UniqueMembers      int               `json:"unique_members"`
	AverageGuestCount  float64           `json:"average_guest_count"`
	CompletionRate     float64           `json:"completion_rate"`
	CancellationRate   float64           `json:"cancellation_rate"`
	ByDate             []DateStatistics  `json:"by_date"`
	ByBatch            []BatchStatistics `json:"by_batch"`
	ByRole             []RoleStatistics  `json:"by_role"`
}
