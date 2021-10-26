/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package org.apache.dubbo;

import java.io.Serializable;

import lombok.Data;

/**
 * Response, the first alpha of field name is Upper case to compatible with golang.
 */
@Data
public final class Response<T> implements Serializable {
    private static final long serialVersionUID = 3727205004706510648L;
    public static final Integer OK = 200;
    public static final Integer ERR = 500;
    private Integer Status;
    private String Err;
    private T Data;

    public static <T> Response<T> ok(T data) {
        Response<T> r = new Response<>();
        r.Status = OK;
        r.Data = data;
        return r;
    }
}