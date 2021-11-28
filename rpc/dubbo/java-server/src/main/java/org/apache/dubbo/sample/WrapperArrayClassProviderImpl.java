package org.apache.dubbo.sample;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Arrays;

public class WrapperArrayClassProviderImpl implements WrapperArrayClassProvider {
    private static final Logger LOG = LoggerFactory.getLogger("WrapperArrayClassProviderImpl"); // Output to com.dubbogo.user-server.log


    @Override
    public Integer[] InvokeWithJavaIntegerArray(Integer[] ia) {
        LOG.info("InvokeWithJavaIntegerArray:?", Arrays.toString(ia));
        return new Integer[]{1, 2, 3};
    }

    @Override
    public Byte[] InvokeWithJavaByteArray(Byte[] ba) {
        LOG.info("InvokeWithJavaByteArray:?", Arrays.toString(ba));
        return new Byte[]{2, 3, 4};
    }

    @Override
    public Short[] InvokeWithJavaShortArray(Short[] ia) {
        LOG.info("InvokeWithJavaShortArray:?", Arrays.toString(ia));
        return new Short[]{4, 5, 6};
    }

    @Override
    public Long[] InvokeWithJavaLongArray(Long[] ia) {
        LOG.info("InvokeWithJavaLongArray:?", Arrays.toString(ia));
        return new Long[]{7L, 8L, 9L};
    }

    @Override
    public Character[] InvokeWithJavaCharacterArray(Character[] ia) {
        LOG.info("InvokeWithJavaCharacterArray:?", Arrays.toString(ia));
        return new Character[]{'A', 'B', 'C'};
    }

    @Override
    public Float[] InvokeWithJavaFloatArray(Float[] ia) {
        LOG.info("InvokeWithJavaFloatArray:?", Arrays.toString(ia));
        return new Float[]{1.2f, 4.5f, 7.8f};
    }

    @Override
    public Double[] InvokeWithJavaDoubleArray(Double[] ia) {
        LOG.info("InvokeWithJavaDoubleArray:?", Arrays.toString(ia));
        return new Double[]{8d, 9d, 10d};
    }

    @Override
    public Boolean[] InvokeWithJavaBooleanArray(Boolean[] ia) {
        LOG.info("InvokeWithJavaBooleanArray:?", Arrays.toString(ia));
        return new Boolean[]{true, false, true};
    }
}
