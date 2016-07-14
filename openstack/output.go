package openstack

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"text/tabwriter"

	"github.com/Sirupsen/logrus"
	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/util"
)

// INDENT is the indentation passed to json.MarshalIndent
const INDENT string = "  "

type output struct {
	writer    io.Writer
	fields    []string
	noHeader  bool
	format    string
	quiet     bool
	logger    *logrus.Logger
	commander lib.Commander
}

var once sync.Once

func (o *output) OutputResult(result interface{}) error {
	o.writer = os.Stdout
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
		//panic(r)
	case DebugMsg:
		o.logger.Debug(r)
	case *ProgressStatus:
		switch o.quiet {
		case false:
			progresser, ok := o.commander.(lib.Progresser)
			if !ok {
				return fmt.Errorf("Command does not allow status updates")
			}
			once.Do(progresser.InitProgress)
			switch r.MsgType {
			case StatusStarted:
				progresser.Started(r)
			case StatusUpdated:
				progresser.Updated(r)
			case StatusCompleted:
				progresser.Completed(r)
				o.OutputResult(r.Result)
			case StatusErrored:
				progresser.Errored(r)
			}
			// o.logger.Info(r)
		}
	case map[string]interface{}, []map[string]interface{}:
		o.LimitFields(r)
		switch o.format {
		case "json":
			o.jsonOut(r)
		default:
			switch s := r.(type) {
			case map[string]interface{}:
				o.singleTable(s)
			case []map[string]interface{}:
				o.listTable(s)
			}
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

	return nil
}

func (o *output) LimitFields(r interface{}) {
	switch len(o.fields) {
	case 0:
		switch r.(type) {
		case []map[string]interface{}:
			fieldser, ok := o.commander.(lib.Fieldser)
			switch ok {
			case false:
				o.logger.Infof("List command has no default fields")
				return
			default:
				o.fields = fieldser.Fields()
			}
		default:
			return
		}
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
	if preTabler, ok := o.commander.(lib.PreTabler); ok {
		err := preTabler.PreTable(many)
		if err != nil {
			fmt.Fprintln(w, fmt.Sprintf("Error formatting table: %s", err))
			return
		}
	}
	if !o.noHeader {
		fmt.Fprintln(w, strings.Join(o.fields, "\t"))
	}
	for _, m := range many {
		f := []string{}
		for _, k := range o.fields {
			f = append(f, fmt.Sprint(m[k]))
		}
		fmt.Fprintln(w, strings.Join(f, "\t"))
	}
	w.Flush()
}

func (o output) singleTable(m map[string]interface{}) {
	w := tabwriter.NewWriter(o.writer, 0, 8, 0, '\t', 0)
	if preTabler, ok := o.commander.(lib.PreTabler); ok {
		err := preTabler.PreTable(m)
		if err != nil {
			fmt.Fprintln(w, fmt.Sprintf("Error formatting table: %s", err))
			return
		}
	}

	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(w, "%s\t%s\n", k, strings.Replace(fmt.Sprint(m[k]), "\n", "\n\t", -1))
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
