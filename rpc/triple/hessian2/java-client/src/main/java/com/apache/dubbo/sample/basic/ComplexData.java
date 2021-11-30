package com.apache.dubbo.sample.basic;

import java.io.Serializable;
import java.util.Arrays;
import java.util.HashMap;
import java.util.StringJoiner;

public class ComplexData implements Serializable {
    boolean booleanData;
    String stringData;

    short int16Data;
    int intData;
    long int64Data;

    User userDefinedData;
    byte [] byteData;
    HashMap<String, String> stringStringHashMap;
    HashMap<String, User> stringUserDefinedPtrMapData;
    String[] arrayListData;
    User[] arrayUserData;

    @Override
    public String toString() {
        return new StringJoiner(", ", ComplexData.class.getSimpleName() + "[", "]")
                .add("booleanData=" + booleanData)
                .add("stringData='" + stringData + "'")
                .add("int16Data=" + int16Data)
                .add("intData=" + intData)
                .add("int64Data=" + int64Data)
                .add("userDefinedData=" + userDefinedData)
                .add("byteData=" + Arrays.toString(byteData))
                .add("stringStringHashMap=" + stringStringHashMap)
                .add("stringUserDefinedPtrMapData=" + stringUserDefinedPtrMapData)
                .add("arrayListData=" + Arrays.toString(arrayListData))
                .add("arrayUserData=" + Arrays.toString(arrayUserData))
                .toString();
    }
}
