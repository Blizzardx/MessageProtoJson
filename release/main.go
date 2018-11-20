package main

import (
	"encoding/json"
	"fmt"
	"github.com/Blizzardx/ConfigProtocol/common"
	"github.com/Blizzardx/ConfigProtocol/define"
	"github.com/Blizzardx/ConfigProtocol/tool/exportTool"
)

type ExportTargetConfig struct {
	Name         string            `json:"Name"`
	Lan          define.SupportLan `json:"Lan"`
	OutPutSuffix string            `json:"OutPutSuffix"`
	PackageName  string            `json:"PackageName"`
}
type ExportConfig struct {
	LanDefine    string                `json:"LanDefine"`
	InputDir     string                `json:"InputDir"`
	OutputDir    string                `json:"OutputDir"`
	ExportTarget []*ExportTargetConfig `json:"ExportTarget"`
}

const configName = "configProto.cfg"

func main() {
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
	var exportTargetList []*exportTool.ExportTarget
	for _, exportTarget := range configInfo.ExportTarget {
		exportTargetList = append(exportTargetList, &exportTool.ExportTarget{
			Name:         exportTarget.Name,
			Lan:          exportTarget.Lan,
			OutPutSuffix: exportTarget.OutPutSuffix,
			PackageName:  exportTarget.PackageName,
		})
	}

	err = exportTool.ExportDirectory(configInfo.InputDir, configInfo.OutputDir, exportTargetList)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("succeed")
	}
}
func createSample() {
	configInfo := &ExportConfig{
		InputDir:  "config",
		OutputDir: "output",
		LanDefine: fmt.Sprint("Go:", define.SupportLan_Go, " C#:", define.SupportLan_Csharp, " Java:", define.SupportLan_Java, " Json:", define.SupportLan_Json),
	}
	configInfo.ExportTarget = append(configInfo.ExportTarget, &ExportTargetConfig{Name: "server", Lan: define.SupportLan_Go, OutPutSuffix: ".bytes", PackageName: "config"})
	configInfo.ExportTarget = append(configInfo.ExportTarget, &ExportTargetConfig{Name: "client", Lan: define.SupportLan_Csharp, OutPutSuffix: ".bytes", PackageName: "config"})
	configInfo.ExportTarget = append(configInfo.ExportTarget, &ExportTargetConfig{Name: "httpServer", Lan: define.SupportLan_Java, OutPutSuffix: ".bytes", PackageName: "config"})
	configInfo.ExportTarget = append(configInfo.ExportTarget, &ExportTargetConfig{Name: "cocosClient", Lan: define.SupportLan_Json, OutPutSuffix: ".json", PackageName: "config"})

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
}
