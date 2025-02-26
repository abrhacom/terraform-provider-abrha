package vm_test

import (
	"reflect"
	"testing"

	"github.com/abrhacom/terraform-provider-abrha/abrha/vm"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAbrhaVmMigrateState(t *testing.T) {
	cases := map[string]struct {
		StateVersion int
		ID           string
		Attributes   map[string]string
		Expected     map[string]string
	}{
		"v0_1_with_values": {
			StateVersion: 0,
			ID:           "id",
			Attributes: map[string]string{
				"backups":    "true",
				"monitoring": "false",
			},
			Expected: map[string]string{
				"backups":    "true",
				"monitoring": "false",
			},
		},
		"v0_1_without_values": {
			StateVersion: 0,
			ID:           "id",
			Attributes:   map[string]string{},
			Expected: map[string]string{
				"backups":    "false",
				"monitoring": "false",
			},
		},
	}

	for tn, tc := range cases {
		is := &terraform.InstanceState{
			ID:         tc.ID,
			Attributes: tc.Attributes,
		}
		is, err := vm.ResourceAbrhaVmMigrateState(tc.StateVersion, is, nil)

		if err != nil {
			t.Fatalf("bad: %q, err: %#v", tn, err)
		}

		if !reflect.DeepEqual(tc.Expected, is.Attributes) {
			t.Fatalf("Bad Vm Migrate\n\n. Got: %+v\n\n expected: %+v", is.Attributes, tc.Expected)
		}
	}
}
