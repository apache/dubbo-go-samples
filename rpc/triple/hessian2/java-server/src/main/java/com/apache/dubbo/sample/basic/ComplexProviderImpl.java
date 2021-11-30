package com.apache.dubbo.sample.basic;

public class ComplexProviderImpl implements ComplexProvider {
    @Override
    public void invokeWithEmptyReq() {
        System.out.println("invokeWithEmptyReq");
    }

    @Override
    public void invokeWithSingleString(String req) {
        System.out.println("invokeWithEmptyReq" + req);
    }

    @Override
    public void invokeWithStringList(String[] req) {
        System.out.println("invokeWithEmptyReq" + req);
    }

    @Override
    public void invokeWithMultiString(String str1, String str2, String str3) {
        System.out.println("invokeWithEmptyReq" + str1 + str2 + str3);
    }

    @Override
    public String invokeWithEmptyReqStringRsp() {
        System.out.println("invokeWithEmptyReq");
        return "invoke success";
    }

    @Override
    public ComplexData invokeWithComplexReqComplexRspPtr(ComplexData complexData) {
        System.out.println("invokeWithComplexReqComplexRspPtr = "+ complexData.GetString());
        return complexData;
    }

    @Override
    public int invokeWithMultiBasicData(String str, byte[] data, int num, boolean boolValue) {
        System.out.println("invokeWithEmptyReq" + str + data + num + boolValue);
        return num;
    }
}
