// unifi_resources/unifi_info.go
package unifi_resources

import (
	uc "github.com/j33pguy/Unifi_Provider/client/unifi"
)

type UnifiInfo struct {
	ApplicationVersion string `json:"applicationVersion"`
}

func GetUnifiInfo(client *uc.Client) (*UnifiInfo, error) {
	var info UnifiInfo
	err := client.DoRequest("GET", "/v1/info", nil, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}
