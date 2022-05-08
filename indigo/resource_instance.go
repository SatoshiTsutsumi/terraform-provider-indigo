package indigo

import (
	"context"
	"strconv"

	"github.com/UndefxDev/indigo-client-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInstanceCreate,
		ReadContext:   resourceInstanceRead,
		UpdateContext: resourceInstanceUpdate,
		DeleteContext: resourceInstanceDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"set_no": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vps_kind": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"sequence_id": {
				Type:     schema.TypeInt,
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
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssh_key_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"snapshot_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"import_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"start_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"host_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"plan": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"plan_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"disk_point": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"mem_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cpus": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"os_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"other_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"uid_gid": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vnc_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vnc_passwd": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"arpa_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"arpa_date": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"started_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"closed_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status_change_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vm_revert": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mac_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"import_instance": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"container_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"daemon_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"out_of_stock": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ip_address_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ve_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceInstanceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*indigo.Client)

	sshKeyIDIf, existsSSHKey := d.GetOk("ssh_key_id")
	regionIDIf, existsRegionID := d.GetOk("region_id")
	osIDIf, existsOSID := d.GetOk("os_id")
	planIDIf, existsPlanID := d.GetOk("plan_id")
	nameIf, existsName := d.GetOk("name")
	snapshotIDIf, existsSnapshot := d.GetOk("snapshot_id")
	importURLIf, existsImportURL := d.GetOk("import_url")

	name := nameIf.(string)
	if !existsName || name == "" {
		return diag.Errorf("invalid arguments")
	}

	if existsSSHKey && existsSnapshot && existsPlanID {
		var err error
		var sshKeyID int
		sshKeyID, err = strconv.Atoi(sshKeyIDIf.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		snapshotID := snapshotIDIf.(int)
		planID := planIDIf.(int)

		instance, err := c.CreateSnapshotInstanceSync(sshKeyID, snapshotID, planID, name)
		if err != nil {
			return diag.FromErr(err)
		}
		setInstance(d, instance)
	} else if existsImportURL && existsRegionID && existsOSID && existsPlanID {
		var err error
		importURL := importURLIf.(string)
		regionID := regionIDIf.(int)
		osID := osIDIf.(int)
		planID := planIDIf.(int)
		instance, err := c.CreateImportInstanceSync(importURL, regionID, osID, planID, name)
		if err != nil {
			return diag.FromErr(err)
		}
		setInstance(d, instance)
	} else if existsSSHKey && existsRegionID && existsOSID && existsPlanID {
		var err error
		var sshKeyID int
		sshKeyID, err = strconv.Atoi(sshKeyIDIf.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		regionID := regionIDIf.(int)
		osID := osIDIf.(int)
		planID := planIDIf.(int)
		instance, err := c.CreateInstanceSync(sshKeyID, regionID, osID, planID, name)
		if err != nil {
			return diag.FromErr(err)
		}
		setInstance(d, instance)
	} else {
		return diag.Errorf("invalid arguments")
	}

	return diags
}

func resourceInstanceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*indigo.Client)
	instanceID, _ := strconv.Atoi(d.Id())
	instance, err := c.GetInstance(instanceID)
	if err != nil {
		return diag.FromErr(err)
	}
	setInstance(d, instance)
	return diags
}

func resourceInstanceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	if d.HasChange("status") {
		c := m.(*indigo.Client)
		instanceID, err := strconv.Atoi(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
		status := d.Get("status").(string)

		err = c.UpdateInstanceStatus(instanceID, status)
		if err != nil {
			return diag.FromErr(err)
		}

		d.Set("status", status)
	}

	return diags
}

func resourceInstanceDelete(tx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(*indigo.Client)
	instanceID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.DeleteInstance(instanceID)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")

	return diags
}

func setInstance(d *schema.ResourceData, instance *indigo.Instance) {
	d.Set("name", instance.InstanceName)
	d.Set("type_id", instance.InstanceTypeID)
	d.Set("set_no", instance.SetNo)
	d.Set("vps_kind", instance.VPSKind)
	d.Set("sequence_id", instance.SequenceID)
	d.Set("user_id", instance.UserID)
	d.Set("service_id", instance.ServiceID)
	d.Set("status", instance.Status)
	d.Set("ssh_key_id", strconv.Itoa(instance.SSHKeyID))
	d.Set("snapshot_id", instance.SnapshotID)
	d.Set("created_at", instance.CreatedAt)
	d.Set("start_date", instance.StartDate)
	d.Set("host_id", instance.HostID)
	d.Set("plan", instance.Plan)
	d.Set("plan_id", instance.PlanID)
	d.Set("disk_point", instance.DiskPoint)
	d.Set("mem_size", instance.MemSize)
	d.Set("cpus", instance.CPUs)
	d.Set("os_id", instance.OSID)
	d.Set("other_status", instance.OtherStatus)
	d.Set("uuid", instance.UUID)
	d.Set("uid_gid", instance.UIDGID)
	d.Set("vnc_port", instance.VNCPort)
	d.Set("vnc_passwd", instance.VNCPasswd)
	d.Set("arpa_name", instance.ARPAName)
	d.Set("arpa_date", instance.ARPADate)
	d.Set("started_at", instance.StartedAt)
	d.Set("closed_at", instance.ClosedAt)
	d.Set("status_change_date", instance.StatusChangeDate)
	d.Set("updated_at", instance.UpdatedAt)
	d.Set("vm_revert", instance.VMRevert)
	d.Set("ip_address", instance.IPAddress)
	d.Set("mac_address", instance.MACAddress)
	d.Set("import_instance", instance.ImportInstance)
	d.Set("container_id", instance.ContainerID)
	d.Set("daemon_status", instance.DaemonStatus)
	d.Set("out_of_stock", instance.OutOfStock)
	d.Set("ip_address_type", instance.IPAddressType)
	d.Set("ve_id", instance.VEID)
	d.Set("os", flattenOS(instance.OS))
	d.SetId(strconv.Itoa(instance.ID))
}

func flattenOS(os *indigo.OS) map[string]interface{} {
	c := make(map[string]interface{})
	if os != nil {
		c["name"] = os.Name
		c["view_name"] = os.ViewName
	}
	return c
}
