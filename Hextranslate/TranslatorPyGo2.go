package Hextranslate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/golang-collections/collections/stack"
	"github.com/hexya-erp/hexya/src/models"

	"github.com/davecgh/go-spew/spew"
)

type modifiedField struct {
	str   string
	label string
}

type hexyaModel struct {
	name           string
	sqlConstraints string
	fields         []*hexyaField
	modFields      []*modifiedField
	funcs          []*hexyaFunc
}

type hexyaField struct {
	name    string
	dump    string
	invalid bool
	comment string
}

type hexyaFunc struct {
	name     string
	override bool
	dump     string
}

var (
	rawText            [][]byte
	conf               *spew.ConfigState
	funcMap            map[string]func(map[string]interface{}) interface{}
	knownImportedWords map[string]string
	knownComputedFuncs []string
	aliasMap           map[string]string
	curModelName       string
	isInInit           bool
	isInAddFields      bool
	fieldFuncMap       map[string][]func(map[string]interface{}) []byte
	odooModelVarNames  []string
	uniqueNamesBank    map[string][]string
	hexyaModels        []*hexyaModel
)

/* utils */

func InitFuncMap() {
	if funcMap == nil {
		funcMap = make(map[string]func(map[string]interface{}) interface{})
		funcMap["ClassDef"] = translateClassDef
		funcMap["ImportFrom"] = translateImportFrom
		funcMap["Attribute"] = translateAttribute
		funcMap["Name"] = translateName
		funcMap["alias"] = translateAlias
		funcMap["Assign"] = translateAssign
		funcMap["FunctionDef"] = translateFunctionDef
		funcMap["Call"] = translateCall
		funcMap["BinOp"] = translateBinOp
		funcMap["BinOpMod"] = translateBinOpMod
		funcMap["Subscript"] = translateSubscript
		funcMap["Index"] = translateIndex
		funcMap["Str"] = translateStr
		funcMap["Num"] = translateNum
	}
}

func getFunctionName(model, methName string) string {
	return fmt.Sprintf("%s_%s", strings.ToLower(model[:1])+model[1:], methName)
}

func Comment(text [][]byte) [][]byte {
	for i, line := range text {
		text[i] = append([]byte("//"), line...)
	}
	return text
}

func IsOdooVar(node map[string]interface{}) bool {
	if odooModelVarNames == nil {
		odooModelVarNames = []string{
			`_auto`, `_register`, `_abstract`, `_transient`, `_name`, `_description`, `_custom`, `_inherit`,
			`_inherits`, `_constraints`, `_table`, `_sequence`, `_sql_constraints`, `_rec_name`, `_order`,
			`_parent_name`, `_parent_store`, `_parent_order`, `_date_name`, `_fold_name`, `_needaction`, `_translate`,
			`_depends`, `_transient_check_count`, `_transient_max_count`, `_transient_max_hours`, `_log_access`,
		}
	}
	target := Node(Sl(node["targets"])[0])
	if Str(target["type"]) == "Name" && !IsNotInStringArray(Str(target["id"]), odooModelVarNames) {
		return true
	}
	return false
}

func GetRawText(l interface{}) [][]byte {
	loc := l.(map[string]interface{})
	startNode := loc["start"].(map[string]interface{})
	endNode := loc["end"].(map[string]interface{})
	return rawText[int(startNode["line"].(float64))-1 : int(endNode["line"].(float64))]
}

func GetRawLinePart(l interface{}) []byte {
	var out []byte
	loc := Node(l)
	startNode := Node(loc["start"])
	endNode := Node(loc["end"])
	lineCur := 0
	lineEnd := int(endNode["line"].(float64))
	colCur := int(startNode["column"].(float64)) + 1
	colEnd := int(endNode["column"].(float64))
	for lineCur = int(startNode["line"].(float64)); lineCur <= lineEnd; lineCur++ {
		if lineCur == lineEnd {
			out = append(out, bytes.TrimSpace(rawText[lineCur-1][colCur-1:colEnd])...)
		} else {
			out = append(out, bytes.TrimSpace(rawText[lineCur-1][colCur-1:])...)
			colCur = 1
		}
	}
	return out

}

