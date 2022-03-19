package defs

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/cucumber/godog"
	"github.com/pawelWritesCode/gdutils"
	"github.com/pawelWritesCode/gdutils/pkg/format"
	"github.com/pawelWritesCode/gdutils/pkg/stringutils"
	"github.com/pawelWritesCode/gdutils/pkg/timeutils"
)

// Scenario is entity that contains utility services and holds methods used behind godog steps.
type Scenario struct {
	// APIContext holds utility services and methods for working with HTTP(s) API.
	APIContext *gdutils.APIContext
}

// IGenerateARandomRunesOfLengthWithCharactersAndSaveItAs creates random runes generator func using provided charset.
// Returned func creates runes from provided range and preserve it under given cacheKey in scenario cache.
func (s *Scenario) IGenerateARandomRunesOfLengthWithCharactersAndSaveItAs(from, to int, charset string, cacheKey string) error {
	var generateWordFunc func(from, to int, cacheKey string) error

	switch strings.ToLower(charset) {
	case "ascii":
		generateWordFunc = s.APIContext.IGenerateARandomRunesInTheRangeToAndSaveItAs(stringutils.CharsetASCII)
	case "unicode":
		generateWordFunc = s.APIContext.IGenerateARandomRunesInTheRangeToAndSaveItAs(stringutils.CharsetUnicode)
	case "polish":
		generateWordFunc = s.APIContext.IGenerateARandomRunesInTheRangeToAndSaveItAs(stringutils.CharsetPolish)
	case "english":
		generateWordFunc = s.APIContext.IGenerateARandomRunesInTheRangeToAndSaveItAs(stringutils.CharsetEnglish)
	case "russian":
		generateWordFunc = s.APIContext.IGenerateARandomRunesInTheRangeToAndSaveItAs(stringutils.CharsetRussian)
	default:
		return fmt.Errorf("unknown charset '%s', available: ascii, unicode, polish, english, russian", charset)
	}

	return generateWordFunc(from, to, cacheKey)
}

// IGenerateARandomNumberInTheRangeFromToAndSaveItAs generates random number from provided range
// and preserve it in scenario cache under provided cacheKey.
func (s *Scenario) IGenerateARandomNumberInTheRangeFromToAndSaveItAs(numberType string, from, to int, cacheKey string) error {
	switch strings.ToLower(numberType) {
	case "float":
		return s.APIContext.IGenerateARandomFloatInTheRangeToAndSaveItAs(from, to, cacheKey)
	case "int":
		return s.APIContext.IGenerateARandomIntInTheRangeToAndSaveItAs(from, to, cacheKey)
	default:
		return fmt.Errorf("unknown type %s, available: int, float", numberType)
	}
}

// IGenerateARandomSentenceInTheRangeFromToWordsAndSaveItAs creates generator func for creating random sentences.
// Each sentence has length from - to as provided in params and is saved in scenario cache under provided cacheKey.
func (s *Scenario) IGenerateARandomSentenceInTheRangeFromToWordsAndSaveItAs(minWordLength, maxWordLength int) func(from, to int, charset string, cacheKey string) error {
	return func(from, to int, charset string, cacheKey string) error {
		var generateSentenceFunc func(from, to int, cacheKey string) error
		switch strings.ToLower(charset) {
		case "ascii":
			generateSentenceFunc = s.APIContext.IGenerateARandomSentenceInTheRangeFromToWordsAndSaveItAs(stringutils.CharsetASCII, minWordLength, maxWordLength)
		case "unicode":
			generateSentenceFunc = s.APIContext.IGenerateARandomSentenceInTheRangeFromToWordsAndSaveItAs(stringutils.CharsetUnicode, minWordLength, maxWordLength)
		case "polish":
			generateSentenceFunc = s.APIContext.IGenerateARandomSentenceInTheRangeFromToWordsAndSaveItAs(stringutils.CharsetPolish, minWordLength, maxWordLength)
		case "english":
			generateSentenceFunc = s.APIContext.IGenerateARandomSentenceInTheRangeFromToWordsAndSaveItAs(stringutils.CharsetEnglish, minWordLength, maxWordLength)
		case "russian":
			generateSentenceFunc = s.APIContext.IGenerateARandomSentenceInTheRangeFromToWordsAndSaveItAs(stringutils.CharsetRussian, minWordLength, maxWordLength)
		default:
			return fmt.Errorf("unknown charset '%s', available: ascii, unicode, polish, english, russian", charset)
		}

		return generateSentenceFunc(from, to, cacheKey)
	}
}

