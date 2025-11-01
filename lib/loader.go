package properties

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/ssongin/core"
)

func (p PropertiesMap) ParseVariables(prefix string, target interface{}) {
	err := p.LoadObject(prefix, target)
	core.CheckError("Failed to parse \""+prefix+"\" properties", err)
	core.GetLogger().Info("Parsed properties", "prefix", prefix)
}

func (p PropertiesMap) LoadObject(prefix string, target interface{}) error {
	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return fmt.Errorf("target must be a non-nil pointer to struct")
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("target must be a pointer to struct")
	}
	return populateStruct(v, p, prefix)
}

func populateStruct(v reflect.Value, props PropertiesMap, prefix string) error {
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if !field.CanSet() {
			continue
		}

		tag := fieldType.Tag.Get("prop")
		if tag == "" {
			continue
		}

		fullKey := tag
		if prefix != "" {
			fullKey = prefix + "." + tag
		}

		switch field.Kind() {
		case reflect.Struct:
			if err := populateStruct(field, props, fullKey); err != nil {
				return err
			}

		case reflect.String:
			if val, ok := props[fullKey]; ok {
				field.SetString(val)
			}

		case reflect.Bool:
			if val, ok := props[fullKey]; ok {
				if boolVal, err := strconv.ParseBool(val); err == nil {
					field.SetBool(boolVal)
				}
			}

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if val, ok := props[fullKey]; ok {
				if intVal, err := strconv.ParseInt(val, 10, 64); err == nil {
					field.SetInt(intVal)
				}
			}

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if val, ok := props[fullKey]; ok {
				if uintVal, err := strconv.ParseUint(val, 10, 64); err == nil {
					field.SetUint(uintVal)
				}
			}

		case reflect.Float32, reflect.Float64:
			if val, ok := props[fullKey]; ok {
				if floatVal, err := strconv.ParseFloat(val, 64); err == nil {
					field.SetFloat(floatVal)
				}
			}
		}
	}
	return nil
}
