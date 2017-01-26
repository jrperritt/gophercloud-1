package openstack

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/gophercloud/gophercloud/cli/lib"
	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/cli/util"
)

type CustomWriterer interface {
	CustomWriter() io.Writer
}

func OutputResults() error {
	for result := range GC.ResultsRunCommand {
		switch r := result.(type) {
		case error:
			outputError(r)
		case map[string]interface{}, []map[string]interface{}:
			outputMap(r)
		case io.Reader:
			outputReader(r)
		default:
			switch GC.GlobalOptions.outputFormat {
			case "json":
				defaultJSON(r)
			default:
				fmt.Fprintf(GC.CommandContext.App.Writer, "%v\n", r)
			}
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
	switch GC.GlobalOptions.outputFormat {
	case "json":
		outputJSON(map[string]interface{}{"error": e.Error()})
	default:
		fmt.Fprintf(GC.CommandContext.App.Writer, "%v\n", e)
	}
}
func outputMapTable(m map[string]interface{}) {
	w := tabwriter.NewWriter(GC.CommandContext.App.Writer, 0, 8, 0, '\t', 0)
	if preTabler, ok := GC.Command.(lib.PreTabler); ok {
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

func outputMapsTable(ms []map[string]interface{}) {
	w := tabwriter.NewWriter(GC.CommandContext.App.Writer, 0, 8, 1, '\t', 0)
	if preTabler, ok := GC.Command.(lib.PreTabler); ok {
		err := preTabler.PreTable(ms)
		if err != nil {
			fmt.Fprintln(w, fmt.Sprintf("Error formatting table: %s", err))
			return
		}
	}
	if !GC.GlobalOptions.noHeader {
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
	switch GC.GlobalOptions.outputFormat {
	case "json":
		outputJSON(i)
	default:
		switch s := i.(type) {
		case map[string]interface{}:
			outputMapTable(s)
		case []map[string]interface{}:
			outputMapsTable(s)
		}
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
	case false:
		writer = GC.CommandContext.App.Writer
	}
	switch GC.GlobalOptions.outputFormat {
	case "json":
		toJSONer, ok := GC.Command.(lib.ToJSONer)
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
	default:
		toTabler, ok := GC.Command.(lib.ToTabler)
		switch ok {
		case true:
			toTabler.ToTable()
		case false:
			_, err := io.Copy(writer, r)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error copying (io.Reader) result: %s\n", err)
			}
		}
	}
}

func LimitFields(r interface{}) {
	if len(GC.GlobalOptions.fields) == 0 {
		if fieldser, ok := GC.Command.(interfaces.DefaultTableFieldser); ok {
			GC.GlobalOptions.fields = fieldser.DefaultTableFields()
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
			for k, _ := range t {
				if !util.Contains(GC.GlobalOptions.fields, k) {
					delete(t, k)
				}
			}
		case []map[string]interface{}:
			for _, i := range t {
				for k, _ := range i {
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
