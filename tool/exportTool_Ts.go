package tool

import (
	"github.com/Blizzardx/MessageProtoJson/common"
	"github.com/Blizzardx/MessageProtoJson/define"
	"strconv"
)

const codeTemplate_TsEnum = `
export enum {{.EnumName}} {
{{range .EnumElemList}}
	 {{.EnumElemName}} = {{.EnumElemValue}} ,{{end}}
}
`
const codeTemplate_TsClass = `
export class {{.ClassName}} {
{{range .Field}}	
	public {{.Name}} : {{.Type}} ;{{end}}

	public GetName():string{
		return "{{.ClassName}}";
	}
}
`
const codeTemplate_Ts = `// Generated by gen-tool
// DO NOT EDIT!
{{range .ImportList}}	
import { {{.FieldName}} } from "./{{.ImportFile}}";{{end}}

{{range $_, $v := .EnumList}}{{$v}}{{end}}
{{range $_, $v := .ClassList}}{{$v}}{{end}}

`

type ExportHandler_TsInfo struct {
	EnumList   []string
	ClassList  []string
	ImportList []*ExportHandler_TsImportInfo
}
type ExportHandler_TsEnumInfo struct {
	EnumName     string
	EnumElemList []*ExportHandler_TsEnumElementInfo
}
type ExportHandler_TsEnumElementInfo struct {
	EnumElemName  string
	EnumElemValue string
}
type ExportHandler_TsClassInfo struct {
	ClassName string
	Field     []*ExportHandler_TsClassElementInfo
}
type ExportHandler_TsClassElementInfo struct {
	Name string
	Type string
}
type ExportHandler_TsImportInfo struct {
	FieldName  string
	ImportFile string
}
type ExportHandler_Ts struct {
}

func (handler *ExportHandler_Ts) DoExportProtoFileOnTarget(fileName string, provisionParserInfo *define.MessageProvisionParserInfo, exportPath string, importfile []*define.MessageProvisionParserInfo) error {
	filePath := exportPath + "/" + fileName + ".ts"

	enumList, err := handler.genEnumContent(provisionParserInfo)
	if nil != err {
		return err
	}
	classList, err := handler.genClassContent(provisionParserInfo)
	if nil != err {
		return err
	}
	template := &ExportHandler_TsInfo{
		EnumList:   enumList,
		ClassList:  classList,
		ImportList: handler.parserImportInfo(provisionParserInfo, importfile),
	}
	content, err := generateCode(codeTemplate_Ts, template, false)
	if nil != err {
		return err
	}

	err = common.WriteFileByName(filePath, []byte(content))
	return nil
}

func (self *ExportHandler_Ts) genEnumContent(provisionParserInfo *define.MessageProvisionParserInfo) ([]string, error) {
	var result []string
	for _, enumProvision := range provisionParserInfo.EnumList {
		template := &ExportHandler_TsEnumInfo{EnumName: enumProvision.Name}
		for _, enumFiledInfo := range enumProvision.EnumInfo {
			enumValue := strconv.Itoa(enumFiledInfo.Value)
			template.EnumElemList = append(template.EnumElemList, &ExportHandler_TsEnumElementInfo{EnumElemName: enumFiledInfo.Name, EnumElemValue: enumValue})
		}
		content, err := generateCode(codeTemplate_TsEnum, template, false)
		if nil != err {
			return nil, err
		}
		result = append(result, content)
	}
	return result, nil
}
func (self *ExportHandler_Ts) genClassContent(provisionParserInfo *define.MessageProvisionParserInfo) ([]string, error) {
	var result []string
	for _, classProvision := range provisionParserInfo.ClassList {
		template := &ExportHandler_TsClassInfo{ClassName: classProvision.Name}
		for _, classFiledInfo := range classProvision.FieldInfo {
			template.Field = append(template.Field, &ExportHandler_TsClassElementInfo{
				Name: classFiledInfo.Name,
				Type: self.getFieldType(classFiledInfo.Type, classFiledInfo.IsList)})
		}
		content, err := generateCode(codeTemplate_TsClass, template, false)
		if nil != err {
			return nil, err
		}
		result = append(result, content)
	}
	return result, nil

}
func (self *ExportHandler_Ts) convertToSelfType(fieldType string) string {
	switch fieldType {
	case "int32":
		return "number"
	case "int64":
		return "number"
	case "float32":
		return "number"
	case "float64":
		return "number"
	case "bool":
		return "boolean"
	case "string":
		return "string"
	default:
		return ""
	}
}
func (handler *ExportHandler_Ts) getFieldType(fieldType string, isList bool) string {
	typeName := handler.convertToSelfType(fieldType)
	if typeName == "" {
		typeName = fieldType
	}
	if !isList {
		return typeName
	}
	return "Array<" + typeName + ">"
}
func (handler *ExportHandler_Ts) parserImportInfo(target *define.MessageProvisionParserInfo, importfile []*define.MessageProvisionParserInfo) []*ExportHandler_TsImportInfo {
	tmpMap := map[string]string{}

	for _, class := range target.ClassList {
		for _, filedInfo := range class.FieldInfo {
			if _, ok := tmpMap[filedInfo.Type]; ok {
				continue
			}
			importFileName := handler.getFieldInfoFromImportFile(filedInfo.Type, importfile)
			if importFileName != "" {
				tmpMap[filedInfo.Type] = importFileName
			}
		}
	}
	var res []*ExportHandler_TsImportInfo

	for field, importFile := range tmpMap {
		res = append(res, &ExportHandler_TsImportInfo{FieldName: field, ImportFile: importFile})
	}
	return res
}
func (handler *ExportHandler_Ts) getFieldInfoFromImportFile(fileType string, importfile []*define.MessageProvisionParserInfo) string {
	for _, file := range importfile {
		for _, class := range file.ClassList {
			if class.Name == fileType {
				return file.Name
			}
		}
		for _, enum := range file.EnumList {
			if enum.Name == fileType {
				return file.Name
			}
		}
	}
	return ""
}
