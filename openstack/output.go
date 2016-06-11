package openstack

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/Sirupsen/logrus"
	"github.com/gophercloud/cli/lib"
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

func (o output) OutputResult(resulter lib.Resulter) error {

	_ = resulter.(*Result)

	switch resulter.(type) {

	}

	resulter.SetType()

	if resulter.GetError() != nil {
		o.writer = os.Stderr
		return nil
	}

	if resulter.GetValue() == nil {
		resulter.SetValue(resulter.GetEmptyValue())
	}

	o.LimitFields(resulter)

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

	switch r := resulter.GetValue().(type) {
	case map[string]interface{}:
		//m = onlyNonNil(r)
		switch o.format {
		case "json":
			MetadataJSON(o.writer, r, o.fields)
		default:
			MetadataTable(o.writer, r, o.fields)
		}
	case []map[string]interface{}:
		for i, m := range r {
			r[i] = onlyNonNil(m)
		}
		switch o.format {
		case "json":
			ListJSON(o.writer, r, o.fields)
		default:
			ListTable(o.writer, r, o.fields, o.noHeader)
		}
	case io.Reader:
		if rc, ok := resulter.GetValue().(io.ReadCloser); ok {
			defer rc.Close()
		}
		_, err := io.Copy(o.writer, r)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error copying (io.Reader) result: %s\n", err)
		}
	default:
		switch o.format {
		case "json":
			DefaultJSON(o.writer, resulter.GetValue())
		default:
			fmt.Fprintf(o.writer, "%v\n", resulter.GetValue())
		}
	}
	return nil
}

func (o output) ToTable() {

}

func (o output) ToJSON() {

}

func (o output) LimitFields(r lib.Resulter) {

}

func limitJSONFields(m map[string]interface{}, keys []string) map[string]interface{} {
	mLimited := make(map[string]interface{})
	for _, key := range keys {
		if v, ok := m[key]; ok {
			mLimited[key] = v
		}
	}
	return mLimited
}

func jsonOut(w io.Writer, i interface{}) {
	j, _ := json.MarshalIndent(i, "", INDENT)
	fmt.Fprintln(w, string(j))
}

func DefaultJSON(w io.Writer, i interface{}) {
	m := map[string]interface{}{"result": i}
	jsonOut(w, m)
}

func MetadataJSON(w io.Writer, m map[string]interface{}, keys []string) {
	mLimited := limitJSONFields(m, keys)
	jsonOut(w, mLimited)
}

func ListJSON(w io.Writer, maps []map[string]interface{}, keys []string) {
	mLimited := make([]map[string]interface{}, len(maps))
	for i, m := range maps {
		mLimited[i] = limitJSONFields(m, keys)
	}
	jsonOut(w, mLimited)
}

// ListTable writes a table composed of keys as the header with values from many
func ListTable(writer io.Writer, many []map[string]interface{}, keys []string, noHeader bool) {
	w := tabwriter.NewWriter(writer, 0, 8, 1, '\t', 0)
	if !noHeader {
		// Write the header
		fmt.Fprintln(w, strings.Join(keys, "\t"))
	}
	for _, m := range many {
		f := []string{}
		for _, key := range keys {
			f = append(f, fmt.Sprint(m[key]))
		}
		fmt.Fprintln(w, strings.Join(f, "\t"))
	}
	w.Flush()
}

// MetadataTable writes a table to the writer composed of keys on the left and
// the associated metadata on the right column from m
func MetadataTable(writer io.Writer, m map[string]interface{}, keys []string) {
	w := tabwriter.NewWriter(writer, 0, 8, 0, '\t', 0)
	for _, key := range keys {
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
