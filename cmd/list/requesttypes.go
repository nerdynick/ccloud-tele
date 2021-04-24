package list

import (
	"github.com/nerdynick/ccloud-go-sdk/telemetry"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/labels"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/metric"
	"github.com/nerdynick/ccloud-tele/cmd/command"
	"github.com/spf13/cobra"
)

type RequestTypes struct {
	AttributeCMD
}

func (am *RequestTypes) Run(cmd *cobra.Command, args []string, context command.CommandContext, client telemetry.TelemetryClient) (bool, error) {
	am.ResouceType = context.ResourceType
	am.ResourceID = context.ResourceID
	res, err := client.GetKafkaRequestTypes(context.ResourceID, context.Interval)

	am.Results = res
	am.Log.Info("Fetched Request Types")

	return (len(res) > 0), err
}

func init() {
	requestTypes := &cobra.Command{
		Use:   "requests",
		Short: "List all request types",
		RunE: command.CobraRunE(&RequestTypes{
			AttributeCMD: AttributeCMD{
				Metric:      metric.Requests,
				MetricLabel: labels.MetricType,
				Log:         command.CMDContext.Log.Named("list.request.type"),
			},
		}),
	}

	command.CMDContext.AddTimeFlags(requestTypes)
	command.CMDContext.AddKafkaID(requestTypes)

	CMDList.AddCommand(requestTypes)
}
