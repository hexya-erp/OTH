package main

import (
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rayman520/GoRayUtils"
)

func EncodeDeps(deps []string) string {
	out := ""
	for _, dep := range deps {
		out += ",\n\t\"" + dep + "\""
	}
	return out
}

func UpdateToml(deps []string) {
	content, err := rayUtils.ReadFileContent("HexyaFilesReference/hexya.toml")
	if err != nil {
		rayUtils.ExitFatal("Could not open file. %v", err)
	}
	//thisModulePath, _ := filepath.Abs(genArgs.outputPath)
	//goPath := os.Getenv("GOPATH") + "/src/"
	//thisModulePath = strings.Replace(thisModulePath, goPath, "", 1)
	thisModulePath := "github.com/hexya-erp/OTH/Launcher/HexyaFiles/module"
	deps = append(deps, thisModulePath)
	content = strings.Replace(content, "$$MODULEPATH$$", EncodeDeps(deps), -1)
	rayUtils.WriteFileContent("HexyaFiles/hexya.toml", content)
}

func UpdateModule() {
	exec.Command("cp", "-r", "./HexyaFilesReference/module", "./HexyaFiles").CombinedOutput()
}

func WriteInputPathFile(inPath string) {
	path, _ := filepath.Abs(inPath)
	rayUtils.WriteFileContent("HexyaFiles/inputPath", path)
}

func WriteOutputPathFile(outPath string) {
	path, _ := filepath.Abs(outPath)
	rayUtils.WriteFileContent("HexyaFiles/outputPath", path)
}

func WriteGoModulePathFile(goModPath string) {
	rayUtils.WriteFileContent("HexyaFiles/goModPath", goModPath)
}

func WriteModuleDir(deps []string, args []string) {
	UpdateToml(deps)
	WriteInputPathFile(args[1])
	WriteOutputPathFile(args[2])
	WriteGoModulePathFile(args[3])
	UpdateModule()
}
