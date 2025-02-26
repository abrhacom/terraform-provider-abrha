package sweep_test

import (
	"testing"

	_ "github.com/abrhacom/terraform-provider-abrha/abrha/app"
	_ "github.com/abrhacom/terraform-provider-abrha/abrha/cdn"
	_ "github.com/abrhacom/terraform-provider-abrha/abrha/certificate"
	_ "github.com/abrhacom/terraform-provider-abrha/abrha/database"
	_ "github.com/abrhacom/terraform-provider-abrha/abrha/domain"
	_ "github.com/abrhacom/terraform-provider-abrha/abrha/firewall"
	_ "github.com/abrhacom/terraform-provider-abrha/abrha/image"
	_ "github.com/abrhacom/terraform-provider-abrha/abrha/kubernetes"
	_ "github.com/abrhacom/terraform-provider-abrha/abrha/loadbalancer"
	_ "github.com/abrhacom/terraform-provider-abrha/abrha/monitoring"
	_ "github.com/abrhacom/terraform-provider-abrha/abrha/project"
	_ "github.com/abrhacom/terraform-provider-abrha/abrha/reservedip"
	_ "github.com/abrhacom/terraform-provider-abrha/abrha/snapshot"
	_ "github.com/abrhacom/terraform-provider-abrha/abrha/spaces"
	_ "github.com/abrhacom/terraform-provider-abrha/abrha/sshkey"
	_ "github.com/abrhacom/terraform-provider-abrha/abrha/uptime"
	_ "github.com/abrhacom/terraform-provider-abrha/abrha/vm"
	_ "github.com/abrhacom/terraform-provider-abrha/abrha/volume"
	_ "github.com/abrhacom/terraform-provider-abrha/abrha/vpc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
