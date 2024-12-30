package filter

type CreateFilterModal struct {
	Name     string `json:"name,omitempty"`
	Type     string `json:"type,omitempty"`
	TypeName string `json:"typeName,omitempty"`
}
