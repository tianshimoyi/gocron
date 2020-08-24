package schema

type NodeRequest struct {
	Name   string `json:"name" description:"node name" validate:"required"`
	Alias  string `json:"alias" description:"alias name" optional:"true"`
	Port   int    `json:"port" description:"node port" validate:"required"`
	Remark string `json:"remark" description:"node remark" optional:"true"`
}
