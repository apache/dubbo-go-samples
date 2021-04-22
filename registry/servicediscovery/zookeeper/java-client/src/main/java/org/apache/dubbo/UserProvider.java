package org.apache.dubbo;

import org.apache.dubbo.User;

/**
 *@author Patrick Jiang
 *@Description
 *@date 2021/4/19 4:57 下午
 */
public interface UserProvider {
    User GetUser(String userId);
}

