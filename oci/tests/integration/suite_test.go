//go:build integration
// +build integration

/*
Copyright 2022 The Flux authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package integration

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	tfjson "github.com/hashicorp/terraform-json"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/fluxcd/test-infra/tftestenv"
)

const (
	// terraformPathAWS is the path to the terraform working directory
	// containing the aws terraform configurations.
	terraformPathAWS = "./terraform/aws"
	// terraformPathAzure is the path to the terraform working directory
	// containing the azure terraform configurations.
	terraformPathAzure = "./terraform/azure"
	// terraformPathGCP is the path to the terraform working directory
	// containing the gcp terraform configurations.
	terraformPathGCP = "./terraform/gcp"
	// kubeconfigPath is the path where the cluster kubeconfig is written to and
	// used from.
	kubeconfigPath = "./build/kubeconfig"

	resultWaitTimeout = 30 * time.Second
)

var (
	// supportedProviders are the providers supported by the test.
	supportedProviders = []string{"aws", "azure", "gcp"}

	// targetProvider is the name of the kubernetes provider to test against.
	targetProvider = flag.String("provider", "", fmt.Sprintf("name of the provider %v", supportedProviders))

	// retain flag to prevent destroy and retaining the created infrastructure.
	retain = flag.Bool("retain", false, "retain the infrastructure for debugging purposes")

	// existing flag to use existing infrastructure terraform state.
	existing = flag.Bool("existing", false, "use existing infrastructure state for debugging purposes")

	// verbose flag to enable output of terraform execution.
	verbose = flag.Bool("verbose", false, "verbose output of the environment setup")

	// testRepos is a map of registry common name and URL of the test
	// repositories. This is used as the test cases to run the tests against.
	// The registry common name need not be the actual registry address but an
	// identifier to identify the test case without logging any sensitive
	// account IDs in the subtest names.
	// For example, map[string]string{"ecr", "xxxxx.dkr.ecr.xxxx.amazonaws.com/foo:v1"}
	// would result in subtest name TestImageRepositoryScanAWS/ecr.
	testRepos map[string]string

	// testEnv is the test environment. It contains test infrastructure and
	// kubernetes client of the created cluster.
	testEnv *tftestenv.Environment

	// testImageTags are the tags used in the test for the generated images.
	testImageTags = []string{"v0.1.0", "v0.1.2", "v0.1.3", "v0.1.4"}

	testAppImage string
)

// registryLoginFunc is used to perform registry login against a provider based
// on the terraform state output values. It returns a map of registry common
// name and test repositories to test against, read from the terraform state
// output.
type registryLoginFunc func(ctx context.Context, output map[string]*tfjson.StateOutput) (map[string]string, error)

// pushTestImages is used to push local flux test images to a remote registry
// after logging in using registryLoginFunc. It takes a map of image name and
// local images and terraform state output. The local images are retagged and
// pushed to a corresponding registry repository for the image.
type pushTestImages func(ctx context.Context, localImgs map[string]string, output map[string]*tfjson.StateOutput) (map[string]string, error)

// ProviderConfig is the test configuration of a supported cloud provider to run
// the tests against.
type ProviderConfig struct {
	// terraformPath is the path to the directory containing the terraform
	// configurations of the provider.
	terraformPath string
	// registryLogin is used to perform registry login.
	registryLogin registryLoginFunc
	// createKubeconfig is used to create kubeconfig of a cluster.
	createKubeconfig tftestenv.CreateKubeconfig
	// pushAppTestImages is used to push flux test images to a remote registry.
	pushAppTestImages pushTestImages
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz1234567890")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func TestMain(m *testing.M) {
	flag.Parse()
	ctx := context.TODO()

	appImg := os.Getenv("TEST_IMG")
	if appImg == "" {
		log.Fatal("TEST_IMG must be set to the test application image, cannot be empty")
	}

	localImgs := map[string]string{
		"app": appImg,
	}

	// Validate the provider.
	if *targetProvider == "" {
		log.Fatalf("-provider flag must be set to one of %v", supportedProviders)
	}
	var supported bool
	for _, p := range supportedProviders {
		if p == *targetProvider {
			supported = true
		}
	}
	if !supported {
		log.Fatalf("Unsupported provider %q, must be one of %v", *targetProvider, supportedProviders)
	}

	providerCfg := getProviderConfig(*targetProvider)
	if providerCfg == nil {
		log.Fatalf("Failed to get provider config for %q", *targetProvider)
	}

	// Construct scheme to be added to the kubeclient.
	scheme := runtime.NewScheme()

	err := batchv1.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}

	// Create environment.
	envOpts := []tftestenv.EnvironmentOption{
		tftestenv.WithVerbose(*verbose),
		tftestenv.WithRetain(*retain),
		tftestenv.WithExisting(*existing),
		tftestenv.WithCreateKubeconfig(providerCfg.createKubeconfig),
	}
	// var err error
	testEnv, err = tftestenv.New(ctx, scheme, providerCfg.terraformPath, kubeconfigPath, envOpts...)
	if err != nil {
		panic(fmt.Sprintf("Failed to provision the test infrastructure: %v", err))
	}

	// Get terraform state output.
	output, err := testEnv.StateOutput(ctx)
	if err != nil {
		panic(fmt.Sprintf("Failed to get the terraform state output: %v", err))
	}

	testRepos, err = providerCfg.registryLogin(ctx, output)
	if err != nil {
		panic(fmt.Sprintf("Failed to log into registry: %v", err))
	}

	pushedImages, err := providerCfg.pushAppTestImages(ctx, localImgs, output)
	if err != nil {
		panic(fmt.Sprintf("Failed to push test images: %v", err))
	}

	if len(pushedImages) != 1 {
		panic(fmt.Sprintf("Unexpected number of app images pushed: %d", len(pushedImages)))
	}

	if appImg, ok := pushedImages["app"]; !ok {
		panic(fmt.Sprintf("Could not find pushed app image in %v", pushedImages))
	} else {
		testAppImage = appImg
	}

	// Create and push test images.
	if err := tftestenv.CreateAndPushImages(testRepos, testImageTags); err != nil {
		panic(fmt.Sprintf("Failed to create and push images: %v", err))
	}

	code := m.Run()

	if err := testEnv.Stop(ctx); err != nil {
		log.Printf("Failed to stop environment: %v", err)
	}
	os.Exit(code)
}

// getProviderConfig returns the test configuration of supported providers.
func getProviderConfig(provider string) *ProviderConfig {
	switch provider {
	case "aws":
		return &ProviderConfig{
			terraformPath:     terraformPathAWS,
			registryLogin:     registryLoginECR,
			pushAppTestImages: pushAppTestImagesECR,
			createKubeconfig:  createKubeconfigEKS,
		}
	case "azure":
		return &ProviderConfig{
			terraformPath:     terraformPathAzure,
			registryLogin:     registryLoginACR,
			pushAppTestImages: pushAppTestImagesACR,
			createKubeconfig:  createKubeConfigAKS,
		}
	case "gcp":
		return &ProviderConfig{
			terraformPath:     terraformPathGCP,
			registryLogin:     registryLoginGCR,
			pushAppTestImages: pushAppTestImagesGCR,
			createKubeconfig:  createKubeconfigGKE,
		}
	}
	return nil
}
