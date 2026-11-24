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

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// mesh is the discovery id configured in conf/dubbo-admin.yaml.
const mesh = "nacos"

// adminClient calls the Dubbo Admin console API under /api/v1.
type adminClient struct {
	base string
	http *http.Client
}

func newAdminClient(address string) (*adminClient, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	return &adminClient{
		base: strings.TrimRight(address, "/") + "/api/v1",
		http: &http.Client{Jar: jar, Timeout: 10 * time.Second},
	}, nil
}

func (a *adminClient) login(user, password string) error {
	resp, err := a.http.PostForm(a.base+"/auth/login", url.Values{"user": {user}, "password": {password}})
	if err != nil {
		return err
	}
	return decodeEnvelope(resp, nil)
}

func (a *adminClient) get(path string, params url.Values, out any) error {
	return a.do(http.MethodGet, path, params, nil, out)
}

func (a *adminClient) post(path string, body, out any) error {
	return a.do(http.MethodPost, path, nil, body, out)
}

func (a *adminClient) put(path string, params url.Values) error {
	return a.do(http.MethodPut, path, params, nil, nil)
}

func (a *adminClient) delete(path string) error {
	return a.do(http.MethodDelete, path, nil, nil, nil)
}

// promQL runs an instant query through Admin, which proxies Prometheus, and
// returns the value of the single resulting sample.
func (a *adminClient) promQL(promQL string) (float64, error) {
	resp, err := a.http.Get(a.base + "/promQL/query?" + url.Values{"query": {promQL}}.Encode())
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	var result struct {
		Status string `json:"status"`
		Data   struct {
			Result []struct {
				Value [2]any `json:"value"`
			} `json:"result"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}
	if result.Status != "success" || len(result.Data.Result) != 1 {
		return 0, fmt.Errorf("query %s: status %q with %d samples", promQL, result.Status, len(result.Data.Result))
	}
	raw, _ := result.Data.Result[0].Value[1].(string)
	return strconv.ParseFloat(raw, 64)
}

func (a *adminClient) do(method, path string, params url.Values, body, out any) error {
	if params == nil {
		params = url.Values{}
	}
	params.Set("mesh", mesh)
	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(payload)
	}
	req, err := http.NewRequest(method, a.base+path+"?"+params.Encode(), reader)
	if err != nil {
		return err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := a.http.Do(req)
	if err != nil {
		return err
	}
	return decodeEnvelope(resp, out)
}

// decodeEnvelope unwraps Admin's {code, message, data} response.
func decodeEnvelope(resp *http.Response, out any) error {
	defer resp.Body.Close()
	var envelope struct {
		Code    string          `json:"code"`
		Message string          `json:"message"`
		Data    json.RawMessage `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return fmt.Errorf("%s %s: HTTP %d: %w", resp.Request.Method, resp.Request.URL.Path, resp.StatusCode, err)
	}
	if envelope.Code != "Success" {
		return &apiError{
			request: resp.Request.Method + " " + resp.Request.URL.Path,
			code:    envelope.Code,
			message: envelope.Message,
		}
	}
	if out == nil {
		return nil
	}
	return json.Unmarshal(envelope.Data, out)
}

// apiError is a non-success code in Admin's response envelope.
type apiError struct {
	request, code, message string
}

func (e *apiError) Error() string {
	return fmt.Sprintf("%s: %s: %s", e.request, e.code, e.message)
}

func isNotFound(err error) bool {
	var apiErr *apiError
	return errors.As(err, &apiErr) && apiErr.code == "NotFoundError"
}

func query(key, value string) url.Values {
	return url.Values{key: {value}}
}

func serviceQuery() url.Values {
	return url.Values{"serviceName": {serviceName}, "version": {serviceVer}, "group": {serviceGroup}}
}
