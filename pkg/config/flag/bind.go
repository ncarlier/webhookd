package configflag

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/ncarlier/webhookd/pkg/helper"
)

// Bind conf struct tags with flags
func Bind(conf interface{}, prefix string) error {
	rv := reflect.ValueOf(conf)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		rv = rv.Elem()
	}
	typ := rv.Type()

	for i := 0; i < typ.NumField(); i++ {
		fieldType := typ.Field(i)
		field := rv.Field(i)

		var key, desc, val string
		// Get field key from struct tag
		if tag, ok := fieldType.Tag.Lookup("flag"); ok {
			key = tag
		} else {
			continue
		}
		// Get field description from struct tag
		if tag, ok := fieldType.Tag.Lookup("desc"); ok {
			desc = tag
		}
		// Get field value from struct tag
		if tag, ok := fieldType.Tag.Lookup("default"); ok {
			val = tag
		}

		// Get field value and description from environment variable
		if fieldType.Type.Kind() == reflect.Slice {
			val = getEnvValue(prefix, key+"s", val)
			desc = getEnvDesc(prefix, key+"s", desc)
		} else {
			val = getEnvValue(prefix, key, val)
			desc = getEnvDesc(prefix, key, desc)
		}

		// Get field value by reflection from struct definition
		// And bind value to command line flag
		switch fieldType.Type.Kind() {
		case reflect.String:
			field.SetString(val)
			ptr, _ := field.Addr().Interface().(*string)
			flag.StringVar(ptr, key, val, desc)
		case reflect.Bool:
			bVal, err := strconv.ParseBool(val)
			if err != nil {
				return fmt.Errorf("invalid boolean value for %s: %v", key, err)
			}
			field.SetBool(bVal)
			ptr, _ := field.Addr().Interface().(*bool)
			flag.BoolVar(ptr, key, bVal, desc)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if field.Kind() == reflect.Int64 && field.Type().PkgPath() == "time" && field.Type().Name() == "Duration" {
				d, err := time.ParseDuration(val)
				if err != nil {
					return fmt.Errorf("invalid duration value for %s: %v", key, err)
				}
				field.SetInt(int64(d))
				ptr, _ := field.Addr().Interface().(*time.Duration)
				flag.DurationVar(ptr, key, d, desc)
			} else {
				i64Val, err := strconv.ParseInt(val, 0, fieldType.Type.Bits())
				if err != nil {
					return fmt.Errorf("invalid number value for %s: %v", key, err)
				}
				field.SetInt(i64Val)
				ptr, _ := field.Addr().Interface().(*int)
				flag.IntVar(ptr, key, int(i64Val), desc)
			}
		case reflect.Slice:
			sliceType := field.Type().Elem()
			if sliceType.Kind() == reflect.String {
				if strings.TrimSpace(val) != "" {
					vals := strings.Split(val, ",")
					sl := make([]string, len(vals))
					copy(sl, vals)
					field.Set(reflect.ValueOf(sl))
					ptr, _ := field.Addr().Interface().(*[]string)
					af := newArrayFlags(ptr)
					flag.Var(af, key, desc)
				}
			}
		}
	}
	return nil
}

func getEnvKey(prefix, key string) string {
	return helper.ToScreamingSnake(prefix + "_" + key)
}

func getEnvValue(prefix, key, fallback string) string {
	if value, ok := os.LookupEnv(getEnvKey(prefix, key)); ok {
		return value
	}
	return fallback
}

func getEnvDesc(prefix, key, desc string) string {
	return fmt.Sprintf("%s (env: %s)", desc, getEnvKey(prefix, key))
}
