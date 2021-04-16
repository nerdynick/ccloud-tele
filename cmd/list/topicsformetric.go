package list

// import (
// 	"encoding/csv"
// 	"encoding/json"
// 	"fmt"
// 	"time"

// 	"github.com/nerdynick/ccloud-tele/cmd/common"
// 	"github.com/nerdynick/confluent-cloud-metrics-go-sdk/ccloudmetrics"
// 	log "github.com/sirupsen/logrus"
// 	"github.com/spf13/cobra"
// )

// type TopicsForMetric struct {
// 	Results []string
// }

// func (am *TopicsForMetric) Run(cmd *cobra.Command, args []string, context common.CommandContext, client ccloudmetrics.MetricsClient) (bool, error) {
// 	metrics := getMetrics(client, args...)
// 	res, err := client.GetTopicsForMetric(context.Cluster, metrics[0], context.getStartTime(), context.getEndTime())

// 	am.Results = res
// 	log.WithFields(log.Fields{
// 		"result":  res,
// 		"err":     err,
// 		"context": context,
// 	}).Info("Fetched Available Topics for Metric")

// 	return (len(res) > 0), err
// }
// func (am TopicsForMetric) outputPlain() error {
// 	log.WithFields(log.Fields{
// 		"result": am.Results,
// 	}).Info("Printing Plain Output")

// 	for _, topic := range am.Results {
// 		fmt.Printf("Topic: %s\n", topic)
// 		fmt.Println()
// 	}
// 	return nil
// }
// func (am TopicsForMetric) outputJSON(encoder *json.Encoder) error {
// 	return encoder.Encode(am.Results)
// }
// func (am TopicsForMetric) outputCSV(writer *csv.Writer) error {
// 	for _, topic := range am.Results {
// 		err := writer.Write([]string{topic})
// 		if err != nil {
// 			return nil
// 		}
// 	}
// 	return nil
// }

// func init() {
// 	topicsForMetric := &cobra.Command{
// 		Use:   "topics [cluster-id] [metric",
// 		Short: "List all available topics for a given metric",
// 		Args:  cobra.ExactArgs(2),
// 		RunE:  runE(&TopicsForMetric{}),
// 	}

// 	topicsForMetric.Flags().StringVar(&requestcontext.StartTime, "start", time.Now().Add(time.Duration(-1)*time.Hour).Format(ccloudmetrics.TimeFormatStr), "Start Time in the format of "+ccloudmetrics.TimeFormatStr)
// 	topicsForMetric.Flags().StringVar(&requestcontext.EndTime, "end", time.Now().Format(ccloudmetrics.TimeFormatStr), "End Time in the format of "+ccloudmetrics.TimeFormatStr)
// 	CMDList.AddCommand(topicsForMetric)
// }
