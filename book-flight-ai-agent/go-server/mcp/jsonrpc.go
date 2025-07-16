/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package mcp

const (
	jsonrpc = "2.0" // jsonrpc version
)

type RequestRPC struct {
	// JsonRPC specifies the JSON-RPC version. It MUST be exactly "2.0".
	JsonRPC string `json:"jsonrpc"` // Explicitly named "jsonrpc" in JSON
	// Method is a string containing the name of the method to be invoked.
	// Method names that begin with a lowercase letter are reserved for
	// system-defined methods and MUST NOT be used for custom methods.
	Method string `json:"method"`
	// Params is a structured value that holds the parameter values to be used
	// during the invocation of the method. This member MAY be omitted.
	Params map[string]any `json:"params,omitempty"` // Omitempty to skip if nil
	// ID is an identifier established by the client that MUST contain a String,
	// Number, or Null value if included. If not included it is assumed to be a
	// notification. The value SHOULD normally not be Null [1] and Numbers SHOULD
	// NOT contain fractional parts [2].
	// [1] Using Null as a value for the id member in a Request object is
	// discouraged, as there are no benefits and it introduces ambiguities
	// when differentiating between Requests, Notifications, and Response
	// objects.
	// [2] Fractional parts SHOULD NOT be used as there is no clear
	// interoperable way to represent them across all systems.
	ID string `json:"id,omitempty"` // Omitempty to skip if empty
}

// NewRequestRPC creates a new RequestRPC with the JsonRPC field set to "2.0".
func NewRequestRPC(method string, params map[string]any, id string) *RequestRPC {
	return &RequestRPC{
		JsonRPC: jsonrpc,
		Method:  method,
		Params:  params,
		ID:      id,
	}
}

type ErrorRPC struct {
	// A Number that indicates the error type that occurred.
	// This MUST be an integer.
	Code int64 `json:"code"`
	// A String providing a short description of the error.
	// The message SHOULD be limited to a concise single sentence.
	Message string `json:"message"`
	// A Primitive or Structured value that contains additional information about the error.
	// This may be omitted.
	// The value of this member is defined by the Server(e.g. detailed error information, nested errors etc.).
	Data map[string]any `json:"data,omitempty"` // Omitempty to skip if nil
}

// NewErrorRPC creates a new ErrorRPC instance.
func NewErrorRPC(code int64, message string, data map[string]any) *ErrorRPC {
	return &ErrorRPC{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

type ResponseRPC struct {
	// JsonRPC specifies the JSON-RPC version. It MUST be exactly "2.0".
	JsonRPC string `json:"jsonrpc"` // Explicitly named "jsonrpc" in JSON
	// Result is the value that holds the result of the method invocation.
	// This member MUST be present when there was no error invoking the
	// method. The value of this member is determined by the method invoked
	// on the Server.
	Result map[string]any `json:"result,omitempty"` // Omitempty to skip if nil
	// Error is an object containing information about the error that occurred
	// during the invocation of the method. This member MUST be present when
	// there was an error invoking the method.
	Error *ErrorRPC `json:"error,omitempty"` // Omitempty to skip if nil
	// ID is the identifier established by the Client that MUST contain a String,
	// Number, or Null value if included in the request. It MUST be the same as
	// the value of the id member in the Request object. If there was an error
	// in detecting the id in the Request object (e.g. Parse error or Invalid
	// Request), it MUST be Null.
	ID string `json:"id,omitempty"` // Omitempty to skip if empty
}

// NewResponseRPC creates a new ResponseRPC with the JsonRPC field set to "2.0".
func NewResponseRPC(result map[string]any, error *ErrorRPC, id string) *ResponseRPC {
	return &ResponseRPC{
		JsonRPC: jsonrpc,
		Result:  result,
		Error:   error,
		ID:      id,
	}
}
