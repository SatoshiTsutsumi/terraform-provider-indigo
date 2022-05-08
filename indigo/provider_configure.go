package indigo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/UndefxDev/indigo-client-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiKey := d.Get("api_key").(string)
	apiSecret := d.Get("api_secret").(string)
	host := d.Get("host").(string)

	var diags diag.Diagnostics
	missingParam := []string{}

	if host == "" {
		missingParam = append(missingParam, "host")
	}
	if apiKey == "" {
		missingParam = append(missingParam, "apiKey")
	}
	if apiSecret == "" {
		missingParam = append(missingParam, "apiSecret")
	}

	if len(missingParam) > 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Indigo client",
			Detail:   fmt.Sprintf("Empty parameters: %s", strings.Join(missingParam, ",")),
		})
		return nil, diags
	}

	c, err := indigo.NewClient(host, apiKey, apiSecret, false)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Indigo client",
			Detail:   err.Error(),
		})
		return nil, diags
	}
	// NOTE: WebArena allows one request / sec and requires more interval actually.
	time.Sleep(time.Second * 3)

	return c, diags
}
