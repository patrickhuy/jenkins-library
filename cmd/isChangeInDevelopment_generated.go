// Code generated by piper's step-generator. DO NOT EDIT.

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/SAP/jenkins-library/pkg/config"
	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/SAP/jenkins-library/pkg/piperenv"
	"github.com/SAP/jenkins-library/pkg/splunk"
	"github.com/SAP/jenkins-library/pkg/telemetry"
	"github.com/SAP/jenkins-library/pkg/validation"
	"github.com/spf13/cobra"
)

type isChangeInDevelopmentOptions struct {
	Endpoint                       string   `json:"endpoint,omitempty"`
	Username                       string   `json:"username,omitempty"`
	Password                       string   `json:"password,omitempty"`
	ChangeDocumentID               string   `json:"changeDocumentId,omitempty"`
	FailIfStatusIsNotInDevelopment bool     `json:"failIfStatusIsNotInDevelopment,omitempty"`
	CmClientOpts                   []string `json:"cmClientOpts,omitempty"`
}

type isChangeInDevelopmentCommonPipelineEnvironment struct {
	custom struct {
		isChangeInDevelopment bool
	}
}

func (p *isChangeInDevelopmentCommonPipelineEnvironment) persist(path, resourceName string) {
	content := []struct {
		category string
		name     string
		value    interface{}
	}{
		{category: "custom", name: "isChangeInDevelopment", value: p.custom.isChangeInDevelopment},
	}

	errCount := 0
	for _, param := range content {
		err := piperenv.SetResourceParameter(path, resourceName, filepath.Join(param.category, param.name), param.value)
		if err != nil {
			log.Entry().WithError(err).Error("Error persisting piper environment.")
			errCount++
		}
	}
	if errCount > 0 {
		log.Entry().Fatal("failed to persist Piper environment")
	}
}

// IsChangeInDevelopmentCommand Checks if a certain change is in status 'in development'
func IsChangeInDevelopmentCommand() *cobra.Command {
	const STEP_NAME = "isChangeInDevelopment"

	metadata := isChangeInDevelopmentMetadata()
	var stepConfig isChangeInDevelopmentOptions
	var startTime time.Time
	var commonPipelineEnvironment isChangeInDevelopmentCommonPipelineEnvironment
	var logCollector *log.CollectorHook
	splunkClient := &splunk.Splunk{}
	telemetryClient := &telemetry.Telemetry{}

	var createIsChangeInDevelopmentCmd = &cobra.Command{
		Use:   STEP_NAME,
		Short: "Checks if a certain change is in status 'in development'",
		Long:  `"Checks if a certain change is in status 'in development'"`,
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
			log.RegisterSecret(stepConfig.Username)
			log.RegisterSecret(stepConfig.Password)

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
			customTelemetryData := telemetry.CustomData{}
			customTelemetryData.ErrorCode = "1"
			handler := func() {
				config.RemoveVaultSecretFiles()
				commonPipelineEnvironment.persist(GeneralConfig.EnvRootPath, "commonPipelineEnvironment")
				customTelemetryData.Duration = fmt.Sprintf("%v", time.Since(startTime).Milliseconds())
				customTelemetryData.ErrorCategory = log.GetErrorCategory().String()
				telemetryClient.SetData(&customTelemetryData)
				telemetryClient.Send()
				if len(GeneralConfig.HookConfig.SplunkConfig.Dsn) > 0 {
					splunkClient.Send(telemetryClient.GetData(), logCollector)
				}
			}
			log.DeferExitHandler(handler)
			defer handler()
			telemetryClient.Initialize(GeneralConfig.NoTelemetry, STEP_NAME)
			if len(GeneralConfig.HookConfig.SplunkConfig.Dsn) > 0 {
				splunkClient.Initialize(GeneralConfig.CorrelationID,
					GeneralConfig.HookConfig.SplunkConfig.Dsn,
					GeneralConfig.HookConfig.SplunkConfig.Token,
					GeneralConfig.HookConfig.SplunkConfig.Index,
					GeneralConfig.HookConfig.SplunkConfig.SendLogs)
			}
			isChangeInDevelopment(stepConfig, &customTelemetryData, &commonPipelineEnvironment)
			customTelemetryData.ErrorCode = "0"
			log.Entry().Info("SUCCESS")
		},
	}

	addIsChangeInDevelopmentFlags(createIsChangeInDevelopmentCmd, &stepConfig)
	return createIsChangeInDevelopmentCmd
}

