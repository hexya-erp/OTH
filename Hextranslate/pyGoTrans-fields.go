package Hextranslate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hexya-addons/web/odooproxy"
)

/* field funcs */

func BuildFieldFuncMap() map[string][]func(map[string]interface{}) []byte {
	o := make(map[string][]func(map[string]interface{}) []byte)
	f := func(v ...func(map[string]interface{}) []byte) []func(map[string]interface{}) []byte {
		return v
	}
	o["Char"] = f(
		encodeFieldString)
	o["Boolean"] = o["Char"]
	o["Binary"] = o["Char"]
	o["Integer"] = o["Char"]
	o["Text"] = o["Char"]
	o["HTML"] = o["Char"]
	o["Date"] = o["Char"]
	o["DateTime"] = o["Char"]
	o["Serialized"] = o["Char"]
	o["ID"] = o["Char"]
	o["Monetary"] = f(
		encodeFieldString,
		encodeFieldCurrency)
	o["Float"] = f(
		encodeFieldString,
		encodeFieldDigits)
	o["Selection"] = f(
		encodeFieldSelection,
		encodeFieldString)
	o["Reference"] = o["Selection"]
	o["Many2One"] = f(
		encodeFieldComodel,
		encodeFieldString)
	o["One2Many"] = f(
		encodeFieldComodel,
		encodeFieldInverse,
		encodeFieldString)
	o["Many2Many"] = f(
		encodeFieldComodel,
		encodeFieldRelation,
		encodeFieldColomn1,
		encodeFieldColomn2,
		encodeFieldString)
	o["string"] = f(encodeFieldString)
	o["compute"] = f(encodeFieldCompute)
	o["required"] = f(encodeFieldRequired)
	o["index"] = f(encodeFieldIndex)
	o["readonly"] = f(encodeFieldReadOnly)
	o["unique"] = f(encodeFieldUnique)
	o["copy"] = f(encodeFieldCopy)
	o["auto_join"] = f(exemptField)
	o["default"] = f(encodeFieldDefault)
	o["defaults"] = f(encodeFieldDefault)
	o["related"] = f(encodeFieldRelated)
	o["ondelete"] = f(encodeFieldOnDelete)
	o["help"] = f(encodeFieldHelp)
	o["domain"] = f(encodeFieldDomain)
	o["translate"] = f(encodeFieldTranslate)
	o["comodel_name"] = f(encodeFieldComodel)
	o["store"] = f(encodeFieldStore)
	o["size"] = f(encodeFieldSize)
	return o
}

func exemptField(node map[string]interface{}) []byte {
	return nil
}

func encodeFieldSize(node map[string]interface{}) []byte {
	return BSl(fmt.Sprintf(`Size: %s,`, translateGeneric(node)))
}

func encodeFieldCurrency(node map[string]interface{}) []byte {
	return []byte("TODO   EncodeFieldsCurrency")
}

func encodeFieldHelp(node map[string]interface{}) []byte {
	switch Str(node["type"]) {
	case "Str":
		spl := SplitSubN(Str(node["s"]), 60)
		for i, s := range spl {
			s = strings.Replace(s, "\n", "\\n", -1)
			s = strings.Replace(s, `"`, `\"`, -1)
			spl[i] = `"` + s + `"`
		}
		return []byte(fmt.Sprintf("Help: %s,", strings.Join(spl, " + \n")))
	}
	return PrintNodeOutput(node)
}

func encodeFieldOnDelete(node map[string]interface{}) []byte {
	s := Str(node["s"])
	return []byte(fmt.Sprintf("OnDelete: `%s`,", s))
}

func encodeFieldRelated(node map[string]interface{}) []byte {
	str := Str(node["s"])
	var splout []string
	for _, s := range strings.Split(str, ".") {
		splout = append(splout, SnakeToCamel(s, true))
	}
	return []byte(fmt.Sprintf("Related: `%s`,", strings.Join(splout, ".")))
}

func encodeFieldDomain(node map[string]interface{}) []byte {
	switch Str(node["type"]) {
	case "Str":
		return []byte(fmt.Sprintf(`Filter: q.%s,`, transDomains(Str(node["s"]))))
	case "Lambda":
		str := Str(GetRawLinePart(Node(node["body"])["loc"]))
		str = ReplaceRegex(str, `self\._name`, fmt.Sprintf(`'%s'`, curModelName))
		str = strings.TrimSpace(transDomains(strings.TrimSpace(str)))
		return []byte(fmt.Sprintf(`Filter: q.%s,`, str))
	case "List":
		str := Str(GetRawLinePart(node["loc"]))
		return []byte(fmt.Sprintf(`Filter: q.%s,`, transDomains(str)))
	default:
		return PrintNodeOutput(node)
	}
}

