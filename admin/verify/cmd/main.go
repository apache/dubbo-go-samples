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

// verify walks through the Admin integration end to end: it reads what Admin
// shows about the sample, changes traffic through Admin's API, and observes
// the effect on real calls. Run it inside the compose network:
//
//	docker compose run --rm verify
package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/client"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/config_center"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/registry"
)

import (
	adminsample "github.com/apache/dubbo-go-samples/admin/proto"
)

const (
	providerApp  = "admin-sample-provider"
	consumerApp  = "admin-sample-consumer"
	serviceName  = "adminsample.GreetService"
	serviceVer   = "1.0.0"
	serviceGroup = "admin-sample"
	providerPort = 20000
	grayTag      = "gray"
	ruleGroup    = "dubbo"

	conditionRuleName = serviceName + ":" + serviceVer + ":" + serviceGroup + ".condition-router"
	tagRuleName       = providerApp + ".tag-router"

	// Admin credentials from conf/dubbo-admin.yaml.
	adminUser     = "admin"
	adminPassword = "dubbo@2025"

	// Rules travel Admin -> Nacos -> SDK listener, and Admin refreshes its view
	// on its discovery watch period.
	ruleTimeout  = 60 * time.Second
	viewTimeout  = 90 * time.Second
	pollInterval = 2 * time.Second
	listenerSync = 6 * time.Second
	callsPerPoll = 10
)

type verifier struct {
	admin    *adminClient
	jaeger   string
	greet    adminsample.GreetService
	ipA, ipB string

	// Learned from Admin by the earlier checks, used to invoke from Admin.
	instanceA, greetSignature string
}

func main() {
	v, err := newVerifier()
	if err != nil {
		fmt.Println("setup failed:", err)
		os.Exit(1)
	}
	checks := []struct {
		name string
		run  func() error
	}{
		{"Admin lists the provider application with two instances", v.applications},
		{"Admin lists both provider instances with their IPs", v.instances},
		{"Admin shows the service with its version and group", v.service},
		{"Admin shows the Greet method from published metadata", v.method},
		{"Admin invokes Greet on provider-a through its debug API", v.invoke},
		{"Without rules, calls reach both providers", v.baseline},
		{"Condition route from Admin pins calls to provider-a", v.conditionRoute},
		{"Tag route from Admin sends gray calls to provider-b", v.tagRoute},
		{"Traffic disable from Admin drains provider-b", v.trafficDisable},
		{"Admin queries consumer and provider request metrics", v.metrics},
		{"Traces reach Jaeger and Admin links the trace dashboard", v.traces},
	}
	failed := 0
	for _, c := range checks {
		if err := c.run(); err != nil {
			failed++
			fmt.Printf("FAIL  %s\n      %v\n", c.name, err)
			continue
		}
		fmt.Printf("PASS  %s\n", c.name)
	}
	v.cleanup()
	if failed > 0 {
		fmt.Printf("%d of %d checks failed\n", failed, len(checks))
		os.Exit(1)
	}
	fmt.Printf("all %d checks passed\n", len(checks))
}

func newVerifier() (*verifier, error) {
	admin, err := newAdminClient(getEnv("ADMIN_ADDRESS", "http://dubbo-admin:8888"))
	if err != nil {
		return nil, err
	}
	if err = admin.login(adminUser, adminPassword); err != nil {
		return nil, fmt.Errorf("login to Admin: %w", err)
	}
	ipA, err := lookup("provider-a")
	if err != nil {
		return nil, err
	}
	ipB, err := lookup("provider-b")
	if err != nil {
		return nil, err
	}

	nacos := getEnv("NACOS_ADDRESS", "nacos:8848")
	ins, err := dubbo.NewInstance(
		dubbo.WithName("admin-sample-verify"),
		dubbo.WithRegistry(registry.WithNacos(), registry.WithAddress(nacos)),
		dubbo.WithConfigCenter(config_center.WithNacos(), config_center.WithAddress(nacos), config_center.WithGroup(ruleGroup)),
	)
	if err != nil {
		return nil, err
	}
	// Round robin makes a batch served by one instance mean a rule applied,
	// never a random streak.
	cli, err := ins.NewClient(client.WithClientVersion(serviceVer), client.WithClientGroup(serviceGroup),
		client.WithClientLoadBalanceRoundRobin())
	if err != nil {
		return nil, err
	}
	greet, err := adminsample.NewGreetService(cli)
	if err != nil {
		return nil, err
	}
	v := &verifier{admin: admin, jaeger: getEnv("JAEGER_ADDRESS", "http://jaeger:16686"), greet: greet, ipA: ipA, ipB: ipB}
	// Leftovers from an interrupted run would skew the routing checks.
	v.cleanup()
	return v, nil
}

