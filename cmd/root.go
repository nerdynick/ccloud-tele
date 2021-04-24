package cmd

import (
	"os"
	"strings"

	"github.com/nerdynick/ccloud-go-sdk/client"
	"github.com/nerdynick/ccloud-go-sdk/telemetry"
	"github.com/nerdynick/ccloud-tele/cmd/command"
	"github.com/nerdynick/ccloud-tele/cmd/list"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "ccloud-tele",
		Short: "Confluent Cloud Telemetry API CLI Tool",
	}
	//All Supported output formats
	AvailableOutputFormats = map[string]command.OutputFormat{
		"plain": command.OutputPlain,
		"json":  command.OutputJSON,
		"csv":   command.OutputCSV,
	}
)

//Global Vars
var (
	verbose           bool
	extraVerbose      bool
	extraExtraVerbose bool
	strOutputFormat   string = string(command.OutputPlain)
	apiKey            string = os.Getenv("API_KEY")
	apiSecret         string = os.Getenv("API_SECRET")
)

func init() {
	apiSecretDefault := ""

	if apiSecret != "" {
		apiSecretDefault = "****Secret***"
	}

	//Root Commands
	rootCmd.Version = "0.1.0"
	rootCmd.PersistentFlags().BoolVar(&verbose, "v", false, "Verbose output")
	rootCmd.PersistentFlags().BoolVar(&extraVerbose, "vv", false, "Extra Verbose output")
	rootCmd.PersistentFlags().BoolVar(&extraExtraVerbose, "vvv", false, "Extra Extra Verbose output")

	rootCmd.PersistentFlags().StringVarP(&apiKey, "api-key", "k", apiKey, "API Key - Optional ENV Var 'API_KEY'")
	rootCmd.PersistentFlags().StringVarP(&apiSecret, "api-secret", "s", apiSecretDefault, "API Secret - Optional ENV Var 'API_SECRET'")

	rootCmd.PersistentFlags().StringVarP(&strOutputFormat, "output", "o", strOutputFormat, "Output Format - Available Options: plain, csv, json")
	rootCmd.PersistentFlags().StringVarP(&command.CMDContext.APIClient.Context.BaseURL, "baseurl", "b", telemetry.DefaultBaseURL, "API Base Url")
	rootCmd.PersistentFlags().StringVarP(&command.CMDContext.APIClient.Context.UserAgent, "agent", "a", "ccloud-go-sdk/go-cli", "HTTP User Agent")

	rootCmd.AddCommand(list.CMDList)

	cobra.OnInitialize(rootInit, func() {
		//Test API Key and Secrets
		if command.CMDContext.APIClient.Context.APIKey == "" || command.CMDContext.APIClient.Context.APISecret == "" {
			// println()
			rootCmd.Usage()
		}
	})
}

func rootInit() {
	command.CMDContext.OutputFormat = AvailableOutputFormats[strings.ToLower(strOutputFormat)]

	//Get API Secret from ENV Vars if MASKED
	if apiSecret == "****Secret***" {
		apiSecret = os.Getenv("API_SECRET")
	}

	command.CMDContext.APIClient.Context.APIKey = apiKey
	command.CMDContext.APIClient.Context.APISecret = client.SecurePassword(apiSecret)

	//Get the level of Logging to preform
	if verbose || extraVerbose || extraExtraVerbose {
		command.CMDContext.LogLevel1()
		if extraVerbose {
			command.CMDContext.LogLevel2()
		}

		if extraExtraVerbose {
			command.CMDContext.LogLevel3()
		}
	} else {
		command.CMDContext.LogLevel0()
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
