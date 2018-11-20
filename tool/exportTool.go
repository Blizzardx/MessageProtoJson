package tool

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Blizzardx/MessageProtoJson/common"
	"github.com/Blizzardx/MessageProtoJson/define"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"strconv"
	"strings"
	"text/template"
)

type ExportTarget struct {
	Lan        define.SupportLan
	ExportPath string
}
type ExportHandler interface {
	DoExportProtoFileOnTarget(fileName string, provisionParserInfo *define.MessageProvisionParserInfo, exportPath string) error
}

var exportHandlerMap = map[define.SupportLan]ExportHandler{}

func init() {
	exportHandlerMap[define.SupportLan_Go] = &ExportHandler_Go{}
	exportHandlerMap[define.SupportLan_Ts] = &ExportHandler_Ts{}
}

//根据输入的路径和目标文件的尾缀 生成对应平台的代码
func ExportProtoFile(provisionFilePath string, provisionFileSuffix string, exPortTarget []*ExportTarget) error {
	allFile := loadAllFile(provisionFilePath)
	if len(allFile) == 0 {
		return errors.New("file not found" + provisionFilePath)
	}
	var targetFile []*FileDetail
	for _, file := range allFile {
		if strings.HasSuffix(file.FileName, provisionFileSuffix) {
			targetFile = append(targetFile, file)
		}
	}
	if len(targetFile) == 0 {
		return errors.New("not found file by suffix" + provisionFileSuffix + " at path " + provisionFileSuffix)
	}
	errorMsg := ""
	for _, file := range targetFile {
		err := exportProtoFileElement(file.FilePath, exPortTarget)
		if nil != err {
			errorMsg += err.Error() + "\n"
		}
	}
	if errorMsg != "" {
		return errors.New(errorMsg)
	}
	return nil
}

