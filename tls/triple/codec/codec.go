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

package codec

import (
	"github.com/dubbogo/grpc-go/encoding"

	msgpack "github.com/ugorji/go/codec"

	triCommon "github.com/dubbogo/triple/pkg/common"
)

func init() {
	triCommon.SetTripleCodec("msgpack", NewMsgpackCodec)
}

func NewMsgpackCodec() encoding.Codec {
	return msgpackCodec{}
}

type msgpackCodec struct{}

func (c msgpackCodec) Name() string {
	return "msgpack"
}

func (c msgpackCodec) Marshal(message any) ([]byte, error) {
	var out []byte
	encoder := msgpack.NewEncoderBytes(&out, new(msgpack.MsgpackHandle))
	return out, encoder.Encode(message)
}

func (c msgpackCodec) Unmarshal(binary []byte, message any) error {
	decoder := msgpack.NewDecoderBytes(binary, new(msgpack.MsgpackHandle))
	return decoder.Decode(message)
}
