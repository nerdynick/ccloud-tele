package list

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nerdynick/ccloud-go-sdk/telemetry"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/metric"
	"github.com/nerdynick/ccloud-tele/cmd/command"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type AvailableMetrics struct {
	Results []metric.Metric
	Log     *zap.Logger
}

func (am *AvailableMetrics) Run(cmd *cobra.Command, args []string, context command.CommandContext, client telemetry.TelemetryClient) (bool, error) {
	res, err := client.GetAvailableMetricsForResource(context.ResourceType)
	am.Results = res
	am.Log.Info("Fetched Available Metrics")

	return (len(res) > 0), err
}
func (am AvailableMetrics) OutputPlain() error {
	for _, metric := range am.Results {
		labels := []string{}
		for _, label := range metric.Labels {
			labels = append(labels, label.Key)
		}

		fmt.Print("============================\n")
		fmt.Printf("Name:      %s\n", metric.Name)
		fmt.Printf("Desc:      %s\n", metric.Desc)
		fmt.Printf("Type:      %s\n", metric.Type)
		fmt.Printf("LifeCycle: %s\n", metric.LifecycleStage)
		fmt.Printf("Labels:    %s\n", strings.Join(labels, ","))
		fmt.Print("============================\n\n")
	}
	return nil
}
func (am AvailableMetrics) OutputJSON(encoder *json.Encoder) error {
	return encoder.Encode(am.Results)
}
func (am AvailableMetrics) OutputCSV(writer *csv.Writer) error {
	for _, metric := range am.Results {
		labels := []string{}
		for _, label := range metric.Labels {
			labels = append(labels, label.Key)
		}
		err := writer.Write([]string{
			metric.Name,
			metric.Desc,
			metric.Type,
			metric.LifecycleStage,
			strings.Join(labels, ";"),
		})
		if err != nil {
			return nil
		}
	}
	return nil
}

func init() {
	metrics := &cobra.Command{
		Use:   "metrics",
		Short: "List currently available metrics.",
		RunE: command.CobraRunE(&AvailableMetrics{
			Log: command.CMDContext.Log.Named("list.metrics"),
		}),
	}
	command.CMDContext.AddResourceTypeFlags(metrics)

	CMDList.AddCommand(metrics)
}
