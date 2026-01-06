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
	"context"
	"testing"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/client"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/config/generic"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	hessian "github.com/apache/dubbo-go-hessian2"

	"github.com/apache/dubbo-go-samples/generic/go-client/pkg"
)

func init() {
	hessian.RegisterPOJO(&pkg.User{})
}

// setupDubboGenericService creates a GenericService using Dubbo protocol
func setupDubboGenericService(t *testing.T) *generic.GenericService {
	ins, err := dubbo.NewInstance(
		dubbo.WithName("generic-test-dubbo"),
	)
	if err != nil {
		t.Fatalf("Failed to create dubbo instance: %v", err)
	}

	cli, err := ins.NewClient(
		client.WithClientProtocolDubbo(),
		client.WithClientSerialization(constant.Hessian2Serialization),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	svc, err := cli.NewGenericService(
		"org.apache.dubbo.samples.UserProvider",
		client.WithURL("dubbo://127.0.0.1:20000"),
		client.WithVersion("1.0.0"),
		client.WithGroup("dubbo"),
	)
	if err != nil {
		t.Fatalf("Failed to create generic service: %v", err)
	}

	return svc
}

// setupTripleGenericService creates a GenericService using Triple protocol
func setupTripleGenericService(t *testing.T) *generic.GenericService {
	ins, err := dubbo.NewInstance(
		dubbo.WithName("generic-test-triple"),
	)
	if err != nil {
		t.Fatalf("Failed to create dubbo instance: %v", err)
	}

	cli, err := ins.NewClient(
		client.WithClientProtocolTriple(),
		client.WithClientSerialization(constant.Hessian2Serialization),
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	svc, err := cli.NewGenericService(
		"org.apache.dubbo.samples.UserProvider",
		client.WithURL("tri://127.0.0.1:50052"),
		client.WithVersion("1.0.0"),
		client.WithGroup("triple"),
	)
	if err != nil {
		t.Fatalf("Failed to create generic service: %v", err)
	}

	return svc
}

// =============================================================================
// Dubbo Protocol Tests
// =============================================================================

func TestDubbo_GenericCall_StringArg(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	svc := setupDubboGenericService(t)

	result, err := svc.Invoke(context.Background(), "GetUser1", []string{"java.lang.String"}, []hessian.Object{"A001"})
	if err != nil {
		t.Fatalf("GetUser1 failed: %v", err)
	}

	resultMap, ok := result.(map[interface{}]interface{})
	if !ok {
		t.Fatalf("Expected map result, got %T", result)
	}

	if resultMap["id"] != "A001" {
		t.Errorf("Expected id=A001, got %v", resultMap["id"])
	}
}

func TestDubbo_GenericCall_MultipleArgs(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	svc := setupDubboGenericService(t)

	result, err := svc.Invoke(context.Background(), "GetUser2", []string{"java.lang.String", "java.lang.String"}, []hessian.Object{"A002", "CustomName"})
	if err != nil {
		t.Fatalf("GetUser2 failed: %v", err)
	}

	resultMap, ok := result.(map[interface{}]interface{})
	if !ok {
		t.Fatalf("Expected map result, got %T", result)
	}

	if resultMap["name"] != "CustomName" {
		t.Errorf("Expected name=CustomName, got %v", resultMap["name"])
	}
}

func TestDubbo_GenericCall_IntArg(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	svc := setupDubboGenericService(t)

	result, err := svc.Invoke(context.Background(), "GetUser3", []string{"int"}, []hessian.Object{int32(1)})
	if err != nil {
		t.Fatalf("GetUser3 failed: %v", err)
	}

	resultMap, ok := result.(map[interface{}]interface{})
	if !ok {
		t.Fatalf("Expected map result, got %T", result)
	}

	if resultMap["id"] != "1" {
		t.Errorf("Expected id=1, got %v", resultMap["id"])
	}
}

func TestDubbo_GenericCall_NoArgs(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	svc := setupDubboGenericService(t)

	result, err := svc.Invoke(context.Background(), "GetOneUser", []string{}, []hessian.Object{})
	if err != nil {
		t.Fatalf("GetOneUser failed: %v", err)
	}

	resultMap, ok := result.(map[interface{}]interface{})
	if !ok {
		t.Fatalf("Expected map result, got %T", result)
	}

	if resultMap["id"] == nil {
		t.Error("Expected non-nil id")
	}
}

func TestDubbo_GenericCall_ArrayArg(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	svc := setupDubboGenericService(t)

	result, err := svc.Invoke(context.Background(), "GetUsers", []string{"[Ljava.lang.String;"}, []hessian.Object{[]string{"001", "002", "003"}})
	if err != nil {
		t.Fatalf("GetUsers failed: %v", err)
	}

	resultList, ok := result.([]interface{})
	if !ok {
		t.Fatalf("Expected list result, got %T", result)
	}

	if len(resultList) != 3 {
		t.Errorf("Expected 3 users, got %d", len(resultList))
	}
}

func TestDubbo_GenericCall_POJOArg(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	svc := setupDubboGenericService(t)

	user := &pkg.User{
		ID:   "test001",
		Name: "TestUser",
		Age:  30,
		Time: time.Now(),
	}

	result, err := svc.Invoke(context.Background(), "QueryUser", []string{"org.apache.dubbo.samples.User"}, []hessian.Object{user})
	if err != nil {
		t.Fatalf("QueryUser failed: %v", err)
	}

	resultMap, ok := result.(map[interface{}]interface{})
	if !ok {
		t.Fatalf("Expected map result, got %T", result)
	}

	if resultMap["name"] != "TestUser" {
		t.Errorf("Expected name=TestUser, got %v", resultMap["name"])
	}
}

// =============================================================================
// Triple Protocol Tests
// =============================================================================

func TestTriple_GenericCall_StringArg(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	svc := setupTripleGenericService(t)

	result, err := svc.Invoke(context.Background(), "GetUser1", []string{"java.lang.String"}, []hessian.Object{"A001"})
	if err != nil {
		t.Fatalf("GetUser1 failed: %v", err)
	}

	resultMap, ok := result.(map[interface{}]interface{})
	if !ok {
		t.Fatalf("Expected map result, got %T", result)
	}

	if resultMap["id"] != "A001" {
		t.Errorf("Expected id=A001, got %v", resultMap["id"])
	}
}

func TestTriple_GenericCall_MultipleArgs(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	svc := setupTripleGenericService(t)

	result, err := svc.Invoke(context.Background(), "GetUser2", []string{"java.lang.String", "java.lang.String"}, []hessian.Object{"A002", "CustomName"})
	if err != nil {
		t.Fatalf("GetUser2 failed: %v", err)
	}

	resultMap, ok := result.(map[interface{}]interface{})
	if !ok {
		t.Fatalf("Expected map result, got %T", result)
	}

	if resultMap["name"] != "CustomName" {
		t.Errorf("Expected name=CustomName, got %v", resultMap["name"])
	}
}

func TestTriple_GenericCall_IntArg(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	svc := setupTripleGenericService(t)

	result, err := svc.Invoke(context.Background(), "GetUser3", []string{"int"}, []hessian.Object{int32(1)})
	if err != nil {
		t.Fatalf("GetUser3 failed: %v", err)
	}

	resultMap, ok := result.(map[interface{}]interface{})
	if !ok {
		t.Fatalf("Expected map result, got %T", result)
	}

	if resultMap["id"] != "1" {
		t.Errorf("Expected id=1, got %v", resultMap["id"])
	}
}

func TestTriple_GenericCall_NoArgs(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	svc := setupTripleGenericService(t)

	result, err := svc.Invoke(context.Background(), "GetOneUser", []string{}, []hessian.Object{})
	if err != nil {
		t.Fatalf("GetOneUser failed: %v", err)
	}

	resultMap, ok := result.(map[interface{}]interface{})
	if !ok {
		t.Fatalf("Expected map result, got %T", result)
	}

	if resultMap["id"] == nil {
		t.Error("Expected non-nil id")
	}
}

func TestTriple_GenericCall_ArrayArg(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	svc := setupTripleGenericService(t)

	result, err := svc.Invoke(context.Background(), "GetUsers", []string{"[Ljava.lang.String;"}, []hessian.Object{[]string{"001", "002", "003"}})
	if err != nil {
		t.Fatalf("GetUsers failed: %v", err)
	}

	resultList, ok := result.([]interface{})
	if !ok {
		t.Fatalf("Expected list result, got %T", result)
	}

	if len(resultList) != 3 {
		t.Errorf("Expected 3 users, got %d", len(resultList))
	}
}

func TestTriple_GenericCall_POJOArg(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	svc := setupTripleGenericService(t)

	user := &pkg.User{
		ID:   "test001",
		Name: "TestUser",
		Age:  30,
		Time: time.Now(),
	}

	result, err := svc.Invoke(context.Background(), "QueryUser", []string{"org.apache.dubbo.samples.User"}, []hessian.Object{user})
	if err != nil {
		t.Fatalf("QueryUser failed: %v", err)
	}

	resultMap, ok := result.(map[interface{}]interface{})
	if !ok {
		t.Fatalf("Expected map result, got %T", result)
	}

	if resultMap["name"] != "TestUser" {
		t.Errorf("Expected name=TestUser, got %v", resultMap["name"])
	}
}

func TestTriple_GenericCall_POJOArrayArg(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	svc := setupTripleGenericService(t)

	users := []hessian.Object{
		&pkg.User{ID: "001", Name: "User1", Age: 20, Time: time.Now()},
		&pkg.User{ID: "002", Name: "User2", Age: 25, Time: time.Now()},
		&pkg.User{ID: "003", Name: "User3", Age: 30, Time: time.Now()},
	}

	result, err := svc.Invoke(context.Background(), "QueryUsers", []string{"java.util.List"}, []hessian.Object{users})
	if err != nil {
		t.Fatalf("QueryUsers failed: %v", err)
	}

	// Result can be []interface{} or []hessian.Object depending on Java server response
	var listLen int
	switch v := result.(type) {
	case []interface{}:
		listLen = len(v)
	case []hessian.Object:
		listLen = len(v)
	default:
		t.Fatalf("Expected list result, got %T", result)
	}

	if listLen != 3 {
		t.Errorf("Expected 3 users, got %d", listLen)
	}
}

// =============================================================================
// Edge Case Tests (Triple Protocol)
// =============================================================================

func TestTriple_GenericCall_EmptyStringArg(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	svc := setupTripleGenericService(t)

	result, err := svc.Invoke(context.Background(), "GetUser1", []string{"java.lang.String"}, []hessian.Object{""})
	if err != nil {
		t.Fatalf("GetUser1 with empty string failed: %v", err)
	}

	resultMap, ok := result.(map[interface{}]interface{})
	if !ok {
		t.Fatalf("Expected map result, got %T", result)
	}

	if resultMap["id"] != "" {
		t.Errorf("Expected empty id, got %v", resultMap["id"])
	}
}

func TestTriple_GenericCall_ZeroIntArg(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	svc := setupTripleGenericService(t)

	result, err := svc.Invoke(context.Background(), "GetUser3", []string{"int"}, []hessian.Object{int32(0)})
	if err != nil {
		t.Fatalf("GetUser3 with zero failed: %v", err)
	}

	resultMap, ok := result.(map[interface{}]interface{})
	if !ok {
		t.Fatalf("Expected map result, got %T", result)
	}

	if resultMap["id"] != "0" {
		t.Errorf("Expected id=0, got %v", resultMap["id"])
	}
}

func TestTriple_GenericCall_EmptyArrayArg(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	svc := setupTripleGenericService(t)

	result, err := svc.Invoke(context.Background(), "GetUsers", []string{"[Ljava.lang.String;"}, []hessian.Object{[]string{}})
	if err != nil {
		t.Fatalf("GetUsers with empty array failed: %v", err)
	}

	resultList, ok := result.([]interface{})
	if !ok {
		t.Fatalf("Expected list result, got %T", result)
	}

	if len(resultList) != 0 {
		t.Errorf("Expected empty list, got %d items", len(resultList))
	}
}

// =============================================================================
// Error Case Tests
// =============================================================================

func TestTriple_GenericCall_NonExistentMethod(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	svc := setupTripleGenericService(t)

	_, err := svc.Invoke(context.Background(), "NonExistentMethod", []string{}, []hessian.Object{})
	if err == nil {
		t.Error("Expected error for non-existent method, got nil")
	}
}

// =============================================================================
// Context Tests
// =============================================================================

func TestTriple_GenericCall_WithTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	svc := setupTripleGenericService(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := svc.Invoke(ctx, "GetUser1", []string{"java.lang.String"}, []hessian.Object{"A001"})
	if err != nil {
		t.Fatalf("GetUser1 with timeout failed: %v", err)
	}
}

func TestTriple_GenericCall_CancelledContext(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	svc := setupTripleGenericService(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := svc.Invoke(ctx, "GetUser1", []string{"java.lang.String"}, []hessian.Object{"A001"})
	if err == nil {
		t.Error("Expected error for cancelled context, got nil")
	}
}

// =============================================================================
// Benchmark Tests
// =============================================================================

func BenchmarkDubbo_GenericCall(b *testing.B) {
	ins, _ := dubbo.NewInstance(dubbo.WithName("benchmark-dubbo"))
	cli, _ := ins.NewClient(
		client.WithClientProtocolDubbo(),
		client.WithClientSerialization(constant.Hessian2Serialization),
	)
	svc, err := cli.NewGenericService(
		"org.apache.dubbo.samples.UserProvider",
		client.WithURL("dubbo://127.0.0.1:20000"),
		client.WithVersion("1.0.0"),
		client.WithGroup("dubbo"),
	)
	if err != nil {
		b.Skip("Dubbo server not available")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = svc.Invoke(context.Background(), "GetUser1", []string{"java.lang.String"}, []hessian.Object{"A001"})
	}
}

func BenchmarkTriple_GenericCall(b *testing.B) {
	ins, _ := dubbo.NewInstance(dubbo.WithName("benchmark-triple"))
	cli, _ := ins.NewClient(
		client.WithClientProtocolTriple(),
		client.WithClientSerialization(constant.Hessian2Serialization),
	)
	svc, err := cli.NewGenericService(
		"org.apache.dubbo.samples.UserProvider",
		client.WithURL("tri://127.0.0.1:50052"),
		client.WithVersion("1.0.0"),
		client.WithGroup("triple"),
	)
	if err != nil {
		b.Skip("Triple server not available")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = svc.Invoke(context.Background(), "GetUser1", []string{"java.lang.String"}, []hessian.Object{"A001"})
	}
}
