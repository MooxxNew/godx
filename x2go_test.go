package x2go

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

func TestPrintStruct(t *testing.T) {
	x2go := New([]byte(xexam1))
	s := x2go.String()

	if s != xstruct1 {
		t.Error("expect\n", xstruct1, "but got\n", s)
	}
}

var xexam1 = `<?xml version="1.0"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:tem="http://tempuri.org/">
<soapenv:Body>
    <tem:Authen>
        <tem:Username>?</tem:Username>
        <tem:Password>?</tem:Password>
        <tem:DomainName>?</tem:DomainName>
        <tem:ClientIP>?</tem:ClientIP>
    </tem:Authen>
</soapenv:Body>
</soapenv:Envelope>`

var xstruct1 = `type Envelope struct {
	XMLName xml.Name ` + "`" + `xml:"soapenv:Envelope"` + "`" + `
	Soapenv string ` + "`" + `xml:"xmlns:soapenv,attr"` + "`" + `
	Tem string ` + "`" + `xml:"xmlns:tem,attr"` + "`" + `
	Body Body ` + "`" + `xml:"soapenv:Body"` + "`" + `
}

type Body struct {
	XMLName xml.Name ` + "`" + `xml:"soapenv:Body"` + "`" + `
	Authen Authen ` + "`" + `xml:"tem:Authen"` + "`" + `
}

type Authen struct {
	XMLName xml.Name ` + "`" + `xml:"tem:Authen"` + "`" + `
	Username string ` + "`" + `xml:"tem:Username"` + "`" + `
	Password string ` + "`" + `xml:"tem:Password"` + "`" + `
	DomainName string ` + "`" + `xml:"tem:DomainName"` + "`" + `
	ClientIP string ` + "`" + `xml:"tem:ClientIP"` + "`" + `
}

`

func TestLayer(t *testing.T) {
	b, err := ioutil.ReadFile("./main/cmp.xml")
	if err != nil {
		t.Error("File cmp.xml not found.")
		return
	}

	x2go := New(b)

	l := x2go.Layer()
	if l != 5 {
		t.Error("It should return 5 but was ", l)
	}
}

func TestSkeleton(t *testing.T) {
	b, err := ioutil.ReadFile("./main/cmp.xml")
	if err != nil {
		t.Error("File cmp.xml not found.")
		return
	}

	x2go := New(b)

	bones := x2go.Skeleton()

	expect := map[string][]string{
		"":                []string{"Envelope"},
		"Envelope":        []string{"Header", "Body"},
		"Body":            []string{"executeBatch"},
		"executeBatch":    []string{"sessionID", "commands"},
		"commands":        []string{"audienceID", "audienceLevel", "debug", "eventParameters", "interactiveChannel", "methodIdentifier", "relyOnExistingSession"},
		"eventParameters": []string{"name", "valueAsString", "valueDataType", "valueAsNumeric"},
	}

	for k, v := range bones.(map[string][]string) {
		if len(v) != len(expect[k]) {
			t.Error("Something went wrong.")
			t.Errorf(">>>%# v", bones)
		}
	}
	// arrange("", bones.(map[string][]string))
}

func TestIdentifyType(t *testing.T) {
	bones := map[string][]string{
		"":                []string{"Envelope"},
		"Envelope":        []string{"Header", "Body"},
		"Body":            []string{"executeBatch"},
		"executeBatch":    []string{"sessionID", "commands"},
		"commands":        []string{"audienceID", "audienceLevel", "debug", "eventParameters", "interactiveChannel", "methodIdentifier", "relyOnExistingSession"},
		"eventParameters": []string{"name", "valueAsString", "valueDataType", "valueAsNumeric"},
	}

	id := Identify(bones)

	expect := map[string]map[string]string{
		"":                map[string]string{"Envelope": "Envelope"},
		"Envelope":        map[string]string{"Header": "Header", "Body": "Body"},
		"Body":            map[string]string{"executeBatch": "ExecuteBatch"},
		"executeBatch":    map[string]string{"sessionID": "string", "commands": "Commands"},
		"commands":        map[string]string{"audienceID": "string", "audienceLevel": "string", "debug": "string", "eventParameters": "EventParameters", "interactiveChannel": "string", "methodIdentifier": "string", "relyOnExistingSession": "string"},
		"eventParameters": map[string]string{"name": "string", "valueAsString": "string", "valueDataType": "string", "valueAsNumeric": "string"},
	}

	for k, v := range expect {
		if len(v) != len(id[k]) {
			t.Error("It might be wrong.")
		}
	}
}

func TestPrint(t *testing.T) {
	id := map[string]map[string]string{
		"":                map[string]string{"Envelope": "Envelope"},
		"Envelope":        map[string]string{"Header": "Header", "Body": "Body"},
		"Body":            map[string]string{"executeBatch": "ExecuteBatch"},
		"executeBatch":    map[string]string{"sessionID": "string", "commands": "Commands"},
		"commands":        map[string]string{"audienceID": "string", "audienceLevel": "string", "debug": "string", "eventParameters": "EventParameters", "interactiveChannel": "string", "methodIdentifier": "string", "relyOnExistingSession": "string"},
		"eventParameters": map[string]string{"name": "string", "valueAsString": "string", "valueDataType": "string", "valueAsNumeric": "string"},
	}

	echo(id)
}

func echo(id map[string]map[string]string) {
	for k, v := range id {
		if k == "" {
			continue
		}

		fmt.Println("type", k, "struct")
		for name, typ := range v {
			fmt.Println(" ", strings.Title(name), typ, "`xml:"+`"`+name+`"`+"`")
		}
	}
}

func arrange(key string, bones map[string][]string) {
	for i := range bones[key] {
		fmt.Println(key + ":" + bones[key][i])
		arrange(bones[key][i], bones)
	}
}
