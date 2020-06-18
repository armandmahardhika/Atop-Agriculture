package router

import (
	"encoding/json"
	"net/http"
	"time"

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
	//CodeInvalidQueryParameter invalid parameter of query
	CodeInvalidQueryParameter = 4
	//CodeInvalidBodyParameter invalid parameter of body
	CodeInvalidBodyParameter = 5
	//CodeInvalidBodyParameter invalid parameter of body
	CodeServerInternalError = 6
	//ResultEmpty empty result
	CodeResultEmpty = 8
	// CodeConvertJSONFail is error status when system
	// marshal data to json fail
	CodeConvertJSONFail = 100
)

// ResponseStatus is a status string
var ResponseStatus = map[int]string{
	CodeOK:                    "OK",
	CodeGeneralFailure:        "Gernal failure",
	CodeAuthenticateFailure:   "Authentic failure",
	CodeInvalidQueryParameter: "Invalid query parameter",
	CodeInvalidBodyParameter:  "Invalid body parameter",
	CodeServerInternalError:   "Server internal error",
	CodeResultEmpty:           "Empty Result",
	CodeConvertJSONFail:       "Convert reponse to JSON format fail",
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

// SendServerError send error response (server error)
func (r *AtopResponse) SendServerError(err error) {
	r.SendErrorResponseWithMessage(CodeServerInternalError, err.Error())
}

// SendAuthenticateError send error response (Body error)
func (r *AtopResponse) SendAuthenticateError(message string) {
	r.SendErrorResponseWithMessage(CodeAuthenticateFailure, message)
}

// SendGeneralError send error response (general error)
func (r *AtopResponse) SendGeneralError(err error) {
	r.SendErrorResponseWithMessage(CodeGeneralFailure, err.Error())
}

// SendBodyError send error response (Body error)
func (r *AtopResponse) SendBodyError(err error) {
	r.SendErrorResponseWithMessage(CodeInvalidBodyParameter, err.Error())
}

// SendQueryError send error response (Body error)
func (r *AtopResponse) SendQueryError(err error) {
	r.SendErrorResponseWithMessage(CodeInvalidQueryParameter, err.Error())
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
