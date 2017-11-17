package noaweb

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"strings"
)

// ViewFunctions struct used to structure package
type ViewFunctions struct{}

// View variable used to structure package
var View ViewFunctions

// ParseTmpl is a function that parses a view template from filePath and inserts
// values from templateMap. Should be called like this:
//  parsedTemplate := ParseTmpl(template.tmpl.html, map[string]interface{}{
//                        "template1": "Content",
//                        "template2": template.HTML("<p>HTML encoded stuff</p>"),
//                    })
func (ViewFunctions) ParseTmpl(filePath string, templateMap interface{}) (string, error) {

	// Read file to string, f
	var dat []byte
	var err error
	if strings.HasPrefix(filePath, "/") {
		dat, err = ioutil.ReadFile(filePath)
	} else {
		dat, err = ioutil.ReadFile(noawebinst.AssetsDir + "/" + filePath)
	}
	if err != nil {
		fmt.Println(err)
	}
	f := string(dat)

	// Parse and execute the template
	t, err := template.New("error").Parse(f)
	if err != nil {
		log.Fatal(err)
	}

	var tpl bytes.Buffer
	t.Execute(&tpl, templateMap)

	return tpl.String(), nil
}

// Parse is a function that parses a view from filePath and returns the
// view as a string.
func (ViewFunctions) Parse(filePath string) (string, error) {

	// Read file to string, f
	var dat []byte
	var err error
	if strings.HasPrefix(filePath, "/") {
		dat, err = ioutil.ReadFile(filePath)
	} else {
		dat, err = ioutil.ReadFile(noawebinst.AssetsDir + "/" + filePath)
	}
	if err != nil {
		fmt.Println(err)
	}

	return string(dat), nil
}
