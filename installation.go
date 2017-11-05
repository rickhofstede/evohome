package evohome

type Installation struct {
    Location Location  `json:"locationInfo"`
    Gateways []Gateway `json:"gateways"`
}
