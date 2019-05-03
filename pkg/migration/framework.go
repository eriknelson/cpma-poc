package migration

import "fmt"

type Migration interface {
	Extract() error
	Validate() error
	Transform() (TransformOutput, error)
}

type MigrationRunner struct {
	Config string
}

// MigrationRunner constructor
func NewMigrationRunner(runnerConfig string) *MigrationRunner {
	fmt.Printf("Building MigrationRunner with RunnerConfig: %s\n", runnerConfig)
	return &MigrationRunner{Config: runnerConfig}
}

func (r MigrationRunner) Run(migrations []Migration) error {
	fmt.Println("MigrationRunner::Run")

	// For each migration, extract the data, validate it, and run the transform.
	// Handle any errors, and finally flush the output to it's desired destination
	// NOTE: This should be parallelized with channels unless the transforms have
	// some dependency on the outputs of others
	for _, mig := range migrations {
		if err := mig.Extract(); err != nil {
			return err
		}

		if err := mig.Validate(); err != nil {
			return err
		}

		output, err := mig.Transform()
		if err != nil {
			return err
		}

		if err := output.Flush(); err != nil {
			return err
		}
	}

	return nil
}
