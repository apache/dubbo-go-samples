package org.apache.dubbo;

import java.util.Arrays;

public class ComplexProviderImpl implements ComplexProvider {
    @Override
    public int InvokeWithMultiBasicData(String str, byte[] data, int num, boolean boolValue) {
        System.out.println("InvokeWithMultiBasicData, str: " + str + ", data: " + Arrays.toString(data) + ", num: " + num + ", boolValue:" +
                " " + boolValue);
        return num;
    }

    @Override
    public void InvokeWithEmptyReq() {
        System.out.println("InvokeWithEmptyReq");
    }

    @Override
    public void InvokeWithSingleString(String req) {
        System.out.println("InvokeWithEmptyReq" + req);
    }

    @Override
    public void InvokeWithStringList(String[] req) {
        System.out.println("InvokeWithEmptyReq" + req);
    }

    @Override
    public void InvokeWithMultiString(String str1, String str2, String str3) {
        System.out.println("InvokeWithEmptyReq" + str1 + str2 + str3);
    }

    @Override
    public String InvokeWithEmptyReqStringRsp() {
        System.out.println("InvokeWithEmptyReq");
        return "invoke success";
    }

    @Override
    public ComplexData InvokeWithComplexReqComplexRspPtr(ComplexData complexData) {
        System.out.println("InvokeWithComplexReqComplexRspPtr = "+ complexData.GetString());
        return complexData;
    }
}
