# private branch for Spring Cloud Service Virtual Network Integration

## Instruction
This feature is based on azurerm provider 2.32.0

to build a private binary, you should 
- checkout this repo and switch to current branch
- `make build`

If you want to cross compile for other OS platform, you could change the file: `GNUmakefile`. In the line 23, add cross compile parameters

## how to use private build binary
please See https://www.terraform.io/docs/extend/how-terraform-works.html#plugin-locations for more details

## verify
The way to verify you are using the private build version is to check the debug log of terraform

- you shall see: `[DEBUG] plugin: starting plugin: path=<plugin location>/terraform-provider-azurerm`
- debug log should contain following information
```
[INFO] *************************************************************************************
[INFO] ***** Using Terraform Azure provider with Spring Cloud Service Vnet Integration *****
[INFO] *************************************************************************************
```

## sample usage (spring cloud vnet, private dns zone and augment route table rules)
```hcl
provider "azurerm" {
  features {}
}

data "azurerm_resource_group" "test" {
  name = "cz"
}

resource "azurerm_virtual_network" "test" {
  name                = "testvnet"
  address_space       = ["10.1.0.0/16"]
  location            = data.azurerm_resource_group.test.location
  resource_group_name = data.azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test1" {
  name                 = "internal1"
  resource_group_name  = data.azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.1.0.0/24"
}

resource "azurerm_subnet" "test2" {
  name                 = "internal2"
  resource_group_name  = data.azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.1.1.0/24"
}

data "azuread_service_principal" "test" {
  display_name = "Azure Spring Cloud Resource Provider"
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_virtual_network.test.id
  role_definition_name = "Owner"
  principal_id         = data.azuread_service_principal.test.object_id
}

resource "azurerm_application_insights" "test" {
  name                = "tf-test-appinsights"
  location            = data.azurerm_resource_group.test.location
  resource_group_name = data.azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_spring_cloud_service" "test" {
  name                = "sc-cz"
  resource_group_name = data.azurerm_resource_group.test.name
  location            = data.azurerm_resource_group.test.location
  
  network {
    app_subnet_id             = azurerm_subnet.test1.id
    service_runtime_subnet_id = azurerm_subnet.test2.id
    cidr_ranges               = ["10.4.0.0/16", "10.5.0.0/16", "10.3.0.1/16"]
  }

  trace {
    instrumentation_key = azurerm_application_insights.test.instrumentation_key
  }

  depends_on = [azurerm_role_assignment.test]
}

output "private_address" {
  value = azurerm_spring_cloud_service.test.service_runtime_lb_private_address
}

output "app_route_table_id" {
  value = azurerm_spring_cloud_service.test.app_subnet_route_table_id
}

output "app_route_table_name" {
  value = azurerm_spring_cloud_service.test.app_subnet_route_table_name
}

output "app_route_table_resource_group" {
  value = azurerm_spring_cloud_service.test.app_subnet_route_table_resource_group
}

output "service_runtime_route_table_id" {
  value = azurerm_spring_cloud_service.test.service_runtime_subnet_route_table_id
}

output "service_runtime_route_table_name" {
  value = azurerm_spring_cloud_service.test.service_runtime_subnet_route_table_name
}

output "service_runtime_route_table_resource_group" {
  value = azurerm_spring_cloud_service.test.service_runtime_subnet_route_table_resource_group
}

resource "azurerm_private_dns_zone" "test" {
  name                = "private.azuremicroservices.io"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_dns_zone_virtual_network_link" "test" {
  name                  = "azure-spring-cloud-dns-link"
  resource_group_name   = azurerm_resource_group.test.name
  private_dns_zone_name = azurerm_private_dns_zone.test.name
  virtual_network_id    = azurerm_virtual_network.test.id
}

resource "azurerm_private_dns_a_record" "test" {
  name                = "*"
  zone_name           = azurerm_private_dns_zone.test.name
  resource_group_name = azurerm_resource_group.test.name
  ttl                 = 300
  records             = [azurerm_spring_cloud_service.test.service_runtime_lb_private_address]
}

resource "azurerm_route" "test" {
  name                = "test"
  resource_group_name = azurerm_spring_cloud_service.test.app_subnet_route_table_resource_group
  route_table_name    = azurerm_spring_cloud_service.test.app_subnet_route_table_name
  address_prefix      = "10.1.0.0/16"
  next_hop_type       = "vnetlocal"
}
```

