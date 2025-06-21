//go:build generate_excel
// +build generate_excel

package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	pb "github.com/gucooing/BaPs/common/server_only"

	"google.golang.org/protobuf/proto"
)

func main() {
	excel := &pb.Excel{}
	folders := []string{"./resources/Excel", "./resources/ExcelDB"}

	for _, folder := range folders {
		err := filepath.Walk(folder, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() || !strings.HasSuffix(info.Name(), ".json") {
				return nil
			}

			raw, err := ioutil.ReadFile(path)
			if err != nil {
				return fmt.Errorf("read %s failed: %w", path, err)
			}

			// Replace None_ → None
			fixedJSON := strings.ReplaceAll(string(raw), `"None_"`, `"None"`)

			// Determine which Excel field to fill
			fieldName := strings.TrimSuffix(info.Name(), ".json")
			v := reflect.ValueOf(excel).Elem()
			field := v.FieldByName(fieldName)
			if !field.IsValid() {
				fmt.Printf("Skipping %s: not a valid field in Excel\n", fieldName)
				return nil
			}

			// Must be a slice
			sliceType := field.Type()
			if sliceType.Kind() != reflect.Slice {
				fmt.Printf("Skipping %s: not a valid slice field in Excel\n", fieldName)
				return nil
			}

			// Element must be struct or pointer-to-struct
			elemType := sliceType.Elem()
			if elemType.Kind() == reflect.Ptr {
				elemType = elemType.Elem()
			}
			if elemType.Kind() != reflect.Struct {
				fmt.Printf("Skipping %s: slice element is not a struct (type: %s)\n", fieldName, elemType.String())
				return nil
			}

			// Unmarshal raw JSON into generic slice
			var rawSlice []map[string]interface{}
			if err := json.Unmarshal([]byte(fixedJSON), &rawSlice); err != nil {
				return fmt.Errorf("json unmarshal %s failed: %w", path, err)
			}

			// Coerce any numeric 0/1 into bool for all bool fields
			for _, item := range rawSlice {
				for key, val := range item {
					if isBoolField(elemType, key) {
						switch v := val.(type) {
						case float64:
							item[key] = v != 0
						case int:
							item[key] = v != 0
						}
					}
				}
			}

			// Marshal cleaned JSON back
			cleanJSON, err := json.Marshal(rawSlice)
			if err != nil {
				return fmt.Errorf("re-marshal failed: %w", err)
			}

			// Unmarshal into the actual slice type (which might be []*T or []T)
			slicePtr := reflect.New(sliceType).Interface()
			if err := json.Unmarshal(cleanJSON, slicePtr); err != nil {
				return fmt.Errorf("final unmarshal %s failed: %w", path, err)
			}

			// Assign to excel struct
			field.Set(reflect.ValueOf(slicePtr).Elem())
			fmt.Printf("Loaded %s from %s\n", fieldName, path)
			return nil
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error walking folder %s: %v\n", folder, err)
			os.Exit(1)
		}
	}

	// Finally, serialize the full Excel struct to protobuf binary
	outPath := "./resources/Excel.bin"
	binData, err := proto.Marshal(excel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal proto: %v\n", err)
		os.Exit(1)
	}
	if err := ioutil.WriteFile(outPath, binData, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write bin file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Excel.bin generated successfully.")
}

// isBoolField reports whether a JSON key corresponds to a bool field in the given struct type.
// It handles pointer-to-struct, uses case-insensitive matching against JSON tags or field names.
func isBoolField(t reflect.Type, jsonKey string) bool {
	// Dereference pointer-to-struct
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return false
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		// Determine JSON name from struct tag or field name
		tag := f.Tag.Get("json")
		name := strings.Split(tag, ",")[0]
		if name == "" {
			name = f.Name
		}
		if strings.EqualFold(jsonKey, name) && f.Type.Kind() == reflect.Bool {
			return true
		}
	}
	return false
}
