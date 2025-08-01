// provider/provider.go
package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uc "github.com/j33pguy/Unifi_Provider/client/unifi"
)

type ProviderConfig struct {
	Username  string
	Password  string
	APIKey    string
	Host      string
	Site      string
	VerifySSL bool
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("UNIFI_USERNAME", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("UNIFI_PASSWORD", nil),
				Sensitive:   true,
			},
			"api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("UNIFI_API_KEY", nil),
				Sensitive:   true,
			},
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("UNIFI_HOST", nil),
			},
			"site": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
			},
			"verify_ssl": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
		ConfigureContextFunc: providerConfigure,
		ResourcesMap:         map[string]*schema.Resource{},
		DataSourcesMap:       map[string]*schema.Resource{},
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config := ProviderConfig{
		Username:  d.Get("username").(string),
		Password:  d.Get("password").(string),
		APIKey:    d.Get("api_key").(string),
		Host:      d.Get("host").(string),
		Site:      d.Get("site").(string),
		VerifySSL: d.Get("verify_ssl").(bool),
	}

	if config.APIKey == "" && (config.Username == "" || config.Password == "") {
		return nil, diag.Errorf("either api_key or username and password must be provided")
	}

	client, err := uc.NewClient(config.Host, config.Username, config.Password, config.APIKey, config.VerifySSL)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return client, nil
}
