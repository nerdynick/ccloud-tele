package list

import (
	"github.com/nerdynick/ccloud-go-sdk/telemetry"
	"github.com/nerdynick/ccloud-tele/cmd/command"
	"github.com/spf13/cobra"
)

type MetricAttribute struct {
	AttributeCMD
}

func (am *MetricAttribute) Run(cmd *cobra.Command, args []string, context command.CommandContext, client telemetry.TelemetryClient) (bool, error) {
	am.ResouceType = context.ResourceType
	am.ResourceID = context.ResourceID
	am.MetricLabel = context.MetricLabel
	am.Metric = context.Metric

	res, err := client.SendAttri(context.ResourceType.Labels[0], context.ResourceID, context.Metric, context.MetricLabel, context.Interval)

	am.Results = res
	am.Log.Info("Fetched Metrics Attributes")

	return (len(res) > 0), err
}

func init() {
	requestTypes := &cobra.Command{
		Use:   "attr",
		Short: "List all values for a given Metric & Metric Label for a given Resource in a given timeframe",
		RunE: command.CobraRunE(&MetricAttribute{
			AttributeCMD: AttributeCMD{
				Log: command.CMDContext.Log.Named("list.request.type"),
			},
		}),
	}

	command.CMDContext.AddTimeFlags(requestTypes)
	command.CMDContext.AddKnownMetricFlags(requestTypes)
	command.CMDContext.AddKnownMetricLabelFlags(requestTypes)
	command.CMDContext.AddResourceTypeFlagsWithID(requestTypes)

	CMDList.AddCommand(requestTypes)
}
