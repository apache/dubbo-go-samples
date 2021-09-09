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

package pkg

import (
	"context"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common/logger"
)

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

func (u *ComplexData) JavaClassName() string {
	return "org.apache.dubbo.ComplexData"
}

type ComplexProvider struct {
}

func (u *ComplexProvider) InvokeWithMultiBasicData(ctx context.Context, str string, data []byte, num int32, boolValue bool) (int32, error) {
	logger.Info("InvokeWithMultiBasicData", str, " ", data, " ", num, " ", boolValue)
	return num, nil
}

func (u *ComplexProvider) InvokeWithEmptyReq(ctx context.Context) error {
	logger.Info("InvokeWithEmptyReq")
	return nil
}

func (u *ComplexProvider) InvokeWithSingleString(ctx context.Context, req string) error {
	logger.Infof("InvokeWithSingleString, req = %s", req)
	return nil
}

func (u *ComplexProvider) InvokeWithMultiString(ctx context.Context, req, req2, req3 string) error {
	logger.Info("InvokeWithMultiString, req = ", req, req2, req3)
	return nil
}

func (u *ComplexProvider) InvokeWithStringList(ctx context.Context, req []string) error {
	logger.Infof("InvokeWithStringList, req = %s", req)
	return nil
}

func (u *ComplexProvider) InvokeWithEmptyReqStringRsp(ctx context.Context) (string, error) {
	logger.Infof("InvokeWithEmptyReqStringRsp")
	return "success rsp", nil
}

func (u *ComplexProvider) InvokeWithEmptyReqMultiStringRsp(ctx context.Context) (string, string, string, error) {
	logger.Infof("InvokeWithEmptyReqMultiStringRsp")
	return "success rsp1", "success rsp2", "success rsp3", nil
}

func (u *ComplexProvider) InvokeWithComplexReqComplexRspPtr(ctx context.Context, req *ComplexData) (*ComplexData, error) {
	logger.Infof("InvokeWithComplexReqComplexRsp req = %+v", req)
	return req, nil
}

func (u *ComplexProvider) Reference() string {
	return "ComplexProvider"
}