func addIsChangeInDevelopmentFlags(cmd *cobra.Command, stepConfig *isChangeInDevelopmentOptions) {
	cmd.Flags().StringVar(&stepConfig.Endpoint, "endpoint", os.Getenv("PIPER_endpoint"), "The service endpoint")
	cmd.Flags().StringVar(&stepConfig.Username, "username", os.Getenv("PIPER_username"), "Service user to authenticate against the ABAP backend")
	cmd.Flags().StringVar(&stepConfig.Password, "password", os.Getenv("PIPER_password"), "Service user password to authenticate against the ABAP backend")
	cmd.Flags().StringVar(&stepConfig.ChangeDocumentID, "changeDocumentId", os.Getenv("PIPER_changeDocumentId"), "ID of the change document to be checked for the status")
	cmd.Flags().BoolVar(&stepConfig.FailIfStatusIsNotInDevelopment, "failIfStatusIsNotInDevelopment", true, "lets the build fail in case the change is not in status 'in developent'. Otherwise a warning is emitted to the log")
	cmd.Flags().StringSliceVar(&stepConfig.CmClientOpts, "cmClientOpts", []string{}, "additional options passed to cm client, e.g. for troubleshooting")

	cmd.MarkFlagRequired("endpoint")
	cmd.MarkFlagRequired("username")
	cmd.MarkFlagRequired("password")
	cmd.MarkFlagRequired("changeDocumentId")
}

// retrieve step metadata
func isChangeInDevelopmentMetadata() config.StepData {
	var theMetaData = config.StepData{
		Metadata: config.StepMetadata{
			Name:        "isChangeInDevelopment",
			Aliases:     []config.Alias{},
			Description: "Checks if a certain change is in status 'in development'",
		},
		Spec: config.StepSpec{
			Inputs: config.StepInputs{
				Secrets: []config.StepSecrets{
					{Name: "credentialsId", Description: "Jenkins 'Username with password' credentials ID containing user and password to authenticate against the ABAP backend", Type: "jenkins", Aliases: []config.Alias{{Name: "changeManagement/credentialsId", Deprecated: false}}},
				},
				Parameters: []config.StepParameters{
					{
						Name:        "endpoint",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS", "GENERAL"},
						Type:        "string",
						Mandatory:   true,
						Aliases:     []config.Alias{{Name: "changeManagement/endpoint"}},
						Default:     os.Getenv("PIPER_endpoint"),
					},
					{
						Name: "username",
						ResourceRef: []config.ResourceReference{
							{
								Name:  "credentialsId",
								Param: "username",
								Type:  "secret",
							},
						},
						Scope:     []string{"PARAMETERS", "STAGES", "STEPS", "GENERAL"},
						Type:      "string",
						Mandatory: true,
						Aliases:   []config.Alias{},
						Default:   os.Getenv("PIPER_username"),
					},
					{
						Name: "password",
						ResourceRef: []config.ResourceReference{
							{
								Name:  "credentialsId",
								Param: "password",
								Type:  "secret",
							},
						},
						Scope:     []string{"PARAMETERS"},
						Type:      "string",
						Mandatory: true,
						Aliases:   []config.Alias{},
						Default:   os.Getenv("PIPER_password"),
					},
					{
						Name: "changeDocumentId",
						ResourceRef: []config.ResourceReference{
							{
								Name:  "commonPipelineEnvironment",
								Param: "custom/changeDocumentId",
							},
						},
						Scope:     []string{"PARAMETERS"},
						Type:      "string",
						Mandatory: true,
						Aliases:   []config.Alias{},
						Default:   os.Getenv("PIPER_changeDocumentId"),
					},
					{
						Name:        "failIfStatusIsNotInDevelopment",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "bool",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     true,
					},
					{
						Name:        "cmClientOpts",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "[]string",
						Mandatory:   false,
						Aliases:     []config.Alias{{Name: "clientOpts"}, {Name: "changeManagement/clientOpts"}},
						Default:     []string{},
					},
				},
			},
			Containers: []config.Container{
				{Name: "cmclient", Image: "ppiper/cm-client"},
			},
			Outputs: config.StepOutputs{
				Resources: []config.StepResources{
					{
						Name: "commonPipelineEnvironment",
						Type: "piperEnvironment",
						Parameters: []map[string]interface{}{
							{"Name": "custom/isChangeInDevelopment"},
						},
					},
				},
			},
		},
	}
	return theMetaData
}
