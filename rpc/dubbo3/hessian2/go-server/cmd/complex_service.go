package main

import (
	"context"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common/logger"
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

func (u *ComplexData) JavaClassName() string {
	return "com.apache.dubbo.sample.basic.ComplexData"
}

type ComplexProvider struct {
}

func (u *ComplexProvider) InvokeWithMultiBasicData(ctx context.Context, str string, data []byte, num int, boolValue bool) (int, error) {
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

//func (u *ComplexProvider) InvokeWithStringMap(ctx context.Context, req map[string]string) (map[string]string,error) {
//	logger.Infof("InvokeWithStringList, req = %s", req)
//	return req, nil
//}

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
