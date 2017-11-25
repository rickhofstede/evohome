package evohome

type ControlSystem interface {
    Zone(string) (*Zone)
    ZoneNames() ([]string)
    ZoneNamesWithOverride() ([]string)
    ZonesMap() (map[string]*Zone)
}
