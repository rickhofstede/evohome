package evohome

type TemperatureStatus struct {
    Temperature float32 `json:"temperature"`
    IsAvailable bool    `json:"isAvailable"`
}
