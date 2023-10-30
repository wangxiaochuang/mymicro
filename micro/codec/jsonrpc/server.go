package jsonrpc

import (
	"encoding/json"
	"io"

	"github.com/wxc/micro/codec"
)

type serverCodec struct {
	dec *json.Decoder // for reading JSON values
	enc *json.Encoder // for writing JSON values
	c   io.Closer

	req  serverRequest
	resp serverResponse
}

type serverRequest struct {
	ID     interface{}      `json:"id"`
	Params *json.RawMessage `json:"params"`
	Method string           `json:"method"`
}

type serverResponse struct {
	ID     interface{} `json:"id"`
	Result interface{} `json:"result"`
	Error  interface{} `json:"error"`
}

func newServerCodec(conn io.ReadWriteCloser) *serverCodec {
	return &serverCodec{
		dec: json.NewDecoder(conn),
		enc: json.NewEncoder(conn),
		c:   conn,
	}
}

func (r *serverRequest) reset() {
	r.Method = ""
	if r.Params != nil {
		*r.Params = (*r.Params)[0:0]
	}
	if r.ID != nil {
		r.ID = nil
	}
}

func (c *serverCodec) ReadHeader(m *codec.Message) error {
	panic(" in ReadHeader")
}

func (c *serverCodec) ReadBody(x interface{}) error {
	if x == nil {
		return nil
	}
	var params [1]interface{}
	params[0] = x
	return json.Unmarshal(*c.req.Params, &params)
}

func (c *serverCodec) Write(m *codec.Message, x interface{}) error {
	var resp serverResponse
	resp.ID = m.Id
	resp.Result = x
	if m.Error == "" {
		resp.Error = nil
	} else {
		resp.Error = m.Error
	}
	return c.enc.Encode(resp)
}

func (c *serverCodec) Close() error {
	return c.c.Close()
}
