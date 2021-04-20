package list

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nerdynick/ccloud-go-sdk/telemetry"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/resourcetype"
	"github.com/nerdynick/ccloud-tele/cmd/command"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type AvailableResources struct {
	Results []resourcetype.ResourceType
	Log     *zap.Logger
}

func (am *AvailableResources) Run(cmd *cobra.Command, args []string, context command.CommandContext, client telemetry.TelemetryClient) (bool, error) {
	res, err := client.GetAvailableResources()
	am.Results = res
	am.Log.Info("Fetched Available Resources")

	return (len(res) > 0), err
}
func (am AvailableResources) OutputPlain() error {
	for _, metric := range am.Results {
		labels := []string{}
		for _, label := range metric.Labels {
			labels = append(labels, label.Key)
		}

		fmt.Print("============================\n")
		fmt.Printf("Type:      %s\n", metric.Type)
		fmt.Printf("Desc:      %s\n", metric.Desc)
		fmt.Printf("Labels:    %s\n", strings.Join(labels, ","))
		fmt.Print("============================\n\n")
	}
	return nil
}
func (am AvailableResources) OutputJSON(encoder *json.Encoder) error {
	return encoder.Encode(am.Results)
}
func (am AvailableResources) OutputCSV(writer *csv.Writer) error {
	for _, metric := range am.Results {
		labels := []string{}
		for _, label := range metric.Labels {
			labels = append(labels, label.Key)
		}
		err := writer.Write([]string{
			metric.Type,
			metric.Desc,
			strings.Join(labels, ";"),
		})
		if err != nil {
			return nil
		}
	}
	return nil
}

func init() {
	resources := &cobra.Command{
		Use:   "resources",
		Short: "List currently available resources.",
		RunE: command.CobraRunE(&AvailableResources{
			Log: command.CMDContext.Log.Named("list.resources"),
		}),
	}

	CMDList.AddCommand(resources)
}
