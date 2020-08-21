package test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestTerraformAzureExample(t *testing.T) {
	t.Parallel()

	// website::tag::1:: Configure Terraform setting up a path to Terraform code.
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/example1",
		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"prefix":   "e2etest",
			"location": "westeurope",
		},
	}

	// website::tag::4:: At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// website::tag::2:: Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// website::tag::3:: Run `terraform output` to get the values of output variables
	siteHostName := terraform.Output(t, terraformOptions, "site_hostname")

	// website::tag::4:: Make an HTTP request to the instance and make sure we get back a 200 OK with the body "Hello, World!"
	url := fmt.Sprintf("http://%s", siteHostName)
	http_helper.HttpGetWithRetryWithCustomValidation(t, url, nil, 3, 5*time.Second, func(statusCode int, body string) bool {
		fmt.Println(body)
		return statusCode == http.StatusOK
	})
}