func ConvertFieldType(input string) string {
	switch input {
	case "Char", "Selection", "Float", "Integer", "Boolean", "Binary", "Date", "Text", "Monetary":
		return input
	case "Many2one":
		return "Many2One"
	case "One2many":
		return "One2Many"
	case "Many2many":
		return "Many2Many"
	case "Datetime":
		return "DateTime"
	case "Html":
		return "HTML"
	}
	return string(append([]byte(input), []byte("_UNKNOWN_FIELD_PLS_FIX")...))
}

func MethodIsExtension(name string) bool {
	model, ok := models.Registry.Get(genArgs.RecordSet.ModelName())
	if !ok {
		return false
	}
	_, ok = model.Methods().Get(name)
	if !ok {
		return false
	}
	return true
}

func GetModelName(b interface{}) string {
	body := b.([]interface{})
	var out string
	for _, p := range body {
		part := p.(map[string]interface{})
		if part["type"].(string) != "Assign" {
			continue
		}
		values := part["value"].(map[string]interface{})
		var val string
		if v, ok := values["s"].(string); !ok {
			val = ""
		} else {
			val = v
		}
		val = CleanName(val)
		for _, t := range part["targets"].([]interface{}) {
			target := t.(map[string]interface{})
			if target["id"].(string) == "_inherit" {
				out = val
			} else if target["id"].(string) == "_name" {
				return val
			}
		}
	}
	return out
}

func GetClassParam(classNode map[string]interface{}) string {
	var out string
	for _, b := range classNode["bases"].([]interface{}) {
		base := b.(map[string]interface{})
		typ := base["type"].(string)
		if _, ok := funcMap[typ]; !ok {
			conf.Dump(base)
			continue
		}
		str := GetFullNameImported(funcMap[typ](base).(string))
		if out == "" {
			out = str
		} else {
			out += "," + str
		}
	}
	return out
}

func GetFullNameImported(str string) string {
	spl := strings.Split(str, ".")
	if knownImportedWords[spl[0]] != "" {
		spl[0] = knownImportedWords[spl[0]]
	}
	return strings.Join(spl, ".")
}

func GetAlias(str string) string {
	if val, ok := aliasMap[str]; ok {
		return val
	}
	return str
}

func Node(n interface{}) map[string]interface{} {
	return n.(map[string]interface{})
}

func Sl(n interface{}) []interface{} {
	return n.([]interface{})
}

func Str(n interface{}) string {
	if IsNil(n) {
		return ""
	}
	t := reflect.TypeOf(n)
	switch t.Kind().String() {
	case "string":
		return n.(string)
	case "slice":
		return string(n.([]byte))
	}
	return n.(string)
}

func BSl(n interface{}) []byte {
	return []byte(Str(n))
}

func SnakeToCamel(str string, firstUpper bool) string {
	splited := strings.Split(strings.TrimSpace(strings.Replace(str, "_", " ", -1)), " ")
	for i, s := range splited {
		if firstUpper {
			splited[i] = strings.Title(s)
		}
		firstUpper = true
	}
	return strings.Join(splited, "")
}

func MakeUnique(str string, bank string) string {
	if uniqueNamesBank == nil {
		uniqueNamesBank = make(map[string][]string)
	}
	out := str
	i := 0
	for !IsNotInStringArray(out, uniqueNamesBank[bank]) {
		i++
		out = str + strconv.Itoa(i)
	}
	uniqueNamesBank[bank] = append(uniqueNamesBank[bank], out)
	return out
}

func PrintNodeOutput(node map[string]interface{}) []byte {
	s, _ := json.Marshal(node)
	return s
}

func IsNil(c interface{}) bool {
	if c == nil || (reflect.ValueOf(c).Kind() == reflect.Ptr && reflect.ValueOf(c).IsNil()) {
		return true
	}
	return false
}

