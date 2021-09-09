package org.apache.dubbo;

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
    HashMap<String, User> stringUserDefinedPtrMapData;
    String[] arrayListData;
    User[] arrayUserData;

    public String GetString(){
        String result = "";
        result += booleanData;
        result += stringData;
        result += int16Data;
        result += intData;
        result += int64Data;
        result += userDefinedData;
        result += byteData;
        result += stringStringHashMap;
        result += stringUserDefinedPtrMapData;
        result += arrayUserData;
        result += arrayListData;
        return result;
    }

}