// applications checks the provider application. Under application-level
// discovery consumers register no instances, so Admin lists providers only.
func (v *verifier) applications() error {
	return eventually(viewTimeout, func() error {
		var page struct {
			List []struct {
				AppName       string `json:"appName"`
				InstanceCount int64  `json:"instanceCount"`
			} `json:"list"`
		}
		if err := v.admin.get("/application/search", query("keywords", providerApp), &page); err != nil {
			return err
		}
		if len(page.List) != 1 || page.List[0].InstanceCount != 2 {
			return fmt.Errorf("want %s with 2 instances, got %+v", providerApp, page.List)
		}
		return nil
	})
}

func (v *verifier) instances() error {
	return eventually(viewTimeout, func() error {
		var page struct {
			List []struct {
				IP   string `json:"ip"`
				Name string `json:"name"`
			} `json:"list"`
		}
		if err := v.admin.get("/instance/search", query("appName", providerApp), &page); err != nil {
			return err
		}
		ips := map[string]bool{}
		for _, instance := range page.List {
			ips[instance.IP] = true
			if instance.IP == v.ipA {
				v.instanceA = instance.Name
			}
		}
		if !ips[v.ipA] || !ips[v.ipB] {
			return fmt.Errorf("want provider IPs %s and %s, Admin lists %v", v.ipA, v.ipB, ips)
		}
		return nil
	})
}

func (v *verifier) service() error {
	return eventually(viewTimeout, func() error {
		var page struct {
			List []struct {
				ServiceName string `json:"serviceName"`
				Version     string `json:"version"`
				Group       string `json:"group"`
			} `json:"list"`
		}
		if err := v.admin.get("/service/search", query("keywords", serviceName), &page); err != nil {
			return err
		}
		for _, s := range page.List {
			if s.ServiceName == serviceName && s.Version == serviceVer && s.Group == serviceGroup {
				return nil
			}
		}
		return fmt.Errorf("service %s:%s:%s not listed, got %+v", serviceName, serviceVer, serviceGroup, page.List)
	})
}

func (v *verifier) method() error {
	return eventually(viewTimeout, func() error {
		var methods []struct {
			MethodName string `json:"methodName"`
			Signature  string `json:"signature"`
		}
		if err := v.admin.get("/service/methods", serviceQuery(), &methods); err != nil {
			return err
		}
		if len(methods) != 1 || methods[0].MethodName != "Greet" {
			return fmt.Errorf("want the Greet method only, got %+v", methods)
		}
		v.greetSignature = methods[0].Signature
		var detail struct {
			ParameterTypes []string `json:"parameterTypes"`
			ReturnType     string   `json:"returnType"`
		}
		params := serviceQuery()
		params.Set("methodName", "Greet")
		params.Set("signature", methods[0].Signature)
		if err := v.admin.get("/service/method/detail", params, &detail); err != nil {
			return err
		}
		if len(detail.ParameterTypes) != 1 || detail.ParameterTypes[0] != "adminsample.GreetRequest" ||
			detail.ReturnType != "adminsample.GreetResponse" {
			return fmt.Errorf("unexpected Greet signature %v -> %s", detail.ParameterTypes, detail.ReturnType)
		}
		return nil
	})
}

// invoke calls Greet through the API behind Admin's debug page. For Triple IDL
// services Admin resolves the message types through the provider's reflection
// service.
func (v *verifier) invoke() error {
	if v.instanceA == "" || v.greetSignature == "" {
		return errors.New("instance or method not found by the previous checks")
	}
	req := map[string]any{
		"mesh":         mesh,
		"serviceName":  serviceName,
		"version":      serviceVer,
		"group":        serviceGroup,
		"instanceName": v.instanceA,
		"methodName":   "Greet",
		"signature":    v.greetSignature,
		"args":         []map[string]string{{"name": "admin"}},
		"timeoutMs":    3000,
	}
	var resp struct {
		RawResult struct {
			Greeting string `json:"greeting"`
			Instance string `json:"instance"`
		} `json:"rawResult"`
	}
	if err := v.admin.post("/service/generic/invoke", req, &resp); err != nil {
		return err
	}
	if resp.RawResult.Greeting != "hello admin" || resp.RawResult.Instance != "provider-a" {
		return fmt.Errorf("unexpected result %+v", resp.RawResult)
	}
	return nil
}

// baseline waits until this process routes to both providers, then a few
// seconds more: a rule added and removed right after this client subscribed
// was seen to stay applied in it.
func (v *verifier) baseline() error {
	if err := v.waitRouting("", both); err != nil {
		return err
	}
	time.Sleep(listenerSync)
	return nil
}

func (v *verifier) conditionRoute() error {
	rule := map[string]any{
		"configVersion": "v3.0",
		"scope":         "service",
		"key":           serviceName,
		"enabled":       true,
		"force":         true,
		"runtime":       true,
		"conditions":    []string{"=> host = " + v.ipA},
	}
	if err := v.admin.post("/condition-rule/"+conditionRuleName, rule, nil); err != nil {
		return err
	}
	if err := v.waitRouting("", only("provider-a")); err != nil {
		return err
	}
	if err := v.admin.delete("/condition-rule/" + conditionRuleName); err != nil {
		return err
	}
	return v.waitRouting("", both)
}

