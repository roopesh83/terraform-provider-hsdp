---
subcategory: "Identity and Access Management (IAM)"
---

# hsdp_iam_service

Provides a resource for managing HSDP IAM services of an application under a proposition.

## Example Usage

The following example creates a service

```hcl
resource "hsdp_iam_service" "testservice" {
  name                = "TESTSERVICE"
  description         = "Test service"
  application_id      = var.app_id

  validity            = 12   # Months
  
  token_validity      = 3600 # Seconds

  scopes              = ["openid"]
  default_scopes      = ["openid"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the service
* `description` - (Required) The description of the service
* `application_id` - (Required) the application ID (GUID) to attach this service to
* `scopes` - (Required) Array. List of supported scopes for this service. Minimum: ["openid"]
* `validity` - (Optional) Integer. Validity of service (in months). Minimum: 1, Maximum: 600 (5 years), Default: 12
* `token_validity` - (Optional) Integer. Access Token Lifetime (in seconds). Default: 1800 (30 minutes), Maximum: 2592000 (30 days)
* `default_scopes` - (Required) Array. Default scopes. You do not have to specify these explicitly when requesting a token. Minimum: ["openid"]
* `self_managed_private_key` - (Optional)  RSA private key in PEM format. When provided, overrides the generated certificate / private key combination of the
  IAM service. This gives you full control over the credentials. When not specified, a private key will be generated by IAM
* `self_managed_expires_on` - (Optional) Sets the certificate validity. When not specified, the certificate will have a validity of 5 years.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The GUID of the client
* `service_id` - (Generated) The service id
* `private_key` - (Generated) The active private of the service
* `expires_on` - (Generated) Sets the certificate validity. When not specified, the certificate will have a validity of 5 years.
* `organization_id` - The organization ID this service belongs to (via application and proposition)

## Import

Existing services can be imported, however they will be missing their private key rendering them pretty much useless. Therefore, we recommend creating them using the provider.
