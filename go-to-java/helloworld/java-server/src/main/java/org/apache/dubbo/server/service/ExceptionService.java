package org.apache.dubbo.server.service;

/**
 * @author zhaoyunxing
 */
public interface ExceptionService {

    String getThrowable(String name) throws Throwable;
}
