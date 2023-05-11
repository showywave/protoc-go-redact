# protoc-go-redact

## Install
- `go install github.com/showywave/protoc-go-redact@latest`

## Usage
Add special comments for sensitive fields in the structure to be processed in the `*.pb.go` file
```go
// @ban
```
Then execute `protoc-go-redact` on the command line

```console
$ protoc-go-redact -input=*.pb.go
```

## Example
eg/eg1.pb.go
```  go
package eg

type LoginRequest1 struct {
	// Username
	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	// Password  @ban
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
}

type LoginRequest2 struct {
	// Username
	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	// Password
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"` // @ban
}
```

```console
$ protoc-go-redact -input=eg/eg1.pb.go
```

will generate `eg/deract_eg1.pb.go` file:
```go
package eg

import "fmt"

func (x *LoginRequest1) Redact () string {
	return fmt.Sprintf("username:%s password:%s", x.Username, "******")
}

func (x *LoginRequest2) Redact () string {
	return fmt.Sprintf("username:%s password:%s", x.Username, "******")
}
```