package triton_test

import (
	"fmt"
	"testing"

	triton "github.com/joyent/triton-go"
)

func getAnyMachineID(t *testing.T, c *triton.Client) (string, error) {
	machines, err := c.Machines().GetMachines()
	if err != nil {
		return "", err
	}

	for _, m := range machines {
		if len(m.ID) > 0 {
			return m.ID, nil
		}
	}

	t.Skip()
	return "", fmt.Errorf("no machines configured")
}

func TestAccMachine_GetMachine(t *testing.T) {
	triton.AccTest(t, triton.TestCase{
		Steps: []triton.Step{
			&triton.StepAPICall{
				StateBagKey: "machine",
				CallFunc: func(client *triton.Client) (interface{}, error) {
					machineID, err := getAnyMachineID(t, client)
					if err != nil {
						return nil, err
					}

					return client.Machines().GetMachine(&triton.GetMachineInput{
						ID: machineID,
					})
				},
			},
			&triton.StepAssertSet{
				StateBagKey: "machine",
				Keys:        []string{"ID", "Name", "Type", "Tags"},
			},
		},
	})
}

// FIXME(seanc@): TestAccMachine_ListMachineTags assumes that any machine ID
// returned from getAnyMachineID will have at least one tag.
func TestAccMachine_ListMachineTags(t *testing.T) {
	triton.AccTest(t, triton.TestCase{
		Steps: []triton.Step{
			&triton.StepAPICall{
				StateBagKey: "machine",
				CallFunc: func(client *triton.Client) (interface{}, error) {
					machineID, err := getAnyMachineID(t, client)
					if err != nil {
						return nil, err
					}

					return client.Machines().ListMachineTags(&triton.ListMachineTagsInput{
						ID: machineID,
					})
				},
			},
			&triton.StepAssertFunc{
				AssertFunc: func(state triton.TritonStateBag) error {
					tagsRaw, found := state.GetOk("machine")
					if !found {
						return fmt.Errorf("State key %q not found", "machines")
					}

					tags := tagsRaw.(map[string]string)
					if len(tags) == 0 {
						return fmt.Errorf("Expected at least one tag on machine")
					}
					return nil
				},
			},
		},
	})
}