func exportProtoFileElement(filePath string, exportTargets []*ExportTarget) error {
	//read file
	fileContent, err := common.LoadFileByName(filePath)
	if nil != err {
		return err
	}
	provisionInfo := &define.MessageProvisionInfo{}
	err = json.Unmarshal(fileContent, provisionInfo)
	if nil != err {
		return errors.New("error on parser file to json " + filePath + " error: " + err.Error())
	}
	fixedProvisionInfo, err := parserProvisionInfo(provisionInfo)
	if nil != err {
		return errors.New("error on parser provision file " + filePath + " error: " + err.Error())
	}
	fileNameWithoutSuffix, _ := common.ParserFileNameByPath(filePath)
	err = doExportProtoFileElement(fileNameWithoutSuffix, fixedProvisionInfo, exportTargets)
	return err
}
func parserClassInfo(fieldInfo []string) ([]*define.MessageProvisionParserClassFieldInfo, error) {
	var result []*define.MessageProvisionParserClassFieldInfo
	for _, field := range fieldInfo {
		fieldInfoList := strings.Split(field, ":")
		if len(fieldInfoList) < 1 {
			return nil, errors.New("解析类成员发生错误，必须用:分割，类型：名字：repeated " + field)
		}
		parserInfo := &define.MessageProvisionParserClassFieldInfo{
			Type:   fieldInfoList[0],
			Name:   fieldInfoList[1],
			IsList: false,
		}
		if len(fieldInfoList) > 2 && fieldInfoList[2] == "repeated" {
			parserInfo.IsList = true
		}
		result = append(result, parserInfo)
	}
	return result, nil
}
func parserEnumInfo(enumInfo []string) ([]*define.MessageProvisionParserEnumFieldInfo, error) {
	var result []*define.MessageProvisionParserEnumFieldInfo
	for _, field := range enumInfo {
		fieldInfoList := strings.Split(field, "=")
		if len(fieldInfoList) != 2 {
			return nil, errors.New("解析枚举成员发生错误，必须用=分割，名字=value// yellow=1 red=2 blue=3 " + field)
		}
		parserInfo := &define.MessageProvisionParserEnumFieldInfo{Name: fieldInfoList[0]}
		value, err := strconv.Atoi(fieldInfoList[1])
		if nil != err {
			return nil, errors.New("解析枚举成员发生错误，枚举值必须是int整形" + field)
		}
		parserInfo.Value = value
		result = append(result, parserInfo)
	}
	return result, nil
}
func parserProvisionInfo(provisionInfo *define.MessageProvisionInfo) (*define.MessageProvisionParserInfo, error) {
	result := &define.MessageProvisionParserInfo{PackageName: provisionInfo.PackageName}
	for _, enum := range provisionInfo.EnumList {
		parserEnum := &define.MessageProvisionParserEnumInfo{Name: enum.Name}
		enumInfoList, err := parserEnumInfo(enum.EnumInfo)
		if nil != err {
			return nil, errors.New(err.Error() + " " + enum.Name + " " + provisionInfo.PackageName)
		}
		parserEnum.EnumInfo = enumInfoList
		result.EnumList = append(result.EnumList, parserEnum)
	}
	for _, class := range provisionInfo.ClassList {
		parserClass := &define.MessageProvisionParserClassInfo{Name: class.Name}
		fieldInfoList, err := parserClassInfo(class.FieldInfo)
		if nil != err {
			return nil, errors.New(err.Error() + " " + class.Name + " " + provisionInfo.PackageName)
		}
		parserClass.FieldInfo = fieldInfoList
		result.ClassList = append(result.ClassList, parserClass)
	}
	return result, nil
}
func doExportProtoFileElement(fileName string, parserInfo *define.MessageProvisionParserInfo, exportTargets []*ExportTarget) error {
	for _, target := range exportTargets {
		realExportPath := target.ExportPath + "/" + parserInfo.PackageName
		common.EnsureFolder(realExportPath)
		doExportProtoFileOnTarget(fileName, parserInfo, realExportPath, target.Lan)
	}

	return nil
}
func doExportProtoFileOnTarget(fileName string, provisionParserInfo *define.MessageProvisionParserInfo, exportPath string, lan define.SupportLan) error {
	if handler, ok := exportHandlerMap[lan]; ok {
		err := handler.DoExportProtoFileOnTarget(fileName, provisionParserInfo, exportPath)
		if nil != err {
			return err
		}
	}
	return errors.New(fmt.Sprint("not support lan ", lan))
}

type FileDetail struct {
	FileName string
	FilePath string
	FileSize int64
}

func loadAllFile(dirPath string) []*FileDetail {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var filePool []*FileDetail
	for _, file := range files {
		if file.IsDir() {
			filePool = append(filePool, loadAllFile(dirPath+"/"+file.Name())...)
		} else {
			filePool = append(filePool, &FileDetail{FileName: file.Name(), FilePath: dirPath + "/" + file.Name(), FileSize: file.Size()})
		}
	}
	return filePool
}

//根据模板生成代码
func generateCode(templateStr string, model interface{}, needFormat bool) (string, error) {

	var err error

	var bf bytes.Buffer

	tpl, err := template.New("Template").Parse(templateStr)
	if err != nil {
		return "", err
	}

	err = tpl.Execute(&bf, model)
	if err != nil {
		return "", err
	}

	if needFormat {
		if err = formatCode(&bf); err != nil {
			fmt.Println("format golang code err", err)
		}
	}

	return string(bf.Bytes()), nil
}

//格式化go文件
func formatCode(bf *bytes.Buffer) error {

	fset := token.NewFileSet()

	ast, err := parser.ParseFile(fset, "", bf, parser.ParseComments)
	if err != nil {
		return err
	}

	bf.Reset()

	err = (&printer.Config{Mode: printer.TabIndent | printer.UseSpaces, Tabwidth: 8}).Fprint(bf, fset, ast)
	if err != nil {
		return err
	}

	return nil
}
