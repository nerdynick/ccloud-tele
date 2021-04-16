package common

import (
	"errors"

	"github.com/nerdynick/ccloud-go-sdk/telemetry"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/labels"
	"github.com/spf13/cobra"
)

const (
	OutputPlain OutputFormat = "plain"
	OutputJSON  OutputFormat = "json"
	OutputCSV   OutputFormat = "csv"
)

//Output Formats
type OutputFormat string

type CommandContext struct {
	APIClient         telemetry.TelemetryClient
	StartTime         string
	EndTime           string
	Topic             string
	Topics            []string
	BlacklistedTopics []string
	IncludePartitions bool
	Granularity       string
	LastXmin          int
	OutputFormat      OutputFormat
	Metric            string

	ResourceID          string
	ResourceType        labels.Resource
	ResourceIsKafka     bool
	ResourceIsConnector bool
	ResourceIsSchemaReg bool
	ResourceIsKSQL      bool

	preRuns []func(*CommandContext, []string) error
}

func (ctx *CommandContext) PreRunE(args []string) error {
	if ctx.preRuns != nil {
		for _, e := range ctx.preRuns {
			err := e(ctx, args)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (ctx *CommandContext) AddPreRun(f func(*CommandContext, []string) error) {
	if ctx.preRuns == nil {
		ctx.preRuns = []func(*CommandContext, []string) error{f}
	} else {
		ctx.preRuns = append(ctx.preRuns, f)
	}
}

func (ctx CommandContext) GetResourceType() (labels.Resource, error) {
	if ctx.ResourceIsKafka {
		return labels.ResourceKafka, nil
	}

	return labels.Resource{}, errors.New("No Valid Resource Type Provided")
}

func AddResourceTypeFlags(cmd *cobra.Command, ctx *CommandContext) {
	ctx.AddPreRun(func(ctx *CommandContext, args []string) error {
		if ctx.ResourceIsConnector {
			ctx.ResourceType = labels.ResourceConnector
		} else if ctx.ResourceIsKSQL {
			ctx.ResourceType = labels.ResourceKSQL
		} else if ctx.ResourceIsKafka {
			ctx.ResourceType = labels.ResourceKafka
		} else if ctx.ResourceIsSchemaReg {
			ctx.ResourceType = labels.ResourceSchemaRegistry
		} else {
			return errors.New("no resource type selected")
		}
		return nil
	})

	cmd.Flags().BoolVar(&ctx.ResourceIsKafka, "kafka", true, "Resource ID refers to a Kafka Cluster")
	cmd.Flags().BoolVar(&ctx.ResourceIsKafka, "connector", true, "Resource ID refers to a Kafka Connector")
	cmd.Flags().BoolVar(&ctx.ResourceIsKafka, "ksql", true, "Resource ID refers to a KSQL Cluster")
	cmd.Flags().BoolVar(&ctx.ResourceIsKafka, "sr", true, "Resource ID refers to a Schema Registry")
}

func AddResourceTypeFlagsWithArg(cmd *cobra.Command, ctx *CommandContext) {
	AddResourceTypeFlags(cmd, ctx)

	ctx.AddPreRun(func(ctx *CommandContext, args []string) error {
		ctx.ResourceID = args[0]
		return nil
	})

	cmd.Args = cobra.ExactArgs(1)
}

// func (r *RequestContext) getStartTime() time.Time {
// 	if r.LastXmin > 0 {
// 		return time.Now().Add(time.Duration(-r.LastXmin) * time.Minute)
// 	}
// 	res, err := time.Parse(ccloudmetrics.TimeFormatStr, r.StartTime)
// 	if err != nil {
// 		log.Panic(fmt.Sprintf("Start Time is invalid. Times must be provided in the %s format. Was given %s", ccloudmetrics.TimeFormatStr, r.StartTime))
// 	}
// 	return res
// }
// func (r *RequestContext) getEndTime() time.Time {
// 	res, err := time.Parse(ccloudmetrics.TimeFormatStr, r.EndTime)
// 	if err != nil {
// 		log.Panic(fmt.Sprintf("End Time is invalid. Times must be provided in the %s format. Was given %s", ccloudmetrics.TimeFormatStr, r.EndTime))
// 	}
// 	return res
// }

// func (r *RequestContext) getGranularity() ccloudmetrics.Granularity {
// 	g := ccloudmetrics.Granularity(r.Granularity)
// 	if !g.IsValid() {
// 		log.Panic(fmt.Sprintf("Granularity is invalid. Was given the value of %s expecting on of %s", r.Granularity, strings.Join(ccloudmetrics.AvailableGranularities, ", ")))
// 	}
// 	return g
// }
