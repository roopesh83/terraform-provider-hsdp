package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserRead,
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"uuid": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "use the id field",
			},
			"email_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

}

func dataSourceUserRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)

	var diags diag.Diagnostics

	client, err := c.IAMClient()
	if err != nil {
		return diag.FromErr(err)
	}

	username := d.Get("username").(string)

	uuid, _, err := client.Users.GetUserIDByLoginID(username)

	if err != nil {
		// Fallback to legacy user find
		uuid, _, err = client.Users.LegacyGetUserIDByLoginID(username)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "user not found",
				Detail:   fmt.Sprintf("user '%s' not found or you don't have 'USER.READ' permission in this users' IAM hierarchy", username),
			})
			d.SetId(fmt.Sprintf("%s-404", username))
			_ = d.Set("uuid", "")
			_ = d.Set("email_address", "")
			return diags
		}
	}
	user, _, err := client.Users.LegacyGetUserByUUID(uuid)
	if err == nil {
		_ = d.Set("email_address", user.Contact.EmailAddress)
	}

	d.SetId(uuid)
	_ = d.Set("uuid", uuid)

	return diags
}
