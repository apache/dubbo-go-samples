package com.apache.dubbo.sample.basic;

public interface ComplexProvider {
    void invokeWithEmptyReq();
    void invokeWithSingleString(String req);
    void invokeWithStringList(String[] req);
    void invokeWithMultiString(String str1, String str2, String str3);
    String invokeWithEmptyReqStringRsp ();
    ComplexData invokeWithComplexReqComplexRspPtr(ComplexData complexData);
}