//
func (s *Scenario) IGenerateRandomBoolValueAndSaveItAs(cacheKey string) error {
	s.APIContext.Cache.Save(cacheKey, rand.Intn(2) == 0)

	return nil
}

// IGenerateCurrentTimeAndTravelByAndSaveItAs creates current time object, move timeDuration in time and
// save it in cache under given cacheKey.
func (s *Scenario) IGenerateCurrentTimeAndTravelByAndSaveItAs(timeDirection, timeDuration, cacheKey string) error {
	duration, err := time.ParseDuration(timeDuration)
	if err != nil {
		return err
	}

	return s.APIContext.IGenerateCurrentTimeAndTravelByAndSaveItAs(timeutils.TimeDirection(timeDirection), duration, cacheKey)
}

/*
	ISendRequestToWithFormatBodyAndHeaders sends HTTP(s) requests with provided body and headers.

	Argument "method" indices HTTP request method for example: "POST", "GET" etc.
 	Argument "urlTemplate" should be full valid URL. May include template values.
	Argument "bodyTemplate" should contain data (may include template values)
	in JSON or YAML format with keys "body" and "headers".
*/
func (s *Scenario) ISendRequestToWithBodyAndHeaders(method, urlTemplate string, reqBody *godog.DocString) error {
	return s.APIContext.ISendRequestToWithBodyAndHeaders(method, urlTemplate, reqBody.Content)
}

// IPrepareNewRequestToAndSaveItAs prepares new request and saves it in cache under cacheKey.
func (s Scenario) IPrepareNewRequestToAndSaveItAs(method, urlTemplate, cacheKey string) error {
	return s.APIContext.IPrepareNewRequestToAndSaveItAs(method, urlTemplate, cacheKey)
}

// ISetFollowingHeadersForPreparedRequest sets provided headers for previously prepared request.
// incoming data should be in format acceptable by injected s.APIContext.Deserializer
func (s Scenario) ISetFollowingHeadersForPreparedRequest(cacheKey string, headersTemplate *godog.DocString) error {
	return s.APIContext.ISetFollowingHeadersForPreparedRequest(cacheKey, headersTemplate.Content)
}

// ISetFollowingCookiesForPreparedRequest sets cookies for previously prepared request
// cookies template should be YAML or JSON deserializable on []http.Cookie
func (s Scenario) ISetFollowingCookiesForPreparedRequest(cacheKey string, cookies *godog.DocString) error {
	return s.APIContext.ISetFollowingCookiesForPreparedRequest(cacheKey, cookies.Content)
}

/*
	ISetFollowingFormForPreparedRequest sets form for previously prepared request.
	Internally method sets proper Content-Type: multipart/form-data header.
	formTemplate should be YAML or JSON deserializable on map[string]string.
*/
func (s Scenario) ISetFollowingFormForPreparedRequest(cacheKey string, formTemplate *godog.DocString) error {
	return s.APIContext.ISetFollowingFormForPreparedRequest(cacheKey, formTemplate.Content)
}

// ISetFollowingBodyForPreparedRequest sets body for previously prepared request.
// bodyTemplate may be in any format and accepts template values.
func (s Scenario) ISetFollowingBodyForPreparedRequest(cacheKey string, bodyTemplate *godog.DocString) error {
	return s.APIContext.ISetFollowingBodyForPreparedRequest(cacheKey, bodyTemplate.Content)
}

// ISendRequest sends previously prepared HTTP(s) request.
func (s Scenario) ISendRequest(cacheKey string) error {
	return s.APIContext.ISendRequest(cacheKey)
}

