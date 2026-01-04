package proto

import (
	"bytes"
	"strings"
	"text/template"
)

const protoTemplate = `
syntax = "proto3";

package {{.Package}};

import "google/api/annotations.proto";
import "buf/validate/validate.proto";
import "errors/errors.proto";

option go_package = "{{.GoPackage}}";

enum ErrorReason {
	// 设置缺省错误码
	option (errors.default_code) = 500;

	ERROR_REASON_UNSPECIFIED          = 0 [(errors.code) = 500];
	// 请求错误，包含参数校验不过
	ERROR_REASON_BAD_REQUEST          = 1 [(errors.code) = 400];
	// 找不到资源
	ERROR_REASON_NOT_FOUND            = 2  [(errors.code) = 404];
	// 服务内部未知
	ERROR_REASON_INTERNAL_ERROR       = 3 [(errors.code) = 500];
	// 数据库异常
	ERROR_REASON_DB_UNAVAILABLE       = 4 [(errors.code) = 503];
	// 缓存异常
	ERROR_REASON_CACHE_UNAVAILABLE    = 5 [(errors.code) = 503];
	// mq异常
	ERROR_REASON_MQ_UNAVAILABLE       = 6 [(errors.code) = 503];
}

service {{.Service}}Service {
	rpc Create{{.Service}} (Create{{.Service}}Request) returns (Create{{.Service}}Response);
	
	rpc Echo(EchoRequest) returns (EchoResponse) {
		option (google.api.http) = {
		  post: "/api/echo/{name}"
		  body: "*"
		};
	}
}

message Create{{.Service}}Request {}
message Create{{.Service}}Response {}

message EchoRequest {
  string name = 1;
  int64 age = 2 [(buf.validate.field).int64 = {gt: 17}];
}
message EchoResponse {
  string message = 1;
  int64 created_at = 2;
}
`

func (p *proto) execute() ([]byte, error) {
	tpl, err := template.New("proto").Parse(strings.TrimSpace(protoTemplate))
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if err := tpl.Execute(buf, p); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
