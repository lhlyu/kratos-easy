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

option go_package = "{{.GoPackage}}";
option java_multiple_files = true;
option java_outer_classname = "{{.Service}}ProtoV1";
option java_package = "{{.JavaPackage}}";

enum ErrorReason {
  ERROR_REASON_UNSPECIFIED = 0;
  ERROR_REASON_NOT_FOUND = 1;
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