func SplitSubN(s string, n int) []string {
	runes := []rune(s)
	for i, r := range runes {
		if r == '\n' {
			return append([]string{string(runes[:i])}, SplitSubN(string(runes[i+1:]), n)...)
		}
		if i >= int(float64(n)*1.2) {
			for j, ru := range runes {
				if j >= int(float64(n)*0.9) {
					if unicode.IsSpace(ru) {
						return append([]string{string(runes[:j])}, SplitSubN(string(runes[j+1:]), n)...)
					}
				}
			}
		}
	}
	return []string{string(runes)}
}

func ReplaceRegex(str string, regex string, repl string) string {
	rx := regexp.MustCompile(regex)
	return rx.ReplaceAllString(str, repl)
}

/* funcMap */

func translateGeneric(node map[string]interface{}) interface{} {
	typ := Str(node["type"])
	if _, ok := funcMap[typ]; !ok {
		s, _ := json.Marshal(node)
		fmt.Printf("Error: type %s not handled. %s\n %s\n\n", typ, conf.Sdump(node["loc"]), s)
		return nil
	}
	return funcMap[typ](node)
}

func translateNum(node map[string]interface{}) interface{} {
	fl := node["n"].(float64)
	str := strconv.FormatFloat(fl, 'f', -1, 64)
	return str
}

func translateStr(node map[string]interface{}) interface{} {
	return `"` + Str(node["s"]) + `"`
}

func translateIndex(node map[string]interface{}) interface{} {
	value := Node(node["value"])
	return Str(translateGeneric(value))
}

func translateSubscript(node map[string]interface{}) interface{} {
	slice := Str(translateGeneric(Node(node["slice"])))
	val := Str(translateGeneric(Node(node["value"])))
	return val + "[" + slice + "]"
}

func translateCall(node map[string]interface{}) interface{} {
	f := GetFullNameImported(Str(translateGeneric(Node(node["func"]))))
	return fmt.Sprintf(`%s()`, f)
}

func translateBinOp(node map[string]interface{}) interface{} {
	node["type"] = Str(node["type"]) + Str(Node(node["op"])["type"])
	return translateGeneric(node)
}

func translateBinOpMod(node map[string]interface{}) interface{} {
	getArgStr := func(node map[string]interface{}) []string {
		typ := Str(node["type"])
		switch typ {
		default:
			return []string{Str(translateGeneric(node))}
		}
	}
	left := Str(Node(node["left"])["s"])
	right := strings.Join(getArgStr(Node(node["right"])), ", ")
	out := []byte(fmt.Sprintf(`fmt.Sprintf("%s", %s)`, left, right))
	return out
}

func translateFunctionDef(node map[string]interface{}) interface{} {
	hf := transHexyaFuncDef(node)
	hexyaModels[len(hexyaModels)-1].funcs = append(hexyaModels[len(hexyaModels)-1].funcs, hf)
	return [][]byte{}
}

func translateAssign(node map[string]interface{}) interface{} {
	switch {
	case isFieldDeclaration(node):
		transFieldDeclaration(node)
		return [][]byte{{}}
	case IsOdooVar(node):
		return [][]byte{{}}
	default:
		return Comment(GetRawText(node["loc"]))
	}
}

func translateClassDef(node map[string]interface{}) interface{} {
	inheritName := GetClassParam(node)
	var out [][]byte
	switch inheritName {
	case "odoo.models.Model", "odoo.models.AbstractModel", "odoo.models.TransientModel":
		transClassDefModel(node)
		return [][]byte{}
	default:
		return append(out, transGenericClass(node)...)
	}
}

func translateImportFrom(node map[string]interface{}) interface{} {
	module := Str(node["module"])
	for _, n := range Sl(node["names"]) {
		name := Str(translateGeneric(Node(n)))
		knownImportedWords[name] = module + "." + name
	}
	return [][]byte{}
}

func translateAttribute(node map[string]interface{}) interface{} {
	attr := Str(node["attr"])
	val := Node(node["value"])
	valval := Str(translateGeneric(val))
	return valval + "." + attr
}

func translateName(node map[string]interface{}) interface{} {
	return GetAlias(Str(node["id"]))
}

