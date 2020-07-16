package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/austinjan/AtopIOTServer/utils"
	"go.mongodb.org/mongo-driver/bson"
)

// APIVersion api version
var APIVersion = "1.0"

const (
	//CodeOK OK
	CodeOK = 0
	//CodeGeneralFailure normal error
	CodeGeneralFailure = 1
	//CodeGeneralFailure normal error
	CodeAuthenticateFailure = 3
	//CodeInvalidQAPIInput invalid parameter of query
	CodeInvalidQAPIInput = 4
	//CodeInvalidBodyParameter invalid parameter of body
	CodeInvalidBodyParameter = 5
	//CodeInvalidBodyParameter invalid parameter of body
	CodeServerInternalError = 6
	//CodeDatabaseError database error
	CodeDatabaseError = 7
	//ResultEmpty empty result
	CodeResultEmpty = 8
	// CodeDuplicateUniqueKey data have duplicate value which should be unique
	CodeDuplicateUniqueValue = 9
	// CodePermissionFailure
	CodePermissionFailure = 11
	// CodeConvertJSONFail is error status when system
	// marshal data to json fail
	CodeConvertJSONFail = 100
)

// ResponseStatus is a status string
var ResponseStatus = map[int]string{
	CodeOK:                   "OK",
	CodeGeneralFailure:       "Gernal failure",
	CodeAuthenticateFailure:  "Authentic failure",
	CodeInvalidQAPIInput:     "Invalid query parameter or URL",
	CodeInvalidBodyParameter: "Invalid body parameter",
	CodeServerInternalError:  "Server internal error",
	CodeDatabaseError:        "Database process error",
	CodeResultEmpty:          "Empty Result",
	CodeDuplicateUniqueValue: "Value should be unique",
	CodePermissionFailure:    "Permission denied",
	CodeConvertJSONFail:      "Convert reponse to JSON format fail",
}

// AtopResponse response object
type AtopResponse struct {
	Body   bson.M
	Writer http.ResponseWriter
}

func (r *AtopResponse) makeTimeStamp() {
	r.Body["responseTS"] = time.Now()
}

func (r *AtopResponse) makeError(code int) {
	r.Body["code"] = code
	r.Body["status"] = ResponseStatus[code]
}

// NewResponse basic response
func NewResponse(r *http.Request, w http.ResponseWriter) AtopResponse {
	var response = bson.M{}
	response["version"] = APIVersion
	response["code"] = 0
	response["status"] = "OK"
	response["command"] = r.URL.RequestURI()
	return AtopResponse{Body: response, Writer: w}
}

// AddPayload add payload into response
func (r *AtopResponse) AddPayload(payload bson.M) {
	r.Body["payload"] = payload
}

// SendResponse sending response
func (r *AtopResponse) SendResponse() {
	w := r.Writer
	// success
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	r.makeTimeStamp()
	j, err := json.Marshal(r.Body)

	if err != nil {
		r.SendErrorResponse(CodeConvertJSONFail)
		return
	}
	w.Write(j)

	return
}

// SendAuthenticateError send error response (Body error)
func (r *AtopResponse) SendAuthenticateError(message string) {
	w := r.Writer
	r.makeError(CodeAuthenticateFailure)
	r.makeTimeStamp()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	j, _ := json.Marshal(r.Body)
	w.Write(j)
	return
}

// SendGeneralError send error response (general error)
func (r *AtopResponse) SendPermissionError(info utils.RequestPermission) {
	msg := fmt.Sprintf("Permission denied: %s %s %s ", info.Name, info.Method, info.Collection)
	r.SendErrorResponseWithMessage(CodePermissionFailure, msg)
}

// SendGeneralError send error response (general error)
//
func (r *AtopResponse) SendGeneralError(err error) {
	r.SendErrorResponseWithMessage(CodeGeneralFailure, err.Error())
}

// SendBodyError send error response (Body error)
func (r *AtopResponse) SendBodyError(err error) {
	r.SendErrorResponseWithMessage(CodeInvalidBodyParameter, err.Error())
}

// SendAPIInputError send error response (Body error)
func (r *AtopResponse) SendAPIInputError(err error) {
	r.SendErrorResponseWithMessage(CodeInvalidQAPIInput, err.Error())
}

// SendServerError send error response (server error)
func (r *AtopResponse) SendServerError(err error) {
	r.SendErrorResponseWithMessage(CodeServerInternalError, err.Error())
}

// SendDatabaseError send error response (database error)
func (r *AtopResponse) SendDatabaseError(err error) {
	r.SendErrorResponseWithMessage(CodeDatabaseError, err.Error())
}

func (r *AtopResponse) SendUniqueValueError(err error) {
	r.SendErrorResponseWithMessage(CodeDuplicateUniqueValue, err.Error())
}

// SendEmptyError send error response (empty result )
func (r *AtopResponse) SendEmptyError() {
	r.SendErrorResponseWithMessage(CodeResultEmpty, "empty result")
}

// SendErrorResponseWithMessage response with error code and reason
func (r *AtopResponse) SendErrorResponseWithMessage(code int, message string) {
	r.AddPayload(bson.M{"reason": message})
	r.SendErrorResponse(code)
}

// SendErrorResponse response with error code
func (r *AtopResponse) SendErrorResponse(code int) {
	w := r.Writer
	r.makeError(code)
	r.makeTimeStamp()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	j, _ := json.Marshal(r.Body)
	w.Write(j)
	return
}
