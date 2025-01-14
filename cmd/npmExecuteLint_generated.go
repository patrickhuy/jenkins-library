// Code generated by piper's step-generator. DO NOT EDIT.

package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/SAP/jenkins-library/pkg/config"
	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/SAP/jenkins-library/pkg/splunk"
	"github.com/SAP/jenkins-library/pkg/telemetry"
	"github.com/SAP/jenkins-library/pkg/validation"
	"github.com/spf13/cobra"
)

type npmExecuteLintOptions struct {
	Install            bool   `json:"install,omitempty"`
	RunScript          string `json:"runScript,omitempty"`
	FailOnError        bool   `json:"failOnError,omitempty"`
	DefaultNpmRegistry string `json:"defaultNpmRegistry,omitempty"`
}

// NpmExecuteLintCommand Execute ci-lint script on all npm packages in a project or execute default linting
func NpmExecuteLintCommand() *cobra.Command {
	const STEP_NAME = "npmExecuteLint"

	metadata := npmExecuteLintMetadata()
	var stepConfig npmExecuteLintOptions
	var startTime time.Time
	var logCollector *log.CollectorHook

	var createNpmExecuteLintCmd = &cobra.Command{
		Use:   STEP_NAME,
		Short: "Execute ci-lint script on all npm packages in a project or execute default linting",
		Long: `Execute ci-lint script for all package json files, if they implement the script. If no ci-lint script is defined,
either use ESLint configurations present in the project or use the provided general purpose configuration to run ESLint.`,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			startTime = time.Now()
			log.SetStepName(STEP_NAME)
			log.SetVerbose(GeneralConfig.Verbose)

			GeneralConfig.GitHubAccessTokens = ResolveAccessTokens(GeneralConfig.GitHubTokens)

			path, _ := os.Getwd()
			fatalHook := &log.FatalHook{CorrelationID: GeneralConfig.CorrelationID, Path: path}
			log.RegisterHook(fatalHook)

			err := PrepareConfig(cmd, &metadata, STEP_NAME, &stepConfig, config.OpenPiperFile)
			if err != nil {
				log.SetErrorCategory(log.ErrorConfiguration)
				return err
			}

			if len(GeneralConfig.HookConfig.SentryConfig.Dsn) > 0 {
				sentryHook := log.NewSentryHook(GeneralConfig.HookConfig.SentryConfig.Dsn, GeneralConfig.CorrelationID)
				log.RegisterHook(&sentryHook)
			}

			if len(GeneralConfig.HookConfig.SplunkConfig.Dsn) > 0 {
				logCollector = &log.CollectorHook{CorrelationID: GeneralConfig.CorrelationID}
				log.RegisterHook(logCollector)
			}

			validation, err := validation.New(validation.WithJSONNamesForStructFields(), validation.WithPredefinedErrorMessages())
			if err != nil {
				return err
			}
			if err = validation.ValidateStruct(stepConfig); err != nil {
				log.SetErrorCategory(log.ErrorConfiguration)
				return err
			}

			return nil
		},
		Run: func(_ *cobra.Command, _ []string) {
			telemetryData := telemetry.CustomData{}
			telemetryData.ErrorCode = "1"
			handler := func() {
				config.RemoveVaultSecretFiles()
				telemetryData.Duration = fmt.Sprintf("%v", time.Since(startTime).Milliseconds())
				telemetryData.ErrorCategory = log.GetErrorCategory().String()
				telemetry.Send(&telemetryData)
				if len(GeneralConfig.HookConfig.SplunkConfig.Dsn) > 0 {
					splunk.Send(&telemetryData, logCollector)
				}
			}
			log.DeferExitHandler(handler)
			defer handler()
			telemetry.Initialize(GeneralConfig.NoTelemetry, STEP_NAME)
			if len(GeneralConfig.HookConfig.SplunkConfig.Dsn) > 0 {
				splunk.Initialize(GeneralConfig.CorrelationID,
					GeneralConfig.HookConfig.SplunkConfig.Dsn,
					GeneralConfig.HookConfig.SplunkConfig.Token,
					GeneralConfig.HookConfig.SplunkConfig.Index,
					GeneralConfig.HookConfig.SplunkConfig.SendLogs)
			}
			npmExecuteLint(stepConfig, &telemetryData)
			telemetryData.ErrorCode = "0"
			log.Entry().Info("SUCCESS")
		},
	}

	addNpmExecuteLintFlags(createNpmExecuteLintCmd, &stepConfig)
	return createNpmExecuteLintCmd
}

func addNpmExecuteLintFlags(cmd *cobra.Command, stepConfig *npmExecuteLintOptions) {
	cmd.Flags().BoolVar(&stepConfig.Install, "install", false, "Run npm install or similar commands depending on the project structure.")
	cmd.Flags().StringVar(&stepConfig.RunScript, "runScript", `ci-lint`, "List of additional run scripts to execute from package.json.")
	cmd.Flags().BoolVar(&stepConfig.FailOnError, "failOnError", false, "Defines the behavior in case linting errors are found.")
	cmd.Flags().StringVar(&stepConfig.DefaultNpmRegistry, "defaultNpmRegistry", os.Getenv("PIPER_defaultNpmRegistry"), "URL of the npm registry to use. Defaults to https://registry.npmjs.org/")

}

// retrieve step metadata
func npmExecuteLintMetadata() config.StepData {
	var theMetaData = config.StepData{
		Metadata: config.StepMetadata{
			Name:        "npmExecuteLint",
			Aliases:     []config.Alias{{Name: "executeNpm", Deprecated: false}},
			Description: "Execute ci-lint script on all npm packages in a project or execute default linting",
		},
		Spec: config.StepSpec{
			Inputs: config.StepInputs{
				Parameters: []config.StepParameters{
					{
						Name:        "install",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "bool",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     false,
					},
					{
						Name:        "runScript",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     `ci-lint`,
					},
					{
						Name:        "failOnError",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "bool",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     false,
					},
					{
						Name:        "defaultNpmRegistry",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "GENERAL", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{{Name: "npm/defaultNpmRegistry"}},
						Default:     os.Getenv("PIPER_defaultNpmRegistry"),
					},
				},
			},
			Containers: []config.Container{
				{Name: "node", Image: "node:lts-stretch"},
			},
		},
	}
	return theMetaData
}