func translateAlias(node map[string]interface{}) interface{} {
	name := node["name"].(string)
	alias := node["asname"]
	if alias != nil {
		aliasMap[alias.(string)] = name
	}
	return name
}

/* other Translations */

func transMethodBody(nodes []interface{}) [][]byte {
	var out [][]byte
	if genArgs.RecordSet.TranslatePython() {
		for _, n := range nodes {
			node := Node(n)
			typ := node["type"].(string)
			if _, ok := funcMap[typ]; !ok {
				s, _ := json.Marshal(node)
				fmt.Printf("Error: type %s not handled. %s\n", typ, conf.Sdump(node["loc"]))
				out = append(out, s)
				continue
			}
			out = append(out, funcMap[typ](node).([][]byte)...)
		}
	} else {
		for _, n := range nodes {
			node := Node(n)
			out = append(out, Comment(GetRawText(node["loc"]))...)
		}
	}
	return out
}

func registerControllers(decorator map[string]interface{}, methName string) [][]byte {
	var out [][]byte
	out = append(out, []byte(`func init() {
root := controllers.Registry`))
	delete(uniqueNamesBank, "controllerVarNames")
	firstArg := Node(Sl(decorator["args"])[0])
	switch Str(firstArg["type"]) {
	case "Str":
		argStr := strings.TrimPrefix(Str(firstArg["s"]), "/")
		out = append(out, registerController(argStr, methName)...)
	case "List":
		for _, e := range Sl(firstArg["elts"]) {
			argStr := strings.TrimPrefix(Str(Node(e)["s"]), "/")
			out = append(out, registerController(argStr, methName)...)
		}
	default:
		s, _ := json.Marshal(firstArg)
		out = append(out, s)
	}
	out = append(out, []byte(`}`))
	return out
}

func registerController(argStr, methName string) [][]byte {
	cleanRelPath := func(str string) string {
		if strings.HasPrefix(str, "<") && strings.Contains(str, ":") {
			str = ":" + strings.Split(str, ":")[1]
		}
		reg := regexp.MustCompile("[^a-zA-Z:_]+")
		out := reg.ReplaceAllString(str, "")
		return out
	}

	var out [][]byte
	prevGrp := "root"
	for i, relPath := range strings.Split(argStr, "/")[:len(strings.Split(argStr, "/"))-1] {
		relPath = cleanRelPath(relPath)
		goVarName := MakeUnique(strings.TrimPrefix(relPath, ":"), "controllerVarNames")

		if i == 0 {
			out = append(out, []byte("var ok bool"))
		}
		out = append(out, []byte(fmt.Sprintf(`var %s *controllers.Group
%s, ok = %s.GetGroup("/%s")
if !ok {
	%s = %s.AddGroup("/%s")
}`, goVarName, goVarName, prevGrp, relPath, goVarName, prevGrp, relPath)))
		prevGrp = goVarName
	}
	relPath := strings.Split(argStr, "/")[len(strings.Split(argStr, "/"))-1]
	relPath = cleanRelPath(relPath)
	out = append(out, []byte(fmt.Sprintf(`if %s.HasController(http.MethodGet, "/%s") {
	%s.ExtendController(http.MethodPost, "/%s", %s)
} else {
	%s.AddController(http.MethodPost, "/%s", %s)
}`, prevGrp, relPath, prevGrp, relPath, methName, prevGrp, relPath, methName)))
	return out
}

func transGenericFuncDef(node map[string]interface{}) [][]byte {
	var out [][]byte
	methName := MakeUnique(SnakeToCamel(Str(node["name"]), true), "GenericFuncDef")
	argsStr, returnStr := transHexyaFuncArgs(Node(node["args"]), "")
	for _, dec := range Sl(node["decorator_list"]) {
		decorator := Node(dec)
		if Str(decorator["type"]) == "Call" && GetFullNameImported(Str(Node(decorator["func"])["id"])) == "odoo.http.route" {
			out = append(out, registerControllers(decorator, methName)...)
		}
	}
	out = append(out, []byte(fmt.Sprintf(`func %s(%s) %s {`, methName, argsStr, returnStr)))
	out = append(out, transMethodBody(Sl(node["body"]))...)
	out = append(out, []byte("}"))
	return out
}

