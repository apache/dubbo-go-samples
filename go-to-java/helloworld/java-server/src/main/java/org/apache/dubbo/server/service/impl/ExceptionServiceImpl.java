package org.apache.dubbo.server.service.impl;

import lombok.extern.slf4j.Slf4j;
import org.apache.dubbo.config.annotation.DubboService;
import org.apache.dubbo.server.service.ExceptionService;

/**
 * @author zhaoyunxing
 */
@Slf4j
@DubboService
public class ExceptionServiceImpl implements ExceptionService {

    @Override
    public String getThrowable(String name) throws Throwable {
        log.info("getThrowable: {}", name);
        throw new Throwable("throwable exception");
    }
}
