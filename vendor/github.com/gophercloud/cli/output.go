package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

// INDENT is the indentation passed to json.MarshalIndent
const INDENT string = "  "

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
