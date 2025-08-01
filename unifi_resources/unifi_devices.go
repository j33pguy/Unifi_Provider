// unifi_resources/unifi_devices.go
package unifi_resources

import (
	"fmt"

	uc "github.com/j33pguy/Unifi_Provider/client/unifi"
)

type UnifiDevices struct {
	Offset     int           `json:"offset"`
	Limit      int           `json:"limit"`
	Count      int           `json:"count"`
	TotalCount int           `json:"totalCount"`
	Data       []UnifiDevice `json:"data"`
}

type UnifiDevice struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Model      string   `json:"model"`
	MacAddress string   `json:"macAddress"`
	IPAddress  string   `json:"ipAddress"`
	State      string   `json:"state"`
	Features   []string `json:"features"`
	Interfaces []string `json:"interfaces"`
}

type UnifiDeviceDetails struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Model             string `json:"model"`
	Supported         bool   `json:"supported"`
	MacAddress        string `json:"macAddress"`
	IPAddress         string `json:"ipAddress"`
	State             string `json:"state"`
	FirmwareVersion   string `json:"firmwareVersion"`
	FirmwareUpdatable bool   `json:"firmwareUpdatable"`
	AdoptedAt         string `json:"adoptedAt"`
	ProvisionedAt     string `json:"provisionedAt"`
	ConfigurationId   string `json:"configurationId"`
	Uplink            struct {
		DeviceID string `json:"deviceId"`
	} `json:"uplink"`
	Features struct {
		Switching   interface{} `json:"switching"`
		AccessPoint interface{} `json:"accessPoint"`
	} `json:"features"`
	Interfaces struct {
		Ports  []UnifiPort  `json:"ports"`
		Radios []UnifiRadio `json:"radios"`
	} `json:"interfaces"`
}

type UnifiPort struct {
	Idx          int    `json:"idx"`
	State        string `json:"state"`
	Connector    string `json:"connector"`
	MaxSpeedMbps int    `json:"maxSpeedMbps"`
	SpeedMbps    int    `json:"speedMbps"`
	Poe          struct {
		Standard string `json:"standard"`
		Type     int    `json:"type"`
		Enabled  bool   `json:"enabled"`
		State    string `json:"state"`
	} `json:"poe"`
}

type UnifiRadio struct {
	WlanStandard    string `json:"wlanStandard"`
	FrequencyGHz    string `json:"frequencyGHz"`
	ChannelWidthMHz int    `json:"channelWidthMHz"`
	Channel         int    `json:"channel"`
}

type UnifiDeviceStats struct {
	UptimeSec            int     `json:"uptimeSec"`
	LastHeartbeatAt      string  `json:"lastHeartbeatAt"`
	NextHeartbeatAt      string  `json:"nextHeartbeatAt"`
	LoadAverage1Min      float64 `json:"loadAverage1Min"`
	LoadAverage5Min      float64 `json:"loadAverage5Min"`
	LoadAverage15Min     float64 `json:"loadAverage15Min"`
	CpuUtilizationPct    float64 `json:"cpuUtilizationPct"`
	MemoryUtilizationPct float64 `json:"memoryUtilizationPct"`
	Uplink               struct {
		TxRateBps int `json:"txRateBps"`
		RxRateBps int `json:"rxRateBps"`
	} `json:"uplink"`
	Interfaces struct {
		Radios []struct {
			FrequencyGHz string  `json:"frequencyGHz"`
			TxRetriesPct float64 `json:"txRetriesPct"`
		} `json:"radios"`
	} `json:"interfaces"`
}

type DeviceAction struct {
	Action string `json:"action"`
}

type PortAction struct {
	Action string `json:"action"`
}

func ListUnifiDevices(client *uc.Client) (*UnifiDevices, error) {
	var devices UnifiDevices
	err := client.DoRequest("GET", client.SitePath("devices"), nil, &devices)
	if err != nil {
		return nil, err
	}
	return &devices, nil
}

func GetUnifiDevice(client *uc.Client, deviceId string) (*UnifiDeviceDetails, error) {
	var device UnifiDeviceDetails
	path := client.SitePath("devices/" + deviceId)
	err := client.DoRequest("GET", path, nil, &device)
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func GetUnifiDeviceStats(client *uc.Client, deviceId string) (*UnifiDeviceStats, error) {
	var stats UnifiDeviceStats
	path := client.SitePath("devices/" + deviceId + "/statistics/latest")
	err := client.DoRequest("GET", path, nil, &stats)
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

func ExecuteUnifiDeviceAction(client *uc.Client, deviceId string, action string) error {
	body := DeviceAction{Action: action}
	path := client.SitePath("devices/" + deviceId + "/actions")
	return client.DoRequest("POST", path, body, nil)
}

func ExecuteUnifiPortAction(client *uc.Client, deviceId string, portIdx int, action string) error {
	body := PortAction{Action: action}
	path := client.SitePath(fmt.Sprintf("devices/%s/interfaces/ports/%d/actions", deviceId, portIdx))
	return client.DoRequest("POST", path, body, nil)
}
