package types

type GetArrivedGuestsResponse struct {
	Name               string `json:"name"`
	AccompanyingGuests int    `json:"accompanying_guests"`
	TimeArrived        string `json:"time_arrived"`
}
