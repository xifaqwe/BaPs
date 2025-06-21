//go:build !debug
// +build !debug

package gdconf

import (
	"os"
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"

	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/pkg/logger"
	pb "google.golang.org/protobuf/proto"
)

var (
    excelPath   = "./resources/Excel"
    excelDbPath = "./resources/ExcelDB"
)

func (g *GameConfig) LoadExcel() {
retry:
    info, err := os.Stat(g.dataPath)
    if err != nil || !info.IsDir() {
        panic("找不到文件夹: " + g.dataPath)
    }
    if !strings.HasSuffix(g.dataPath, "/") {
        g.dataPath += "/"
    }

    binFile := g.dataPath + "Excel.bin"
    data, err := os.ReadFile(binFile)
    if err != nil {
        if os.IsNotExist(err) {
            logger.Info("Excel.bin not found; generating from dumped JSONs...")
            if genErr := GenerateExcelBin(excelPath, excelDbPath, binFile); genErr != nil {
                logger.Error("生成 Excel.bin 失败: %v", genErr)
                return
            }
            logger.Info("Excel.bin 生成成功，重新加载...")
            goto retry
        }
        logger.Error("无法读取 Excel.bin: %v", err)
        return
    }

    g.Excel = new(sro.Excel)
    if err := pb.Unmarshal(data, g.Excel); err != nil {
        panic("解析 Excel.bin 失败，请检查版本是否匹配")
    }
}

func GenerateExcelBin(excelPath, excelDbPath, outBinPath string) error {
    excel := &sro.Excel{}
    folders := []string{excelPath, excelDbPath}

    for _, folder := range folders {
        err := filepath.Walk(folder, func(path string, info fs.FileInfo, err error) error {
            if err != nil {
                logger.Error("walk %s failed: %v", path, err)
                return err
            }
            if info.IsDir() || !strings.HasSuffix(info.Name(), ".json") {
                return nil
            }

            raw, err := ioutil.ReadFile(path)
            if err != nil {
                logger.Error("read %s failed: %v", path, err)
                return err
            }

            fixedJSON := strings.ReplaceAll(string(raw), `"None_"`, `"None"`)

            fieldName := strings.TrimSuffix(info.Name(), ".json")
            v := reflect.ValueOf(excel).Elem()
            field := v.FieldByName(fieldName)
            if !field.IsValid() {
                return nil
            }

            sliceType := field.Type()
            if sliceType.Kind() != reflect.Slice {
                return nil
            }

            elemType := sliceType.Elem()
            if elemType.Kind() == reflect.Ptr {
                elemType = elemType.Elem()
            }
            if elemType.Kind() != reflect.Struct {
                return nil
            }

            var rawSlice []map[string]interface{}
            if err := json.Unmarshal([]byte(fixedJSON), &rawSlice); err != nil {
                logger.Error("json unmarshal %s failed: %v", path, err)
                return err
            }

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

            cleanJSON, err := json.Marshal(rawSlice)
            if err != nil {
                logger.Error("re-marshal failed: %v", err)
                return err
            }

            slicePtr := reflect.New(sliceType).Interface()
            if err := json.Unmarshal(cleanJSON, slicePtr); err != nil {
                logger.Error("final unmarshal %s failed: %v", path, err)
                return err
            }

            field.Set(reflect.ValueOf(slicePtr).Elem())
            logger.Info("Loaded %s", fieldName)
            return nil
        })

        if err != nil {
            logger.Error("error walking folder %s: %v", folder, err)
            return err
        }
    }

    binData, err := pb.Marshal(excel)
    if err != nil {
        logger.Error("failed to marshal proto: %v", err)
        return err
    }
    if err := ioutil.WriteFile(outBinPath, binData, 0644); err != nil {
        logger.Error("failed to write bin file: %v", err)
        return err
    }
    return nil
}

func isBoolField(t reflect.Type, jsonKey string) bool {
    if t.Kind() == reflect.Ptr {
        t = t.Elem()
    }
    if t.Kind() != reflect.Struct {
        return false
    }

    for i := 0; i < t.NumField(); i++ {
        f := t.Field(i)
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

func loadExcelFile[T any](path string, table *[]*T) {
	*table = make([]*T, 0)
}
