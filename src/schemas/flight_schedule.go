package schemas

import "time"

type FlightScheduleCreateRequest struct {
	FlightNo          string     `json:"flight_no" validate:"required"`
	DepartureAirport  string     `json:"departure_airport" validate:"required"`
	ArrivalAirport    string     `json:"arrival_airport" validate:"required"`
	ScheduledDepart   time.Time  `json:"scheduled_depart" validate:"required"`
	ScheduledArrive   time.Time  `json:"scheduled_arrive" validate:"required"`
	ActualDepart      *time.Time `json:"actual_depart"`
	ActualArrive      *time.Time `json:"actual_arrive"`
	FlightStatus      string     `json:"flight_status"`
	VipLoungeCapacity int        `json:"vip_lounge_capacity"`
	VipLoungeUsed     int        `json:"vip_lounge_used"`
	Status            string     `json:"status"`
	ResponsiblePerson string     `json:"responsible_person"`
	BatchNo           string     `json:"batch_no"`
	Remarks           string     `json:"remarks"`
}

type FlightScheduleUpdateRequest struct {
	DepartureAirport  string     `json:"departure_airport"`
	ArrivalAirport    string     `json:"arrival_airport"`
	ScheduledDepart   *time.Time `json:"scheduled_depart"`
	ScheduledArrive   *time.Time `json:"scheduled_arrive"`
	ActualDepart      *time.Time `json:"actual_depart"`
	ActualArrive      *time.Time `json:"actual_arrive"`
	FlightStatus      string     `json:"flight_status"`
	VipLoungeCapacity *int       `json:"vip_lounge_capacity"`
	VipLoungeUsed     *int       `json:"vip_lounge_used"`
	Status            string     `json:"status"`
	ResponsiblePerson string     `json:"responsible_person"`
	BatchNo           string     `json:"batch_no"`
	Remarks           string     `json:"remarks"`
}

type FlightScheduleResponse struct {
	ID                int64      `json:"id"`
	FlightNo          string     `json:"flight_no"`
	DepartureAirport  string     `json:"departure_airport"`
	ArrivalAirport    string     `json:"arrival_airport"`
	ScheduledDepart   time.Time  `json:"scheduled_depart"`
	ScheduledArrive   time.Time  `json:"scheduled_arrive"`
	ActualDepart      *time.Time `json:"actual_depart"`
	ActualArrive      *time.Time `json:"actual_arrive"`
	FlightStatus      string     `json:"flight_status"`
	VipLoungeCapacity int        `json:"vip_lounge_capacity"`
	VipLoungeUsed     int        `json:"vip_lounge_used"`
	Status            string     `json:"status"`
	ResponsiblePerson string     `json:"responsible_person"`
	BatchNo           string     `json:"batch_no"`
	Remarks           string     `json:"remarks"`
	SnapshotData      string     `json:"snapshot_data"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type FlightScheduleListResponse struct {
	Total   int64                    `json:"total"`
	Page    int                      `json:"page"`
	PerPage int                      `json:"per_page"`
	Data    []FlightScheduleResponse `json:"data"`
}

type FlightArchiveRequest struct {
	FlightScheduleIDs []int64 `json:"flight_schedule_ids" validate:"required"`
	CreateSnapshot    bool    `json:"create_snapshot"`
	Operator          string  `json:"operator"`
}

type FlightSnapshotResponse struct {
	FlightScheduleID int64     `json:"flight_schedule_id"`
	SnapshotData     string    `json:"snapshot_data"`
	ArchivedAt       time.Time `json:"archived_at"`
}
