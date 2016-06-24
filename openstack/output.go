package openstack

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/Sirupsen/logrus"
	"github.com/jrperritt/rack/util"
)

// INDENT is the indentation passed to json.MarshalIndent
const INDENT string = "  "

type output struct {
	writer   io.Writer
	fields   []string
	noHeader bool
	format   string
	logger   *logrus.Logger
}

// FormatOptions satisfies the Outputter.FormatOptions method
func (o output) GetFormatOptions() []string {
	return []string{
		"json",
		"table",
	}
}

func (o output) OutputResult(result interface{}) error {

	switch r := result.(type) {
	case error:
		o.writer = os.Stderr
		switch o.format {
		case "json":
			o.jsonOut(map[string]interface{}{"error": r.Error()})
		default:
			fmt.Fprintf(o.writer, "%v\n", r)
		}
		return nil
	case DebugMsg:
		o.logger.Debug(r)
	case ProgressStatus:
		o.logger.Info(r)
	case map[string]interface{}:
		o.LimitFields(r)
		switch o.format {
		case "json":
			o.jsonOut(r)
		default:
			o.singleTable(r)
		}
	case []map[string]interface{}:
		o.LimitFields(r)
		switch o.format {
		case "json":
			o.jsonOut(r)
		default:
			o.listTable(r)
		}
	case io.Reader:
		if rc, ok := r.(io.ReadCloser); ok {
			defer rc.Close()
		}
		_, err := io.Copy(o.writer, r)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error copying (io.Reader) result: %s\n", err)
		}
	default:
		switch o.format {
		case "json":
			o.defaultJSON(r)
		default:
			fmt.Fprintf(o.writer, "%v\n", r)
		}
	}

	//_ = resulter.(*Result)

	/*
		switch o.format {
		case "json":
			if jsoner, ok := command.(PreJSONer); ok {
				err = jsoner.PreJSON(resource)
			}
		default:
			if tabler, ok := command.(PreTabler); ok {
				err = tabler.PreTable(resource)
			}
		}
		if err != nil {
			resource.Keys = []string{"error"}
			resource.Result = map[string]interface{}{"error": err.Error()}
		}
	*/

	return nil
}

func (o output) ToTable() {

}

func (o output) ToJSON() {

}

func (o output) LimitFields(r interface{}) {
	switch len(o.fields) {
	case 0:
		return
	}
	switch t := r.(type) {
	case map[string]interface{}:
		for k, _ := range t {
			if !util.Contains(o.fields, k) {
				delete(t, k)
			}
		}
	case []map[string]interface{}:
		for _, i := range t {
			for k, _ := range i {
				if !util.Contains(o.fields, k) {
					delete(i, k)
				}
			}
		}
	}
}

func (o output) jsonOut(i interface{}) {
	j, _ := json.MarshalIndent(i, "", INDENT)
	fmt.Fprintln(o.writer, string(j))
}

func (o output) defaultJSON(i interface{}) {
	m := map[string]interface{}{"result": i}
	o.jsonOut(m)
}

func (o output) listTable(many []map[string]interface{}) {
	w := tabwriter.NewWriter(o.writer, 0, 8, 1, '\t', 0)
	if !o.noHeader {
		// Write the header
		fmt.Fprintln(w, strings.Join(o.fields, "\t"))
	}
	for _, m := range many {
		f := []string{}
		for _, key := range o.fields {
			f = append(f, fmt.Sprint(m[key]))
		}
		fmt.Fprintln(w, strings.Join(f, "\t"))
	}
	w.Flush()
}

func (o output) singleTable(m map[string]interface{}) {
	w := tabwriter.NewWriter(o.writer, 0, 8, 0, '\t', 0)
	for _, key := range o.fields {
		val := fmt.Sprint(m[key])
		fmt.Fprintf(w, "%s\t%s\n", key, strings.Replace(val, "\n", "\n\t", -1))
	}
	w.Flush()
}

func onlyNonNil(m map[string]interface{}) map[string]interface{} {
	for k, v := range m {
		if v == nil {
			m[k] = ""
		}
	}
	return m
}
