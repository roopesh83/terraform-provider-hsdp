package hsdp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/cdl"
)

func dataSourceCDLResearchStudy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCDLResearchStudyRead,
		Schema: map[string]*schema.Schema{
			"study_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cdl_endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"title": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ends_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"study_owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"uploaders": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"monitors": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"data_scientists": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"study_managers": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}

}

func dataSourceCDLResearchStudyRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)

	var diags diag.Diagnostics

	endpoint := d.Get("cdl_endpoint").(string)
	studyID := d.Get("study_id").(string)

	client, err := config.getCDLClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	study, resp, err := client.Study.GetStudyByID(studyID)
	if err != nil {
		return diag.FromErr(err)
	}
	if resp.StatusCode == http.StatusForbidden {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  fmt.Sprintf("permission denied ready study %s", studyID),
			Detail:   "not enough permissions to read study details",
		})
	}

	permissions, resp, err := client.Study.GetPermissions(cdl.Study{ID: studyID}, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	if resp.StatusCode == http.StatusForbidden {
		return diag.FromErr(fmt.Errorf("permission denied reading study %s", studyID))
	}
	var monitors []string
	var uploaders []string
	var dataScientists []string
	var studyManagers []string
	for _, p := range permissions {
		for _, r := range p.Roles {
			switch r.Role {
			case cdl.ROLE_DATA_SCIENTIST:
				dataScientists = append(dataScientists, p.IAMUserUUID)
			case cdl.ROLE_STUDYMANAGER:
				studyManagers = append(studyManagers, p.IAMUserUUID)
			case cdl.ROLE_UPLOADER:
				uploaders = append(uploaders, p.IAMUserUUID)
			case cdl.ROLE_MONITOR:
				monitors = append(monitors, p.IAMUserUUID)
			}
		}
	}

	d.SetId(study.ID)
	_ = d.Set("title", study.Title)
	_ = d.Set("description", study.Description)
	_ = d.Set("ends_at", study.Period.End)
	_ = d.Set("study_owner", study.StudyOwner)

	_ = d.Set("monitors", monitors)
	_ = d.Set("data_scientists", dataScientists)
	_ = d.Set("uploaders", uploaders)
	_ = d.Set("study_managers", studyManagers)

	return diags
}