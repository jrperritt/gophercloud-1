package openstack

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"

	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/cli/util"
)

// CustomWriterer is an interface implemented by commands that offer
// custom output destinations
type CustomWriterer interface {
	CustomWriter() io.Writer
}

// OutputResults prints the results of the command
func OutputResults() error {
	for result := range GC.ResultsRunCommand {
		switch r := result.(type) {
		case error:
			outputError(r)
		case map[string]interface{}:
			outputMap(r)
		case []map[string]interface{}:
			outputMap(r)
		case io.Reader:
			outputReader(r)
		case string:
			fmt.Fprintf(GC.CommandContext.App.Writer, "%v\n", r)
		default:
			defaultJSON(r)
		}
	}
	return nil
}

func outputJSON(i interface{}) {
	j, _ := json.MarshalIndent(i, "", "  ")
	fmt.Fprintln(GC.CommandContext.App.Writer, string(j))
}

func outputError(e error) {
	GC.CommandContext.App.Writer = os.Stderr
	outputJSON(map[string]interface{}{"error": e.Error()})
}

func outputTable(i interface{}) {
	ms, ok := i.([]map[string]interface{})
	if !ok {
		fmt.Fprintln(GC.CommandContext.App.Writer, fmt.Sprintf("Don't know how to properly print type (%T)", i))
		fmt.Fprintln(GC.CommandContext.App.Writer, fmt.Sprintf("%v", i))
	}

	w := tabwriter.NewWriter(GC.CommandContext.App.Writer, 0, 8, 1, '\t', 0)
	if preTabler, ok := GC.Command.(interfaces.PreTabler); ok {
		err := preTabler.PreTable(ms)
		if err != nil {
			fmt.Fprintln(w, fmt.Sprintf("Error formatting table: %s", err))
			return
		}
	}
	if GC.Command.(interfaces.Tabler).ShouldHeader() {
		fmt.Fprintln(w, strings.Join(GC.GlobalOptions.fields, "\t"))
	}
	for _, m := range ms {
		f := []string{}
		for _, k := range GC.GlobalOptions.fields {
			f = append(f, fmt.Sprint(m[k]))
		}
		fmt.Fprintln(w, strings.Join(f, "\t"))
	}
	w.Flush()
}

func outputMap(i interface{}) {
	LimitFields(i)
	if t, ok := GC.Command.(interfaces.Tabler); ok && t.ShouldTable() {
		outputTable(i)
	} else {
		outputJSON(i)
	}
}

func outputReader(r io.Reader) {
	if rc, ok := r.(io.ReadCloser); ok {
		defer rc.Close()
	}
	var writer io.Writer
	customWriterer, ok := GC.Command.(CustomWriterer)
	switch ok {
	case true:
		writer = customWriterer.CustomWriter()
		toTabler, ok := GC.Command.(interfaces.ToTabler)
		switch ok {
		case true:
			toTabler.ToTable()
		case false:
			_, err := io.Copy(writer, r)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error copying (io.Reader) result: %s\n", err)
			}
		}
	case false:
		toJSONer, ok := GC.Command.(interfaces.ToJSONer)
		switch ok {
		case true:
			toJSONer.ToJSON()
		case false:
			bytes, err := ioutil.ReadAll(r)
			if err != nil {
				//return err
			}
			defaultJSON(string(bytes))
		}
	}
}

// LimitFields reduces the number of fields in the output
func LimitFields(r interface{}) {
	if len(GC.GlobalOptions.fields) == 0 {
		if tabler, ok := GC.Command.(interfaces.Tabler); ok {
			GC.GlobalOptions.fields = tabler.DefaultTableFields()
		} else {
			switch t := r.(type) {
			case map[string]interface{}:
				for k, v := range t {
					switch reflect.ValueOf(v).Kind() {
					case reflect.Map, reflect.Slice, reflect.Struct:
						delete(t, k)
					}
				}
			case []map[string]interface{}:
				for _, i := range t {
					for k, v := range i {
						switch reflect.ValueOf(v).Kind() {
						case reflect.Map, reflect.Slice, reflect.Struct:
							delete(i, k)
						}
					}
				}
			}
		}
	} else {
		switch t := r.(type) {
		case map[string]interface{}:
			for k := range t {
				if !util.Contains(GC.GlobalOptions.fields, k) {
					delete(t, k)
				}
			}
		case []map[string]interface{}:
			for _, i := range t {
				for k := range i {
					if !util.Contains(GC.GlobalOptions.fields, k) {
						delete(i, k)
					}
				}
			}
		}
	}
}

func defaultJSON(i interface{}) {
	m := map[string]interface{}{"result": i}
	outputJSON(m)
}
