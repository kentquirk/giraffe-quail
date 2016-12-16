package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kentquirk/giraffe-quail/types"
)

type Status map[string]interface{}

func (st Status) GetStr(name string) (string, bool) {
	v, ok := st[name]
	if !ok {
		return "", false
	}
	switch s := v.(type) {
	case string:
		return s, true
	case int:
		return strconv.Itoa(s), true
	default:
		fmt.Printf("GetStr: Can't interpret status value %s (%#v)\n", name, s)
		return "", false
	}
}

func (st Status) GetInt(name string) (int, bool) {
	v, ok := st[name]
	if !ok {
		return 0, false
	}
	switch i := v.(type) {
	case string:
		return 0, false
	case int:
		return i, true
	default:
		fmt.Printf("GetInt: Can't interpret status value %s (%#v)\n", name, i)
		return 0, false
	}
}

type StatusHandler struct {
	RawData []Status `json:"data"`
	Data    map[string]interface{}
}

func NewStatusHandler() *StatusHandler {
	return &StatusHandler{Data: make(map[string]interface{})}
}

func (h *StatusHandler) Operate(field types.QueryField, global, local *types.Scope) error {
	if h.RawData == nil {
		resp, err := http.Get("http://atool2-dev.achievementnetwork.org:8082/status/detail")
		if err != nil {
			return err
		}
		err = json.NewDecoder(resp.Body).Decode(&h.RawData)
		if err != nil {
			return err
		}
		// fmt.Println(field)
		// fmt.Println("global=", global)
		// fmt.Println("local=", local)
		result := make([]interface{}, 0)
		for _, s := range h.RawData {
			d := make(map[string]interface{})
			for _, subfield := range field.SelectionSet.Fields {
				name := subfield.Alias
				if name == "" {
					name = subfield.Name
				}
				if v, ok := s[subfield.Name]; ok {
					d[name] = v
				} else {
					d[name] = nil
				}
			}
			result = append(result, d)
		}
		h.Data["Data"] = result
	}
	return nil
}
