package godog_example_setup

import (
	"context"
	"log"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/joho/godotenv"
	"github.com/pawelWritesCode/gdutils"
	"github.com/spf13/pflag"

	"github.com/pawelWritesCode/godog-example-setup/defs"
)

const (
	//envDebug describes environment variable responsible for debug mode - (true/false)
	envDebug = "GODOG_DEBUG"

	// envMyAppURL describes URL to "My app" - should be valid URL
	envMyAppURL = "GODOG_MY_APP_URL"

	// envJsonSchemaDir path to JSON schemas dir - relative path from project root
	envJsonSchemaDir = "GODOG_JSON_SCHEMA_DIR"
)

// opt defines options for godog CLI while running tests from "go test" command.
var opt = godog.Options{Output: colors.Colored(os.Stdout), Format: "progress", Randomize: time.Now().UTC().UnixNano()}

func init() {
	godog.BindCommandLineFlags("godog.", &opt)
	checkErr(godotenv.Load()) // loading environment variables from .env file
}

func TestMain(m *testing.M) {
	pflag.Parse()
	opt.Paths = pflag.Args()

	status := godog.TestSuite{Name: "godogs", ScenarioInitializer: InitializeScenario, Options: &opt}.Run()

	os.Exit(status)
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	isDebug := strings.ToLower(os.Getenv(envDebug)) == "true"
	wd, err := os.Getwd()
	checkErr(err)

	/*
		Variable scenario represents godog scenario.
		Scenario has State field with plenty of utility services useful to create godog steps.

		If you would like to replace default state service with your own, do following:

		Let's assume, you want to replace default debugger with your own. Field scenario.State.Debugger has interface
		type of debugger.Debugger. Create your own struct, implement on it interface debugger.Debugger and then, use
		proper setter method to inject it into scenario.State. In this example it would be: scenario.State.SetDebugger.

		In another example you may want to use your own http.Client. So you have to create your own implementation of
		httpctx.HttpContext interface with your custom http.Client and then inject it to State with scenario.State.SetHttpContext
	*/
	scenario := defs.Scenario{State: gdutils.NewDefaultState(isDebug, path.Join(wd, os.Getenv(envJsonSchemaDir)))}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		scenario.State.ResetState(isDebug)

		// Here you can define more scenario-scoped values using scenario.State.Cache.Save() method
		scenario.State.Cache.Save("MY_APP_URL", os.Getenv(envMyAppURL))

		return ctx, nil
	})

	/*
	   |--------------------------------------------------------------------------
	   | Random data generation
	   |--------------------------------------------------------------------------
	   |
	   | This section contains utility methods for random data generation.
	   | Those methods contains creation of
	   | - random length runes of ASCII/UNICODE characters
	   | - random length sentence of ASCII/UNICODE words
	   | - int/float from provided range.
	   |
	   | Every method saves its output in state's cache under provided key
	*/
	ctx.Step(`^I generate a random word having from "(\d+)" to "(\d+)" of "(ASCII|UNICODE)" characters and save it as "([^"]*)"$`, scenario.IGenerateARandomRunesOfLengthWithCharactersAndSaveItAs)
	ctx.Step(`^I generate a random sentence having from "(\d+)" to "(\d+)" of "(ASCII|UNICODE)" words and save it as "([^"]*)"$`, scenario.IGenerateARandomSentenceInTheRangeFromToWordsAndSaveItAs(3, 10))
	ctx.Step(`^I generate a random "(int|float)" in the range from "(\d+)" to "(\d+)" and save it as "([^"]*)"$`, scenario.IGenerateARandomNumberInTheRangeFromToAndSaveItAs)

	/*
	   |--------------------------------------------------------------------------
	   | Sending HTTP(s) requests
	   |--------------------------------------------------------------------------
	   |
	   | This section contains methods for preparing and sending HTTP(s) requests.
	   |
	   | Methods that start with "I set following ..." accept docstring in form of JSON.
	*/

	ctx.Step(`^I prepare new "(GET|POST|PUT|PATCH|DELETE|HEAD)" request to "([^"]*)" and save it as "([^"]*)"$`, scenario.IPrepareNewRequestToAndSaveItAs)
	ctx.Step(`^I set following headers for prepared request "([^"]*)":$`, scenario.ISetFollowingHeadersForPreparedRequest)
	ctx.Step(`^I set following body for prepared request "([^"]*)":$`, scenario.ISetFollowingBodyForPreparedRequest)
	ctx.Step(`^I send request "([^"]*)"$`, scenario.ISendRequest)

	// this method accepts docstring in form of JSON with two keys: "body" and "headers"
	ctx.Step(`^I send "(GET|POST|PUT|PATCH|DELETE|HEAD)" request to "([^"]*)" with body and headers:$`, scenario.ISendRequestToWithBodyAndHeaders)

	/*
	   |--------------------------------------------------------------------------
	   | Assertions
	   |--------------------------------------------------------------------------
	   |
	   | This section contains assertions against last HTTP(s) responses.
	   | Those include assertions against:
	   | - response body JSON nodes,
	   | - HTTP(s) headers,
	   | - status code.
	   |
	   | Every argument following immediately after word "node" or "nodes"
	   | should have syntax acceptable by one of json-path libraries:
	   | https://github.com/pawelWritesCode/qjson
	   | https://github.com/oliveagle/jsonpath
	   |
	   | Method "the JSON response should have nodes" accepts list of nodes,
	   | separated with comma ",". For example: "data[0].user, $.data[1].user, data".
	   |
	   | Every method, that ends with 'of value "([^"]*)"' accepts fixed values or
	   | template values. Template values are references to previously saved values
	   | in scenario cache.
	   | For example, after generating random string and saving it under "USER_NAME",
	   | argument of this method may be: "{{.USER_NAME}}"
	   |
	   | Argument in method starting with 'time between ...' should be string valid for
	   | golang standard library time.ParseDuration func, for example: 3s, 1h, 30ms
	*/
	ctx.Step(`^the response should have header "([^"]*)"$`, scenario.TheResponseShouldHaveHeader)
	ctx.Step(`^the response should have header "([^"]*)" of value "([^"]*)"$`, scenario.TheResponseShouldHaveHeaderOfValue)

	ctx.Step(`^the response status code should be (\d+)$`, scenario.TheResponseStatusCodeShouldBe)

	ctx.Step(`^the JSON response should have nodes "([^"]*)"$`, scenario.TheJSONResponseShouldHaveNodes)
	ctx.Step(`^the JSON response should have node "([^"]*)"$`, scenario.TheJSONResponseShouldHaveNodes)

	ctx.Step(`^the JSON node "([^"]*)" should be "(string|int|float|bool)" of value "([^"]*)"$`, scenario.TheJSONNodeShouldBeOfValue)
	ctx.Step(`^the JSON node "([^"]*)" should be slice of length "(\d+)"$`, scenario.TheJSONNodeShouldBeSliceOfLength)
	ctx.Step(`^the JSON node "([^"]*)" should be "(nil|string|int|float|bool|map|slice)"$`, scenario.TheJSONNodeShouldBe)
	ctx.Step(`^the JSON node "([^"]*)" should not be "(nil|string|int|float|bool|map|slice)"$`, scenario.TheJSONNodeShouldNotBe)

	ctx.Step(`^the response body should have type "(JSON)"$`, scenario.TheResponseBodyShouldHaveType)

	ctx.Step(`^the response body should be valid according to JSON schema "([^"]*)"$`, scenario.IValidateLastResponseBodyWithSchema)
	ctx.Step(`^the response body should be valid according to JSON schema:$`, scenario.IValidateLastResponseBodyWithFollowingSchema)

	ctx.Step(`^time between last request and response should be less than or equal to "([^"]*)"$`, scenario.TimeBetweenLastHTTPRequestResponseShouldBeLessThanOrEqualTo)

	/*
	   |--------------------------------------------------------------------------
	   | Preserving data
	   |--------------------------------------------------------------------------
	   |
	   | This section contains method for preserving data
	   |
	   | Argument following immediately after word "node"
	   | should have syntax acceptable by one of json-path libraries:
	   | https://github.com/pawelWritesCode/qjson
	   | https://github.com/oliveagle/jsonpath
	*/
	ctx.Step(`^I save "([^"]*)" as "([^"]*)"$`, scenario.ISaveAs)
	ctx.Step(`^I save from the last response JSON node "([^"]*)" as "([^"]*)"$`, scenario.ISaveFromTheLastResponseJSONNodeAs)

	/*
	   |--------------------------------------------------------------------------
	   | Debugging
	   |--------------------------------------------------------------------------
	   |
	   | This section contains methods that are useful during test creation
	*/
	ctx.Step(`^I print last response body$`, scenario.IPrintLastResponseBody)

	ctx.Step(`^I start debug mode$`, scenario.IStartDebugMode)
	ctx.Step(`^I stop debug mode$`, scenario.IStopDebugMode)

	/*
	   |--------------------------------------------------------------------------
	   | Flow control
	   |--------------------------------------------------------------------------
	   |
	   | This section contains methods for control scenario flow
	   |
	   | Argument in method 'I wait ([^"]*)"' should be string valid for
	   | golang standard library time.ParseDuration func, for example: 3s, 1h, 30ms
	*/
	ctx.Step(`^I wait "([^"]*)"`, scenario.IWait)
}

// checkErr checks error and log if found.
func checkErr(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
