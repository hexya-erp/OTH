package main

import (
	"os"

	"path/filepath"

	"strings"

	"sort"

	"github.com/hexya-erp/hexya/src/tools/po"
	"github.com/rayman520/GoRayUtils"
)

var logger rayUtils.Logger
var globals Globals

type Globals struct {
	inPath      string
	outPath     string
	inFiles     []string
	outFiles    []string
	forceUpdate bool
}

func listPoFiles(path string) []string {
	var out []string
	err := filepath.Walk(path, func(p string, info os.FileInfo, er error) error {
		if er != nil {
			return er
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".po") {
			out = append(out, info.Name())
		} else if info.IsDir() && path != p {
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		logger.LogFatal("Error while walking though %s.  err: %v", path, err)
		os.Exit(0)
	}
	return out
}

func loadPoFile(path string) *po.File {
	file, err := po.Load(path)
	if err != nil {
		logger.LogError("error while trying to load %s. %v", path, err)
	}
	return file
}

type poMsgMap map[string]*po.Message

func mapifizeFile(file *po.File) poMsgMap {
	out := make(poMsgMap)
	for i, msg := range file.Messages {
		out[msg.MsgId] = &file.Messages[i]
	}
	return out
}

func (msgMap poMsgMap) updateMessageMap(inFile *po.File) {
	for _, msg := range inFile.Messages {
		if msgMap[msg.MsgId] != nil {
			if msgMap[msg.MsgId].MsgStr == "" || globals.forceUpdate {
				msgMap[msg.MsgId].MsgStr = msg.MsgStr
			}
		}
	}
}

func (msgMap poMsgMap) deMapifize() []po.Message {
	var out []po.Message
	for _, val := range msgMap {
		out = append(out, *val)
	}
	sort.Slice(out, func(i, j int) bool {
		cmp := strings.Compare(out[i].MsgId, out[j].MsgId)
		return cmp < 0
	})
	return out
}

func updatePoFile(fn string) {
	logger.LogInfo("Starting %s update", fn)
	inPath := filepath.Join(globals.inPath, fn)
	inFile := loadPoFile(inPath)
	outPath := filepath.Join(globals.outPath, fn)
	outFile := loadPoFile(outPath)
	outMessageMap := mapifizeFile(outFile)
	outMessageMap.updateMessageMap(inFile)
	outFile.Messages = outMessageMap.deMapifize()
	outFile.Save(outPath)
	logger.LogInfo("%s update: Done", fn)
}

func updatePoFiles() {
	for _, outFile := range globals.outFiles {
		if rayUtils.IsContainedInStringArray(outFile, globals.inFiles) {
			updatePoFile(outFile)
		} else {
			logger.LogWarn("%s not found in the input path. no update will be done", outFile)
		}
	}
}

func main() {
	logger = rayUtils.NewLogger("OTH-PO")
	if len(os.Args) < 3 {
		logger.LogFatal("Missing args. Usage: ./poUpdate [input_path] [output_path]")
		os.Exit(0)
	}
	globals.outPath = os.Args[2]
	globals.outFiles = listPoFiles(globals.outPath)
	if len(globals.outFiles) == 0 {
		logger.LogFatal("No po files found in output path %s", globals.outPath)
		os.Exit(0)
	}
	if len(os.Args) > 3 && os.Args[3] == "force" {
		globals.forceUpdate = true
	}
	for _, inP := range strings.Split(os.Args[1], ",") {
		ps, err := filepath.Glob(inP)
		if err != nil {
			logger.LogError("Error while parsing input %s. %v", inP, err)
			continue
		}
		for _, inPath := range ps {
			globals.inPath = inPath
			globals.inFiles = listPoFiles(globals.inPath)
			if len(globals.inFiles) == 0 {
				logger.LogError("No po files found in input path %s. this path will be ignored", globals.inPath)
				continue
			}
			logger.LogInfo("Reading input dir %s", globals.inPath)
			updatePoFiles()
		}
	}
	logger.LogInfo("Done")
}