func transDomains(str string) string {
	// string to ast
	jsonRaw, err := exec.Command(`python`, `../../python-parse-to-json/parse_python_to_json.py`, str).CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error while jsonize python: %v, %s", err, jsonRaw)
	}
	var pyFile map[string]interface{}
	err = json.Unmarshal(jsonRaw, &pyFile)
	if err != nil {
		fmt.Println(fmt.Sprint(err))
	}
	list := Node(Node(Sl(pyFile["body"])[0])["value"])
	// single domain
	domainOperators := GetDomainOperators()
	transDomain := func(d map[string]interface{}) string {
		domainSl := Sl(Node(d)["elts"])
		strA := Str(translateGeneric(Node(domainSl[0])))
		strA = strings.Replace(strA, `"`, ``, -1)
		strAspl := strings.Split(strA, ".")
		for i, s := range strAspl {
			strAspl[i] = SnakeToCamel(s, true)
		}
		strA = strings.Join(strAspl, "().")
		strB := Str(translateGeneric(Node(domainSl[1])))
		if val, ok := domainOperators[strings.Replace(strB, `"`, ``, -1)]; ok {
			strB = val
		}
		strC := translateGeneric(Node(domainSl[2]))
		return fmt.Sprintf("%s().%s(%s)", strA, strB, strC)
	}
	// main translator
	operandStack := stack.New()
	operatorStack := stack.New()
	pending := false
	if _, ok := list["elts"]; !ok {
		return ""
	}
	for _, d := range Sl(list["elts"]) {
		domain := Node(d)
		if Str(domain["type"]) == "Tuple" {
			operand := transDomain(domain)
			if pending {
				for operandStack.Len() > 0 {
					operandOld := Str(operandStack.Pop())
					operator := ".And()."
					if operatorStack.Len() > 0 {
						operator = Str(operatorStack.Pop())
					}
					operand = operandOld + operator + operand
				}
			}
			operandStack.Push(operand)
			pending = true
		} else {
			operatorStack.Push(domainOperators[Str(domain["s"])])
			pending = false
		}
	}
	return Str(operandStack.Pop())
}

func findMethodDocString(methName string, body []interface{}) (string, []interface{}) {
	first := Node(body[0])
	if Str(first["type"]) == "Expr" {
		value := Node(first["value"])
		if Str(value["type"]) == "Str" {
			spl := SplitSubN(Str(value["s"]), 60)
			return strings.Join(spl, "\n"), body[1:]
		}
	}
	return methName, body
}

func transHexyaFuncDef(node map[string]interface{}) *hexyaFunc {
	res := new(hexyaFunc)
	methName := SnakeToCamel(Str(node["name"]), true)
	res.name = methName
	res.override = MethodIsExtension(methName)
	var out bytes.Buffer

	docString, body := findMethodDocString(methName, Sl(node["body"]))
	docLines := strings.Split(docString, "\n")
	for _, dl := range docLines {
		out.WriteString("// ")
		out.WriteString(dl)
		out.WriteByte('\n')
	}

	argsStr, returnStr := transHexyaFuncArgs(Node(node["args"]), methName)
	funcName := getFunctionName(curModelName, methName)
	out.WriteString(fmt.Sprintf(`func %s(%s) %s {`, funcName, argsStr, returnStr))
	out.WriteByte('\n')
	out.Write(bytes.Join(transMethodBody(body), []byte{'\n'}))
	out.WriteByte('\n')
	out.WriteString("}")
	res.dump = out.String()
	return res
}

