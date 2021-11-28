package org.apache.dubbo.sample;

public interface WrapperArrayClassProvider {
    Integer[] InvokeWithJavaIntegerArray(Integer[] ia);
    Byte[] InvokeWithJavaByteArray(Byte[] ba);
    Short[] InvokeWithJavaShortArray(Short[] ia);
    Long[] InvokeWithJavaLongArray(Long[] ia);
    Character[] InvokeWithJavaCharacterArray(Character[] ia);
    Float[] InvokeWithJavaFloatArray(Float[] ia);
    Double[] InvokeWithJavaDoubleArray(Double[] ia);
    Boolean[] InvokeWithJavaBooleanArray(Boolean[] ia);
}