func encodeFieldDefault(node map[string]interface{}) []byte {
	filterRegexes := func(line string) string {
		line = ReplaceRegex(line, `odoo\.fields\.Datetime\.now`, `dates.Now()`)
		line = ReplaceRegex(line, `self\.env\.user`, `env.Uid()`)
		line = ReplaceRegex(line, `self\.env`, `env`)
		return line
	}
	encodeFieldDefaultLambdaAttribute := func(node map[string]interface{}) []byte {
		line := GetFullNameImported(Str(translateGeneric(node)))
		line = filterRegexes(line)
		return []byte(fmt.Sprintf(`Default: func (env models.Environment) interface{} { return %s },`, line))
	}
	encodeFieldDefaultLambda := func(node map[string]interface{}) []byte {
		body := Node(node["body"])
		line := Str(translateGeneric(body))
		line = filterRegexes(line)
		return []byte(fmt.Sprintf(`Default: func (env models.Environment) interface{} { return %s },`, line))
	}

	var line string
	switch Str(node["type"]) {
	case "Lambda":
		return encodeFieldDefaultLambda(node)
	case "Attribute":
		return encodeFieldDefaultLambdaAttribute(node)
	default:
		line = Str(translateGeneric(node))
		line = strings.Replace(line, "True", `true`, 1)
		line = strings.Replace(line, "False", `false`, 1)
	}
	if line == "" {
		b, _ := json.Marshal(node)
		line = string(b)
	}
	return []byte(fmt.Sprintf(`Default: models.DefaultValue(%s),`, line))
}

func encodeFieldRequired(node map[string]interface{}) []byte {
	node["fieldArgLabel"] = "Required"
	return encodeFieldBool(node)
}

func encodeFieldStore(node map[string]interface{}) []byte {
	node["fieldArgLabel"] = "Stored"
	return encodeFieldBool(node)
}

func encodeFieldTranslate(node map[string]interface{}) []byte {
	node["fieldArgLabel"] = "Translate"
	return encodeFieldBool(node)
}

func encodeFieldIndex(node map[string]interface{}) []byte {
	node["fieldArgLabel"] = "Index"
	return encodeFieldBool(node)
}

func encodeFieldReadOnly(node map[string]interface{}) []byte {
	node["fieldArgLabel"] = "ReadOnly"
	return encodeFieldBool(node)
}

func encodeFieldUnique(node map[string]interface{}) []byte {
	node["fieldArgLabel"] = "Unique"
	return encodeFieldBool(node)
}

func encodeFieldCopy(node map[string]interface{}) []byte {
	if !IsNotInStringArray(Str(node["parentField"]), []string{"One2Many"}) {
		node["fieldArgLabel"] = "Copy"
		return encodeFieldBool(node)
	}
	node["fieldArgLabel"] = "NoCopy"
	node["inverted"] = true
	return encodeFieldBool(node)
}

func encodeFieldInverse(node map[string]interface{}) []byte {
	str := SnakeToCamel(Str(node["Str"]), true)
	str = `"` + str + `"`
	return []byte(fmt.Sprintf(`ReverseFK: %s,`, str))
}

func encodeFieldRelation(node map[string]interface{}) []byte {
	str := SnakeToCamel(Str(node["Str"]), true)
	str = `"` + str + `"`
	return []byte(fmt.Sprintf(`M2MLinkModelName: %s,`, str))
}

func encodeFieldColomn1(node map[string]interface{}) []byte {
	str := SnakeToCamel(Str(node["Str"]), true)
	str = `"` + str + `"`
	return []byte(fmt.Sprintf(`M2MOurField: %s,`, str))
}

func encodeFieldColomn2(node map[string]interface{}) []byte {
	str := SnakeToCamel(Str(node["Str"]), true)
	str = `"` + str + `"`
	return []byte(fmt.Sprintf(`M2MTheirField: %s,`, str))
}

func encodeFieldComodel(node map[string]interface{}) []byte {
	return []byte(fmt.Sprintf(`RelationModel: h.%s(),`, odooproxy.ConvertModelName(Str(node["s"]))))
}

func encodeFieldDigits(node map[string]interface{}) []byte {
	return []byte("TODO   EncodeFieldsDigits")
}

func encodeFieldSelection(node map[string]interface{}) []byte {
	var out [][]byte
	if Str(node["type"]) == "List" {
		out = append(out, []byte("Selection: types.Selection{"))
		for _, e := range Sl(node["elts"]) {
			elem := Node(e)
			first := Str(Node(Sl(elem["elts"])[0])["s"])
			sec := Str(Node(Sl(elem["elts"])[1])["s"])
			out = append(out, []byte(fmt.Sprintf(`"%s": "%s",`, first, sec)))
		}
		out = append(out, []byte("},"))
		return bytes.Join(out, []byte("\n"))
	}
	return []byte(fmt.Sprintf(`Selection: %s,`, Str(translateGeneric(node))))
}

func encodeFieldCompute(node map[string]interface{}) []byte {
	name := SnakeToCamel(Str(node["s"]), true)
	knownComputedFuncs = append(knownComputedFuncs, name)
	return []byte(fmt.Sprintf(`Compute: h.%s().Methods().%s(),`, curModelName, name))
}

func encodeFieldString(node map[string]interface{}) []byte {
	return []byte(fmt.Sprintf(`String: "%s",`, Str(node["s"])))
}

/* generic encodeField */

func encodeFieldBool(node map[string]interface{}) []byte {
	outbool := false
	str := strings.Replace(Str(translateGeneric(node)), `"`, ``, -1)
	if !IsNotInStringArray(str, []string{`true`, `True`, `1`}) {
		outbool = true
	}
	if val, ok := node["inverted"]; ok && val == true {
		outbool = !outbool
	}
	outboolStr := "false"
	if outbool {
		outboolStr = "true"
	}
	return []byte(fmt.Sprintf(`%s: %s,`, Str(node["fieldArgLabel"]), outboolStr))
}