// TheResponseShouldHaveHeader checks whether last HTTP response has given header.
func (s *Scenario) TheResponseShouldHaveHeader(name string) error {
	return s.APIContext.TheResponseShouldHaveHeader(name)
}

// TheResponseShouldHaveHeaderOfValue checks whether last HTTP response has given header with provided value.
func (s *Scenario) TheResponseShouldHaveHeaderOfValue(name, value string) error {
	return s.APIContext.TheResponseShouldHaveHeaderOfValue(name, value)
}

// TheResponseStatusCodeShouldBe checks last response status code.
func (s *Scenario) TheResponseStatusCodeShouldBe(code int) error {
	return s.APIContext.TheResponseStatusCodeShouldBe(code)
}

// TheResponseShouldHaveNode checks whether last response body contains given node.
// expr should be valid according to injected PathFinder for given data format
func (s *Scenario) TheResponseShouldHaveNode(dataFormat, exprTemplate string) error {
	return s.APIContext.TheResponseShouldHaveNode(format.DataFormat(dataFormat), exprTemplate)
}

// TheNodeShouldBeOfValue compares json node value from expression to expected by user dataValue of given by user dataType
// Available data types are listed in switch section in each case directive.
// expr should be valid according to injected PathFinder for provided dataFormat.
func (s *Scenario) TheNodeShouldBeOfValue(dataFormat, exprTemplate, dataType, dataValue string) error {
	return s.APIContext.TheNodeShouldBeOfValue(format.DataFormat(dataFormat), exprTemplate, dataType, dataValue)
}

// TheNodeShouldBeSliceOfLength checks whether given key is slice and has given length
// expr should be valid according to injected PathFinder for provided dataFormat
func (s *Scenario) TheNodeShouldBeSliceOfLength(dataFormat, exprTemplate string, length int) error {
	return s.APIContext.TheNodeShouldBeSliceOfLength(format.DataFormat(dataFormat), exprTemplate, length)
}

// TheNodeShouldBe checks whether node from last response body is of provided type
// goType may be one of: nil, string, int, float, bool, map, slice
// expr should be valid according to injected PathResolver
func (s *Scenario) TheNodeShouldBe(dataFormat, exprTemplate, goType string) error {
	return s.APIContext.TheNodeShouldBe(format.DataFormat(dataFormat), exprTemplate, goType)
}

// TheNodeShouldNotBe checks whether node from last response body is not of provided type.
// goType may be one of: nil, string, int, float, bool, map, slice,
// expr should be valid according to injected PathFinder for given data format.
func (s *Scenario) TheNodeShouldNotBe(dataFormat, exprTemplate, goType string) error {
	return s.APIContext.TheNodeShouldNotBe(format.DataFormat(dataFormat), exprTemplate, goType)
}

// TheResponseShouldHaveNodes checks whether last request body has keys defined in string separated by comma
// nodeExprs should be valid according to injected PathFinder expressions separated by comma (,)
func (s *Scenario) TheResponseShouldHaveNodes(dataFormat, nodesExpr string) error {
	return s.APIContext.TheResponseShouldHaveNodes(format.DataFormat(dataFormat), nodesExpr)
}

// TheNodeShouldMatchRegExp checks whether last response body node matches provided regExp.
func (s *Scenario) TheNodeShouldMatchRegExp(dataFormat, exprTemplate, regExpTemplate string) error {
	return s.APIContext.TheNodeShouldMatchRegExp(format.DataFormat(dataFormat), exprTemplate, regExpTemplate)
}

// TheResponseBodyShouldHaveFormat checks whether last response body has given data format.
// Available data formats are listed in format package.
func (s *Scenario) TheResponseBodyShouldHaveFormat(dataFormat string) error {
	return s.APIContext.TheResponseBodyShouldHaveFormat(format.DataFormat(dataFormat))
}

