package define

import (
	"bytes"
	"encoding/json"
	"github.com/Blizzardx/MessageProtoJson/common"
	"strconv"
)

type MessageProvisionInfo struct {
	PackageName string                       `json:"pkg"`
	ImportList  []string                     `json:"import"`
	EnumList    []*MessageProvisionEnumInfo  `json:"enums"`
	ClassList   []*MessageProvisionClassInfo `json:"classes"`
}
type MessageProvisionEnumInfo struct {
	Name     string   `json:"name"`
	EnumInfo []string `json:"enums"` // yellow=1 red=2
}
type MessageProvisionClassInfo struct {
	Name      string   `json:"name"`
	FieldInfo []string `json:"fields"` //int32:id:isList name:string:repeated
}

type MessageProvisionParserInfo struct {
	PackageName string
	ImportList  []string
	Name        string
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

func GenSampleFile() {
	sampleFile := &MessageProvisionInfo{PackageName: "message"}
	sampleFile.ImportList = append(sampleFile.ImportList, "sampleFile")
	enumSample1 := &MessageProvisionEnumInfo{Name: "Color"}
	enumSample1.EnumInfo = append(enumSample1.EnumInfo, "yellow=1")
	enumSample1.EnumInfo = append(enumSample1.EnumInfo, "blue=2")
	enumSample1.EnumInfo = append(enumSample1.EnumInfo, "red=3")
	enumSample1.EnumInfo = append(enumSample1.EnumInfo, "green=4")
	sampleFile.EnumList = append(sampleFile.EnumList, enumSample1)

	enumSample := &MessageProvisionEnumInfo{Name: "Quality"}
	enumSample.EnumInfo = append(enumSample.EnumInfo, "yellow=1")
	enumSample.EnumInfo = append(enumSample.EnumInfo, "blue=2")
	enumSample.EnumInfo = append(enumSample.EnumInfo, "red=3")
	enumSample.EnumInfo = append(enumSample.EnumInfo, "green=4")
	sampleFile.EnumList = append(sampleFile.EnumList, enumSample)

	for i := 0; i < 10; i++ {
		classSample := &MessageProvisionClassInfo{Name: "class" + strconv.Itoa(i)}
		classSample.FieldInfo = append(classSample.FieldInfo, "int32"+":field1")
		classSample.FieldInfo = append(classSample.FieldInfo, "int64"+":field2"+":repeated")
		classSample.FieldInfo = append(classSample.FieldInfo, "float32"+":field3")
		classSample.FieldInfo = append(classSample.FieldInfo, "float64"+":field4"+":repeated")
		classSample.FieldInfo = append(classSample.FieldInfo, "bool"+":field5")
		classSample.FieldInfo = append(classSample.FieldInfo, "string"+":field6"+":repeated")
		if i > 0 {
			classSample.FieldInfo = append(classSample.FieldInfo, "class"+strconv.Itoa(i-1)+":field7:repeated")
		}
		sampleFile.ClassList = append(sampleFile.ClassList, classSample)
	}
	fileName := "sampleFile.json"
	content, _ := json.Marshal(sampleFile)
	var out bytes.Buffer
	json.Indent(&out, content, "", "\t")
	common.WriteFileByName(fileName, out.Bytes())
}
