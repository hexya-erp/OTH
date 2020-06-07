package Hextranslate

import (
	"encoding/xml"
	"regexp"
)

var (
	genArgs               GenArgs
	GreservedKeywordArray []string
	GDomainOperators      map[string]string
	GTemplateIdTranslator map[string]string
	GEntryTypeTranslator  map[string]string
	GKnownMethodNames     []string
	laterCSVs             []OdooRecord
)

const (
	RECTYPEVIEW int = iota
	RECTYPEACTION
	RECTYPEMENU
)

type Manifest map[string][]string

type SecurityEntry struct {
	ModelId string
	GroupId string
}

type ComplexReplacementRegex struct {
	rx   *regexp.Regexp
	rstr string
	fnc  func(gr []string) []string
}

type FileData struct {
	Imports []string
	Raw     []string
	Classes []ClassData
	Methods []MethodData
}

type ClassData struct {
	ClassName string
	VarDef    map[string]string
	SqlDef    []string
	Raw       []string
	Fields    []FieldData
	Classes   []ClassData
	Methods   []MethodData
}

type FieldData struct {
	Name           string
	Type           string
	Args           map[string]string
	ArgsRaw        string
	AlreadyCreated bool
}

type ApiData struct {
	ApiName string
	ApiArgs string
}

type MethodInfoData struct {
	Path string
	Line int
	Raw  string
}

type MethodData struct {
	Info       MethodInfoData
	Apis       []ApiData
	MethodName string
	MethodArgs string
	Methods    []MethodData
	Raw        []string
}

type HexyaHexyaXml struct {
	Data HexyaData `xml:"data"`
}

type HexyaData struct {
	Records []HexyaRecord `xml:"HEXYA_HEXYA_TMP_RECORD"`
}

type HexyaRecord struct {
	Id       string `xml:"id,attr"`
	Type     string `xml:"type,attr"`
	Name     string `xml:"name,attr"`
	Model    string `xml:"model,attr"`
	ViewMode string `xml:"view_mode,attr"`
	ViewId   string `xml:"view_id,attr"`
	InnerXML string `xml:",innerxml"`
}

type GenericRecord struct {
	Type     string
	ID       string
	Name     string
	Model    string
	ViewMode string
	ViewId   string
	Data     string
}

type OdooOdoo struct {
	XMLNAME  xml.Name       `xml:"odoo"`
	Record   []OdooRecord   `xml:"record"`
	Data     []OdooData     `xml:"data"`
	Template []OdooTemplate `xml:"template"`
}

type OdooScriptEntry struct {
	XMLNAME xml.Name `xml:"script"`
	Type    string   `xml:"type,attr"`
	Src     string   ` xml:"src,attr"`
}

type OdooLinkEntry struct {
	XMLNAME xml.Name `xml:"link"`
	Rel     string   `xml:"rel,attr"`
	Type    string   `xml:"type,attr"`
	Href    string   `xml:"href,attr"`
}

type GenericEntry struct {
	Type string
	Path string
}

type OdooXpath struct {
	XMLNAME  xml.Name          `xml:"xpath"`
	Expr     string            `xml:"expr,attr"`
	Position string            `xml:"position,attr"`
	Script   []OdooScriptEntry `xml:"script"`
	Link     []OdooLinkEntry   `xml:"link"`
}

type OdooTemplate struct {
	XMLNAME   xml.Name          `xml:"template"`
	Id        string            `xml:"id,attr"`
	Name      string            `xml:"name,attr"`
	InheritId string            `xml:"inherit_id,attr"`
	Script    []OdooScriptEntry `xml:"script"`
	Link      []OdooLinkEntry   `xml:"link"`
	Xpath     []OdooXpath       `xml:"xpath"`
}

type OdooData struct {
	XMLNAME  xml.Name       `xml:"data"`
	Template []OdooTemplate `xml:"template"`
	Record   []OdooRecord   `xml:"record"`
}

type OdooRecord struct {
	XMLNAME xml.Name    `xml:"record"`
	Field   []OdooField `xml:"field"`
	Id      string      `xml:"id,attr"`
	Model   string      `xml:"model,attr"`
}

type OdooField struct {
	XMLNAME  xml.Name `xml:"field"`
	InnerXML string   `xml:",innerxml"`
	Name     string   `xml:"name,attr"`
	Ref      string   `xml:"ref,attr"`
}

type SecuOdoo struct {
	XMLNAME xml.Name   `xml:"odoo"`
	Data    []SecuData `xml:"data"`
}

