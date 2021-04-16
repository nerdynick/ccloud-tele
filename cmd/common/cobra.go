package common

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nerdynick/ccloud-go-sdk/telemetry"
	"github.com/spf13/cobra"
)

var (
	CMDContext = CommandContext{
		APIClient: telemetry.New("", ""),
	}
)

//CobraRunFunc is a struct to handle all CMDs in a uniform and common manor to reduce duplicate code
type CobraRunFunc interface {
	Run(*cobra.Command, []string, CommandContext, telemetry.TelemetryClient) (bool, error)
	OutputPlain() error
	OutputCSV(*csv.Writer) error
	OutputJSON(*json.Encoder) error
}

func CobraRunE(run CobraRunFunc) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		err := CMDContext.PreRunE(args)
		if err != nil {
			return err
		}

		results, err := run.Run(cmd, args, CMDContext, CMDContext.APIClient)
		if err != nil {
			log.Panic(fmt.Sprintf("Failed to get full results. Error: %s", err.Error()))
			return err
		}

		outputErrs := []error{}
		if err != nil {
			outputErrs = append(outputErrs, err)
		}

		if results {
			switch CMDContext.OutputFormat {
			case OutputCSV:
				writer := csv.NewWriter(os.Stdout)
				defer writer.Flush()
				err := run.OutputCSV(writer)
				if err != nil {
					outputErrs = append(outputErrs, err)
				}
				break
			case OutputJSON:
				encoder := json.NewEncoder(os.Stdout)
				err := run.OutputJSON(encoder)
				if err != nil {
					outputErrs = append(outputErrs, err)
				}
				break
			case OutputPlain:
				err := run.OutputPlain()
				if err != nil {
					outputErrs = append(outputErrs, err)
				}
				break
			}
		}

		finalErrors := []string{}
		for _, err := range outputErrs {
			finalErrors = append(finalErrors, err.Error())
		}
		if len(finalErrors) > 0 {
			return errors.New(strings.Join(finalErrors, "\n\n"))
		}
		return nil
	}
}
