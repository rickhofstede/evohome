package evohome

type HeatSetPointStatus struct {
    TargetTemperature float32 `json:"targetTemperature"`
    SetPointMode string       `json:"setpointMode"`
}