type SecuData struct {
	XMLNAME xml.Name     `xml:"data"`
	Record  []SecuRecord `xml:"record"`
}

type SecuRecord struct {
	XMLNAME xml.Name    `xml:"record"`
	Field   []SecuField `xml:"field"`
	Id      string      `xml:"id,attr"`
	Model   string      `xml:"model,attr"`

	Name string
}

type SecuField struct {
	XMLNAME  xml.Name `xml:"field"`
	InnerXML string   `xml:",innerxml"`
	Name     string   `xml:"name,attr"`
	Ref      string   `xml:"ref,attr"`
}

type CsvData struct {
	Models          []string
	ModelEntryCount map[string]int
	ModelHeads      map[string][]string
	ModelContent    map[string][]CsvValue
}

type CsvValue struct {
	Value map[string]string
}

type XmlDataOdoo struct {
	XMLNAME xml.Name      `xml:"odoo"`
	Data    []XmlDataData `xml:"data"`
}

type XmlDataData struct {
	XMLNAME xml.Name        `xml:"data"`
	Record  []XmlDataRecord `xml:"record"`
}

type XmlDataRecord struct {
	XMLNAME xml.Name       `xml:"record"`
	Field   []XmlDataField `xml:"field"`
	Id      string         `xml:"id,attr"`
	Model   string         `xml:"model,attr"`
}

type XmlDataField struct {
	XMLNAME  xml.Name `xml:"field"`
	InnerXML string   `xml:",innerxml"`
	Name     string   `xml:"name,attr"`
	Ref      string   `xml:"ref,attr"`
}

func GetReservedKeywordArray() []string {
	if len(GreservedKeywordArray) < 4 {
		GreservedKeywordArray = []string{
			"break", "default", "func", "interface",
			"select", "case", "defer", "go",
			"map", "struct", "chan", "else",
			"goto", "package", "switch", "const",
			"return", "var", "if", "range",
			"type", "continue", "for", "import",
			"fallthrough",
		}
	}
	return GreservedKeywordArray
}

func GetDomainOperators() map[string]string {
	if len(GDomainOperators) < 4 {
		out := make(map[string]string)
		out["|"] = ".Or()."
		out["&"] = ".And()."
		out["="] = "Equals"
		out["!="] = "NotEquals"
		out[">"] = "Greater"
		out[">="] = "GreaterOrEqual"
		out["<"] = "Lower"
		out["<="] = "LowerOrEqual"
		out["in"] = "In"
		out["not in"] = "NotIn"
		out["=like"] = "Like"
		out["=ilike"] = "ILike"
		out["like"] = "Contains"
		out["ilike"] = "IContains"
		out["not like"] = "NotContains"
		out["not ilike"] = "NotIContains"
		out["child_of"] = "ChildOf"
		GDomainOperators = out
	}
	return GDomainOperators
}

func GetKnownMethodNames() []string {
	if len(GKnownMethodNames) < 4 {
		GKnownMethodNames = []string{
			"SearchByName", "NameGet", "Search",
			"SearchAll", "Load", "ForceLoad",
			"Write", "Copy", "Read",
			"Unlink", "Create", "ComputeLastUpdate",
			"ComputeDisplayName", "FieldsGet", "FieldGet",
			"DefaultGet", "CheckRecursion", "Onchange",
			"Browse", "SearchCount", "Fetch",
			"SearchAll", "GroupBy", "Aggregates",
			"Limit", "Offset", "OrderBy",
			"Union", "Subtract", "Intersect",
			"CartesianProduct", "Equals", "Sorted",
			"SortedDefault", "SortedByField", "Filtered",
			"GetRecord", "WithEnv", "WithNewContext",
			"Sudo",
		}
	}
	return GKnownMethodNames
}

func GetTemplateIdTranslator() map[string]string {
	if len(GTemplateIdTranslator) < 2 {
		out := make(map[string]string)
		out["web.assets_common"] = "Common"
		out["web.assets_backend"] = "Backend"
		out["web.assets_frontend"] = "Frontend"
		out["web.less_helpers"] = "LessHelpers"
		GTemplateIdTranslator = out
	}
	return GTemplateIdTranslator
}

func GetEntryTypeTranslator() map[string]string {

	if len(GEntryTypeTranslator) < 2 {
		out := make(map[string]string)
		out["text/less"] = "Less"
		out["text/css"] = "CSS"
		out["text/javascript"] = "JS"
		GEntryTypeTranslator = out
	}
	return GEntryTypeTranslator
}
