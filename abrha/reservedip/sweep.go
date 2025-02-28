package reservedip

import (
	"context"

	"github.com/abrhacom/terraform-provider-abrha/abrha/config"
	"github.com/abrhacom/terraform-provider-abrha/abrha/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("abrha_reserved_ip", &resource.Sweeper{
		Name: "abrha_reserved_ip",
		F:    sweepReservedIPs,
	})

	resource.AddTestSweepers("abrha_floating_ip", &resource.Sweeper{
		Name: "abrha_floating_ip",
		F:    testSweepFloatingIps,
	})
}

func sweepReservedIPs(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	ips, _, err := client.ReservedIPs.List(context.Background(), nil)
	if err != nil {
		return err
	}

	for _, ip := range ips {
		if _, err := client.ReservedIPs.Delete(context.Background(), ip.IP); err != nil {
			return err
		}
	}

	return nil
}

func testSweepFloatingIps(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GoApiAbrhaClient()

	ips, _, err := client.FloatingIPs.List(context.Background(), nil)
	if err != nil {
		return err
	}

	for _, ip := range ips {
		if _, err := client.FloatingIPs.Delete(context.Background(), ip.IP); err != nil {
			return err
		}
	}

	return nil
}
