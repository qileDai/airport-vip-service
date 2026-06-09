package models

import "time"

type FlightSchedule struct {
	ID                int64     `json:"id"`
	FlightNo          string    `json:"flight_no"`
	DepartureAirport  string    `json:"departure_airport"`
	ArrivalAirport    string    `json:"arrival_airport"`
	ScheduledDepart   time.Time `json:"scheduled_depart"`
	ScheduledArrive   time.Time `json:"scheduled_arrive"`
	ActualDepart      *time.Time `json:"actual_depart"`
	ActualArrive      *time.Time `json:"actual_arrive"`
	FlightStatus      string    `json:"flight_status"`
	VipLoungeCapacity int       `json:"vip_lounge_capacity"`
	VipLoungeUsed     int       `json:"vip_lounge_used"`
	Status            string    `json:"status"`
	ResponsiblePerson string    `json:"responsible_person"`
	BatchNo           string    `json:"batch_no"`
	Remarks           string    `json:"remarks"`
	SnapshotData      string    `json:"snapshot_data"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

const (
	FlightStatusScheduled   = "scheduled"
	FlightStatusBoarding    = "boarding"
	FlightStatusDeparted    = "departed"
	FlightStatusArrived     = "arrived"
	FlightStatusDelayed     = "delayed"
	FlightStatusCancelled   = "cancelled"
)

const (
	FlightRecordStatusActive   = "active"
	FlightRecordStatusArchived = "archived"
)