func transFuncArgs(node map[string]interface{}, methName string) ([]byte, []byte) {
	var argsStr []byte
	var returnStr []byte
	args := node["args"]
	if args != nil {
		for i, a := range Sl(args) {
			if i > 0 {
				argsStr = append(argsStr, []byte(", ")...)
			}
			arg := Node(a)
			if Str(arg["arg"]) == "self" && methName != "" {
				argsStr = append(argsStr, []byte(fmt.Sprintf(`rs m.%sSet`, curModelName))...)
			} else {
				argName := Str(arg["arg"])
				if !IsNotInStringArray(argName, GetReservedKeywordArray()) {
					argName = argName + "Name"
				}
				argsStr = append(argsStr, []byte(fmt.Sprintf(`%s %s`, argName, findArgType(methName, i)))...)
			}
		}
	}
	returns := node["returns"]
	if returns != nil {
		if len(Sl(returns)) > 1 {
			returnStr = append(returnStr, []byte("(")...)
		}
		for i, r := range Sl(returns) {
			if i > 0 {
				returnStr = append(returnStr, []byte(", ")...)
			}
			fmt.Printf("Error: return string not handled. %s\n", conf.Sdump(Node(r)["loc"]))
		}
		if len(Sl(returns)) > 1 {
			returnStr = append(returnStr, []byte(")")...)
		}
	}
	return argsStr, returnStr
}

func transHexyaFuncArgs(node map[string]interface{}, methName string) ([]byte, []byte) {
	if !IsNotInStringArray(methName, knownComputedFuncs) {
		return []byte(fmt.Sprintf("rs m.%sSet", curModelName)), []byte(fmt.Sprintf("m.%sData", curModelName))
	}
	return transFuncArgs(node, methName)
}

func findArgType(methName string, argNb int) string {
	model, ok := models.Registry.Get(genArgs.RecordSet.ModelName())
	if !ok {
		return "interface{}"
	}
	meth, ok := model.Methods().Get(methName)
	if !ok || argNb >= meth.MethodType().NumIn() {
		return "interface{}"
	}
	return meth.MethodType().In(argNb).String()
}

func isFieldAlreadyCreated(label string) bool {
	mod, ok := models.Registry.Get(curModelName)
	if !ok {
		if IsNotInStringArray(label, []string{
			"CreateDate", "CreateUID", "DisplayName", "HexyaExternalID",
			"HexyaVersion", "LastUpdate", "WriteDate", "WriteUID"}) {
			return false
		}
		return true
	}
	_, ok = mod.Fields().Get(label)
	if !ok {
		return false
	}
	return true
}

func transFieldDeclaration(node map[string]interface{}) {
	fieldLabel := SnakeToCamel(Str(Node(Sl(node["targets"])[0])["id"]), true)
	fieldType := ConvertFieldType(Node(Node(node["value"])["func"])["attr"].(string))
	if len(fieldFuncMap[fieldType]) <= 0 {
		// FieldType is unknown. Create an invalid field
		comment := bytes.Join(Comment([][]byte{GetRawLinePart(node["loc"])}), []byte{'\n'})
		hexyaModels[len(hexyaModels)-1].fields = append(hexyaModels[len(hexyaModels)-1].fields, &hexyaField{
			invalid: true,
			comment: string(comment),
		})
		return
	}
	createdField := isFieldAlreadyCreated(fieldLabel)
	var fieldDump [][]byte
	if !createdField {
		fieldDump = append(fieldDump, []byte(fmt.Sprintf(`fields.%s{`, fieldType)))
	}
	for i, a := range Sl(Node(node["value"])["args"]) {
		arg := Node(a)
		switch createdField {
		case true:
			hexyaModels[len(hexyaModels)-1].modFields = append(hexyaModels[len(hexyaModels)-1].modFields, &modifiedField{
				str:   string(fieldFuncMap[fieldType][i](arg)),
				label: fieldLabel,
			})
		case false:
			fieldDump = append(fieldDump, fieldFuncMap[fieldType][i](arg))
		}
	}
	for _, a := range Sl(Node(node["value"])["keywords"]) {
		arg := Node(a)
		arg["AlreadyCreated"] = createdField
		f, ok := fieldFuncMap[Str(arg["arg"])]
		if !ok {
			fmt.Printf("Error: field keyword %s not handled.\n", Str(arg["arg"]))
			fieldDump = append(fieldDump, Comment([][]byte{GetRawLinePart(arg["loc"])})...)
			continue
		}
		switch createdField {
		case true:
			hexyaModels[len(hexyaModels)-1].modFields = append(hexyaModels[len(hexyaModels)-1].modFields, &modifiedField{
				str:   string(f[0](Node(arg["value"]))),
				label: fieldLabel,
			})
		case false:
			fieldDump = append(fieldDump, f[0](Node(arg["value"])))
		}
	}
	if !createdField {
		hexyaModels[len(hexyaModels)-1].fields = append(hexyaModels[len(hexyaModels)-1].fields, &hexyaField{
			name: fieldLabel,
			dump: string(bytes.Join(fieldDump, []byte{'\n'})) + "},\n",
		})
	}
}

