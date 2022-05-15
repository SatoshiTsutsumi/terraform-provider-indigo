package indigo

import (
	"context"
	"strconv"

	"github.com/UndefxDev/indigo-client-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFirewall() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFirewallCreate,
		ReadContext:   resourceFirewallRead,
		UpdateContext: resourceFirewallUpdate,
		DeleteContext: resourceFirewallDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"inbound": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Required: true,
						},
						"port": {
							Type:     schema.TypeString,
							Required: true,
						},
						"source": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"outbound": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Required: true,
						},
						"port": {
							Type:     schema.TypeString,
							Required: true,
						},
						"source": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"instance_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceFirewallCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return createOrUpdateFirewall(true, ctx, d, m)
}

func resourceFirewallRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*indigo.Client)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("invalid firewall id %s", d.Id())
	}
	firewall, err := c.RetrieveFirewall(id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", firewall.Name)
	d.Set("inbound", flattenRules(firewall.Inbound))
	d.Set("outbound", flattenRules(firewall.Outbound))
	return diags
}

func resourceFirewallUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return createOrUpdateFirewall(false, ctx, d, m)
}

func resourceFirewallDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*indigo.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("invalid firewall id %s", d.Id())
	}

	err = c.DeleteFirewall(id)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")

	return diags
}

func createOrUpdateFirewall(create bool, ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*indigo.Client)

	name := d.Get("name").(string)

	inboundIf, inboundOk := d.GetOk("inbound")
	inbound := []*indigo.Rule{}
	if inboundOk {
		for _, ruleIf := range inboundIf.(([]interface{})) {
			rule := ruleIf.(map[string]interface{})
			inbound = append(inbound, &indigo.Rule{
				Type:     rule["type"].(string),
				Protocol: rule["protocol"].(string),
				Port:     rule["port"].(string),
				Source:   rule["source"].(string),
			})
		}
	}
	outboundIf, outboundOk := d.GetOk("outbound")
	outbound := []*indigo.Rule{}
	if outboundOk {
		for _, ruleIf := range outboundIf.(([]interface{})) {
			rule := ruleIf.(map[string]interface{})
			outbound = append(outbound, &indigo.Rule{
				Type:     rule["type"].(string),
				Protocol: rule["protocol"].(string),
				Port:     rule["port"].(string),
				Source:   rule["source"].(string),
			})
		}
	}
	if !inboundOk && !outboundOk {
		return diag.Errorf("either inbound or outbound required")
	}
	instanceIDsIf, instanceIDsOk := d.GetOk("instance_ids")
	instanceIDs := []int{}
	if instanceIDsOk {
		for _, instanceID := range instanceIDsIf.([]interface{}) {
			instanceIDs = append(instanceIDs, instanceID.(int))
		}
	}

	id, err := c.CreateFirewall(name, inbound, outbound, instanceIDs)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(id))
	d.Set("instance_ids", instanceIDs) // NOTE: API never returns instanceIDs

	return resourceFirewallRead(ctx, d, m)
}

func flattenRules(rules []*indigo.Rule) []map[string]interface{} {
	rs := make([]map[string]interface{}, len(rules))
	for i, rule := range rules {
		r := make(map[string]interface{})
		r["type"] = rule.Type
		r["protocol"] = rule.Protocol
		r["port"] = rule.Port
		r["source"] = rule.Source
		rs[i] = r
	}
	return rs
}
