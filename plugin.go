package main

import (
	"fmt"
	"os"
	"os/exec"

	goplugin "github.com/hashicorp/go-plugin"

	"aspect.build/cli/bazel/buildeventstream"
	"aspect.build/cli/pkg/plugin/sdk/v1alpha3/config"
	aspectplugin "aspect.build/cli/pkg/plugin/sdk/v1alpha3/plugin"
)

// main starts up the plugin as a child process of the CLI and connects the gRPC communication.
func main() {
	goplugin.Serve(config.NewConfigFor(&DefaultPlugin{}))
}

// DefaultPlugin declares the fields on an instance of the plugin.
type DefaultPlugin struct {
	// Base gives default implementations of the plugin methods, so implementing them below is optional.
	// See the definition of aspectplugin.Base for more methods that can be implemented by the plugin.
	aspectplugin.Base
}

var hasFirstFailureOccured = false
const fristFailureMetaDateKey = "dd_tags.ci.first_failure.name"

// BEPEventCallback subscribes to all Build Events, and lets our logic react to ones we care about.
func (plugin *DefaultPlugin) BEPEventCallback(event *buildeventstream.BuildEvent) error {
	// If the first failure has not reported yet.
	if isCIEnvironment() && !hasFirstFailureOccured {
		switch event.Payload.(type) {
			case *buildeventstream.BuildEvent_TestResult:
				testResult := *event.GetTestResult()
				if testResult.Status == buildeventstream.TestStatus_FAILED {
					hasFirstFailureOccured = true
					switch event.Id.Id.(type) {
						case *buildeventstream.BuildEventId_TestResult:
							testResultId := *event.Id.GetTestResult()
							reportFirstFailure(testResultId.Label)
					}
				}
		}
	}
	return nil
}

func reportFirstFailure(failedTarget string) {
	cmd := exec.Command("buildkite-agent", "meta-data", "exists", fristFailureMetaDateKey)
	err := cmd.Run()
	if err != nil {
		cmd := exec.Command("buildkite-agent", "meta-data", "set", fristFailureMetaDateKey, failedTarget)
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("Buildkite first failure meta-data already exists")
	}
}

func isCIEnvironment() bool {
	ci := os.Getenv("CI")
	return ci != ""
}