package main

import (
	"fmt"
	"github.com/eriknelson/cpma-poc/pkg/migration"
)

// Mock data
const MasterConfigFile = "/tmp/openshift.local.clusterup/kube-apiserver/master-config.yaml"

// Simple config type to hold all the user input
type Config struct {
	MasterConfigFile string
	// RunnerConfig - just using a string as an example, likely would be a struct with more fields
	RunnerConfig string
}

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
	migrationRunner := migration.NewMigrationRunner(config.RunnerConfig)

	if err := migrationRunner.Run([]migration.Migration{
		migration.IdentityProviderMigration{
			MasterConfigFileName: config.MasterConfigFile,
		},
	}); err != nil {
		HandleError(err)
	}

	fmt.Println("fin")
}

func HandleError(err error) error {
	return fmt.Errorf("An error has occurred: %s", err)
}
