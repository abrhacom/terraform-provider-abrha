package vm_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/abrhacom/terraform-provider-abrha/abrha/acceptance"
	"github.com/abrhacom/terraform-provider-abrha/abrha/util"
	"github.com/abrhacom/terraform-provider-abrha/abrha/vm"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	defaultSize  = "s-1vcpu-1gb"
	defaultImage = "ubuntu-22-04-x64"

	gpuSize      = "gpu-h100x1-80gb"
	gpuImage     = "gpu-h100x1-base"
	runGPUEnvVar = "DO_RUN_GPU_TESTS"
)

func TestAccAbrhaVm_Basic(t *testing.T) {
	var vm goApiAbrha.Vm
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckAbrhaVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: acceptance.TestAccCheckAbrhaVmConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
					testAccCheckAbrhaVmAttributes(&vm),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "size", defaultSize),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "price_hourly", "0.00893"),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "price_monthly", "6"),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "image", defaultImage),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "user_data", util.HashString("foobar")),
					resource.TestCheckResourceAttrSet(
						"abrha_vm.foobar", "ipv4_address_private"),
					resource.TestCheckResourceAttrSet(
						"abrha_vm.foobar", "vpc_uuid"),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "ipv6_address", ""),
					resource.TestCheckResourceAttrSet("abrha_vm.foobar", "urn"),
					resource.TestCheckResourceAttrSet("abrha_vm.foobar", "created_at"),
				),
			},
		},
	})
}

func TestAccAbrhaVm_WithID(t *testing.T) {
	var vm goApiAbrha.Vm
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckAbrhaVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaVmConfig_withID(name, defaultImage),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
				),
			},
		},
	})
}

func TestAccAbrhaVm_withSSH(t *testing.T) {
	var vm goApiAbrha.Vm
	name := acceptance.RandomTestName()
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("abrha@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckAbrhaVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaVmConfig_withSSH(name, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
					testAccCheckAbrhaVmAttributes(&vm),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "size", defaultSize),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "image", defaultImage),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "user_data", util.HashString("foobar")),
				),
			},
		},
	})
}

func TestAccAbrhaVm_Update(t *testing.T) {
	var vm goApiAbrha.Vm
	name := acceptance.RandomTestName()
	newName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckAbrhaVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: acceptance.TestAccCheckAbrhaVmConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
					testAccCheckAbrhaVmAttributes(&vm),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "name", name),
				),
			},

			{
				Config: testAccCheckAbrhaVmConfig_RenameAndResize(newName),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
					testAccCheckAbrhaVmRenamedAndResized(&vm),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "name", newName),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "size", "s-1vcpu-2gb"),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "disk", "50"),
				),
			},
		},
	})
}

func TestAccAbrhaVm_ResizeWithOutDisk(t *testing.T) {
	var vm goApiAbrha.Vm
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckAbrhaVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: acceptance.TestAccCheckAbrhaVmConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
					testAccCheckAbrhaVmAttributes(&vm),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "name", name),
				),
			},

			{
				Config: testAccCheckAbrhaVmConfig_resize_without_disk(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
					testAccCheckAbrhaVmResizeWithOutDisk(&vm),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "size", "s-1vcpu-2gb"),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "disk", "25"),
				),
			},
		},
	})
}

func TestAccAbrhaVm_ResizeSmaller(t *testing.T) {
	var vm goApiAbrha.Vm
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckAbrhaVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: acceptance.TestAccCheckAbrhaVmConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
					testAccCheckAbrhaVmAttributes(&vm),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "name", name),
				),
			},
			// Test moving to larger plan with resize_disk = false only increases RAM, not disk.
			{
				Config: testAccCheckAbrhaVmConfig_resize_without_disk(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
					testAccCheckAbrhaVmResizeWithOutDisk(&vm),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "size", "s-1vcpu-2gb"),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "disk", "25"),
				),
			},
			// Test that we can downgrade a Vm plan as long as the disk remains the same
			{
				Config: acceptance.TestAccCheckAbrhaVmConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
					testAccCheckAbrhaVmAttributes(&vm),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "size", "s-1vcpu-1gb"),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "disk", "25"),
				),
			},
			// Test that resizing resize_disk = true increases the disk
			{
				Config: testAccCheckAbrhaVmConfig_resize(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
					testAccCheckAbrhaVmResizeSmaller(&vm),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "size", "s-1vcpu-2gb"),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "disk", "50"),
				),
			},
		},
	})
}

