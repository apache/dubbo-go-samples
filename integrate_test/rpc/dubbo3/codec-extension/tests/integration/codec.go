package integration

import (
	"encoding/json"
)

import (
	triCommon "github.com/dubbogo/triple/pkg/common"
)

func init() {
	triCommon.SetTripleCodec("json", NewJSONCodec)
}

func NewJSONCodec() triCommon.Codec {
	return &JSONCodec{}
}

type JSONCodec struct {
}

func (j *JSONCodec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (j *JSONCodec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
