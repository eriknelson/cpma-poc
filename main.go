package main

import (
	"fmt"
)

// Mock data
const MasterConfigFile = "/tmp/openshift.local.clusterup/kube-apiserver/master-config.yaml"
const NodeConfigFile = "/tmp/openshift.local.clusterup/node/node-config.yaml"

////////////////////////////////////////////////////////////////////////////////
// Loader part of the ETL, abstracted because it's possible we may want a
// different type of output in the future
type TransformOutput interface {
	Flush() error
}

// TransformOutput specific for flushing data to a file, in reality this will
// probably be a set of output CRs, but could even have a TransformOutput
// that flushes directly into the OCP4 cluster with the APIs.
type FileTransformOutput struct {
	FileData string
}

func (f FileTransformOutput) Flush() error {
	fmt.Println("Writing file data:")
	fmt.Printf("%s", f.FileData)
	return nil
}

// Simple config type to hold all the user input
type Config struct {
	MasterConfigFile string
	NodeConfigFile   string
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
type Transform interface {
	Run(string) (TransformOutput, error)
	Validate(string) error
	Extract() string
}

type MasterConfigTransform struct {
	RemoteFileName string
}

func (m MasterConfigTransform) Run(extraction string) (TransformOutput, error) {
	fmt.Println("MasterConfigTransform::Run")
	return FileTransformOutput{
		FileData: fmt.Sprintf("[MasterConfigTransform output file contents]", m.RemoteFileName),
	}, nil
}

func (m MasterConfigTransform) Extract() string {
	fmt.Println("MasterConfigTransform::Extract")
	return "[Master config source contents]"
}

func (m MasterConfigTransform) Validate(extraction string) error {
	return nil
}

type NodeConfigTransform struct {
	RemoteFileName string
}

func (n NodeConfigTransform) Run(extraction string) (TransformOutput, error) {
	fmt.Println("NodeConfigTransform::Run")
	return FileTransformOutput{
		FileData: fmt.Sprintf("[NodeConfigTransform output file contents from: %s]", n.RemoteFileName),
	}, nil
}

func (m NodeConfigTransform) Validate(extraction string) error {
	return nil // Simulate fine
}

func (m NodeConfigTransform) Extract() string {
	fmt.Println("NodeConfigTransform::Extract")
	return "[Node config source contents]"
}

////////////////////////////////////////////////////////////////////////////////
// Type responsible for running transforms, possible this could need additional config
type TransformRunner struct {
	Config string
}

// TransformRunner constructor
func NewTransformRunner(config Config) *TransformRunner {
	fmt.Printf("Building TransformRunner with RunnerConfig: %s", config.RunnerConfig)
	return &TransformRunner{Config: config.RunnerConfig}
}

func (r TransformRunner) Run(transforms []Transform) error {
	fmt.Println("TransformRunner::Run")

	// For each transform, extract the data, validate it, and run the transform.
	// Handle any errors, and finally flush the output to it's desired destination
	for _, transform := range transforms {
		extraction := transform.Extract()

		if err := transform.Validate(extraction); err != nil {
			return HandleError(err)
		}

		output, err := transform.Run(extraction)
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
		NodeConfigFile:   NodeConfigFile,
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
		NodeConfigTransform{
			RemoteFileName: NodeConfigFile,
		},
	}); err != nil {
		fmt.Printf("%s", err.Error())
	}
}

func HandleError(err error) error {
	return fmt.Errorf("An error has occurred: %s", err)
}