func TestAccAbrhaVm_UpdateUserData(t *testing.T) {
	var afterCreate, afterUpdate goApiAbrha.Vm
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckAbrhaVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: acceptance.TestAccCheckAbrhaVmConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &afterCreate),
					testAccCheckAbrhaVmAttributes(&afterCreate),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "name", name),
				),
			},

			{
				Config: testAccCheckAbrhaVmConfig_userdata_update(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &afterUpdate),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar",
						"user_data",
						util.HashString("foobar foobar")),
					testAccCheckAbrhaVmRecreated(
						t, &afterCreate, &afterUpdate),
				),
			},
		},
	})
}

func TestAccAbrhaVm_UpdateTags(t *testing.T) {
	var afterCreate, afterUpdate goApiAbrha.Vm
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckAbrhaVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: acceptance.TestAccCheckAbrhaVmConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &afterCreate),
					testAccCheckAbrhaVmAttributes(&afterCreate),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "name", name),
				),
			},

			{
				Config: testAccCheckAbrhaVmConfig_tag_update(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &afterUpdate),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar",
						"tags.#",
						"1"),
				),
			},
		},
	})
}

func TestAccAbrhaVm_VPCAndIpv6(t *testing.T) {
	var vm goApiAbrha.Vm
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckAbrhaVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaVmConfig_VPCAndIpv6(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
					testAccCheckAbrhaVmAttributes_PrivateNetworkingIpv6(&vm),
					resource.TestCheckResourceAttrSet(
						"abrha_vm.foobar", "vpc_uuid"),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "ipv6", "true"),
				),
			},
		},
	})
}

func TestAccAbrhaVm_UpdatePrivateNetworkingIpv6(t *testing.T) {
	var afterCreate, afterUpdate goApiAbrha.Vm
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckAbrhaVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: acceptance.TestAccCheckAbrhaVmConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &afterCreate),
					testAccCheckAbrhaVmAttributes(&afterCreate),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "name", name),
				),
			},
			// For "private_networking," this is now a effectively a no-opt only updating state.
			// All Vms are assigned to a VPC by default. The API should still respond successfully.
			{
				Config: testAccCheckAbrhaVmConfig_PrivateNetworkingIpv6(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &afterUpdate),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "private_networking", "true"),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "ipv6", "true"),
				),
			},
		},
	})
}

func TestAccAbrhaVm_Monitoring(t *testing.T) {
	var vm goApiAbrha.Vm
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckAbrhaVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaVmConfig_Monitoring(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "monitoring", "true"),
				),
			},
		},
	})
}

func TestAccAbrhaVm_conditionalVolumes(t *testing.T) {
	var firstVm goApiAbrha.Vm
	var secondVm goApiAbrha.Vm
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckAbrhaVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaVmConfig_conditionalVolumes(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar.0", &firstVm),
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar.1", &secondVm),
					resource.TestCheckResourceAttr("abrha_vm.foobar.0", "volume_ids.#", "1"),

					// This could be improved in core/HCL to make it less confusing
					// but it's the only way to use conditionals in this context for now and "it works"
					resource.TestCheckResourceAttr("abrha_vm.foobar.1", "volume_ids.#", "1"),
				),
			},
		},
	})
}

func TestAccAbrhaVm_EnableAndDisableBackups(t *testing.T) {
	var vm goApiAbrha.Vm
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckAbrhaVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: acceptance.TestAccCheckAbrhaVmConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
					testAccCheckAbrhaVmAttributes(&vm),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "backups", "false"),
				),
			},

			{
				Config: testAccCheckAbrhaVmConfig_EnableBackups(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "backups", "true"),
				),
			},

			{
				Config: testAccCheckAbrhaVmConfig_DisableBackups(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "backups", "false"),
				),
			},
		},
	})
}

