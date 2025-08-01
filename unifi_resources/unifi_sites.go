// unifi_resources/unifi_sites.go
package unifi_resources

import (
	uc "github.com/j33pguy/Unifi_Provider/client/unifi"
)

type UnifiSites struct {
	Offset     int `json:"offset"`
	Limit      int `json:"limit"`
	Count      int `json:"count"`
	TotalCount int `json:"totalCount"`
	Data       []struct {
		ID                string `json:"id"`
		InternalReference string `json:"internalReference"`
		Name              string `json:"name"`
	} `json:"data"`
}

func GetUnifiSites(client *uc.Client) (*UnifiSites, error) {
	var sites UnifiSites
	err := client.DoRequest("GET", "/v1/sites", nil, &sites)
	if err != nil {
		return nil, err
	}
	return &sites, nil
}
