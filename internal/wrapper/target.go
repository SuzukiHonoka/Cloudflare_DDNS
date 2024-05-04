package wrapper

import "fmt"

type Target struct {
	Domain  string `json:"domain"`
	SubName string `json:"sub"`
}

func (t *Target) String() string {
	return fmt.Sprintf("%s.%s", t.SubName, t.Domain)
}