func getSqlConstraints(sl []interface{}) string {
	var out bytes.Buffer
	var constraintsValNode map[string]interface{}
	for _, p := range sl {
		part := Node(p)
		if Str(part["type"]) != "Assign" {
			continue
		}
		if Str(Node(Sl(part["targets"])[0])["id"]) == "_sql_constraints" {
			constraintsValNode = Node(part["value"])
		}
	}
	if constraintsValNode != nil {
		for _, e := range Sl(constraintsValNode["elts"]) {
			elem := Node(e)
			elets := Sl(elem["elts"])
			out.WriteString(fmt.Sprintf(`h.%s().AddSQLConstraint("%s", "%s", "%s")`, curModelName, Str(Node(elets[0])["s"]), Str(Node(elets[1])["s"]), Str(Node(elets[2])["s"])))
			out.WriteByte('\n')
		}
	}
	return out.String()
}

func transGenericClass(node map[string]interface{}) [][]byte {
	var out [][]byte
	if isInInit {
		out = append(out, []byte("}"))
		isInInit = false
	}
	for _, p := range Sl(node["body"]) {
		part := Node(p)
		typ := part["type"].(string)
		if _, ok := funcMap[typ]; !ok {
			fmt.Printf("Error: type %s not handled. %s\n", typ, conf.Sdump(part["loc"]))
			s, _ := json.Marshal(part)
			out = append(out, []byte("/*"))
			out = append(out, []byte(s))
			out = append(out, []byte("*/"))
			continue
		}
		out = append(out, funcMap[typ](part).([][]byte)...)
	}
	return out
}

func transClassDefModel(node map[string]interface{}) {
	// model name
	curModelName = GetModelName(node["body"])
	hexyaModels = append(hexyaModels, &hexyaModel{
		name:           curModelName,
		sqlConstraints: getSqlConstraints(Sl(node["body"])),
	})
	// body
	for _, p := range Sl(node["body"]) {
		part := Node(p)
		typ := part["type"].(string)
		if _, ok := funcMap[typ]; !ok {
			continue
		}
		funcMap[typ](part)
	}
}

func isFieldDeclaration(node map[string]interface{}) bool {
	value := Node(node["value"])
	if value["type"].(string) == "Call" {
		fnc := Node(value["func"])
		if fnc["type"].(string) == "Attribute" {
			value = Node(fnc["value"])
			if value["type"].(string) == "Name" {
				fullName := GetFullNameImported(funcMap["Name"](value).(string))
				if fullName == "odoo.fields" {
					return true
				}
			}
		}
	}
	return false
}

/* Main */

