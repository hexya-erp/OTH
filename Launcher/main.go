package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/rayman520/GoRayUtils"
)

var logger rayUtils.Logger

func GetKnownPathMap() map[string]string {
	out := make(map[string]string)
	out["base"] = "github.com/hexya-addons/base"
	out["web"] = "github.com/hexya-addons/web"
	out["webKanban"] = "github.com/hexya-addons/webKanban"
	// out["account"] = "github.com/hexya-addons/account"        			// duplicate heyxa_external_id constraint
	// out["analytic"] = "github.com/hexya-addons/analytic"      			// duplicate heyxa_external_id constraint
	out["decimalPrecision"] = "github.com/hexya-addons/decimalPrecision"
	// out["procurement"] = "github.com/hexya-addons/procurement"			// outdated
	// out["product"] = "github.com/hexya-addons/product"        			// duplicate heyxa_external_id constraint
	// out["sale"] = "github.com/hexya-addons/sale"              			// outdated
	out["saleTeams"] = "github.com/hexya-addons/saleTeams"
	return out
}

func ReadDependencies(content string) []string {
	rx := regexp.MustCompile(`'depends'\s*:\s*\[([^~]*?)],`)
	gr := rx.FindStringSubmatch(content)
	rawDeps := rayUtils.SplitIgnoreBrackets(gr[1], ',')
	for i, dep := range rawDeps {
		dep = strings.TrimSpace(dep)
		rawDeps[i] = rayUtils.RemoveSubstringsFromString(dep, []string{"'", "\""})
	}
	return rawDeps
}

func HandleUnknDep() (f func(dep string) string) {
	unknownDependanceMetOnce := false
	var givenpath string
	f = func(dep string) string {
		if !unknownDependanceMetOnce {
			unknownDependanceMetOnce = true
			logger.LogWarn("Some dependencies found are not recognized.\n" +
				"\t\tplease specify for each one its hexya module path\n" +
				"\t\tLeave empty for default hexya-addons folder\n" +
				"\t\tType 'skip' to ignore thre requirement\n" +
				"\t\tExemple: web -> github.com/hexya-addons/web\n\n")
		}
		logger.LogWarn("please specify hexya module path for:   %s", dep)
		givenpath = ""
		fmt.Scanf("%s", &givenpath)
		if len(givenpath) > 0 {
			return givenpath
		} else {
			return "github.com/hexya-addons/" + dep
		}
	}
	return f
}

func ChangePaths(fullPaths []string) []string {
	logger.LogWarn("you're about to change module paths. Leave empty to keep old path")
	var changed string
	for i, entry := range fullPaths {
		changed = ""
		logger.LogWarn("modification for:   %s", entry)
		fmt.Scanf("%s", &changed)
		if len(changed) > 1 {
			fullPaths[i] = changed
		}
	}
	return fullPaths
}

func ConfirmDeps(depFullPaths []string) []string {
	logger.LogInfo("These module paths will attempt to be loaded:")
	for _, entry := range depFullPaths {
		logger.LogInfo("\t%s", entry)
	}
	var confirmation string
	logger.LogInfo("Confirm?   [Y/N]")
	fmt.Scanf("%s", &confirmation)
	if rayUtils.BooleanTranslate(confirmation) || len(confirmation) == 0 {
		return depFullPaths
	} else {
		logger.LogInfo("Modify?   [Y/N]")
		fmt.Scanf("%s", &confirmation)
		if rayUtils.BooleanTranslate(confirmation) || len(confirmation) == 0 {
			depFullPaths = ConfirmDeps(ChangePaths(depFullPaths))
		} else {
			rayUtils.Exit()
		}
	}
	return depFullPaths
}

func HandleDeps(dependencies []string) []string {
	var depFullPaths []string
	knownPathsMap := GetKnownPathMap()
	HandleUnknownDeps := HandleUnknDep()
	for _, dep := range dependencies {
		dep = rayUtils.CamelCase(dep, false)
		if len(knownPathsMap[dep]) > 0 {
			depFullPaths = append(depFullPaths, knownPathsMap[dep])
		} else {
			newdep := HandleUnknownDeps(dep)
			if newdep != "skip" {
				depFullPaths = append(depFullPaths, newdep)
			}
		}
	}
	return ConfirmDeps(depFullPaths)
}