func TestAccAbrhaVm_ChangeBackupPolicy(t *testing.T) {
	var vm goApiAbrha.Vm
	name := acceptance.RandomTestName()
	backupsEnabled := `backups = true`
	backupsDisabled := `backups = false`
	dailyPolicy := `  backup_policy {
		plan    = "daily"
		hour    = 4
	}`
	weeklyPolicy := `  backup_policy {
		plan    = "weekly"
		weekday = "MON"
		hour    = 0
	}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckAbrhaVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaVmConfig_ChangeBackupPolicy(name, backupsEnabled, ""),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "backups", "true"),
				),
			},
			{
				Config: testAccCheckAbrhaVmConfig_ChangeBackupPolicy(name, backupsEnabled, weeklyPolicy),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "backup_policy.0.plan", "weekly"),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "backup_policy.0.weekday", "MON"),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "backup_policy.0.hour", "0"),
				),
			},
			{
				Config: testAccCheckAbrhaVmConfig_ChangeBackupPolicy(name, backupsEnabled, dailyPolicy),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "backup_policy.0.plan", "daily"),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "backup_policy.0.hour", "4"),
				),
			},
			// Verify specified backup policy is applied after re-enabling, and default policy is not used.
			{
				Config: testAccCheckAbrhaVmConfig_ChangeBackupPolicy(name, backupsDisabled, ""),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "backups", "false"),
				),
			},
			{
				Config: testAccCheckAbrhaVmConfig_ChangeBackupPolicy(name, backupsEnabled, weeklyPolicy),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "backup_policy.0.plan", "weekly"),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "backup_policy.0.weekday", "MON"),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "backup_policy.0.hour", "0"),
				),
			},
		},
	})
}

func TestAccAbrhaVm_WithBackupPolicy(t *testing.T) {
	var vm goApiAbrha.Vm
	name := acceptance.RandomTestName()
	backupsEnabled := `backups = true`
	backupPolicy := `  backup_policy {
	   plan    = "weekly"
	   weekday = "MON"
	   hour    = 0
	 }`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckAbrhaVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaVmConfig_ChangeBackupPolicy(name, backupsEnabled, backupPolicy),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "backup_policy.0.plan", "weekly"),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "backup_policy.0.weekday", "MON"),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "backup_policy.0.hour", "0"),
				),
			},
		},
	})
}

func TestAccAbrhaVm_EnableAndDisableGracefulShutdown(t *testing.T) {
	var vm goApiAbrha.Vm
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckAbrhaVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: acceptance.TestAccCheckAbrhaVmConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
					testAccCheckAbrhaVmAttributes(&vm),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "graceful_shutdown", "false"),
				),
			},

			{
				Config: testAccCheckAbrhaVmConfig_EnableGracefulShutdown(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "graceful_shutdown", "true"),
				),
			},

			{
				Config: testAccCheckAbrhaVmConfig_DisableGracefulShutdown(name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "graceful_shutdown", "false"),
				),
			},
		},
	})
}

// TestAccAbrhaVm_withVmAgentSetTrue tests that no error is returned
// from the API when creating a Vm using an OS that supports the agent
// if the `vm_agent` field is explicitly set to true.
func TestAccAbrhaVm_withVmAgentSetTrue(t *testing.T) {
	var vm goApiAbrha.Vm
	keyName := acceptance.RandomTestName()
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("abrha@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}
	vmName := acceptance.RandomTestName()
	agent := "vm_agent = true"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckAbrhaVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaVmConfig_VmAgent(keyName, publicKeyMaterial, vmName, "ubuntu-20-04-x64", agent),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "name", vmName),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "vm_agent", "true"),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "image", "ubuntu-20-04-x64"),
				),
			},
		},
	})
}

// TestAccAbrhaVm_withVmAgentSetFalse tests that no error is returned
// from the API when creating a Vm using an OS that does not support the agent
// if the `vm_agent` field is explicitly set to false.
func TestAccAbrhaVm_withVmAgentSetFalse(t *testing.T) {
	t.Skip("All Vm OSes currently support the Vm agent")

	var vm goApiAbrha.Vm
	keyName := acceptance.RandomTestName()
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("abrha@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}
	vmName := acceptance.RandomTestName()
	agent := "vm_agent = false"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckAbrhaVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaVmConfig_VmAgent(keyName, publicKeyMaterial, vmName, "rancheros", agent),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "name", vmName),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "vm_agent", "false"),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "image", "rancheros"),
				),
			},
		},
	})
}

// TestAccAbrhaVm_withVmAgentNotSet tests that no error is returned
// from the API when creating a Vm using an OS that does not support the agent
// if the `vm_agent` field is not explicitly set.
func TestAccAbrhaVm_withVmAgentNotSet(t *testing.T) {
	t.Skip("All Vm OSes currently support the Vm agent")

	var vm goApiAbrha.Vm
	keyName := acceptance.RandomTestName()
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("abrha@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}
	vmName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckAbrhaVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAbrhaVmConfig_VmAgent(keyName, publicKeyMaterial, vmName, "rancheros", ""),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "name", vmName),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "size", "s-1vcpu-1gb"),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "image", "rancheros"),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "region", "nyc3"),
				),
			},
		},
	})
}

// TestAccAbrhaVm_withVmAgentExpectError tests that an error is returned
// from the API when creating a Vm using an OS that does not support the agent
// if the `vm_agent` field is explicitly set to true.
func TestAccAbrhaVm_withVmAgentExpectError(t *testing.T) {
	t.Skip("All Vm OSes currently support the Vm agent")

	keyName := acceptance.RandomTestName()
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("abrha@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}
	vmName := acceptance.RandomTestName()
	agent := "vm_agent = true"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckAbrhaVmDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckAbrhaVmConfig_VmAgent(keyName, publicKeyMaterial, vmName, "rancheros", agent),
				ExpectError: regexp.MustCompile(`is not supported`),
			},
		},
	})
}

func TestAccAbrhaVm_withTimeout(t *testing.T) {
	vmName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckAbrhaVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`resource "abrha_vm" "foobar" {
  name   = "%s"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc3"
  timeouts {
    create = "5s"
  }
}`, vmName),
				ExpectError: regexp.MustCompile(`timeout while waiting for state`),
			},
		},
	})
}

func TestAccAbrhaVm_Regionless(t *testing.T) {
	var vm goApiAbrha.Vm
	vmName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      acceptance.TestAccCheckAbrhaVmDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "abrha_vm" "foobar" {
  name  = "%s"
  size  = "s-1vcpu-1gb"
  image = "ubuntu-22-04-x64"
}`, vmName),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckAbrhaVmExists("abrha_vm.foobar", &vm),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "name", vmName),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "size", "s-1vcpu-1gb"),
					resource.TestCheckResourceAttr(
						"abrha_vm.foobar", "image", "ubuntu-22-04-x64"),
					resource.TestCheckResourceAttrSet(
						"abrha_vm.foobar", "region"),
				),
			},
		},
	})
}

