package evohome

import (
    "encoding/json"
)

type Gateway struct {
    GatewayInfo *json.RawMessage                         `json:"gatewayInfo"`
    TemperatureControlSystems []TemperatureControlSystem `json:"temperatureControlSystems"`
}
