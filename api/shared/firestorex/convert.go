package firestorex

import (
	"fmt"
	"reflect"
	"time"
)

// SnapshotDataTo populates out (pointer to struct) from the given snapshot using `firestore` tags.
// Centralized helper used by fakes/adapters or services if needed.
func SnapshotDataTo(ds map[string]interface{}, out interface{}) error {
	v := reflect.ValueOf(out)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("out must be non-nil pointer")
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("out must be pointer to struct")
	}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("firestore")
		if tag == "" {
			tag = field.Name
		}
		if val, ok := ds[tag]; ok {
			fv := v.Field(i)
			if fv.CanSet() {
				switch fv.Kind() {
				case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
					if n, ok := val.(int64); ok {
						fv.SetUint(uint64(n))
					} else if n, ok := val.(int); ok {
						fv.SetUint(uint64(n))
					} else if n, ok := val.(uint64); ok {
						fv.SetUint(n)
					} else if n, ok := val.(uint); ok {
						fv.SetUint(uint64(n))
					}
				case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
					if n, ok := val.(int64); ok {
						fv.SetInt(n)
					} else if n, ok := val.(int); ok {
						fv.SetInt(int64(n))
					} else if n, ok := val.(uint64); ok {
						fv.SetInt(int64(n))
					} else if n, ok := val.(uint); ok {
						fv.SetInt(int64(n))
					}
				case reflect.String:
					if s, ok := val.(string); ok {
						fv.SetString(s)
					}
				case reflect.Struct:
					// Support time.Time fields
					if fv.Type() == reflect.TypeOf(time.Time{}) {
						switch tval := val.(type) {
						case time.Time:
							fv.Set(reflect.ValueOf(tval))
						case string:
							if parsed, err := time.Parse(time.RFC3339, tval); err == nil {
								fv.Set(reflect.ValueOf(parsed))
							}
						}
					}
				}
			}
		}
	}
	return nil
}
