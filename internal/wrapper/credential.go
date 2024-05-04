package wrapper

type Credential struct {
	Email    string `json:"email,omitempty"`
	IsGlobal bool   `json:"global,omitempty"`
	Key      string `json:"key"`
}