func (v *verifier) tagRoute() error {
	rule := map[string]any{
		"configVersion": "v3.0",
		"scope":         "application",
		"key":           providerApp,
		"enabled":       true,
		"force":         false,
		"runtime":       true,
		"tags": []map[string]any{{
			"name":      grayTag,
			"addresses": []string{net.JoinHostPort(v.ipB, strconv.Itoa(providerPort))},
		}},
	}
	if err := v.admin.post("/tag-rule/"+tagRuleName, rule, nil); err != nil {
		return err
	}
	if err := v.waitRouting(grayTag, only("provider-b")); err != nil {
		return err
	}
	// Untagged calls avoid addresses claimed by a tag.
	if err := v.waitRouting("", only("provider-a")); err != nil {
		return err
	}
	if err := v.admin.delete("/tag-rule/" + tagRuleName); err != nil {
		return err
	}
	return v.waitRouting("", both)
}

func (v *verifier) trafficDisable() error {
	if err := v.setTraffic(v.ipB, true); err != nil {
		return err
	}
	if err := v.waitRouting("", only("provider-a")); err != nil {
		return err
	}
	if err := v.setTraffic(v.ipB, false); err != nil {
		return err
	}
	return v.waitRouting("", both)
}

func (v *verifier) metrics() error {
	return eventually(viewTimeout, func() error {
		for metric, app := range map[string]string{
			"dubbo_consumer_requests_total": consumerApp,
			"dubbo_provider_requests_total": providerApp,
		} {
			promQL := fmt.Sprintf(`sum(%s{application_name=%q,interface=%q})`, metric, app, serviceName)
			value, err := v.admin.promQL(promQL)
			if err != nil {
				return err
			}
			if value <= 0 {
				return fmt.Errorf("%s returned %v", promQL, value)
			}
		}
		return nil
	})
}

func (v *verifier) traces() error {
	var dashboard struct {
		FullURL string `json:"fullURL"`
	}
	if err := v.admin.get("/application/trace-dashboard", query("appName", providerApp), &dashboard); err != nil {
		return err
	}
	if dashboard.FullURL == "" {
		return errors.New("no trace dashboard URL configured in Admin")
	}
	return eventually(viewTimeout, func() error {
		return jaegerHasCrossServiceTrace(v.jaeger, consumerApp, providerApp, serviceName)
	})
}

func (v *verifier) setTraffic(ip string, disable bool) error {
	params := query("appName", providerApp)
	params.Set("instanceIP", ip)
	params.Set("trafficDisable", strconv.FormatBool(disable))
	return v.admin.put("/instance/config/trafficDisable", params)
}

// cleanup removes what the checks create so the sample returns to normal.
// Rules that are already gone are fine; anything else is reported, because
// it can leave traffic altered.
func (v *verifier) cleanup() {
	for _, err := range []error{
		v.admin.delete("/condition-rule/" + conditionRuleName),
		v.admin.delete("/tag-rule/" + tagRuleName),
		v.setTraffic(v.ipB, false),
	} {
		if err != nil && !isNotFound(err) {
			fmt.Println("cleanup:", err)
		}
	}
}

// waitRouting polls batches of calls until every batch satisfies want.
func (v *verifier) waitRouting(tag string, want func(map[string]int) bool) error {
	var last map[string]int
	err := eventually(ruleTimeout, func() error {
		last = v.callBatch(tag)
		if !want(last) {
			return errors.New("routing not applied yet")
		}
		return nil
	})
	if err != nil {
		label := "untagged"
		if tag != "" {
			label = "tag=" + tag
		}
		return fmt.Errorf("%s calls still served by %v", label, last)
	}
	return nil
}

func (v *verifier) callBatch(tag string) map[string]int {
	served := map[string]int{}
	for i := 0; i < callsPerPoll; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		if tag != "" {
			ctx = context.WithValue(ctx, constant.AttachmentKey, map[string]string{constant.Tagkey: tag})
		}
		resp, err := v.greet.Greet(ctx, &adminsample.GreetRequest{Name: "verify"})
		cancel()
		if err != nil {
			served["error"]++
			continue
		}
		served[resp.GetInstance()]++
	}
	return served
}

func only(instance string) func(map[string]int) bool {
	return func(served map[string]int) bool {
		return len(served) == 1 && served[instance] == callsPerPoll
	}
}

func both(served map[string]int) bool {
	return len(served) == 2 && served["provider-a"] > 0 && served["provider-b"] > 0
}

func eventually(timeout time.Duration, fn func() error) error {
	deadline := time.Now().Add(timeout)
	for {
		err := fn()
		if err == nil || time.Now().After(deadline) {
			return err
		}
		time.Sleep(pollInterval)
	}
}

func lookup(host string) (string, error) {
	addrs, err := net.LookupHost(host)
	if err != nil || len(addrs) == 0 {
		return "", fmt.Errorf("resolve %s: %v", host, err)
	}
	return addrs[0], nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
