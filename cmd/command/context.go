package command

import (
	"errors"

	"github.com/nerdynick/ccloud-go-sdk/telemetry"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/resourcetype"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	CMDContext = New()
)

const (
	OutputPlain OutputFormat = "plain"
	OutputJSON  OutputFormat = "json"
	OutputCSV   OutputFormat = "csv"
)

//Output Formats
type OutputFormat string

type CommandContext struct {
	APIClient telemetry.TelemetryClient
	logConfig zap.Config
	Log       *zap.Logger

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
	ResourceType        resourcetype.ResourceType
	ResourceIsKafka     bool
	ResourceIsConnector bool
	ResourceIsSchemaReg bool
	ResourceIsKSQL      bool
}

func New() CommandContext {
	lConf := zap.NewProductionConfig()
	lConf.Level.SetLevel(zap.ErrorLevel)
	logger, _ := lConf.Build()

	return CommandContext{
		APIClient: telemetry.New("", ""),
		Log:       logger.Named("ccloud-tele-cli"),
		logConfig: lConf,
	}
}
func (ctx *CommandContext) LogLevel1() {
	ctx.logConfig.Level.SetLevel(zap.WarnLevel)
	ctx.APIClient.SetLogLevel(zap.WarnLevel)
}
func (ctx *CommandContext) LogLevel2() {
	ctx.logConfig.Level.SetLevel(zap.InfoLevel)
	ctx.APIClient.SetLogLevel(zap.InfoLevel)
}
func (ctx *CommandContext) LogLevel3() {
	ctx.logConfig.Level.SetLevel(zap.DebugLevel)
	ctx.APIClient.SetLogLevel(zap.DebugLevel)
}

func AddResourceTypeFlags(cmd *cobra.Command, ctx *CommandContext) {
	currentPreRunE := cmd.PreRunE
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if ctx.ResourceIsConnector {
			ctx.ResourceType = resourcetype.ResourceTypeConnector
		} else if ctx.ResourceIsKSQL {
			ctx.ResourceType = resourcetype.ResourceTypeKSQL
		} else if ctx.ResourceIsKafka {
			ctx.ResourceType = resourcetype.ResourceTypeKafka
		} else if ctx.ResourceIsSchemaReg {
			ctx.ResourceType = resourcetype.ResourceTypeSchemaRegistry
		} else {
			return errors.New("no resource type selected")
		}
		if currentPreRunE != nil {
			return currentPreRunE(cmd, args)
		}
		return nil
	}

	cmd.Flags().BoolVar(&ctx.ResourceIsKafka, "kafka", false, "Resource ID refers to a Kafka Cluster")
	cmd.Flags().BoolVar(&ctx.ResourceIsKafka, "connector", false, "Resource ID refers to a Kafka Connector")
	cmd.Flags().BoolVar(&ctx.ResourceIsKafka, "ksql", false, "Resource ID refers to a KSQL Cluster")
	cmd.Flags().BoolVar(&ctx.ResourceIsKafka, "sr", false, "Resource ID refers to a Schema Registry")
}

func AddResourceTypeFlagsWithArg(cmd *cobra.Command, ctx *CommandContext) {
	AddResourceTypeFlags(cmd, ctx)

	currentPreRunE := cmd.PreRunE
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		ctx.ResourceID = args[0]
		if currentPreRunE != nil {
			return currentPreRunE(cmd, args)
		}
		return nil
	}

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
