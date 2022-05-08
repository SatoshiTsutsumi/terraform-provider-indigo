package indigo

import (
	"context"
	"strconv"

	"github.com/UndefxDev/indigo-client-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSnapshot() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSnapshotCreate,
		ReadContext:   resourceSnapshotRead,
		UpdateContext: resourceSnapshotUpdate,
		DeleteContext: resourceSnapshotDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"service_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"disk_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"volume": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"slot_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deleted": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deleted_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"retry": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"restore": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceSnapshotCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*indigo.Client)

	name := d.Get("name").(string)
	instanceID := d.Get("instance_id").(int)

	snapshot, err := c.CreateSnapshotSync(name, instanceID)
	if err != nil {
		return diag.FromErr(err)
	}

	setSnapshot(d, snapshot)

	return diags
}

func resourceSnapshotRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*indigo.Client)
	instanceID := d.Get("instance_id").(int)
	id, _ := strconv.Atoi(d.Id())
	snapshot, err := c.GetSnapshot(instanceID, id)
	if err != nil {
		return diag.FromErr(err)
	}
	setSnapshot(d, snapshot)

	return diags
}

func resourceSnapshotUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*indigo.Client)
	instanceID := d.Get("instance_id").(int)
	snapshotID, _ := strconv.Atoi(d.Id())

	if d.HasChange("restore") {
		restoreIf, existsRestore := d.GetOk("restore")
		if !existsRestore || !restoreIf.(bool) {
			return diags
		}

		err := c.RestoreSnapshot(instanceID, snapshotID)
		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("restore", false)
	} else if d.HasChange("retry") {
		retryIf, existsRetry := d.GetOk("retry")
		if !existsRetry || !retryIf.(bool) {
			return diags
		}

		snapshot, err := c.RecreateSnapshotSync(instanceID, snapshotID)
		if err != nil {
			return diag.FromErr(err)
		}
		setSnapshot(d, snapshot)
		d.Set("retry", false)
	}

	return diags
}

func resourceSnapshotDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*indigo.Client)
	snapshotID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.DeleteSnapshot(snapshotID)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")

	return diags
}

func setSnapshot(d *schema.ResourceData, snapshot *indigo.Snapshot) {
	d.Set("service_id", snapshot.ServiceID)
	d.Set("user_id", snapshot.UserID)
	d.Set("disk_id", snapshot.DiskID)
	d.Set("volume", snapshot.Volume)
	d.Set("slot_number", snapshot.SlotNumber)
	d.Set("status", snapshot.Status)
	d.Set("size", snapshot.Size)
	d.Set("deleted", snapshot.Deleted)
	d.Set("created_at", snapshot.CreatedAt)
	d.Set("deleted_at", snapshot.Deleted)
	d.SetId(strconv.Itoa(snapshot.ID))
}
