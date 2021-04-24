package list

import (
	"encoding/csv"
	"encoding/json"
	"fmt"

	"github.com/nerdynick/ccloud-go-sdk/telemetry/labels"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/metric"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/resourcetype"
	"go.uber.org/zap"
)

type AttributeCMD struct {
	Results     []string
	Log         *zap.Logger
	ResouceType resourcetype.ResourceType
	ResourceID  string
	Metric      metric.Metric
	MetricLabel labels.Metric
}

func (am AttributeCMD) OutputPlain() error {
	fmt.Printf("=== %s(%s) - %s(%s) ===\n", am.ResouceType.Type, am.ResourceID, am.Metric.Name, am.MetricLabel.Key)
	for _, result := range am.Results {
		fmt.Printf("%s\n", result)
	}
	return nil
}
func (am AttributeCMD) OutputJSON(encoder *json.Encoder) error {
	return encoder.Encode(am.Results)
}
func (am AttributeCMD) OutputCSV(writer *csv.Writer) error {
	for _, topic := range am.Results {
		err := writer.Write([]string{topic})
		if err != nil {
			return nil
		}
	}
	return nil
}
