package publisher

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// CheckConfig checks publisher config
func CheckConfig(pub Publisher, settings map[string]interface{}) error {
	cfg := pub.GetConfigDef()
	errs := []string{}

	for k, v := range cfg {
		if _, ok := settings[k]; !ok && v.Required {
			errs = append(errs, fmt.Sprintf("field '%s' not provided in settings", k))
			continue
		}

		switch v.Type {

		case JSON_STRING:
			value, ok := settings[k].(string)
			if !ok {
				errs = append(errs, fmt.Sprintf("field '%s' must be %s", k, v.Type))
				continue
			}

			m := make(map[string]interface{})
			err := json.Unmarshal([]byte(value), &m)
			if err != nil {
				errs = append(errs, fmt.Sprintf("field '%s' is not a valid json", k))
				continue
			}
			break

		case BOOL:
			_, ok := settings[k].(bool)
			if !ok {
				errs = append(errs, fmt.Sprintf("field '%s' must be %s", k, v.Type))
				continue
			}

		case STRING:
			value, ok := settings[k].(string)
			if !ok {
				errs = append(errs, fmt.Sprintf("field '%s' must be %s", k, v.Type))
				continue
			}

			if v.Possible != nil {
				possible := (v.Possible).([]string)
				if !findStringInSlice(possible, value) {
					errs = append(errs, fmt.Sprintf("field '%s' value must be %s", k, strings.Join(possible, " or ")))
					continue
				}
			}
			break

		case INT:
			value, ok := settings[k].(int)
			if !ok {
				errs = append(errs, fmt.Sprintf("field '%s' must be %s", k, v.Type))
				continue
			}

			if v.Possible != nil {
				possible := (v.Possible).([]int)
				if !findIntInSlice(possible, value) {
					errs = append(errs, fmt.Sprintf("field '%s' value must be %s", k, arrayToString(possible, ", ")))
					continue
				}
			}
			break

		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ", "))
	}

	return nil
}

func findStringInSlice(sl []string, str string) bool {
	for _, v := range sl {
		if v == str {
			return true
		}
	}
	return false
}

func findIntInSlice(sl []int, i int) bool {
	for _, v := range sl {
		if v == i {
			return true
		}
	}
	return false
}

func arrayToString(a []int, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}

func checkPluginValidity(pub Publisher) error {
	m := pub.GetConfigDef()
	errs := []string{}

	for k, v := range m {
		switch v.Type {

		case STRING:
			_, ok := v.Default.(string)
			if !ok {
				errs = append(errs, fmt.Sprintf("default value for field '%s' cannot be casted into %s", k, STRING))
			}

			if v.Possible != nil {
				_, ok = v.Possible.([]string)
				if !ok {
					errs = append(errs, fmt.Sprintf("possible values for field '%s' cannot be casted into %s array", k, STRING))
				}
			}
			break

		case JSON_STRING:
			_, ok := v.Default.(string)
			if !ok {
				errs = append(errs, fmt.Sprintf("default value for field '%s' cannot be casted into %s", k, JSON_STRING))
			}
			break

		case INT:
			_, ok := v.Default.(int)
			if !ok {
				errs = append(errs, fmt.Sprintf("default value for field '%s' cannot be casted into %s", k, INT))
			}

			if v.Possible != nil {
				_, ok = v.Possible.([]int)
				if !ok {
					errs = append(errs, fmt.Sprintf("possible values for field '%s' cannot be casted into %s array", k, INT))
				}
			}
			break

		case BOOL:
			_, ok := v.Default.(bool)
			if !ok {
				errs = append(errs, fmt.Sprintf("default value for field '%s' cannot be casted into %s", k, BOOL))
			}
			break

		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ", "))
	}

	return nil
}
