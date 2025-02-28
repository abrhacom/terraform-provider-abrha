package util

import (
	"context"
	"time"

	goApiAbrha "github.com/abrhacom/go-api-abrha"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// WaitForAction waits for the action to finish using the resource.StateChangeConf.
func WaitForAction(client *goApiAbrha.Client, action *goApiAbrha.Action) error {
	var (
		pending   = "in-progress"
		target    = "completed"
		refreshfn = func() (result interface{}, state string, err error) {
			if action.ID == 0 {
				return action, target, nil
			}
			a, _, err := client.Actions.Get(context.Background(), action.ID)
			if err != nil {
				return nil, "", err
			}
			if a.Status == "errored" {
				return a, "errored", nil
			}
			if a.CompletedAt != nil {
				return a, target, nil
			}
			return a, pending, nil
		}
	)
	_, err := (&retry.StateChangeConf{
		Pending: []string{pending},
		Refresh: refreshfn,
		Target:  []string{target},

		Delay:      10 * time.Second,
		Timeout:    60 * time.Minute,
		MinTimeout: 3 * time.Second,

		// This is a hack around DO API strangeness.
		// https://github.com/hashicorp/terraform/issues/481
		//
		NotFoundChecks: 60,
	}).WaitForStateContext(context.Background())
	return err
}
