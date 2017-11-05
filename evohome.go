package evohome

import (
    "bytes"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net/http"
    // "net/http/httputil"
    "net/url"
    // "strings"
    "time"
)

var accessToken string
var refreshToken string

const (
    temperatureRefreshInterval time.Duration = 2 * time.Second
)

type Evohome struct {
    initialized bool
    account Account
    installations []Installation
}

// Create new Evohome instance.
func NewEvohome(username string, password string) (*Evohome) {
    var err error
    accessToken, refreshToken, err = tokens(username, password)
    if err != nil {
        return nil
    }

    e := &Evohome {
        initialized: true,
        account: account(),
    }

    e.installations = installations(e.account.UserId)
    go e.UpdateTemperatures()
    return e
}

// Update zone temperatures
func (e *Evohome) UpdateTemperatures() () {
    for {
        zones := e.TemperatureControlSystem().Zones
        zoneTemperatures := zoneTemperatures(e.installations[0].Location.Id)

        // Merge zone temperatures into Zone objects
        for i, outerZone := range zones {
            for _, innerZone := range zoneTemperatures {
                if outerZone.Id == innerZone.Id {
                    zones[i].TemperatureStatus.IsAvailable = innerZone.TemperatureStatus.IsAvailable
                    zones[i].TemperatureStatus.Temperature = innerZone.TemperatureStatus.Temperature
                    zones[i].HeatSetPointStatus.TargetTemperature = innerZone.HeatSetPointStatus.TargetTemperature
                    zones[i].HeatSetPointStatus.SetPointMode = innerZone.HeatSetPointStatus.SetPointMode
                    break
                }
            }
        }
        time.Sleep(temperatureRefreshInterval)
    }
}

func (e *Evohome) Initialized() (bool) {
    return e.initialized
}

func (e *Evohome) Account() (Account) {
    return e.account
}

// Try to retrieve the first registered temperature control system.
func (e *Evohome) TemperatureControlSystem() (TemperatureControlSystem) {
    if !e.Initialized() {
        panic("Evohome not initialized")
    }
    if len(e.installations) == 0 ||
            len(e.installations[0].Gateways) == 0 ||
            len(e.installations[0].Gateways[0].TemperatureControlSystems) == 0 {
        panic("Cannot retrieve temperature control system")
    }
    return e.installations[0].Gateways[0].TemperatureControlSystems[0]
}

func tokens(username string, password string) (accessToken string, refreshToken string, err error) {
    data := url.Values{
        "Content-Type":     { "application/x-www-form-urlencoded; charset=utf-8'" },
        "Host":             { "rs.alarmnet.com/" },
        "Cache-Control":    { "no-store no-cache" },
        "Pragma":           { "no-cache" },
        "grant_type":       { "password" },
        "scope":            { "EMEA-V1-Basic EMEA-V1-Anonymous EMEA-V1-Get-Current-User-Account" },
        "Username":         { username },
        "Password":         { password },
        "Connection":       { "Keep-Alive" },
    }

    url := "https://tccna.honeywell.com/Auth/OAuth/Token"
    req, _ := http.NewRequest("POST", url, bytes.NewBufferString(data.Encode()))
    req.Header.Set("Authorization", "Basic YjAxM2FhMjYtOTcyNC00ZGJkLTg4OTctMDQ4YjlhYWRhMjQ5OnRlc3Q=")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        panic(err)
    }

    defer resp.Body.Close()

    var respData = make(map[string]interface{})
    err = json.NewDecoder(resp.Body).Decode(&respData)
    if err != nil {
        panic(err)
    }

    var ok bool
    if accessToken, ok = respData["access_token"].(string); !ok {
        err = errors.New("Could not retrieve token(s)")
    }
    if refreshToken, ok = respData["refresh_token"].(string); !ok {
        err = errors.New("Could not retrieve token(s)")
    }

    return accessToken, refreshToken, err
}

func request(requestType string, data io.Reader, path string, pathVars ...interface{}) (responseBody io.ReadCloser, err error) {
    if requestType != "GET" && requestType != "PUT" {
        panic(fmt.Sprintf("Invalid HTTP request type: '%s'", requestType))
    }
    if accessToken == "" {
        panic("No access token available")
    }

    url := "https://tccna.honeywell.com/WebAPI/emea/api/v1/" + path
    req, _ := http.NewRequest(requestType, fmt.Sprintf(url, pathVars...), data)

    req.Header.Set("Authorization", "bearer " + accessToken)
    req.Header.Set("ApplicationId", "b013aa26-9724-4dbd-8897-048b9aada249")
    req.Header.Set("Content-Type", "application/json")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return responseBody, err
    }

    // reqDump, _ := httputil.DumpRequest(req, true)
    // fmt.Printf("request: %q\n", reqDump)

    // defer resp.Body.Close()
    return resp.Body, err
}

func account() (Account) {
    body, err := request("GET", nil, "userAccount")
    if err != nil {
        panic(err)
    }

    var acc Account
    err = json.NewDecoder(body).Decode(&acc)
    if err != nil {
        panic(err)
    }

    defer body.Close()
    return acc
}

func installations(userId string) ([]Installation) {
    body, err := request("GET", nil,
            "location/installationInfo?userId=%s&includeTemperatureControlSystems=True", userId)
    if err != nil {
        panic(err)
    }

    var installations []Installation
    err = json.NewDecoder(body).Decode(&installations)
    if err != nil {
        panic(err)
    }

    defer body.Close()
    return installations
}

func zoneTemperatures(locationId string) ([]Zone) {
    body, err := request("GET", nil,
            "location/%s/status?includeTemperatureControlSystems=True", locationId)
    if err != nil {
        panic(err)
    }

    var installation Installation
    err = json.NewDecoder(body).Decode(&installation)
    if err != nil {
        panic(err)
    }

    defer body.Close()
    return installation.Gateways[0].TemperatureControlSystems[0].Zones
}
