package tcvectordb

import (
	"encoding/json"
	"log"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/olama"
)

func ConvertDbType(dataType olama.DataType) string {
	switch dataType {
	case olama.DataType_BASE:
		return DbTypeBase
	case olama.DataType_AI_DOC:
		return AIDOCDbType
	default:
		return ""
	}
}

func ConvertField2Grpc(field *Field) (result *olama.Field) {
	switch field.Type() {
	case Double:
		result = &olama.Field{OneofVal: &olama.Field_ValDouble{ValDouble: field.Float()}}
	case Uint64:
		result = &olama.Field{OneofVal: &olama.Field_ValU64{ValU64: field.Uint64()}}
	case String:
		result = &olama.Field{OneofVal: &olama.Field_ValStr{ValStr: []byte(field.String())}}
	case Array:
		stringArray := field.StringArray()
		byteArray := make([][]byte, 0, len(stringArray))
		for _, s := range stringArray {
			byteArray = append(byteArray, []byte(s))
		}
		result = &olama.Field{OneofVal: &olama.Field_ValStrArr{ValStrArr: &olama.Field_StringArray{StrArr: byteArray}}}
	case Json:
		jsonData, err := json.Marshal(field.Val)
		if err != nil {
			log.Printf("[Error] marshal failed when converting field to rpc request body. err: %v", err.Error())
			return
		}
		result = &olama.Field{OneofVal: &olama.Field_ValJson{ValJson: jsonData}}
	}
	return
}

func ConvertGrpc2Field(field *olama.Field) (result *Field) {
	result = &Field{}
	switch v := field.GetOneofVal().(type) {
	case *olama.Field_ValStr:
		result.Val = string(v.ValStr)
	case *olama.Field_ValU64:
		result.Val = v.ValU64
	case *olama.Field_ValDouble:
		result.Val = v.ValDouble
	case *olama.Field_ValStrArr:
		result.Val = ConvertByte2StringSlice(v.ValStrArr.StrArr)
	case *olama.Field_ValJson:
		result.Val = make(map[string]interface{}, 0)
		err := json.Unmarshal(v.ValJson, &result.Val)
		if err != nil {
			log.Printf("[Error] Unmarshal failed when converting rpc request body to field. err: %v", err.Error())
			return
		}
	}
	return
}

func ConvertByte2StringSlice(bytes [][]byte) []string {
	strings := make([]string, len(bytes))
	for i := 0; i < len(bytes); i++ {
		strings[i] = string(bytes[i])
	}
	return strings
}
