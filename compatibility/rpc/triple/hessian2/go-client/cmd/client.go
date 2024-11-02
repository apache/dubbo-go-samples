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
	"encoding/json"
	"os"
)

import (
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	hessian "github.com/apache/dubbo-go-hessian2"

	gxlog "github.com/dubbogo/gost/log"
)

var userProvider = new(UserProvider)
var complexProvider = new(ComplexProvider)

func init() {
	config.SetConsumerService(userProvider)
	config.SetConsumerService(complexProvider)
	hessian.RegisterPOJO(&User{})
	hessian.RegisterPOJO(&ComplexData{})
}

// need to setup environment variable "DUBBO_GO_CONFIG_PATH" to "conf/dubbogo.yml" before run
func main() {
	if err := config.Load(); err != nil {
		panic(err)
	}

	gxlog.CInfo("\n\n\nstart to test dubbo")
	testNormalService()

	testComplexService()
}

func testNormalService() {
	user, err := userProvider.GetUser(context.TODO(), &User{Name: "laurence"})
	if err != nil {
		gxlog.CError("error: %v\n", err)
		os.Exit(1)
		return
	}
	gxlog.CInfo("response result: %v\n", user)
}

func testComplexService() {
	// test with normal data

	//test without rsp and request
	err := complexProvider.InvokeWithEmptyReq(context.TODO())
	if err != nil {
		gxlog.CError("error: %v\n", err)
		os.Exit(1)
		return
	}

	// test without response
	err = complexProvider.InvokeWithSingleString(context.TODO(), "request string")
	if err != nil {
		gxlog.CError("error: %v\n", err)
		os.Exit(1)
		return
	}

	err = complexProvider.InvokeWithStringList(context.TODO(), []string{"myfirststring", "mysecondstring"})
	if err != nil {
		gxlog.CError("error: %v\n", err)
		os.Exit(1)
		return
	}

	err = complexProvider.InvokeWithMultiString(context.TODO(), "first string", "secondString", "third str")
	if err != nil {
		gxlog.CError("error: %v\n", err)
		os.Exit(1)
		return
	}

	// test without request
	rsp, err := complexProvider.InvokeWithEmptyReqStringRsp(context.TODO())
	if err != nil {
		gxlog.CError("error: %v\n", err)
		os.Exit(1)
		return
	}
	gxlog.CInfo("get InvokeWithEmptyReqStringRsp rsp = %+v", rsp)

	// complex data
	stringIntMapData := make(map[string]int)
	stringIntMapData["test1"] = 1
	stringIntMapData["test2"] = 2

	stringStringMapData := make(map[string]string)
	stringStringMapData["test1"] = "1"
	stringStringMapData["test2"] = "2"

	stringUserMapData := make(map[string]User)
	stringUserMapData["test1"] = User{Name: "1"}
	stringUserMapData["test2"] = User{Name: "2"}

	stringUintMapData := make(map[string]uint32)
	stringUintMapData["test1"] = 1
	stringUintMapData["test2"] = 2

	stringUserPtrMapData := make(map[string]*User)
	stringUserPtrMapData["test1"] = &User{Name: "1"}
	stringUserPtrMapData["test2"] = &User{Name: "2"}

	intStringMapData := make(map[int]string)
	intStringMapData[1] = "1"
	intStringMapData[2] = "2"

	data, _ := json.Marshal(User{Name: "myJson", Age: 19, Id: "jsonID"})

	cplexData := &ComplexData{
		BooleanData: true,
		StringData:  "testString",
		//UIntData: 8,
		UInt8Data:  8,
		UInt16Data: 16,
		UInt32Data: 32,
		UInt64Data: 64,
		Int8Data:   8,
		Int16Data:  16,
		Int32Data:  32,
		Int64Data:  64,
		IntData:    8,
		//StringIntMapData: stringIntMapData,
		StringStringMapData: stringStringMapData,
		//StringUserDefinedMapData:stringUserMapData,
		//StringUIntMapData: stringUintMapData,
		StringUserDefinedPtrMapData: stringUserPtrMapData,
		//IntStringMapData: intStringMapData,
		UserDefinedData:         User{Name: "myuser", Age: 18, Id: "testid"},
		UserDefinedDataPtr:      &User{Name: "myuserPtr", Age: 18, Id: "testid"},
		ByteData:                data,
		ArrayListData:           []string{"string1", "string2", "string3"},
		ArrayUserDefinedData:    []User{{Name: "name1", Id: "id1", Age: 19}, {Name: "name1", Id: "id1", Age: 19}, {Name: "name1", Id: "id1", Age: 19}},
		ArrayUserDefinedPtrData: []*User{{Name: "name1", Id: "id1", Age: 19}, {Name: "name1", Id: "id1", Age: 19}, {Name: "name1", Id: "id1", Age: 19}},
	}

	cplxRsp, err := complexProvider.InvokeWithComplexReqComplexRspPtr(context.TODO(), cplexData)
	if err != nil {
		gxlog.CError("error: %v\n", err)
		os.Exit(1)
		return
	}
	gxlog.CInfo("get InvokeWithComplexReqComplexRspPtr rsp = %+v", cplxRsp)

	intRsp, err := complexProvider.InvokeWithMultiBasicData(context.TODO(), "reqstr", []byte{1, 2, 4}, 32, true)
	if err != nil {
		panic(err)
	}
	gxlog.CInfo("get InvokeWithMultiBasicData rsp = %d", intRsp)
}