## sample usage (create app, deployment, active_deployment)
```
resource "azurerm_spring_cloud_app" "test" {
  name                = "app1"
  resource_group_name = azurerm_spring_cloud_service.test.resource_group_name
  service_name        = azurerm_spring_cloud_service.test.name
  is_public           = true
  https_only          = true

  identity {
    type = "SystemAssigned"
  }

  persistent_disk {
    size_in_gb = 50
    mount_path = "/persistent"
  }
}

resource "azurerm_spring_cloud_java_deployment" "test" {
  name                = "deploy1"
  spring_cloud_app_id = azurerm_spring_cloud_app.test.id
  cpu                 = 2
  memory_in_gb        = 4
  instance_count      = 2
  jvm_options         = "-XX:+PrintGC"
  runtime_version     = "Java_8"

  env = {
    "Env" : "Staging"
  }
}

resource "azurerm_spring_cloud_active_deployment" "test" {
  spring_cloud_app_id = azurerm_spring_cloud_app.test.id
  deployment_name     = azurerm_spring_cloud_java_deployment.test.name
}
```

## sample usage (app access key vault secret)
```
data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                = "key-vault-test-cz"
  location            = data.azurerm_resource_group.test.location
  resource_group_name = data.azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  secret_permissions      = ["get", "set", "delete", "list"]
  certificate_permissions = ["create", "delete", "get", "update", "list"]
}

resource "azurerm_key_vault_access_policy" "test1" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_spring_cloud_app.test.identity.0.principal_id

  secret_permissions = ["get", "list"]
}
```

## sample usage (spring cloud certificate)
```
data "azurerm_client_config" "current" {}

data "azuread_service_principal" "test" {
  display_name = "Azure Spring Cloud Domain-Management"
}

resource "azurerm_key_vault" "test" {
  name                = "key-vault-test-cz"
  location            = data.azurerm_resource_group.test.location
  resource_group_name = data.azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"
}

resource "azurerm_key_vault_access_policy" "test" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  secret_permissions      = ["get", "set", "delete", "list"]
  certificate_permissions = ["create", "delete", "get", "update", "list"]
}

resource "azurerm_key_vault_access_policy" "test1" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azuread_service_principal.test.object_id

  secret_permissions = ["get", "list"]
  certificate_permissions = ["get", "list"]
}

resource "azurerm_key_vault_certificate" "test" {
  name         = "acctestcertcz"
  key_vault_id = azurerm_key_vault.test.id

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = true
    }

    lifetime_action {
      action {
        action_type = "AutoRenew"
      }

      trigger {
        days_before_expiry = 30
      }
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }

    x509_certificate_properties {
      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]

      subject            = "CN=*.azdmss-test.net"
      subject_alternative_names {
        dns_names = ["*.azdmss-test.net"]
      }
      validity_in_months = 12
    }
  }
}

resource "azurerm_spring_cloud_certificate" "test" {
  name                     = "acctest-scc-cz"
  resource_group_name      = azurerm_spring_cloud_service.test.resource_group_name
  service_name             = azurerm_spring_cloud_service.test.name
  key_vault_certificate_id = azurerm_key_vault_certificate.test.id
}
```

## sample usage (custom domain)
```
resource "azurerm_spring_cloud_custom_domain" "test" {
  name                = "tf-test.azdmss-test.net"
  spring_cloud_app_id = azurerm_spring_cloud_app.test.id
  cert_name           = azurerm_spring_cloud_certificate.test.name
  thumbprint          = azurerm_spring_cloud_certificate.test.thumbprint
}
```

