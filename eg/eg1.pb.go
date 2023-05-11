package eg

type LoginRequest1 struct {
	// Username
	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	// Password  @ban
	Password  string  `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	Intcase   int     `protobuf:"bytes,2,opt,name=intcase,proto3" json:"intcase,omitempty"`
	Floatcase float32 `protobuf:"bytes,2,opt,name=floattcase,proto3" json:"floatcase,omitempty"`
}

type LoginRequest2 struct {
	// Username
	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	// Password
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"` // @ban
}
