package com.apache.dubbo.sample.basic;

import java.io.Serializable;
import java.util.HashMap;

public class ComplexData implements Serializable {
    boolean booleanData;
    String stringData;

    short int16Data;
    int intData;
    long int64Data;

    User userDefinedData;
    byte [] byteData;
    HashMap<String, String> stringStringHashMap;
//    HashMap<String, User> stringUserDefinedPtrMapData;
    String[] arrayListData;
//    User[] arrayUserData;

}
