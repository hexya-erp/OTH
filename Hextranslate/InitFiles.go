package Hextranslate

import (
	"container/list"
	"fmt"
	"os/exec"
	"strings"
)

func GetHexyaFileBase() string {
	return "" +
		"package $MODNAME$												\n" +
		"																\n" +
		"import (														\n" +
		"	\"github.com/hexya-addons/web/controllers\"			\n" +
		"	$MODULELIST$												\n" +
		")																\n" +
		"																\n" +
		"const MODULE_NAME string = \"$MODNAME$\"						\n" +
		"																\n" +
		"func init() {													\n" +
		"	server.RegisterModule(&server.Module{						\n" +
		"		Name: MODULE_NAME,										\n" +
		"		PreInit: func() {},										\n" +
		"		PostInit: func() {},									\n" +
		"	})															\n" +
		"	$CONTROLLERASSETS$											\n" +
		"}																\n"
}

func GetHexyaModuleList() string {
	return "" +
		"\"github.com/hexya-erp/hexya/src/server\"\n"
}

func GetTomlFileBase() string {
	return `
Modules = [
    "github.com/hexya-addons/web"$$MODULEPATH$$
]
LogStdout = true
#LogFile = ""
#LogLevel = "info"
#Debug = false
#DataDir = ""

[DB]
#Driver = "postgres"
#Host = "/var/run/postgresql"
#Port = 5432
#User = ""
#Password = ""
#Name = "Hexya"

[Server]
#Languages = ["fr"]
#Interface = ""
#Port = 8080`
}

func GetTomlModuleList() string {
	// thisModulePath, _ := filepath.Abs(genArgs.outputPath)
	// goPath := os.Getenv("GOPATH") + "/src/"
	// thisModulePath = strings.Replace(thisModulePath, goPath, "", 1)
	thisModulePath := "../HexyaFiles/modulelu"
	return "\"" + thisModulePath + "\",\n" +
		"\t\"github.com/hexya-addons/web\""
}

func GetControllerAssetsList() string {
	out := ""
	for key, val := range genArgs.ControllerAssets {
		out = string(append([]byte(out), []byte(fmt.Sprintf("controllers.%s = append(controllers.%s,\n", key, key))...))
		for _, v := range val {
			out = string(append([]byte(out), []byte("\t\t\""+v+"\",\n")...))
		}
		out = string(append([]byte(out), []byte("\t)\n")...))
	}
	return out
}

func CreateHexyaFile() {
	fmt.Println("creating hexya.go")
	fileContent := strings.Replace(GetHexyaFileBase(), "$MODNAME$", genArgs.moduleName, -1)
	fileContent = strings.Replace(fileContent, "$MODULELIST$", GetHexyaModuleList(), -1)
	fileContent = strings.Replace(fileContent, "$CONTROLLERASSETS$", GetControllerAssetsList(), -1)
	fileContentList := list.New()
	fileContentList.PushBack(fileContent)
	path := genArgs.outputPath + "/hexya.go"
	WriteFile(path, fileContentList)
	out, err := exec.Command("goimports", "-w", fmt.Sprintf("%s", path)).CombinedOutput()
	if err != nil {
		fmt.Printf("could not execute command. Error: %v\n%s\n", err, out)
	}
}

func CreateGoModFile() {
	fmt.Println("creating go.mod")
	cmd := exec.Command("go", "mod", "init", genArgs.goModPath)
	cmd.Dir = genArgs.outputPath
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func CreateToml() {
	fmt.Println("creating hexya.toml")
	fileContent := GetTomlFileBase()
	fileContent = strings.Replace(fileContent, "$MODULELIST$", GetTomlModuleList(), -1)
	fileContentList := list.New()
	fileContentList.PushBack(RemoveExtraneousSpaces(fileContent))
	path := genArgs.outputPath + "/../hexya.toml"
	WriteFile(path, fileContentList)
}

func WriteInitFiles() {
	CreateHexyaFile()
	CreateGoModFile()
	CreateToml()
}
