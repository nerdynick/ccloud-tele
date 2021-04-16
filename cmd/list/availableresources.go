package list

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nerdynick/ccloud-go-sdk/telemetry"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/response"
	"github.com/nerdynick/ccloud-tele/cmd/common"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type AvailableResources struct {
	Results []response.ResourceType
}

func (am *AvailableResources) Run(cmd *cobra.Command, args []string, context common.CommandContext, client telemetry.TelemetryClient) (bool, error) {
	res, err := client.GetAvailableResources()
	am.Results = res
	log.WithFields(log.Fields{
		"result": res,
		"err":    err,
	}).Info("Fetched Available Resources")

	return (len(res) > 0), err
}
func (am AvailableResources) OutputPlain() error {
	log.WithFields(log.Fields{
		"result": am.Results,
	}).Info("Printing Plain Output")

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
	metrics := &cobra.Command{
		Use:   "resources",
		Short: "List currently available resources.",
		RunE:  common.CobraRunE(&AvailableResources{}),
	}

	CMDList.AddCommand(metrics)
}