func testAccCheckAbrhaVmAttributes(vm *goApiAbrha.Vm) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if vm.URN() != fmt.Sprintf("do:vm:%s", vm.ID) {
			return fmt.Errorf("Bad URN: %s", vm.URN())
		}

		if vm.Image.Slug != "ubuntu-22-04-x64" {
			return fmt.Errorf("Bad image_slug: %s", vm.Image.Slug)
		}

		if vm.Size.Slug != "s-1vcpu-1gb" {
			return fmt.Errorf("Bad size_slug: %s", vm.Size.Slug)
		}

		if vm.Size.PriceHourly != 0.00893 {
			return fmt.Errorf("Bad price_hourly: %v", vm.Size.PriceHourly)
		}

		if vm.Size.PriceMonthly != 6.0 {
			return fmt.Errorf("Bad price_monthly: %v", vm.Size.PriceMonthly)
		}

		if vm.Region.Slug != "nyc3" {
			return fmt.Errorf("Bad region_slug: %s", vm.Region.Slug)
		}

		return nil
	}
}

func testAccCheckAbrhaVmRenamedAndResized(vm *goApiAbrha.Vm) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if vm.Size.Slug != "s-1vcpu-2gb" {
			return fmt.Errorf("Bad size_slug: %s", vm.SizeSlug)
		}

		if vm.Disk != 50 {
			return fmt.Errorf("Bad disk: %d", vm.Disk)
		}

		return nil
	}
}

func testAccCheckAbrhaVmResizeWithOutDisk(vm *goApiAbrha.Vm) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if vm.Size.Slug != "s-1vcpu-2gb" {
			return fmt.Errorf("Bad size_slug: %s", vm.SizeSlug)
		}

		if vm.Disk != 25 {
			return fmt.Errorf("Bad disk: %d", vm.Disk)
		}

		return nil
	}
}

func testAccCheckAbrhaVmResizeSmaller(vm *goApiAbrha.Vm) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if vm.Size.Slug != "s-1vcpu-2gb" {
			return fmt.Errorf("Bad size_slug: %s", vm.SizeSlug)
		}

		if vm.Disk != 50 {
			return fmt.Errorf("Bad disk: %d", vm.Disk)
		}

		return nil
	}
}

