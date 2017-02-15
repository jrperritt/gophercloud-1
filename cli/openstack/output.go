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

	"github.com/gophercloud/gophercloud/cli/lib"
	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/cli/util"
)

var outreg = os.Stdout
var outerr = os.Stderr

// outres prints the results of the command
func outres(cmd interfaces.Commander, in chan interface{}) error {
	for result := range in {
		lib.Log.Debugf("rcvd result on in: %+v", result)
		switch r := result.(type) {
		case error:
			errout(r)
		case map[string]interface{}:
			mapout(cmd, r)
		case []map[string]interface{}:
			mapout(cmd, r)
		case string:
			fmt.Fprintf(outreg, "%v\n", r)
		case io.Reader:
			readerout(cmd, r)
		default:
			defaultjson(r)
		}
	}
	return nil
}

func jsonout(i interface{}) {
	j, _ := json.MarshalIndent(i, "", "  ")
	fmt.Fprintln(outreg, string(j))
}

func errout(e error) {
	jsonout(map[string]interface{}{"error": e.Error()})
}

func tableout(tabler interfaces.Tabler, i interface{}) {
	ms, ok := i.([]map[string]interface{})
	if !ok {
		fmt.Fprintln(outreg, fmt.Sprintf("Don't know how to properly print type (%T)", i))
		fmt.Fprintln(outreg, fmt.Sprintf("%v", i))
	}

	w := tabwriter.NewWriter(outreg, 0, 8, 1, '\t', 0)
	if preTabler, ok := tabler.(interfaces.PreTabler); ok {
		err := preTabler.PreTable(ms)
		if err != nil {
			fmt.Fprintln(w, fmt.Sprintf("Error formatting table: %s", err))
			return
		}
	}

	var fields []string
	if f, ok := tabler.(interfaces.Fieldser); ok {
		fields = f.Fields()
	} else {
		fields = tabler.DefaultTableFields()
	}

	if tabler.ShouldHeader() {
		fmt.Fprintln(w, strings.Join(fields, "\t"))
	}

	for _, m := range ms {
		s := []string{}
		for _, k := range fields {
			s = append(s, fmt.Sprint(m[k]))
		}
		fmt.Fprintln(w, strings.Join(s, "\t"))
	}
	w.Flush()

}

func mapout(cmd interfaces.Commander, i interface{}) {
	LimitFields(cmd, i)
	if tabler, ok := cmd.(interfaces.Tabler); ok && tabler.ShouldTable() {
		tableout(tabler, i)
		return
	}
	jsonout(i)
}

func readerout(cmd interfaces.Commander, r io.Reader) {
	if rc, ok := r.(io.ReadCloser); ok {
		defer rc.Close()
	}
	if customWriterer, ok := cmd.(interfaces.CustomWriterer); ok {
		writer, err := customWriterer.CustomWriter()
		if err != nil {
			fmt.Fprintf(outerr, "Error creating custom writer: %s\n", err)
			return
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

	defaultjson(string(bytes))
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

func defaultjson(i interface{}) {
	m := map[string]interface{}{"result": i}
	jsonout(m)
}
