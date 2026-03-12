package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// You normally want to run this under a separate "Testing" subscription
// For lab purposes you will use your assigned subscription under the Cloud Dev/Ops program tenant
var subscriptionID string = "d41fda54-dd2c-48eb-9aaf-c0135dc1dfd3"

func TestAzureLinuxVMCreation(t *testing.T) {
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../",
		// Override the default terraform variables
		Vars: map[string]interface{}{
			"labelPrefix": "dunn0203",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of output variable
	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	nicName := terraform.Output(t, terraformOptions, "nic_name")

	// Confirm VM exists
	assert.True(t, azure.VirtualMachineExists(t, vmName, resourceGroupName, subscriptionID))

	// Confirm NIC exists
	assert.True(t, azure.NetworkInterfaceExists(t, nicName, resourceGroupName, subscriptionID))

	// Confirm NIC is attached to the VM
	vm := azure.GetVirtualMachine(t, vmName, resourceGroupName, subscriptionID)

	nics := *vm.NetworkProfile.NetworkInterfaces
	assert.Contains(t, *nics[0].ID, nicName)

	// Confirm VM is running Ubuntu
	imagePublisher := *vm.StorageProfile.ImageReference.Publisher
	imageOffer := *vm.StorageProfile.ImageReference.Offer
	imageSku := *vm.StorageProfile.ImageReference.Sku

	assert.Equal(t, "Canonical", imagePublisher)
	assert.Contains(t, imageOffer, "ubuntu")
	assert.Contains(t, imageSku, "22_04")
}
