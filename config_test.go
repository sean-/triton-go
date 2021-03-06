package triton

import "testing"

func TestAccConfig_Get(t *testing.T) {
	AccTest(t, TestCase{
		Steps: []Step{
			&StepAPICall{
				StateBagKey: "config",
				CallFunc: func(client *Client) (interface{}, error) {
					return client.Config().GetConfig(&GetConfigInput{})
				},
			},
			&StepAssertSet{
				StateBagKey: "config",
				Keys:        []string{"DefaultNetwork"},
			},
		},
	})
}