func translatePyFile(pyFile map[string]interface{}) []byte {
	var out bytes.Buffer
	isInInit = false
	isInAddFields = false
	curModelName = ""
	if fieldFuncMap == nil {
		fieldFuncMap = BuildFieldFuncMap()
	}
	aliasMap = make(map[string]string)
	knownImportedWords = make(map[string]string)
	knownComputedFuncs = []string{}
	out.WriteString(fmt.Sprintf(`package %s`, genArgs.moduleName))
	out.WriteByte('\n')
	out.WriteString(fmt.Sprintf(`
	import (
		"net/http"

		"github.com/hexya-erp/hexya/src/controllers"
		"github.com/hexya-erp/hexya/src/models"
		"github.com/hexya-erp/hexya/src/models/types"
		"github.com/hexya-erp/hexya/src/models/types/dates"
		"github.com/hexya-erp/pool/h"
		"github.com/hexya-erp/pool/q"
	)
	`))
	for _, n := range Sl(pyFile["body"]) {
		node := Node(n)
		typ := node["type"].(string)
		if _, ok := funcMap[typ]; !ok {
			out.Write(bytes.Join(Comment(GetRawText(node["loc"])), []byte{'\n'}))
			continue
		}
		funcMap[typ](node)
	}
	for _, hm := range hexyaModels {
		out.WriteString(fmt.Sprintf("var fields_%s = map[string]models.FieldDefinition {\n", hm.name))
		for _, f := range hm.fields {
			switch f.invalid {
			case true:
				out.WriteString(f.comment)
			case false:
				out.WriteString(fmt.Sprintf(`"%s": %s`, f.name, f.dump))
				out.WriteByte('\n')
			}
		}
		out.WriteString("}\n\n")
		for _, hf := range hm.funcs {
			out.WriteString(hf.dump)
			out.WriteString("\n")
		}
	}
	out.WriteString("func init() {\n")
	for _, hm := range hexyaModels {
		out.WriteString(fmt.Sprintf("models.NewModel(\"%s\")\n", hm.name))
		out.WriteString(fmt.Sprintf("h.%s().AddFields(fields_%s)\n", hm.name, hm.name))
		for _, mf := range hm.modFields {
			spl := strings.Split(mf.str, ":")
			out.WriteString(fmt.Sprintf("h.%s().Fields().%s().Set%s(%s)\n", hm.name, mf.label, spl[0], spl[1]))
		}
		for _, hf := range hm.funcs {
			var format string
			switch hf.override {
			case false:
				format = "h.%s().NewMethod(\"%s\", %s)\n"
			case true:
				format = "h.%s().Methods().%s().Extend(%s)\n"
			}
			out.WriteString(fmt.Sprintf(format, hm.name, hf.name, getFunctionName(hm.name, hf.name)))
		}
		out.WriteByte('\n')
	}
	out.WriteString("}")
	return out.Bytes()
}

func TranslatePy2(filePath string) {
	// fix pyfile mistakes
	text, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	rx := regexp.MustCompile(`,(\s(#.*\n)?)*\)`)
	text = rx.ReplaceAll(text, []byte(`)`))
	err = ioutil.WriteFile(filePath, text, 0777)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	// autopep8
	o, err := exec.Command("autopep8", fmt.Sprintf("%s", filePath), "--in-place").CombinedOutput()
	if err != nil {
		FatalError("could not execute command. Error: %v\n%s", err, o)
	}
	// init
	conf = spew.NewDefaultConfig()
	conf.MaxDepth = 2
	InitFuncMap()
	GKnownMethodNames = GetKnownMethodNames()
	if strings.HasPrefix(filepath.Base(filePath), "__") {
		NoTranslation(filePath)
		return
	}
	// read file
	text, err = ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	rawText = bytes.Split(text, []byte("\n"))
	cmd := exec.Command(`python`, `../../python-parse-to-json/parse_python_to_json.py`, `--pyfile=`+filePath)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
	// marshall
	var pyFile map[string]interface{}
	err = json.Unmarshal(out.Bytes(), &pyFile)
	if err != nil {
		fmt.Println(fmt.Sprint(err))
	}
	// translate
	output := translatePyFile(pyFile)
	// write
	newPath := genArgs.outputPath + "/" + strings.Replace(strings.Replace(strings.Replace(filePath, genArgs.inputPath+"/", "", -1), "/", "_", -1), ".py", ".go", -1)
	err = ioutil.WriteFile(newPath, output, 0777)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	// goimports
	o, err = exec.Command("goimports", "-w", fmt.Sprintf("%s", newPath)).CombinedOutput()
	if err != nil {
		fmt.Printf("could not execute command. Error: %v\n%s", err, o)
	}
	fmt.Printf("File Written: %s\n", newPath)
}