## sample usage (autoscale)
```
resource "azurerm_monitor_autoscale_setting" "test" {
  name                = "acctestautoscale-cz"
  resource_group_name = data.azurerm_resource_group.test.name
  location            = data.azurerm_resource_group.test.location
  target_resource_id  = azurerm_spring_cloud_java_deployment.test.id
  enabled             = true
  profile {
    name = "metricRules"
    capacity {
      default = 1
      minimum = 1
      maximum = 2
    }
    rule {
      metric_trigger {
        dimensions {
          name     = "AppName"
          operator = "Equals"
          values   = [azurerm_spring_cloud_app.test.name]
        }

        dimensions {
          name     = "Deployment"
          operator = "Equals"
          values   = [azurerm_spring_cloud_java_deployment.test.name]
        }

        metric_name        = "AppCpuUsage"
        metric_namespace   = "microsoft.appplatform/spring"
        metric_resource_id = azurerm_spring_cloud_service.test.id
        time_grain         = "PT1M"
        statistic          = "Average"
        time_window        = "PT5M"
        time_aggregation   = "Average"
        operator           = "GreaterThan"
        threshold          = 75
      }
      scale_action {
        direction = "Increase"
        type      = "ChangeCount"
        value     = 1
        cooldown  = "PT1M"
      }
    }
  }
}
```

## sample usage (config server)
```
resource "azurerm_spring_cloud_config_server" "test" {
  spring_cloud_service_id = azurerm_spring_cloud_service.test.id
  uri                     = "https://github.com/Azure-Samples/piggymetrics"
  label                   = "config"
  search_paths            = ["dir1", "dir4"]

  repository {
    name         = "repo1"
    uri          = "https://github.com/Azure-Samples/piggymetrics"
    label        = "config"
    search_paths = ["dir1", "dir2"]
    http_basic_auth {
      username = "username"
      password = "password"
    }
  }

  repository {
    name         = "repo2"
    uri          = "git@bitbucket.org:Azure-Samples/piggymetrics.git"
    label        = "config"
    search_paths = ["dir1", "dir2"]

    ssh_auth {
      private_key                      = file("testdata/private_key")
      host_key                         = file("testdata/host_key")
      host_key_algorithm               = "ssh-rsa"
      strict_host_key_checking_enabled = false
    }
  }
}
```

# Terraform Provider for Azure (Resource Manager)

Version 2.x of the AzureRM Provider requires Terraform 0.12.x and later.

