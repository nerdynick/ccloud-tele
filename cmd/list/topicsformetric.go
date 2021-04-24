package list

import (
	"github.com/nerdynick/ccloud-go-sdk/telemetry"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/labels"
	"github.com/nerdynick/ccloud-tele/cmd/command"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type TopicsForMetric struct {
	AttributeCMD
}

func (am *TopicsForMetric) Run(cmd *cobra.Command, args []string, context command.CommandContext, client telemetry.TelemetryClient) (bool, error) {
	am.ResouceType = context.ResourceType
	am.ResourceID = context.ResourceID
	am.Metric = context.Metric
	res, err := client.GetKafkaTopicsForMetric(context.ResourceID, context.Metric, context.Interval)

	am.Results = res
	log.WithFields(log.Fields{
		"result":  res,
		"err":     err,
		"context": context,
	}).Info("Fetched Available Topics for Metric")

	return (len(res) > 0), err
}

func init() {
	topicsForMetric := &cobra.Command{
		Use:   "topics",
		Short: "List all available topics for a given metric",
		RunE: command.CobraRunE(&TopicsForMetric{
			AttributeCMD: AttributeCMD{
				MetricLabel: labels.MetricType,
				Log:         command.CMDContext.Log.Named("list.topics"),
			},
		}),
	}

	command.CMDContext.AddTimeFlags(topicsForMetric)
	command.CMDContext.AddKnownMetricFlags(topicsForMetric)
	command.CMDContext.AddKafkaID(topicsForMetric)

	CMDList.AddCommand(topicsForMetric)
}
