package main

import "fmt"

type TransformOutput interface {
	Flush()
}

type TestOutput struct {
	Output string
}

func (t TestOutput) Flush() {
	fmt.Printf("Flushing output %s\n", t.Output)
}

type Extraction interface {
	Transform() (TransformOutput, error)
	Validate() error
}

type Migration interface {
	Extract() Extraction
}

type IdentityProviderExtraction struct {
	Data string
}

func (e IdentityProviderExtraction) Transform() (TransformOutput, error) {
	fmt.Println("IdentityProviderExtraction::Transform")
	// Transform the extraction in some way and created an output
	return TestOutput{Output: e.Data + " output"}, nil
}

type IdentityProviderMigration struct {
	RemoteMasterConfig string
}

// The following doesn't work, must use generic extraction and a type assertion
func (e IdentityProviderExtraction) Validate() error {
	fmt.Printf("Validate remote data: %s\n", e.Data)
	return nil // Passes validation
}

func (m IdentityProviderMigration) Extract() Extraction {
	fmt.Println("IdentityProviderMigration::Extract")
	return IdentityProviderExtraction{
		Data: "some remote file contents",
	}
}

func main() {
	mig := IdentityProviderMigration{
		RemoteMasterConfig: "/tmp/master-config.yaml",
	}

	extraction := mig.Extract()
	if err := extraction.Validate(); err != nil {
		fmt.Println(err)
	}

	output, err := extraction.Transform()
	if err != nil {
		fmt.Println(err)
	}
	output.Flush()
}
