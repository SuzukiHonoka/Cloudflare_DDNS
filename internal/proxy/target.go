package proxy

import "fmt"

type Target struct {
	Domain  string `json:"domain"`
	SubName string `json:"sub"`
}

func (t *Target) Name() string {
	return fmt.Sprintf("%s.%s", t.SubName, t.Domain)
}
