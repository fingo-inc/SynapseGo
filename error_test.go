package synapse

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var errorData map[string]interface{}

/********** METHODS **********/

func init() {
	data, err := readFile("error_responses")

	if err != nil {
		panic(err)
	}

	errorData = data
}

/********** TESTS **********/

func Test_HandleHTTPError(t *testing.T) {
	assert := assert.New(t)

	for k := range errorData {

		testErrRes, _ := json.Marshal(errorData[k])
		testErr := handleHTTPError(testErrRes)

		errData := errorData[k].(map[string]interface{})
		httpCode := errData["http_code"].(string)
		errCode := errData["error_code"].(string)
		msg := errData["error"].(map[string]interface{})["en"].(string)
		responseMsg := "http_code " + httpCode + " error_code " + errCode + " " + msg

		// error message should be an error and print error code plus original API message
		assert.EqualError(testErr, responseMsg)
	}
}

func Test_HandleInvalidHTTPError(t *testing.T) {
	assert := assert.New(t)

	tests := []struct{
		res map[string]interface{}
		msg string
	}{
		// empty response
		{
			res: map[string]interface{}{},
			msg: "http_code 999 error_code 999 unknown error",
		},
		// missing "http_code" attribute
		{
			res: map[string]interface{}{"error_code": "10", "error" : map[string]interface{}{"en": "error message"}},
			msg: "http_code 999 error_code 10 error message",
		},
		// missing "error_code" attribute
		{
			res: map[string]interface{}{"http_code": "202", "error" : map[string]interface{}{"en": "error message"}},
			msg: "http_code 202 error_code 999 error message",
		},
		// missing "error" attribute
		{
			res: map[string]interface{}{"http_code": "202", "error_code": "10"},
			msg: "http_code 202 error_code 10 unknown error",
		},
	}
	for _, test := range tests {
		testErrRes, _ := json.Marshal(test.res)
		testErr := handleHTTPError(testErrRes)
		assert.EqualError(testErr, test.msg)
	}
}
