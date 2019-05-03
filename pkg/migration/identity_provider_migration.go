package migration

import (
	"fmt"
	"github.com/eriknelson/cpma-poc/pkg/util"
)

type IdentityProviderMigration struct {
	MasterConfigFileName       string
	ExtractedIdentityProviders []string
}

func (i IdentityProviderMigration) Extract() error {
	fmt.Println("IdentityProviderMigration::Extract")
	// Here is where you download all the data you need by crawling the remote files
	// Then build up the struct with all the downloaded data
	masterConfig := util.GetMasterConfig(i.MasterConfigFileName)

	fmt.Printf("%s", masterConfig)
	fmt.Println("Reading identity provider related data based on downloaded masterconfig data")
	i.ExtractedIdentityProviders = []string{"foo", "bar"}

	return nil
}

func (i IdentityProviderMigration) Validate() error {
	fmt.Println("Validating migration struct before executing Transform")
	return nil // Simulate fine
}

func (i IdentityProviderMigration) Transform() (TransformOutput, error) {
	fmt.Println("IdentityProviderMigration::Transform")
	return OCP4FileTransformOutput{
		OCP4Files: []string{
			fmt.Sprintf("Secret from IdentityProviders: %s", i.ExtractedIdentityProviders),
		},
	}, nil
}