/*
	IValidateLastResponseBodyWithSchema validates last response body against JSON schema under provided reference.
	reference may be:
		- full OS path to JSON schema
		- relative path from JSON schema's dir which was passed in main_test to initialize *Scenario struct instance,
		- URL
*/
func (s *Scenario) IValidateLastResponseBodyWithSchema(referenceTemplate string) error {
	return s.APIContext.IValidateLastResponseBodyWithSchemaReference(referenceTemplate)
}

// IValidateLastResponseBodyWithFollowingSchema validates last response body against JSON schema provided by user.
func (s *Scenario) IValidateLastResponseBodyWithFollowingSchema(schemaBytes *godog.DocString) error {
	return s.APIContext.IValidateLastResponseBodyWithSchemaString(schemaBytes.Content)
}

/*
	TimeBetweenLastHTTPRequestResponseShouldBeLessThanOrEqualTo asserts that last HTTP request-response time
	is <= than expected timeInterval.
	timeInterval should be string acceptable by time.ParseDuration func
*/
func (s *Scenario) TimeBetweenLastHTTPRequestResponseShouldBeLessThanOrEqualTo(timeInterval string) error {
	duration, err := time.ParseDuration(timeInterval)
	if err != nil {
		return err
	}

	return s.APIContext.TimeBetweenLastHTTPRequestResponseShouldBeLessThanOrEqualTo(duration)
}

// TheResponseShouldHaveCookie checks whether last HTTP(s) response has cookie of given name.
func (s *Scenario) TheResponseShouldHaveCookie(name string) error {
	return s.APIContext.TheResponseShouldHaveCookie(name)
}

// TheResponseShouldHaveCookieOfValue checks whether last HTTP(s) response has cookie of given name and value.
func (s *Scenario) TheResponseShouldHaveCookieOfValue(name, valueTemplate string) error {
	return s.APIContext.TheResponseShouldHaveCookieOfValue(name, valueTemplate)
}

// IValidateNodeWithSchemaReference validates last response body node against schema as provided in reference
func (s *Scenario) IValidateNodeWithSchemaReference(dataFormat, exprTemplate, referenceTemplate string) error {
	return s.APIContext.IValidateNodeWithSchemaReference(format.DataFormat(dataFormat), exprTemplate, referenceTemplate)
}

// IValidateNodeWithSchemaString validates last response body JSON node against schema
func (s *Scenario) IValidateNodeWithSchemaString(dataFormat, exprTemplate string, schemaTemplate *godog.DocString) error {
	return s.APIContext.IValidateNodeWithSchemaString(format.DataFormat(dataFormat), exprTemplate, schemaTemplate.Content)
}

// ISaveAs saves into cache arbitrary passed value
func (s *Scenario) ISaveAs(valueTemplate, cacheKey string) error {
	return s.APIContext.ISaveAs(valueTemplate, cacheKey)
}

// ISaveFromTheLastResponseNodeAs saves from last response json node under given cache key.
func (s *Scenario) ISaveFromTheLastResponseNodeAs(dataFormat, exprTemplate, cacheKey string) error {
	return s.APIContext.ISaveFromTheLastResponseNodeAs(format.DataFormat(dataFormat), exprTemplate, cacheKey)
}

// IPrintLastResponseBody prints response body from last scenario request
func (s *Scenario) IPrintLastResponseBody() error {
	return s.APIContext.IPrintLastResponseBody()
}

/*
	IWait waits for provided time interval amount of time
	timeInterval should be string valid for time.ParseDuration func,
	for example: 3s, 1h, 30ms
*/
func (s *Scenario) IWait(timeInterval string) error {
	duration, err := time.ParseDuration(timeInterval)
	if err != nil {
		return err
	}

	return s.APIContext.IWait(duration)
}

// IStartDebugMode starts debugging mode
func (s *Scenario) IStartDebugMode() error {
	return s.APIContext.IStartDebugMode()
}

// IStopDebugMode stops debugging mode
func (s *Scenario) IStopDebugMode() error {
	return s.APIContext.IStopDebugMode()
}

// IStopScenarioExecution stops scenario execution
func (s *Scenario) IStopScenarioExecution() error {
	return errors.New("scenario stopped")
}
