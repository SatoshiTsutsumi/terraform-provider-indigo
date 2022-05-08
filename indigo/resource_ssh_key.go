package indigo

import (
	"context"
	"strconv"

	"github.com/UndefxDev/indigo-client-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSSHKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSSHKeyCreate,
		ReadContext:   resourceSSHKeyRead,
		UpdateContext: resourceSSHKeyUpdate,
		DeleteContext: resourceSSHKeyDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"service_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceSSHKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*indigo.Client)
	name := d.Get("name").(string)
	sshKeyStr := d.Get("key").(string)

	sshKey, err := c.CreateSSHKey(name, sshKeyStr)
	if err != nil {
		return diag.FromErr(err)
	}
	setSSHKey(d, sshKey)

	return diags
}

func resourceSSHKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*indigo.Client)

	sshKeyID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	sshKey, err := c.RetrieveSSHKey(sshKeyID)
	if err != nil {
		return diag.FromErr(err)
	}
	setSSHKey(d, sshKey)

	return diags
}

func resourceSSHKeyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*indigo.Client)
	sshKeyID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	sshName := d.Get("name").(string)
	sshKey := d.Get("key").(string)
	sshKeyStatus := d.Get("status").(string)

	err = c.UpdateSSHKey(sshKeyID, sshName, sshKey, sshKeyStatus)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceSSHKeyDelete(tx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*indigo.Client)
	sshKeyID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.DeleteSSHKey(sshKeyID)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")

	return diags
}

func setSSHKey(d *schema.ResourceData, sshKey *indigo.SSHKey) {
	d.Set("name", sshKey.Name)
	d.Set("key", sshKey.Key)
	d.Set("status", sshKey.Status)
	d.Set("user_id", sshKey.UserID)
	d.Set("service_id", sshKey.ServiceID)
	d.Set("updated_at", sshKey.UpdatedAt)
	d.Set("created_at", sshKey.CreatedAt)
	d.SetId(strconv.Itoa(sshKey.ID))
}
