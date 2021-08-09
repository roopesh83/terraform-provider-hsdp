package hsdp

import (
	"context"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/inference"
)

func resourceInferenceComputeTarget() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CreateContext: resourceInferenceComputeTargetCreate,
		ReadContext:   resourceInferenceComputeTargetRead,
		DeleteContext: resourceInferenceComputeTargetDelete,

		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"storage": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"is_factory": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceInferenceComputeTargetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	endpoint := d.Get("endpoint").(string)
	client, err := config.getInferenceClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	instanceType := d.Get("instance_type").(string)
	storage := d.Get("storage").(int)

	var createdTarget *inference.ComputeTarget
	// Do initial boarding
	operation := func() error {
		var resp *inference.Response
		createdTarget, resp, err = client.ComputeTarget.CreateComputeTarget(inference.ComputeTarget{
			ResourceType: "ComputeEnvironment",
			Name:         name,
			Description:  description,
			InstanceType: instanceType,
			Storage:      storage,
		})
		if resp == nil {
			resp = &inference.Response{}
		}
		return checkForIAMPermissionErrors(client, resp.Response, err)
	}
	err = backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 8))
	if err != nil {
		return diag.FromErr(err)
	}

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(createdTarget.ID)
	return resourceInferenceComputeTargetRead(ctx, d, m)
}

func resourceInferenceComputeTargetRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(*Config)
	endpoint := d.Get("endpoint").(string)
	client, err := config.getInferenceClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	id := d.Id()

	target, _, err := client.ComputeTarget.GetComputeTargetByID(id)
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("name", target.Name)
	_ = d.Set("description", target.Description)
	_ = d.Set("instance_type", target.InstanceType)
	_ = d.Set("storage", target.Storage)
	_ = d.Set("is_factory", target.IsFactory)
	_ = d.Set("created", target.Created)
	_ = d.Set("created_by", target.CreatedBy)

	return diags
}

func resourceInferenceComputeTargetDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(*Config)
	endpoint := d.Get("endpoint").(string)
	client, err := config.getInferenceClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	id := d.Id()

	_, err = client.ComputeTarget.DeleteComputeTarget(inference.ComputeTarget{
		ID: id,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
