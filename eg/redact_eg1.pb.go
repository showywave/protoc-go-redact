package eg

import "fmt"

func (x *LoginRequest1) Redact () string {
	return fmt.Sprintf("username:%s password:%s intcase:%d floattcase:%v", x.Username, "******", x.Intcase, x.Floatcase)
}

func (x *LoginRequest2) Redact () string {
	return fmt.Sprintf("username:%s password:%s", x.Username, "******")
}

