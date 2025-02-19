package notification_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/philips-software/terraform-provider-hsdp/internal/acc"
)

func TestAccResourceNotificationSubscriber_basic(t *testing.T) {
	t.Parallel()

	resourceName := "hsdp_notification_subscriber.subscriber"
	randomName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	iamOrgID := acc.AccIAMOrgGUID()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acc.PreCheck(t) },
		ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNotificationSubscriber(randomName, iamOrgID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", fmt.Sprintf("acc subscriber %s", randomName)),
				),
			},
		},
	})
}

func testAccResourceNotificationSubscriber(random, parentId string) string {
	return fmt.Sprintf(`
resource "hsdp_iam_org" "test" {
    name = "ACCTest-%s"
    description = "ACCResourceNotificationProducer %s"

	parent_org_id = "%s"
    wait_for_delete = false
}

resource "hsdp_iam_proposition" "test" {
   name = "ACC-%s"
   description = "IAM Service Test %s"
   
   organization_id = hsdp_iam_org.test.id
}

resource "hsdp_iam_application" "test" {
    name = "ACC-%s"
    description = "IAM Service Test %s"
    proposition_id = hsdp_iam_proposition.test.id
}

resource "hsdp_iam_service" "test" {
  name        = "ACC-%s"
  description = "IAM Service Test %s"

  validity = 12
  token_validity = 3600
  scopes = ["openid"]
  default_scopes = ["openid"]
 
  application_id = hsdp_iam_application.test.id
}

resource "hsdp_iam_role" "producer_admin" {
  name = "PRODUCER_ADMIN_TF"
  permissions = [
    "ORGANIZATION.READ",
    "NS_PRODUCER.CREATE",
    "NS_PRODUCER.READ",
    "NS_PRODUCER.DELETE",
    "NS_TOPIC_SCOPE.READ",
    "NS_TOPIC.CREATE",
    "NS_TOPIC.READ",
    "NS_TOPIC.UPDATE",
    "NS_TOPIC.DELETE",
    "NS_SUBSCRIBER.READ",
    "NS_SUBSCRIPTION.READ"
  ]
  managing_organization = hsdp_iam_org.test.id
}

resource "hsdp_iam_role" "publisher" {
  name = "PUBLISHER_TF"
  permissions = [
    "NS_PUBLISH.CREATE",
    "NS_TOPIC.READ"
  ]
  managing_organization = hsdp_iam_org.test.id
}

resource "hsdp_iam_role" "subscriber_admin" {
  name = "SUBSCRIBER_ADMIN_TF"
  permissions = [
    "ORGANIZATION.READ",
    "NS_SUBSCRIBER.CREATE",
    "NS_SUBSCRIBER.DELETE",
    "NS_PRODUCER.READ",
    "NS_SUBSCRIPTION.READ"
  ]
  managing_organization = hsdp_iam_org.test.id
}

resource "hsdp_iam_role" "subscriber" {
  name = "SUBSCRIBER_TF"
  permissions = [
    "NS_SUBSCRIPTION.CREATE",
    "NS_SUBSCRIPTION.READ",
    "NS_SUBSCRIPTION.DELETE",
    "NS_SUBSCRIPTION.CONFIRM",
    "NS_SUBSCRIPTION.SYNC",
    "NS_TOPIC.READ"
  ]
  managing_organization = hsdp_iam_org.test.id
}

resource "hsdp_iam_group" "producer_admins" {
  name                  = "PRODUCER_ADMINS_TF"
  roles                 = [hsdp_iam_role.producer_admin.id]
  services              = [hsdp_iam_service.test.id]
  users                 = []
  managing_organization = hsdp_iam_org.test.id
}

resource "hsdp_iam_group" "publishers" {
  name                  = "PUBLISHERS_TF"
  roles                 = [hsdp_iam_role.publisher.id]
  users                 = []
  services              = [hsdp_iam_service.test.id]
  managing_organization = hsdp_iam_org.test.id
}

resource "hsdp_iam_group" "subscriber_admins" {
  name                  = "SUBSCRIBER_ADMINS_TF"
  roles                 = [hsdp_iam_role.subscriber_admin.id]
  users                 = []
  services              = [hsdp_iam_service.test.id]
  managing_organization = hsdp_iam_org.test.id
}

resource "hsdp_iam_group" "subscribers" {
  name                  = "SUBSCRIBERS_TF"
  roles                 = [hsdp_iam_role.subscriber.id]
  users                 = []
  services              = [hsdp_iam_service.test.id]
  managing_organization = hsdp_iam_org.test.id
}

resource "hsdp_notification_subscriber" "subscriber" {
  principal {
    service_id          = hsdp_iam_service.test.service_id
    service_private_key = hsdp_iam_service.test.private_key
  }

  managing_organization_id       = hsdp_iam_org.test.id
  managing_organization          = hsdp_iam_org.test.name

  subscriber_product_name          = "subscriberProd%s"
  subscriber_service_name          = "subscriberService%s"
  subscriber_service_instance_name = "userEvent-%s"
  subscriber_service_base_url      = "https://notification-receiver.cloud.pcftest.com/"
  subscriber_service_path_url      = "notification/UserEvent-%s"
  description                      = "acc subscriber %s"

  depends_on = [hsdp_iam_group.subscriber_admins]
}
`,
		// IAM TEST ORG
		random,
		random,
		parentId,

		// IAM Proposition
		strings.ToUpper(random),
		random,

		// IAM Application
		strings.ToUpper(random),
		random,

		// IAM Service
		strings.ToUpper(random),
		random,

		// IAM GROUP producer_admins

		// IAM GROUP publishers

		// IAM GROUP subscriber_admins

		// IAM GROUP subscribers

		// NOTIFICATION PRODUCER producer
		random,
		random,
		random,
		random,
		random,
	)
}
