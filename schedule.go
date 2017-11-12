package evohome

type ZoneSchedule struct {
    DailySchedules []DailySchedule `json:"dailySchedules"`
}

type DailySchedule struct {
    DayOfWeek string           `json:"dayOfWeek"`
    SwitchPoints []SwitchPoint `json:"switchpoints"`
}

type SwitchPoint struct {
    Temperature float32 `json:"temperature"`
    Time string         `json:"timeOfDay"`
}
