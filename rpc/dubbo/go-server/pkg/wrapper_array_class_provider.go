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

package pkg

import (
	"context"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common/logger"

	hessian "github.com/apache/dubbo-go-hessian2"
)

type WrapperArrayClassProvider struct {
}

func (u *WrapperArrayClassProvider) InvokeWithJavaIntegerArray(ctx context.Context, ia *hessian.IntegerArray) (*hessian.IntegerArray, error) {
	logger.Info("InvokeWithJavaIntegerArray:", ia.Values)

	newV := []int32{1, 2, 3}
	ia.Values = newV
	return ia, nil
}

func (u *WrapperArrayClassProvider) InvokeWithJavaBooleanArray(ctx context.Context, ba *hessian.BooleanArray) (*hessian.BooleanArray, error) {
	logger.Info("InvokeWithJavaBooleanArray:", ba.Values)

	newV := []bool{true, false, true}
	ba.Values = newV
	return ba, nil
}

func (u *WrapperArrayClassProvider) InvokeWithJavaShortArray(ctx context.Context, ba *hessian.ShortArray) (*hessian.ShortArray, error) {
	logger.Info("InvokeWithJavaShortArray:", ba.Values)

	newV := []int16{1, 2, 1000}
	ba.Values = newV
	return ba, nil
}

func (u *WrapperArrayClassProvider) InvokeWithJavaByteArray(ctx context.Context, ia *hessian.ByteArray) (*hessian.ByteArray, error) {
	logger.Info("InvokeWithJavaByteArray:", ia.Values)

	newV := []uint8{1, 2, 200}
	ia.Values = newV
	return ia, nil
}

func (u *WrapperArrayClassProvider) InvokeWithJavaFloatArray(ctx context.Context, ia *hessian.FloatArray) (*hessian.FloatArray, error) {
	logger.Info("InvokeWithJavaFloatArray:", ia.Values)

	newV := []float32{1, 2, 3}
	ia.Values = newV
	return ia, nil
}

func (u *WrapperArrayClassProvider) InvokeWithJavaDoubleArray(ctx context.Context, ia *hessian.DoubleArray) (*hessian.DoubleArray, error) {
	logger.Info("InvokeWithJavaDoubleArray:", ia.Values)

	newV := []float64{1, 2, 3}
	ia.Values = newV
	return ia, nil
}

func (u *WrapperArrayClassProvider) InvokeWithJavaLongArray(ctx context.Context, ia *hessian.LongArray) (*hessian.LongArray, error) {
	logger.Info("InvokeWithJavaLongArray:", ia.Values)

	newV := []int64{1, 2, 3}
	ia.Values = newV
	return ia, nil
}

func (u *WrapperArrayClassProvider) InvokeWithJavaCharacterArray(ctx context.Context, ia *hessian.CharacterArray) (*hessian.CharacterArray, error) {
	logger.Info("InvokeWithJavaCharacterArray:", ia.Values)

	newV := "hello world"
	ia.Values = newV
	return ia, nil
}
