package testsupport

import "github.com/didil/kubexcloud/kxc-api/lib"

// BootstrapTests loads config and db for tests
func BootstrapTests(envFilePath string) {
	err := lib.LoadEnv(envFilePath)
	if err != nil {
		panic(err)
	}
}
