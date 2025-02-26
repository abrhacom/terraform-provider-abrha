package vm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func ResourceAbrhaVmMigrateState(v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {

	switch v {
	case 0:
		log.Println("[INFO] Found Abrha Vm State v0; migrating to v1")
		return migrateAbrhaVmStateV0toV1(is)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func migrateAbrhaVmStateV0toV1(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	if is.Empty() {
		log.Println("[DEBUG] Empty InstanceState; nothing to migrate.")
		return is, nil
	}

	log.Printf("[DEBUG] Abrha Vm Attributes before Migration: %#v", is.Attributes)

	if _, ok := is.Attributes["backups"]; !ok {
		is.Attributes["backups"] = "false"
	}
	if _, ok := is.Attributes["monitoring"]; !ok {
		is.Attributes["monitoring"] = "false"
	}

	log.Printf("[DEBUG] Abrha Vm Attributes after State Migration: %#v", is.Attributes)

	return is, nil
}
