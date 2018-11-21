package main

import (
	"encoding/json"
	"fmt"
	"github.com/Blizzardx/MessageProtoJson/common"
	"github.com/Blizzardx/MessageProtoJson/define"
	"github.com/Blizzardx/MessageProtoJson/tool"
)

type ExportTargetConfig struct {
	Lan define.SupportLan `json:"Lan"`
}
type ExportConfig struct {
	LanDefine        string                `json:"LanDefine"`
	InputDir         string                `json:"InputDir"`
	OutputDir        string                `json:"OutputDir"`
	TargetFileSuffix string                `json:"TargetFileSuffix"`
	ExportTarget     []*ExportTargetConfig `json:"ExportTarget"`
}

const configName = "messageProtoJson.cfg"

func main() {
	Exec()
}
func Exec() {
	content, err := common.LoadFileByName(configName)
	if err != nil {
		//
		createSample()
		return
	}

	configInfo := &ExportConfig{}
	err = json.Unmarshal(content, configInfo)
	if err != nil {
		fmt.Println(err)
		return
	}
	var exportTargetList []*tool.ExportTarget
	for _, exportTarget := range configInfo.ExportTarget {
		exportTargetList = append(exportTargetList, &tool.ExportTarget{
			ExportPath: configInfo.OutputDir + "/" + getPathByLan(exportTarget.Lan),
			Lan:        exportTarget.Lan,
		})
	}
	if configInfo.InputDir == "" {
		configInfo.InputDir = common.GetCurrentPath()
	}

	err = tool.ExportProtoFile(configInfo.InputDir, configInfo.TargetFileSuffix, exportTargetList)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("succeed")
	}
}
func getPathByLan(lan define.SupportLan) string {
	switch lan {
	case define.SupportLan_Go:
		return "Go"
	case define.SupportLan_Csharp:
		return "C#"
	case define.SupportLan_Java:
		return "Java"
	case define.SupportLan_Ts:
		return "Ts"
	}
	return "unknown lan"
}
func createSample() {
	configInfo := &ExportConfig{
		InputDir:         "msg",
		OutputDir:        "output",
		TargetFileSuffix: "json",
		LanDefine:        fmt.Sprint("Go:", define.SupportLan_Go, " C#:", define.SupportLan_Csharp, " Java:", define.SupportLan_Java, " Ts:", define.SupportLan_Ts),
	}
	configInfo.ExportTarget = append(configInfo.ExportTarget, &ExportTargetConfig{Lan: define.SupportLan_Go})
	configInfo.ExportTarget = append(configInfo.ExportTarget, &ExportTargetConfig{Lan: define.SupportLan_Csharp})
	configInfo.ExportTarget = append(configInfo.ExportTarget, &ExportTargetConfig{Lan: define.SupportLan_Java})
	configInfo.ExportTarget = append(configInfo.ExportTarget, &ExportTargetConfig{Lan: define.SupportLan_Ts})

	configContent, err := json.Marshal(configInfo)
	if nil != err {
		fmt.Println("error on create sample ", err)
		return
	}
	err = common.WriteFileByName(configName, configContent)
	if err != nil {

		fmt.Println("error on create sample ", err)
		return
	}
	fmt.Println("finished config first ", configName)
	define.GenSampleFile()
}
