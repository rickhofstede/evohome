package evohome

import (
    "bytes"
    "encoding/json"
    "time"
    "github.com/jehiah/go-strftime"
)

type Zone struct {
    Id string                             `json:"zoneId"`
    Name string                           `json:"name"`
    ModelType string                      `json:"modelType"`
    ZoneType string                       `json:"zoneType"`
    TemperatureStatus TemperatureStatus   `json:"temperatureStatus"`
    HeatSetPointStatus HeatSetPointStatus `json:"heatSetpointStatus"`
}

type ZoneTemperatureSetting struct {
    HeatSetpointValue float32
    SetpointMode int
    TimeUntil *string
}

// Set the zone temperature to the specified value. If `until` is `nil`,
// a permanent override will be set.
func (z *Zone) SetTemperature(temp float32, until time.Time) (err error) {
    tempSetting := ZoneTemperatureSetting {
        HeatSetpointValue: temp,
    }

    if until.IsZero() {
        tempSetting.SetpointMode = 1
    } else {
        tempSetting.SetpointMode = 2
        untilString := strftime.Format("%Y-%m-%dT%H:%M:%SZ", until.UTC())
        tempSetting.TimeUntil = &untilString
    }

    b := new(bytes.Buffer)
    json.NewEncoder(b).Encode(tempSetting)
    body, err := request("PUT", b, "temperatureZone/%s/heatSetpoint", z.Id)

    defer body.Close()
    return err
}

// Cancel the configured temperature/schedule override.
func (z *Zone) CancelTemperatureOverride() (err error) {
    tempSetting := ZoneTemperatureSetting {
        HeatSetpointValue: 0.0,
        SetpointMode: 0,
    }

    b := new(bytes.Buffer)
    json.NewEncoder(b).Encode(tempSetting)
    body, err := request("PUT", b, "temperatureZone/%s/heatSetpoint", z.Id)

    defer body.Close()
    return err
}
