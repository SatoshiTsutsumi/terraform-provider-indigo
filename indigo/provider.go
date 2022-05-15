package indigo

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("INDIGO_API_HOST", nil),
			},
			"api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("INDIGO_API_KEY", nil),
			},
			"api_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("INDIGO_API_SECRET", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"indigo_ssh_key":  resourceSSHKey(),
			"indigo_instance": resourceInstance(),
			"indigo_snapshot": resourceSnapshot(),
			"indigo_firewall": resourceFirewall(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}
