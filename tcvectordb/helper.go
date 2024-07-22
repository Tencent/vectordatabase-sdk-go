package tcvectordb

import (
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

func ConvertField(field *Field) (result *olama.Field) {
	switch field.Type() {
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
	}
	return
}
