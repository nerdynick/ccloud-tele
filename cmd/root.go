package cmd

import (
	"os"
	"strings"

	"github.com/nerdynick/ccloud-go-sdk/telemetry"
	"github.com/nerdynick/ccloud-tele/cmd/common"
	"github.com/nerdynick/ccloud-tele/cmd/list"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "ccloud-tele",
		Short: "Confluent Cloud Telemetry API CLI Tool",
	}
	//All Supported output formats
	AvailableOutputFormats = map[string]common.OutputFormat{
		"plain": common.OutputPlain,
		"json":  common.OutputJSON,
		"csv":   common.OutputCSV,
	}
)

//Global Vars
var (
	verbose           bool
	extraVerbose      bool
	extraExtraVerbose bool
	strOutputFormat   string = string(common.OutputPlain)
)

func init() {
	apiKey := os.Getenv("API_KEY")
	apiSecret := os.Getenv("API_SECRET")
	apiSecretDefault := ""

	if apiSecret != "" {
		apiSecretDefault = "****"
	}

	//Root Commands
	rootCmd.Version = "0.1.0"
	rootCmd.PersistentFlags().BoolVar(&verbose, "v", false, "Verbose output")
	rootCmd.PersistentFlags().BoolVar(&extraVerbose, "vv", false, "Extra Verbose output")
	rootCmd.PersistentFlags().BoolVar(&extraExtraVerbose, "vvv", false, "Extra Extra Verbose output")

	rootCmd.PersistentFlags().StringVarP(&common.CMDContext.APIClient.Context.APIKey, "api-key", "k", apiKey, "API Key - Optional ENV Var 'API_KEY'")
	rootCmd.PersistentFlags().StringVarP(&common.CMDContext.APIClient.Context.APISecret, "api-secret", "s", apiSecretDefault, "API Secret - Optional ENV Var 'API_SECRET'")

	rootCmd.PersistentFlags().StringVarP(&strOutputFormat, "output", "o", strOutputFormat, "Output Format - Available Options: plain, csv, json")
	rootCmd.PersistentFlags().StringVarP(&common.CMDContext.APIClient.BaseURL, "baseurl", "b", telemetry.DefaultBaseURL, "API Base Url")
	rootCmd.PersistentFlags().StringVarP(&common.CMDContext.APIClient.Context.UserAgent, "agent", "a", "ccloud-go-sdk/go-cli", "HTTP User Agent")

	rootCmd.AddCommand(list.CMDList)

	cobra.OnInitialize(rootInit, func() {
		//Test API Key and Secrets
		if common.CMDContext.APIClient.Context.APIKey == "" || common.CMDContext.APIClient.Context.APISecret == "" {
			// println()
			rootCmd.Usage()
		}
	})
}

func rootInit() {
	common.CMDContext.OutputFormat = AvailableOutputFormats[strings.ToLower(strOutputFormat)]

	//Get API Secret from ENV Vars if MASKED
	if common.CMDContext.APIClient.Context.APISecret == "****" {
		common.CMDContext.APIClient.Context.APISecret = os.Getenv("API_SECRET")
	}

	//Get the level of Logging to preform
	if verbose || extraVerbose || extraExtraVerbose {
		log.SetLevel(log.InfoLevel)

		if extraVerbose {
			log.SetLevel(log.DebugLevel)
		}

		if extraExtraVerbose {
			log.SetReportCaller(true)
			log.SetLevel(log.TraceLevel)
		}
	} else {
		log.SetLevel(log.WarnLevel)
	}

}

func Execute() error {
	return rootCmd.Execute()
}

// func getMetrics(client ccloudmetrics.MetricsClient, metrics ...string) []ccloudmetrics.Metric {
// 	availMetrics, err := client.GetAvailableMetrics()
// 	if err != nil {
// 		log.Panic(fmt.Sprintf("Failed to get all Available Metrics. Got error %s", err.Error()))
// 	}
// 	validMetrics := []ccloudmetrics.Metric{}

// 	for _, metric := range availMetrics {
// 		for _, m := range metrics {
// 			if metric.Matches(m) {
// 				validMetrics = append(validMetrics, metric)
// 			}
// 		}
// 	}
// 	return validMetrics
// }
