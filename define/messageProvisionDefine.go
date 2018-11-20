package define

type MessageProvisionInfo struct {
	PackageName string
	EnumList    []*MessageProvisionEnumInfo  `json:"enums"`
	ClassList   []*MessageProvisionClassInfo `json:"classes"`
}
type MessageProvisionEnumInfo struct {
	Name     string
	EnumInfo []string `json:"enums"` // yellow=1 red=2
}
type MessageProvisionClassInfo struct {
	Name      string
	FieldInfo []string `json:"fields"` //int32:id:isList name:string:repeated
}

type MessageProvisionParserInfo struct {
	PackageName string
	EnumList    []*MessageProvisionParserEnumInfo
	ClassList   []*MessageProvisionParserClassInfo
}
type MessageProvisionParserEnumInfo struct {
	Name     string
	EnumInfo []*MessageProvisionParserEnumFieldInfo
}
type MessageProvisionParserEnumFieldInfo struct {
	Name  string
	Value int
}
type MessageProvisionParserClassInfo struct {
	Name      string
	FieldInfo []*MessageProvisionParserClassFieldInfo
}
type MessageProvisionParserClassFieldInfo struct {
	Type   string
	Name   string
	IsList bool
}