func GenerateHexya() []string {
	logger.LogInfo("Generating hexya")
	cmdGen := exec.Command("hexya", "generate", "-o")
	HexFilesDir, _ := filepath.Abs("HexyaFiles")
	cmdGen.Dir = HexFilesDir
	stdoutLog, _ := rayUtils.ExecPrintCmd(cmdGen)
	line := stdoutLog.Back()
	if strings.HasPrefix(line.Value.(string), "couldn't load packages due to errors:") {
		out := strings.Split(strings.Replace(line.Value.(string), "couldn't load packages due to errors:", "", -1), ",")
		for i, o := range out {
			out[i] = strings.TrimSpace(o)
		}
		return out
	}
	return []string{}
}

func RemoveHexFiles() {
	dir, _ := filepath.Abs("HexyaFiles ")
	rayUtils.ExecCmd(exec.Command("rm", "-rf", dir))
}

func Retry(deps, outDeps []string, args []string) (bool, []string, []string) {
	for i := len(deps) - 1; i >= 0; i-- {
		if rayUtils.IsContainedInStringArray(deps[i], outDeps) {
			deps = append(deps[:i], deps[i+1:]...)
		}
	}
	for _, dep := range deps {
		fmt.Println(dep)
	}
	WriteModuleDir(deps, args)
	outDeps = GenerateHexya()
	if len(outDeps) > 0 {
		confirmation := ""
		logger.LogWarn("Modify: those dependencies will be ignored: %v", outDeps)
		logger.LogWarn("Modify?   [Y/N]")
		fmt.Scanf("%s", &confirmation)
		retry := rayUtils.BooleanTranslate(confirmation) || len(confirmation) == 0
		return retry, deps, outDeps
	}
	return false, deps, outDeps
}

func StartHex() {
	hexFileDir, _ := filepath.Abs("HexyaFiles")

	cmd := exec.Command("hexya", "generate", "-o", ".")
	cmd.Dir = hexFileDir
	rayUtils.ExecPrintCmd(cmd)

	cmd = exec.Command("dropdb", "hexyaoth")
	cmd.Dir = hexFileDir
	rayUtils.ExecPrintCmd(cmd)

	cmd = exec.Command("createdb", "hexyaoth")
	cmd.Dir = hexFileDir
	rayUtils.ExecPrintCmd(cmd)

	cmd = exec.Command("hexya", "updatedb", "--db-name", "hexyaoth")
	cmd.Dir = hexFileDir
	rayUtils.ExecPrintCmd(cmd)

	cmd = exec.Command("hexya", "server", "-o", "--db-name", "hexyaoth")
	cmd.Dir = hexFileDir
	rayUtils.ExecPrintCmd(cmd)
}

func Finish() {
	StartHex()
	RemoveHexFiles()
	fmt.Println("finish")
}

func main() {
	time.Sleep(10 * time.Millisecond)
	logger = rayUtils.NewLogger("HexTranslate")
	if len(os.Args) < 4 {
		logger.LogFatal("No parameters given.")
		logger.LogFatal("Usage: Launcher <input_path> <output_path> <go_module_path>")
		logger.LogFatal("")
		logger.LogFatal("example: Launcher ./input/product ./output/product github.com/hexya-addons/product")
		rayUtils.Exit()
	}
	content, err := ioutil.ReadFile(os.Args[1] + "/__manifest__.py")
	if err != nil {
		logger.LogFatal("Could not open file. %v", err)
		logger.LogFatal("Specified path is not considered as a translatable module\n")
		rayUtils.Exit()
	}
	RemoveHexFiles()
	dependencies := ReadDependencies(string(content))
	dependencies = HandleDeps(dependencies)
	WriteModuleDir(dependencies, os.Args)
	outDependencies := GenerateHexya()
	if len(outDependencies) > 0 {
		confirmation := ""
		logger.LogWarn("Modify: those dependencies will be ignored: %v", outDependencies)
		logger.LogWarn("Modify?   [Y/N]")
		fmt.Scanf("%s", &confirmation)
		retry := rayUtils.BooleanTranslate(confirmation) || len(confirmation) == 0
		for retry {
			retry, dependencies, outDependencies = Retry(dependencies, outDependencies, os.Args)
		}
		if len(outDependencies) == 0 {
			Finish()
		}
		rayUtils.Exit()
	}
	Finish()
}
