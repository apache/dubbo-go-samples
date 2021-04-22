package org.apache.dubbo;

import org.apache.dubbo.common.logger.Logger;
import org.apache.dubbo.common.logger.LoggerFactory;
import org.apache.dubbo.rpc.RpcContext;

/**
 *@author Patrick Jiang
 *@Description
 *@date 2021/4/19 5:06 下午
 */
public class UserProviderImpl implements UserProvider {
    private static final Logger logger = LoggerFactory.getLogger(UserProviderImpl.class);
    
    @Override
    public User GetUser(String userId) {
        logger.info("Hello " + userId + ", request from consumer: " + RpcContext.getContext().getRemoteAddress());
        return new User(userId, "Alex Stocks", 18);
    }
}

