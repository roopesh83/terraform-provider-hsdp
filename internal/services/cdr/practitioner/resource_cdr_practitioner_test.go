package practitioner_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/philips-software/terraform-provider-hsdp/internal/acc"
)

func TestAccResourceCDRPractitioner_basic(t *testing.T) {
	t.Parallel()

	resourceName := "hsdp_cdr_practitioner.test"
	parentOrgID := acc.AccIAMOrgGUID()
	cdrURL := acc.AccCDRURL()
	now := time.Now().Format(time.RFC3339)

	randomNameSTU3 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	randomNameR4 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	// STU3
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: resourceName,
				Config:       testAccResourceCDRPractitioner(cdrURL, parentOrgID, randomNameSTU3, now, "stu3"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identifier.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "identifier.0.value", fmt.Sprintf("ron.swanson@%s.hsdp.io", randomNameSTU3)),
				),
			},
		},
	})

	// R4
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: resourceName,
				Config:       testAccResourceCDRPractitioner(cdrURL, parentOrgID, randomNameR4, now, "r4"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identifier.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "identifier.0.value", fmt.Sprintf("ron.swanson@%s.hsdp.io", randomNameR4)),
				),
			},
		},
	})
}

func testAccResourceCDRPractitioner(cdrURL, parentOrgID, name, now, version string) string {
	return fmt.Sprintf(`

data "hsdp_cdr_fhir_store" "sandbox" {
  base_url = "%s"
  fhir_org_id = hsdp_iam_org.test.id
}

data "hsdp_iam_introspect" "terraform" {
}

data "hsdp_iam_service" "service" {
  service_id = data.hsdp_iam_introspect.terraform.subject
}

resource "hsdp_iam_org" "test" {
  name  = "%s"
  description = "Acceptance Test CDR %s %s"
  parent_org_id = "%s"
}

resource "hsdp_iam_role" "cdr_admin" {
  managing_organization = hsdp_iam_org.test.id
  name                  = "TF_CDR_ADMIN"
  permissions = [
    "ALL.READ",
    "ALL.WRITE"
  ]
}

resource "hsdp_iam_group" "cdr_admins" {
  managing_organization = hsdp_iam_org.test.id
  name                  = "TF_CDR_ADMIN"
  description           = "CDR Admins"
  roles                 = [hsdp_iam_role.cdr_admin.id]
  users                 = []
  services              = [data.hsdp_iam_service.service.id]
}

resource "hsdp_cdr_org" "test" {
  fhir_store  = data.hsdp_cdr_fhir_store.sandbox.endpoint

  name        = "Practitioner Resource Test %s"
  org_id      = hsdp_iam_org.test.id

  version     = "%s"

  depends_on = [hsdp_iam_group.cdr_admins]
}

resource "hsdp_cdr_practitioner" "test" {
  fhir_store  = hsdp_cdr_org.test.fhir_store

  identifier {
    system = "https://provider.terrakube.com"
    value  = "ron.swanson@%s.hsdp.io"
    use    = "usual"
  }

  name {
     text   = "Ron Swanson"
     family = "Swanson"
     given  = ["Ron", "%s"]
  }
 version = "%s"

 soft_delete = true

 depends_on = [hsdp_iam_group.cdr_admins]
}`,
		// DATA SOURCE
		cdrURL,

		// IAM ORG
		name,
		now,
		name,
		parentOrgID,

		// IAM GROUP

		// CDR ORG
		name,
		version,

		// CDR PRACTITIONER
		name,
		name,
		version,
	)
}
