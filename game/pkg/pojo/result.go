package pojo

type Result struct {
    Code int32                  `json:"code"` // 0: success, [1, *): error
    Msg  string                 `json:"msg,omitempty"`
    Data map[string]interface{} `json:"data,omitempty"`
}

func (m Result) JavaClassName() string {
    return "org.apache.dubbo.pojo.Result"
}

func (m *Result) Success() bool {
    return m.Code == 0
}
