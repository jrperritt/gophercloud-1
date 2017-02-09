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

var outreg = os.Stdout
var outerr = os.Stderr

// outres prints the results of the command
func outres(cmd interfaces.Commander, alldonech chan interface{}) error {
	for result := range alldonech {
		switch r := result.(type) {
		case error:
			outputError(r)
		case map[string]interface{}:
			outputMap(cmd, r)
		case []map[string]interface{}:
			outputMap(cmd, r)
		case io.Reader:
			outputReader(cmd, r)
		case string:
			fmt.Fprintf(outreg, "%v\n", r)
		default:
			defaultJSON(r)
		}
	}
	return nil
}

func outputJSON(i interface{}) {
	j, _ := json.MarshalIndent(i, "", "  ")
	fmt.Fprintln(outreg, string(j))
}

func outputError(e error) {
	outputJSON(map[string]interface{}{"error": e.Error()})
}

func outputTable(cmd interfaces.Commander, i interface{}) {
	ms, ok := i.([]map[string]interface{})
	if !ok {
		fmt.Fprintln(outreg, fmt.Sprintf("Don't know how to properly print type (%T)", i))
		fmt.Fprintln(outreg, fmt.Sprintf("%v", i))
	}

	w := tabwriter.NewWriter(outreg, 0, 8, 1, '\t', 0)
	if preTabler, ok := cmd.(interfaces.PreTabler); ok {
		err := preTabler.PreTable(ms)
		if err != nil {
			fmt.Fprintln(w, fmt.Sprintf("Error formatting table: %s", err))
			return
		}
	}
	if f, ok := cmd.(interfaces.Fieldser); ok {
		if cmd.(interfaces.Tabler).ShouldHeader() {
			fmt.Fprintln(w, strings.Join(f.Fields(), "\t"))
		}
		for _, m := range ms {
			s := []string{}
			for _, k := range f.Fields() {
				s = append(s, fmt.Sprint(m[k]))
			}
			fmt.Fprintln(w, strings.Join(s, "\t"))
		}
		w.Flush()
	}
}

func outputMap(cmd interfaces.Commander, i interface{}) {
	LimitFields(cmd, i)
	if t, ok := cmd.(interfaces.Tabler); ok && t.ShouldTable() {
		outputTable(cmd, i)
	} else {
		outputJSON(i)
	}
}

func outputReader(cmd interfaces.Commander, r io.Reader) {
	if rc, ok := r.(io.ReadCloser); ok {
		defer rc.Close()
	}
	if customWriterer, ok := cmd.(interfaces.CustomWriterer); ok {
		writer, err := customWriterer.CustomWriter()
		if err != nil {

		}
		_, err = io.Copy(writer, r)
		if err != nil {
			fmt.Fprintf(outerr, "Error copying (io.Reader) result: %s\n", err)
		}
		return
	}

	if toJSONer, ok := cmd.(interfaces.ToJSONer); ok {
		toJSONer.ToJSON()
		return
	}

	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Fprintf(outerr, "Error reading (io.Reader) result: %s\n", err)
		return
	}
	defaultJSON(string(bytes))
}

// LimitFields reduces the number of fields in the output
func LimitFields(cmd interfaces.Commander, r interface{}) {
	if f, ok := cmd.(interfaces.Fieldser); ok {
		if len(f.Fields()) == 0 {
			if tabler, ok := cmd.(interfaces.Tabler); ok {
				f.SetFields(tabler.DefaultTableFields())
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
					if !util.Contains(f.Fields(), k) {
						delete(t, k)
					}
				}
			case []map[string]interface{}:
				for _, i := range t {
					for k := range i {
						if !util.Contains(f.Fields(), k) {
							delete(i, k)
						}
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
