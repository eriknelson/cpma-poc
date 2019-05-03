package migration

import "fmt"

////////////////////////////////////////////////////////////////////////////////
// Loader part of the ETL, abstracted because it's possible we may want a
// different type of output in the future, it's been suggested HTML
// reports, error reports, in addition to just CRs.
type TransformOutput interface {
	Flush() error
}

// TransformOutput specific for flushing data to a file, in reality this will
// probably be a set of output CRs, but could even have a TransformOutput
// that flushes directly into the OCP4 cluster with the APIs.
type OCP4FileTransformOutput struct {
	OCP4Files []string // Multiple output CRs
}

func (f OCP4FileTransformOutput) Flush() error {
	// TODO: Will need access to the config to know where to write to disk
	fmt.Println("Writing all OCP4Files to disk")
	for _, file := range f.OCP4Files {
		fmt.Printf("%s\n", file)
	}
	return nil
}
