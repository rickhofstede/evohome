package evohome

type Account struct {
    UserId string        `json:"userId"`
    Username string      `json:"username"`
    Language string      `json:"language"`
    FirstName string     `json:"firstname"`
    LastName string      `json:"lastname"`
    Address string       `json:"streetAddress"`
    ZipCode string       `json:"postcode"`
    City string          `json:"city"`
    Country string       `json:"country"`
}
