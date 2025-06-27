package integration

import (
	"flag"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

var testContext *TestContext
var opts = godog.Options{Output: colors.Colored(os.Stdout)}

func init() {
	godog.BindCommandLineFlags("godog.", &opts)
}

func TestMain(m *testing.M) {
	flag.Parse()
	opts.Paths = flag.Args()

	status := godog.TestSuite{
		Name:                 "fizz-and-buzz",
		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer:  InitializeScenario,
		Options:              &opts,
	}.Run()

	os.Exit(status)
}

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	testContext = &TestContext{}
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(testContext.Before)
	ctx.After(testContext.After)

	// Given
	ctx.Step(`^now is the "([^"]*)"$`, nowIs)

	// When
	ctx.Step(`^Fizz and Buzz service has a zero context of "([^"]*)"$`, fizzAndBuzzServiceHasAZeroOf)
	ctx.Step(`^I send a "([^"]*)" request to "([^"]*)"$`, iSendARequestTo)
	ctx.Step(`^I send a "([^"]*)" request to "([^"]*)" with the given params:$`, iSendARequestToWithTheGivenParams)
	ctx.Step(`^I send "([^"]*)" requests to "([^"]*)" with the given params:$`, iSendRequestsToWithTheGivenParams)

	// Then
	ctx.Step(`^the response code should be (\d+)$`, theResponseCodeShouldBe)
	ctx.Step(`^the response codes should be (\d+)$`, theResponseCodesShouldBe)
	ctx.Step(`^the response should match json:$`, theResponseShouldMatchJSON)
}
