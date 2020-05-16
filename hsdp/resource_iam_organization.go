package hsdp

import (
	"errors"
	"fmt"
	"github.com/philips-software/go-hsdp-api/iam"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIAMOrg() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Create: resourceIAMOrgCreate,
		Read:   resourceIAMOrgRead,
		Update: resourceIAMOrgUpdate,
		Delete: resourceIAMOrgDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"distinct_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"org_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_root_org": &schema.Schema{
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"parent_org_id"},
			},
			"parent_org_id": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"is_root_org"},
			},
		},
	}
}

func resourceIAMOrgCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client := config.IAMClient()

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	isRootOrg := d.Get("is_root_org").(bool)
	if isRootOrg {
		return errors.New("cannot create root orgs")
	}
	parentOrgID, ok := d.Get("parent_org_id").(string)
	if !ok {
		return errors.New("non root orgs must specify a `parent_org_id`")
	}
	var newOrg iam.Organization
	newOrg.Name = name
	newOrg.Description = description
	newOrg.Parent.Value = parentOrgID
	org, resp, err := client.Organizations.CreateOrganization(newOrg)
	if err != nil {
		return err
	}
	if org == nil {
		return fmt.Errorf("failed to create organization: %d", resp.StatusCode)
	}
	d.SetId(org.ID)
	return nil
}

func resourceIAMOrgRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client := config.IAMClient()

	id := d.Id()
	org, resp, err := client.Organizations.GetOrganizationByID(id)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return err
	}
	d.Set("org_id", org.ID)
	d.Set("description", org.Description)
	d.Set("name", org.Name)
	return nil
}

func resourceIAMOrgUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	client := config.IAMClient()

	id := d.Id()
	org, _, err := client.Organizations.GetOrganizationByID(id)
	if err != nil {
		return err
	}

	if d.HasChange("description") {
		description := d.Get("description").(string)
		org.Description = description
		_, _, err = client.Organizations.UpdateOrganization(*org)
	}
	return err
}

func resourceIAMOrgDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
