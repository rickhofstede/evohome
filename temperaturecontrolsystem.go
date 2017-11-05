package evohome

type TemperatureControlSystem struct {
    Id string    `json:"systemId"`
    Type string  `json:"modelType"`
    Zones []Zone `json:"zones"`
}

// Retrieve a Zone by its name.
func (t *TemperatureControlSystem) Zone(name string) (*Zone) {
    var foundZone *Zone
    for _, zone := range t.Zones {
        if zone.Name == name {
            foundZone = &zone
            break
        }
    }
    return foundZone
}

// Returns a list of all zone names within this temperatureControlSystem.
func (t *TemperatureControlSystem) ZoneNames() ([]string) {
    var names []string
    for _, zone := range t.Zones {
        names = append(names, zone.Name)
    }

    return names
}

// Returns a list of all zone names with a (temperature) override.
func (t *TemperatureControlSystem) ZoneNamesWithOverride() ([]string) {
    var names []string

    for _, zone := range t.Zones {
        if zone.HeatSetPointStatus.SetPointMode != "FollowSchedule" {
            names = append(names, zone.Name)
        }
    }

    return names
}
