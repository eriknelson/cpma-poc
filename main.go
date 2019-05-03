package main

import (
	"fmt"
)

// Mock data
const MasterConfigFile = "/tmp/openshift.local.clusterup/kube-apiserver/master-config.yaml"

////////////////////////////////////////////////////////////////////////////////
// Loader part of the ETL, abstracted because it's possible we may want a
// different type of output in the future
type TransformOutput interface {
	Flush() error
}

// TransformOutput specific for flushing data to a file, in reality this will
// probably be a set of output CRs, but could even have a TransformOutput
// that flushes directly into the OCP4 cluster with the APIs.
type OCP4FileTransformOutput struct {
	OCP4Files []string
}

func (f OCP4FileTransformOutput) Flush() error {
	// Probably need access to the config to know where to write
	fmt.Println("Writing all OCP4Files to disk...")
	for _, file := range f.OCP4Files {
		fmt.Printf("%s\n", file)
	}
	return nil
}

// Simple config type to hold all the user input
type Config struct {
	MasterConfigFile string
	// RunnerConfig - just using a string as an example, likely would be a struct with more fields
	RunnerConfig string
}

////////////////////////////////////////////////////////////////////////////////
// Transform interface and implementations - this is the bulk of the work
// I could see all the transforms that have file sources being composed of
// something that implements the shared behavior
//
// A Transform encapsulates the details of transforming a single remote source with the
// intention that the bulk of the work will be in defining the details of
// individual transforms. Makes test relatively easy as well, and codifies
// the transforms to strong types.
// Undecided if extract's return type should be different here, I'm not seeing
// a case where it's anything but a string. Could be a struct with more fields.
////////////////////////////////////////////////////////////////////////////////
type Extraction interface {
	Transform() (TransformOutput, error)
}

type Transform interface {
	Validate(Extraction) error
	Extract() Extraction
}

type MasterConfigExtraction struct {
	MasterConfigFileContents string
	IdentityProviders        []string
	OtherDownloadedData      string
}

type MasterConfigTransform struct {
	RemoteFileName         string
	MasterConfigExtraction MasterConfigExtraction
}

func (m MasterConfigTransform) Extract() Extraction {
	fmt.Println("MasterConfigTransform::Extract")
	// Here is where you download all the data you need by crawling the remote files
	// Then build up an Extraction with all the downloaded data
	fmt.Printf("Downloading the master config file from: %s\n", m.RemoteFileName)
	fmt.Printf("Reading identity providers and downloading relevant files now that we have the master config content\n")

	return MasterConfigExtraction{
		MasterConfigFileContents: "[Master config file contents]",
		IdentityProviders:        []string{"HTPassword"},
		OtherDownloadedData:      "other downloaded data",
	}
}

func (m MasterConfigExtraction) Transform() (TransformOutput, error) {
	fmt.Println("MasterConfigExtraction::Transform")
	return OCP4FileTransformOutput{
		OCP4Files: []string{
			// Could forsee functions called here that take master contents and turn it into ocp4 file
			// Ex: transformMasterConfigFile(ocp3FileContent) ocp4CRFile
			fmt.Sprintf("CR from master config contents: %s", m.MasterConfigFileContents),
			fmt.Sprintf("Secret from IdentityProviders: %s", m.IdentityProviders),
		},
	}, nil
}

func (m MasterConfigTransform) Validate(extraction Extraction) error {
	fmt.Println("Looking at the MasterConfigExtraction and validating")
	return nil // Simulate fine
}

////////////////////////////////////////////////////////////////////////////////
// Type responsible for running transforms, possible this could need additional config
type TransformRunner struct {
	Config string
}

// TransformRunner constructor
func NewTransformRunner(config Config) *TransformRunner {
	fmt.Printf("Building TransformRunner with RunnerConfig: %s\n", config.RunnerConfig)
	return &TransformRunner{Config: config.RunnerConfig}
}

func (r TransformRunner) Run(transforms []Transform) error {
	fmt.Println("TransformRunner::Run")

	// For each transform, extract the data, validate it, and run the transform.
	// Handle any errors, and finally flush the output to it's desired destination
	// NOTE: This should be parallelized with channels unless the transforms have
	// some dependency on the outputs of others
	for _, transform := range transforms {
		extraction := transform.Extract()

		if err := transform.Validate(extraction); err != nil {
			return HandleError(err)
		}

		output, err := extraction.Transform()
		if err != nil {
			HandleError(err)
		}

		if err := output.Flush(); err != nil {
			HandleError(err)
		}
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func LoadConfig() Config {
	// Mocking out the details of collecting cli input and file input
	config := Config{
		MasterConfigFile: MasterConfigFile,
		RunnerConfig:     "some_runner_config",
	}

	fmt.Println("Loaded config")
	return config
}

func main() {
	config := LoadConfig()
	transformRunner := NewTransformRunner(config)

	if err := transformRunner.Run([]Transform{
		MasterConfigTransform{
			RemoteFileName: MasterConfigFile,
		},
	}); err != nil {
		fmt.Printf("%s", err.Error())
	}

	fmt.Println("fin")
}

func HandleError(err error) error {
	return fmt.Errorf("An error has occurred: %s", err)
}
