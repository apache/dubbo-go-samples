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
)

import (
	"dubbo.apache.org/dubbo-go/v3/config"

	hessian "github.com/apache/dubbo-go-hessian2"
)

func init() {
	// ------for hessian2------
	hessian.RegisterPOJO(&User{})
	config.SetProviderService(new(UserProvider))
}

type ComplexData struct {
	BooleanData bool

	StringData string

	//UIntData   uint
	UInt8Data  uint8
	UInt16Data uint16
	UInt32Data uint32
	UInt64Data uint64

	IntData   int
	Int8Data  int8
	Int16Data int16
	Int32Data int32
	Int64Data int64

	StringStringMapData map[string]string
	//StringIntMapData            map[string]int
	//StringUIntMapData           map[string]uint32
	//IntStringMapData            map[int]string
	//StringUserDefinedMapData    map[string]User
	StringUserDefinedPtrMapData map[string]*User

	UserDefinedData    User
	UserDefinedDataPtr *User

	ByteData []byte

	ArrayListData           []string
	ArrayUserDefinedData    []User
	ArrayUserDefinedPtrData []*User
}

type ComplexProvider struct {
	InvokeWithEmptyReq                func(ctx context.Context) error
	InvokeWithSingleString            func(ctx context.Context, req string) error
	InvokeWithMultiString             func(ctx context.Context, req, req2, req3 string) error
	InvokeWithStringList              func(ctx context.Context, req []string) error
	InvokeWithEmptyReqStringRsp       func(ctx context.Context) (string, error)
	InvokeWithComplexReqComplexRspPtr func(ctx context.Context, req *ComplexData) (*ComplexData, error)
	InvokeWithMultiBasicData          func(ctx context.Context, str string, data []byte, num int, boolValue bool) (int, error)
	//InvokeWithStringMap func(ctx context.Context, req map[string]string) (map[string]string,error)
}

func (u *ComplexData) JavaClassName() string {
	return "com.apache.dubbo.sample.basic.ComplexData"
}