func testAccCheckAbrhaVmAttributes_PrivateNetworkingIpv6(d *goApiAbrha.Vm) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if d.Image.Slug != "ubuntu-22-04-x64" {
			return fmt.Errorf("Bad image_slug: %s", d.Image.Slug)
		}

		if d.Size.Slug != "s-1vcpu-1gb" {
			return fmt.Errorf("Bad size_slug: %s", d.Size.Slug)
		}

		if d.Region.Slug != "nyc3" {
			return fmt.Errorf("Bad region_slug: %s", d.Region.Slug)
		}

		if vm.FindIPv4AddrByType(d, "private") == "" {
			return fmt.Errorf("No ipv4 private: %s", vm.FindIPv4AddrByType(d, "private"))
		}

		if vm.FindIPv4AddrByType(d, "public") == "" {
			return fmt.Errorf("No ipv4 public: %s", vm.FindIPv4AddrByType(d, "public"))
		}

		if vm.FindIPv6AddrByType(d, "public") == "" {
			return fmt.Errorf("No ipv6 public: %s", vm.FindIPv6AddrByType(d, "public"))
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "abrha_vm" {
				continue
			}
			if rs.Primary.Attributes["ipv6_address"] != strings.ToLower(vm.FindIPv6AddrByType(d, "public")) {
				return fmt.Errorf("IPV6 Address should be lowercase")
			}

		}

		return nil
	}
}

func testAccCheckAbrhaVmRecreated(t *testing.T,
	before, after *goApiAbrha.Vm) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if before.ID == after.ID {
			t.Fatalf("Expected change of vm IDs, but both were %v", before.ID)
		}
		return nil
	}
}

func testAccCheckAbrhaVmConfig_withID(name string, slug string) string {
	return fmt.Sprintf(`
data "abrha_image" "foobar" {
  slug = "%s"
}

resource "abrha_vm" "foobar" {
  name      = "%s"
  size      = "%s"
  image     = data.abrha_image.foobar.id
  region    = "nyc3"
  user_data = "foobar"
}`, slug, name, defaultSize)
}

func testAccCheckAbrhaVmConfig_withSSH(name string, testAccValidPublicKey string) string {
	return fmt.Sprintf(`
resource "abrha_ssh_key" "foobar" {
  name       = "%s-key"
  public_key = "%s"
}

resource "abrha_vm" "foobar" {
  name      = "%s"
  size      = "%s"
  image     = "%s"
  region    = "nyc3"
  user_data = "foobar"
  ssh_keys  = [abrha_ssh_key.foobar.id]
}`, name, testAccValidPublicKey, name, defaultSize, defaultImage)
}

func testAccCheckAbrhaVmConfig_tag_update(name string) string {
	return fmt.Sprintf(`
resource "abrha_tag" "barbaz" {
  name = "barbaz"
}

resource "abrha_vm" "foobar" {
  name      = "%s"
  size      = "%s"
  image     = "%s"
  region    = "nyc3"
  user_data = "foobar"
  tags      = [abrha_tag.barbaz.id]
}
`, name, defaultSize, defaultImage)
}

func testAccCheckAbrhaVmConfig_userdata_update(name string) string {
	return fmt.Sprintf(`
resource "abrha_vm" "foobar" {
  name      = "%s"
  size      = "%s"
  image     = "%s"
  region    = "nyc3"
  user_data = "foobar foobar"
}
`, name, defaultSize, defaultImage)
}

func testAccCheckAbrhaVmConfig_RenameAndResize(newName string) string {
	return fmt.Sprintf(`
resource "abrha_vm" "foobar" {
  name   = "%s"
  size   = "s-1vcpu-2gb"
  image  = "%s"
  region = "nyc3"
}
`, newName, defaultImage)
}

func testAccCheckAbrhaVmConfig_resize_without_disk(name string) string {
	return fmt.Sprintf(`
resource "abrha_vm" "foobar" {
  name        = "%s"
  size        = "s-1vcpu-2gb"
  image       = "%s"
  region      = "nyc3"
  user_data   = "foobar"
  resize_disk = false
}
`, name, defaultImage)
}

func testAccCheckAbrhaVmConfig_resize(name string) string {
	return fmt.Sprintf(`
resource "abrha_vm" "foobar" {
  name        = "%s"
  size        = "s-1vcpu-2gb"
  image       = "%s"
  region      = "nyc3"
  user_data   = "foobar"
  resize_disk = true
}
`, name, defaultImage)
}