* [Terraform Website](https://www.terraform.io)
* [AzureRM Provider Documentation](https://www.terraform.io/docs/providers/azurerm/index.html)
* [AzureRM Provider Usage Examples](https://github.com/terraform-providers/terraform-provider-azurerm/tree/master/examples)
* [Slack Workspace for Contributors](https://terraform-azure.slack.com) ([Request Invite](https://join.slack.com/t/terraform-azure/shared_invite/enQtNDMzNjQ5NzcxMDc3LWNiY2ZhNThhNDgzNmY0MTM0N2MwZjE4ZGU0MjcxYjUyMzRmN2E5NjZhZmQ0ZTA1OTExMGNjYzA4ZDkwZDYxNDE))

## Usage Example

```
# Configure the Microsoft Azure Provider
provider "azurerm" {
  # We recommend pinning to the specific version of the Azure Provider you're using
  # since new versions are released frequently
  version = "=2.40.0"

  features {}

  # More information on the authentication methods supported by
  # the AzureRM Provider can be found here:
  # https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs

  # subscription_id = "..."
  # client_id       = "..."
  # client_secret   = "..."
  # tenant_id       = "..."
}

# Create a resource group
resource "azurerm_resource_group" "example" {
  name     = "production-resources"
  location = "West US"
}

# Create a virtual network in the production-resources resource group
resource "azurerm_virtual_network" "test" {
  name                = "production-network"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  address_space       = ["10.0.0.0/16"]
}
```

Further [usage documentation is available on the Terraform website](https://www.terraform.io/docs/providers/azurerm/index.html).

## Developer Requirements

* [Terraform](https://www.terraform.io/downloads.html) version 0.12.x +
* [Go](https://golang.org/doc/install) version 1.15.x (to build the provider plugin)

### On Windows

If you're on Windows you'll also need:
* [Git Bash for Windows](https://git-scm.com/download/win)
* [Make for Windows](http://gnuwin32.sourceforge.net/packages/make.htm)

For *GNU32 Make*, make sure its bin path is added to PATH environment variable.*

For *Git Bash for Windows*, at the step of "Adjusting your PATH environment", please choose "Use Git and optional Unix tools from Windows Command Prompt".*

Or install via [Chocolatey](https://chocolatey.org/install) (`Git Bash for Windows` must be installed per steps above)
```powershell
choco install make golang terraform -y
refreshenv
```

You must run  `Developing the Provider` commands in `bash` because `sh` scrips are invoked as part of these.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.15+ is **required**). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

First clone the repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-azurerm`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/terraform-providers
$ git clone git@github.com:terraform-providers/terraform-provider-azurerm
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-azurerm
```

Once inside the provider directory, you can run `make tools` to install the dependent tooling required to compile the provider.

At this point you can compile the provider by running `make build`, which will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-azurerm
...
```

You can also cross-compile if necessary:

```sh
GOOS=windows GOARCH=amd64 make build
```

In order to run the `Unit Tests` for the provider, you can run:

```sh
$ make test
```

The majority of tests in the provider are `Acceptance Tests` - which provisions real resources in Azure. It's possible to run the entire acceptance test suite by running `make testacc` - however it's likely you'll want to run a subset, which you can do using a prefix, by running:

```sh
make acctests SERVICE='resource' TESTARGS='-run=TestAccAzureRMResourceGroup' TESTTIMEOUT='60m'
```

The following Environment Variables must be set in your shell prior to running acceptance tests:

- `ARM_CLIENT_ID`
- `ARM_CLIENT_SECRET`
- `ARM_SUBSCRIPTION_ID`
- `ARM_TENANT_ID`
- `ARM_ENVIRONMENT`
- `ARM_METADATA_HOST`
- `ARM_TEST_LOCATION`
- `ARM_TEST_LOCATION_ALT`
- `ARM_TEST_LOCATION_ALT2`

**Note:** Acceptance tests create real resources in Azure which often cost money to run.

---

## Developer: Using the locally compiled Azure Provider binary

When using Terraform 0.14 and later, after successfully compiling the Azure Provider, you must [instruct Terraform to use your locally compiled provider binary](https://www.terraform.io/docs/commands/cli-config.html#development-overrides-for-provider-developers) instead of the official binary from the Terraform Registry.

For example, add the following to `~/.terraformrc` for a provider binary located in `/home/developer/go/bin`:

```hcl
provider_installation {

  # Use /home/developer/go/bin as an overridden package directory
  # for the hashicorp/azurerm provider. This disables the version and checksum
  # verifications for this provider and forces Terraform to look for the
  # azurerm provider plugin in the given directory.
  dev_overrides {
    "hashicorp/azurerm" = "/home/developer/go/bin"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

---

## Developer: Generating Resource ID Formatters, Parsers and Validators

You can generate a Resource ID Formatter, Parser and Validator by adding the following line to a `resourceids.go` within each Service Package (for example `./azurerm/internal/services/someservice/resourceids.go`):

```go
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Server -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AnalysisServices/servers/Server1
```

Where `name` is the name of the Resource ID Type - and `id` is an example Resource ID with placeholder data.

When `make generate` is run, this will then generate the following for this Resource ID:

* Resource ID Struct, containing the fields and a Formatter to convert this into a string - and the associated Unit Tests.
* Resource ID Parser (`./parse/{name}.go`) - to be able to parse a Resource ID into said struct - and the associated Unit Tests.
* Resource ID Validator (`./validate/{name}_id.go`) - to validate the Resource ID is what's expected (and not for a different resource) - and the associated Unit Tests.

---

## Developer: Scaffolding the Website Documentation

You can scaffold the documentation for a Data Source by running:

```sh
$ make scaffold-website BRAND_NAME="Resource Group" RESOURCE_NAME="azurerm_resource_group" RESOURCE_TYPE="data"
```

You can scaffold the documentation for a Resource by running:

```sh
$ make scaffold-website BRAND_NAME="Resource Group" RESOURCE_NAME="azurerm_resource_group" RESOURCE_TYPE="resource" RESOURCE_ID="/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1"
```
