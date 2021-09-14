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

	newV := []int32{1, 2, 3}
	ba.Values = newV
	return ba, nil
}

func (u *WrapperArrayClassProvider) InvokeWithJavaByteArray(ctx context.Context, ia *hessian.ByteArray) (*hessian.ByteArray, error) {
	logger.Info("InvokeWithJavaByteArray:", ia.Values)

	newV := []int32{1, 2, 3}
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
