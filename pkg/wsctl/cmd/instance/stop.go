package instance

import (
	"encoding/json"
	"fmt"
	"github.com/machinefi/w3bstream/pkg/wsctl/client"
	"github.com/machinefi/w3bstream/pkg/wsctl/cmd/types"
	"github.com/machinefi/w3bstream/pkg/wsctl/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"net/http"
)

var (
	_instanceStopUse = map[config.Language]string{
		config.English: "stop INSTANCE_ID",
		config.Chinese: "stop INSTANCE_ID",
	}
	_instanceStopCmdShorts = map[config.Language]string{
		config.English: "Stop a instance",
		config.Chinese: "通过 INSTANCE_ID 停止 INSTANCE",
	}
)

// newInstanceStopCmd is a command to stop instance
func newInstanceStopCmd(client client.Client) *cobra.Command {
	return &cobra.Command{
		Use:   client.SelectTranslation(_instanceStopUse),
		Short: client.SelectTranslation(_instanceStopCmdShorts),
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("accepts 1 arg(s), received %d", len(args))
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			result, err := stop(client, args)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("problem stop instance %+v", args))
			}
			cmd.Println(result)
			return nil
		},
	}
}

func stop(client client.Client, args []string) (string, error) {
	insId := args[0]
	url := fmt.Sprintf("%s/srv-applet-mgr/v0/deploy/%s/STOP", client.Config().Endpoint, insId)
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to stop instance request")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Call(url, req)
	if err != nil {
		return "", errors.Wrap(err, "failed to stop instance")
	}
	defer resp.Body.Close()

	cr := types.InstanceResp{}
	if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil {
		errors.Wrap(err, "failed to decode instance response")
	}
	fmt.Println(cr)
	if cr.Code != 0 {
		return "", fmt.Errorf("failed to stop instance, error code: %d, error message: %s", cr.Code, cr.Desc)
	}

	return cases.Title(language.Und).String(args[0]) + " instance stopped successfully ", nil
}
