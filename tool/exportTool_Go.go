package tool

import (
	"github.com/Blizzardx/MessageProtoJson/common"
	"github.com/Blizzardx/MessageProtoJson/define"
	"strconv"
)

const codeTemplate_GoEnum = `

type {{.EnumName}} int32

const (
{{range .EnumElemList}}
	{{.EnumName}}_{{.EnumElemName}} {{.EnumName}} = {{.EnumElemValue}}
{{end}}
)


`
const codeTemplate_GoClass = `

type {{.ClassName}} struct {
{{range .Field}}
	{{.Name}} {{.Type}} {{.Decorate}}
{{end}}
}

`
const codeTemplate_Go = `// Generated by gen-tool
// DO NOT EDIT!
package {{.PackageName}}

{{range $_, $v := .EnumList}}
{{$v}}
{{end}}
{{range $_, $v := .ClassList}}
{{$v}}
{{end}}

`

type ExportHandler_GoInfo struct {
	PackageName string
	EnumList    []string
	ClassList   []string
}
type ExportHandler_GoEnumInfo struct {
	EnumName     string
	EnumElemList []*ExportHandler_GoEnumElementInfo
}
type ExportHandler_GoEnumElementInfo struct {
	EnumName      string
	EnumElemName  string
	EnumElemValue string
}
type ExportHandler_GoClassInfo struct {
	ClassName string //首字母 必须大写
	Field     []*ExportHandler_GoClassElementInfo
}
type ExportHandler_GoClassElementInfo struct {
	Name     string //首字母 必须大写
	Type     string
	Decorate string
}
type ExportHandler_Go struct {
}

func (handler *ExportHandler_Go) DoExportProtoFileOnTarget(fileName string, provisionParserInfo *define.MessageProvisionParserInfo, exportPath string) error {
	filePath := exportPath + "/" + fileName + ".go"

	enumList, err := handler.genEnumContent(provisionParserInfo)
	if nil != err {
		return err
	}
	classList, err := handler.genClassContent(provisionParserInfo)
	if nil != err {
		return err
	}
	template := &ExportHandler_GoInfo{
		PackageName: provisionParserInfo.PackageName,
		EnumList:    enumList,
		ClassList:   classList,
	}
	content, err := generateCode(codeTemplate_Go, template, true)
	if nil != err {
		return err
	}

	err = common.WriteFileByName(filePath, []byte(content))
	return nil
}
func (self *ExportHandler_Go) genEnumContent(provisionParserInfo *define.MessageProvisionParserInfo) ([]string, error) {
	var result []string
	for _, enumProvision := range provisionParserInfo.EnumList {
		template := &ExportHandler_GoEnumInfo{EnumName: enumProvision.Name}
		for _, enumFiledInfo := range enumProvision.EnumInfo {
			enumValue := strconv.Itoa(enumFiledInfo.Value)
			template.EnumElemList = append(template.EnumElemList, &ExportHandler_GoEnumElementInfo{EnumElemName: enumFiledInfo.Name, EnumElemValue: enumValue, EnumName: enumProvision.Name})
		}
		content, err := generateCode(codeTemplate_GoEnum, template, false)
		if nil != err {
			return nil, err
		}
		result = append(result, content)
	}
	return result, nil
}
func (self *ExportHandler_Go) genClassContent(provisionParserInfo *define.MessageProvisionParserInfo) ([]string, error) {
	var result []string
	for _, classProvision := range provisionParserInfo.ClassList {

		template := &ExportHandler_GoClassInfo{ClassName: common.FirstLetterToUpper(classProvision.Name)}
		for _, classFiledInfo := range classProvision.FieldInfo {
			realFieldName := common.FirstLetterToUpper(classFiledInfo.Name)
			template.Field = append(template.Field, &ExportHandler_GoClassElementInfo{
				Name:     realFieldName,
				Type:     self.getFieldType(provisionParserInfo, classFiledInfo.Type, classFiledInfo.IsList),
				Decorate: self.getDecorateByName(classFiledInfo.Name),
			})
		}
		content, err := generateCode(codeTemplate_GoClass, template, false)
		if nil != err {
			return nil, err
		}
		result = append(result, content)
	}
	return result, nil

}
func (self *ExportHandler_Go) getDecorateByName(name string) string {
	return "`" + "json:" + "\"" + name + "\"" + "`"
}
func (self *ExportHandler_Go) convertToSelfType(fieldType string) string {
	switch fieldType {
	case "int32":
		return "int32"
	case "int64":
		return "int64"
	case "float32":
		return "float32"
	case "float64":
		return "float64"
	case "bool":
		return "bool"
	case "string":
		return "string"
	default:
		return ""
	}
}
func (self *ExportHandler_Go) isEnum(provisionParserInfo *define.MessageProvisionParserInfo, typeName string) bool {
	for _, enum := range provisionParserInfo.EnumList {
		if enum.Name == typeName {
			return true
		}
	}
	return false
}
func (handler *ExportHandler_Go) getFieldType(provisionParserInfo *define.MessageProvisionParserInfo, fieldType string, isList bool) string {
	typeName := handler.convertToSelfType(fieldType)
	if isList {
		if typeName == "" {
			if handler.isEnum(provisionParserInfo, fieldType) {
				return "[]" + common.FirstLetterToUpper(fieldType)
			}
			return "[]*" + common.FirstLetterToUpper(fieldType)
		}
		return "[]" + typeName
	}
	if typeName == "" {
		if handler.isEnum(provisionParserInfo, fieldType) {
			return common.FirstLetterToUpper(fieldType)
		}
		return "*" + common.FirstLetterToUpper(fieldType)
	}
	return typeName
}
