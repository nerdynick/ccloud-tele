package list

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nerdynick/ccloud-go-sdk/telemetry"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/metric"
	"github.com/nerdynick/ccloud-tele/cmd/common"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type AvailableMetrics struct {
	Results []metric.Metric
}

func (am *AvailableMetrics) Run(cmd *cobra.Command, args []string, context common.CommandContext, client telemetry.TelemetryClient) (bool, error) {
	resType, err := context.GetResourceType()
	if err != nil {
		return false, err
	}

	res, err := client.GetAvailableMetricsForResource(resType)
	am.Results = res
	log.WithFields(log.Fields{
		"result": res,
		"err":    err,
	}).Info("Fetched Available Metrics")

	return (len(res) > 0), err
}
func (am AvailableMetrics) OutputPlain() error {
	log.WithFields(log.Fields{
		"result": am.Results,
	}).Info("Printing Plain Output")

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
		RunE:  common.CobraRunE(&AvailableMetrics{}),
	}
	common.AddResourceTypeFlags(metrics, &common.CMDContext)

	CMDList.AddCommand(metrics)
}
