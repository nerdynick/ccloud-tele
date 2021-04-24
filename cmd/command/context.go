package command

import (
	"errors"
	"strings"
	"time"

	"github.com/nerdynick/ccloud-go-sdk/telemetry"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/labels"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/metric"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/query/interval"
	"github.com/nerdynick/ccloud-go-sdk/telemetry/resourcetype"
	"github.com/nerdynick/ccloud-tele/pflag"
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

	IntervalStr string
	Interval    interval.Interval

	Topics            []string
	BlacklistedTopics []string
	IncludePartitions bool

	OutputFormat OutputFormat

	MetricStr string
	Metric    metric.Metric

	ResourceID          string
	ResourceType        resourcetype.ResourceType
	ResourceIsKafka     bool
	ResourceIsConnector bool
	ResourceIsSchemaReg bool
	ResourceIsKSQL      bool

	MetricLabelStr string
	MetricLabel    labels.Metric
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
func (ctx *CommandContext) LogLevel0() {
	ctx.logConfig.Level.SetLevel(zap.ErrorLevel)
	ctx.APIClient.SetLogLevel(zap.ErrorLevel)
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

func wrapPreRunE(cmd *cobra.Command, fnc func(cmd *cobra.Command, args []string) error) {
	currentPreRunE := cmd.PreRunE
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		err := fnc(cmd, args)
		if err != nil {
			return err
		}

		if currentPreRunE != nil {
			return currentPreRunE(cmd, args)
		}
		return nil
	}
}

func (ctx *CommandContext) AddTimeFlags(cmd *cobra.Command) {
	wrapPreRunE(cmd, func(cmd *cobra.Command, args []string) error {
		i, err := interval.Parse(ctx.IntervalStr)
		if err != nil {
			return err
		}
		ctx.Interval = i

		return nil
	})
	cmd.Flags().StringVar(&ctx.IntervalStr, "interval", interval.EndingAt(2*time.Hour, time.Now().Round(time.Hour)).String(), "Time Interval in the form of ISO-8601  (START_TIME/END_TIME, START_TIME/DURATION)")
}

func metricParser(ctx *CommandContext) error {
	if ctx.MetricStr == "" {
		return errors.New("no Metric was provided")
	}
	ctx.Metric = metric.New(ctx.MetricStr)
	return nil
}

func (ctx *CommandContext) addKafkaServerMetrics(cmd *cobra.Command) {
	for _, m := range metric.KnownKafkaServerMetrics {
		if strings.HasPrefix(m.Name, "io.confluent.kafka.server") {
			flagName := strings.ReplaceAll(m.Name, "_", "-")
			flagName = strings.ReplaceAll(flagName, "/", "-")
			flagName = strings.ReplaceAll(flagName, "io.confluent.", "")
			flagName = strings.ReplaceAll(flagName, ".", "-")

			pflag.BoolStringVar(cmd.Flags(), &ctx.MetricStr, "metric-"+flagName, m.Name, "Show Topics for "+m.Name)
		}
	}
}

func (ctx *CommandContext) AddKnownKafkaServerMetricFlags(cmd *cobra.Command) {
	wrapPreRunE(cmd, func(cmd *cobra.Command, args []string) error {
		return metricParser(ctx)
	})

	ctx.addKafkaServerMetrics(cmd)
	cmd.Flags().StringVar(&ctx.MetricStr, "metric-other", "", "Provide metric maybe not yet known")
}

func (ctx *CommandContext) AddKnownMetricFlags(cmd *cobra.Command) {
	wrapPreRunE(cmd, func(cmd *cobra.Command, args []string) error {
		return metricParser(ctx)
	})

	ctx.addKafkaServerMetrics(cmd)
	cmd.Flags().StringVar(&ctx.MetricStr, "metric-other", "", "Provide metric maybe not yet known")
}

func (ctx *CommandContext) AddKnownMetricLabelFlags(cmd *cobra.Command) {
	wrapPreRunE(cmd, func(cmd *cobra.Command, args []string) error {
		if ctx.MetricLabelStr == "" {
			return errors.New("no Metric Label was provided")
		}
		ctx.MetricLabel = labels.NewMetric(ctx.MetricLabelStr)
		return nil
	})

	for _, m := range labels.KnownMetrics {
		flagName := strings.ReplaceAll(m.Key, ".", "-")
		pflag.BoolStringVar(cmd.Flags(), &ctx.MetricLabelStr, "label-"+flagName, m.Key, m.Key+" Label/Attribute")
	}
	cmd.Flags().StringVar(&ctx.MetricLabelStr, "label-other", "", "Provide your own label value")
}

func (ctx *CommandContext) AddResourceTypeFlags(cmd *cobra.Command) {
	wrapPreRunE(cmd, func(cmd *cobra.Command, args []string) error {
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
		return nil
	})

	cmd.Flags().BoolVar(&ctx.ResourceIsKafka, "resource-kafka", false, "Resource ID refers to a Kafka Cluster")
	cmd.Flags().BoolVar(&ctx.ResourceIsKafka, "resource-connector", false, "Resource ID refers to a Kafka Connector")
	cmd.Flags().BoolVar(&ctx.ResourceIsKafka, "resource-ksql", false, "Resource ID refers to a KSQL Cluster")
	cmd.Flags().BoolVar(&ctx.ResourceIsKafka, "resource-sr", false, "Resource ID refers to a Schema Registry")
}

func (ctx *CommandContext) AddResourceTypeFlagsWithID(cmd *cobra.Command) {
	ctx.AddResourceTypeFlags(cmd)

	cmd.Flags().StringVar(&ctx.ResourceID, "resource-id", "", "The actual Resource ID. (ex: LKC-XXXXX)")
	cmd.MarkFlagRequired("resource-id")
}

func (ctx *CommandContext) AddKafkaID(cmd *cobra.Command) {
	ctx.ResourceType = resourcetype.ResourceTypeKafka
	cmd.Flags().StringVar(&ctx.ResourceID, "resource-id", "", "Resource ID for for Cluster in question. (ex: LKC-XXXXX)")
	cmd.MarkFlagRequired("resource-id")
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
