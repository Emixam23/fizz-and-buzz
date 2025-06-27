package integration

import (
	"bou.ke/monkey"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cucumber/godog"
	"gitlab.com/emixam23/fizz-and-buzz/internal"
	"net/http"
	"net/http/httptest"
	"reflect"
	"time"

	"github.com/rdumont/assistdog"
)

var assist = assistdog.NewDefault()

/*****************************************/
/*************** GIVEN *******************/
/*****************************************/

func nowIs(timeStr string) error {
	t, err := time.Parse("2006-01-02 15:04:05+02", timeStr)

	if err != nil {
		return fmt.Errorf("could not parse date: %w", err)
	}

	testContext.mockedTime = monkey.Patch(time.Now, func() time.Time { return t })

	return nil
}

/*****************************************/
/****************** WHEN *****************/
/*****************************************/

func fizzAndBuzzServiceHasAZeroOf(zero int) error {

	if zero < 0 {
		return errors.New("negative number as zero are not supported (uint)")
	}

	testContext.conf.FnbServiceConfig.Zero = uint64(zero)

	app, err := internal.NewWithArgs(testContext.conf, testContext.dbClient)
	if err != nil {
		return fmt.Errorf("an error happened while creating/initializing new app with args: %w", err)
	}
	testContext.app = app

	return nil
}

func iSendARequestTo(method, url string) error {
	req, _ := http.NewRequest(method, url, nil)

	resp := httptest.NewRecorder()
	testContext.responsesRecorder = append(testContext.responsesRecorder, resp)
	testContext.app.GetRouter().ServeHTTP(resp, req)

	return nil
}

func iSendARequestToWithTheGivenParams(method, url string, table *godog.Table) error {
	rows, _ := assist.ParseSlice(table)
	for _, row := range rows {
		var urlWithQuery string
		switch url {
		case "/history":
			urlWithQuery = fmt.Sprintf("%s?limit=%s", url, row["limit"])
			break
		case "/stats":
			urlWithQuery = fmt.Sprintf("%s?sorted=%s", url, row["sorted"])
			break
		default:
			return godog.ErrUndefined
		}
		if err := iSendARequestTo(method, urlWithQuery); err != nil {
			return err
		}
	}

	return nil
}

func iSendRequestsToWithTheGivenParams(method, url string, table *godog.Table) error {
	rows, _ := assist.ParseSlice(table)
	for _, row := range rows {
		var urlWithQuery string
		switch url {
		case "/fizz-and-buzz":
			urlWithQuery = fmt.Sprintf("%s?n1=%s&s1=%s&n2=%s&s2=%s&limit=%s", url, row["n1"], row["s1"], row["n2"], row["s2"], row["limit"])
			break
		default:
			return godog.ErrUndefined
		}
		if err := iSendARequestTo(method, urlWithQuery); err != nil {
			return err
		}
	}

	return nil
}

/*****************************************/
/**************** THEN *******************/
/*****************************************/

func theResponseCodeShouldBe(statusCode int) error {

	if len(testContext.responsesRecorder) == 0 {
		return errors.New("no response caught or available")
	}

	resp := testContext.responsesRecorder[len(testContext.responsesRecorder)-1]
	if statusCode != resp.Code {
		return fmt.Errorf("expected http response code to be: %d, but actual is: %d", statusCode, resp.Code)
	}

	return nil
}

func theResponseCodesShouldBe(statusCode int) error {

	if len(testContext.responsesRecorder) == 0 {
		return errors.New("no response caught or available")
	}

	for _, resp := range testContext.responsesRecorder {
		if statusCode != resp.Code {
			return fmt.Errorf("expected http response code to be: %d, but actual is: %d", statusCode, resp.Code)
		}
	}

	return nil
}

func theResponseShouldMatchJSON(body *godog.DocString) (err error) {

	if len(testContext.responsesRecorder) == 0 {
		return errors.New("no response caught or available")
	}
	resp := testContext.responsesRecorder[len(testContext.responsesRecorder)-1]

	var expected, actual interface{}

	if err = json.Unmarshal([]byte(body.Content), &expected); err != nil {
		fmt.Println("1")
		return
	}

	if err = json.Unmarshal(resp.Body.Bytes(), &actual); err != nil {
		fmt.Println("2")
		return
	}

	if !reflect.DeepEqual(expected, actual) {
		return fmt.Errorf("expected JSON does not match actual, expected: %v, actual: %v", expected, actual)
	}
	return nil
}