func testAccCheckAbrhaVmConfig_PrivateNetworkingIpv6(name string) string {
	return fmt.Sprintf(`
resource "abrha_vm" "foobar" {
  name               = "%s"
  size               = "%s"
  image              = "%s"
  region             = "nyc3"
  ipv6               = true
  private_networking = true
}
`, name, defaultSize, defaultImage)
}

func testAccCheckAbrhaVmConfig_VPCAndIpv6(name string) string {
	return fmt.Sprintf(`
resource "abrha_vpc" "foobar" {
  name   = "%s"
  region = "nyc3"
}

resource "abrha_vm" "foobar" {
  name     = "%s"
  size     = "%s"
  image    = "%s"
  region   = "nyc3"
  ipv6     = true
  vpc_uuid = abrha_vpc.foobar.id
}
`, acceptance.RandomTestName(), name, defaultSize, defaultImage)
}

func testAccCheckAbrhaVmConfig_Monitoring(name string) string {
	return fmt.Sprintf(`
resource "abrha_vm" "foobar" {
  name       = "%s"
  size       = "%s"
  image      = "%s"
  region     = "nyc3"
  monitoring = true
}
 `, name, defaultSize, defaultImage)
}

func testAccCheckAbrhaVmConfig_conditionalVolumes(name string) string {
	return fmt.Sprintf(`
resource "abrha_volume" "myvol-01" {
  region      = "sfo3"
  name        = "%s-01"
  size        = 1
  description = "an example volume"
}

resource "abrha_volume" "myvol-02" {
  region      = "sfo3"
  name        = "%s-02"
  size        = 1
  description = "an example volume"
}

resource "abrha_vm" "foobar" {
  count      = 2
  name       = "%s-${count.index}"
  region     = "sfo3"
  image      = "%s"
  size       = "%s"
  volume_ids = [count.index == 0 ? abrha_volume.myvol-01.id : abrha_volume.myvol-02.id]
}
`, name, name, name, defaultImage, defaultSize)
}

func testAccCheckAbrhaVmConfig_EnableBackups(name string) string {
	return fmt.Sprintf(`
resource "abrha_vm" "foobar" {
  name      = "%s"
  size      = "%s"
  image     = "%s"
  region    = "nyc3"
  user_data = "foobar"
  backups   = true
}`, name, defaultSize, defaultImage)
}

func testAccCheckAbrhaVmConfig_DisableBackups(name string) string {
	return fmt.Sprintf(`
resource "abrha_vm" "foobar" {
  name      = "%s"
  size      = "%s"
  image     = "%s"
  region    = "nyc3"
  user_data = "foobar"
  backups   = false
}`, name, defaultSize, defaultImage)
}

func testAccCheckAbrhaVmConfig_ChangeBackupPolicy(name, backups, backupPolicy string) string {
	return fmt.Sprintf(`
resource "abrha_vm" "foobar" {
  name      = "%s"
  size      = "%s"
  image     = "%s"
  region    = "nyc3"
  user_data = "foobar"
  %s
  %s
}`, name, defaultSize, defaultImage, backups, backupPolicy)
}

func testAccCheckAbrhaVmConfig_VmAgent(keyName, testAccValidPublicKey, vmName, image, agent string) string {
	return fmt.Sprintf(`
resource "abrha_ssh_key" "foobar" {
  name       = "%s"
  public_key = "%s"
}

resource "abrha_vm" "foobar" {
  name     = "%s"
  size     = "%s"
  image    = "%s"
  region   = "nyc3"
  ssh_keys = [abrha_ssh_key.foobar.id]
  %s
}`, keyName, testAccValidPublicKey, vmName, defaultSize, image, agent)
}

func testAccCheckAbrhaVmConfig_EnableGracefulShutdown(name string) string {
	return fmt.Sprintf(`
resource "abrha_vm" "foobar" {
  name              = "%s"
  size              = "%s"
  image             = "%s"
  region            = "nyc3"
  user_data         = "foobar"
  graceful_shutdown = true
}`, name, defaultSize, defaultImage)
}

func testAccCheckAbrhaVmConfig_DisableGracefulShutdown(name string) string {
	return fmt.Sprintf(`
resource "abrha_vm" "foobar" {
  name              = "%s"
  size              = "%s"
  image             = "%s"
  region            = "nyc3"
  user_data         = "foobar"
  graceful_shutdown = false
}`, name, defaultSize, defaultImage)
}
