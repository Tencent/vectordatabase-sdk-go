package ai_document

import (
	"encoding/json"
	"reflect"
	"strings"
)

// Document document struct for document api
type QueryDocumentSet struct {
	DocumentSetId   string                 `json:"documentSetId"`
	DocumentSetName string                 `json:"documentSetName"`
	TextPrefix      string                 `json:"textPrefix"`
	DocumentSetInfo DocumentSetInfo        `json:"documentSetInfo"`
	ScalarFields    map[string]interface{} `json:"-"`
}

type DocumentSetInfo struct {
	TextLength      uint64 `json:"textLength"`
	ByteLength      uint64 `json:"byteLength"`
	IndexedProgress uint64 `json:"indexedProgress"`
	IndexedStatus   string `json:"indexedStatus"` // Ready | New | Loading | Failure
	CreateTime      string `json:"createTime"`
	LastUpdateTime  string `json:"lastUpdateTime"`
}

func (d QueryDocumentSet) MarshalJSON() ([]byte, error) {
	type Alias QueryDocumentSet
	res, err := json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(&d),
	})
	if err != nil {
		return nil, err
	}
	if len(d.ScalarFields) != 0 {
		field, err := json.Marshal(d.ScalarFields)
		if err != nil {
			return nil, err
		}
		if len(field) == 0 {
			return res, nil
		}
		// res = {}
		if len(res) == 2 {
			res = append(res[:1], field[1:]...)
		} else {
			res[len(res)-1] = ','
			res = append(res, field[1:]...)
		}
	}
	return res, nil
}

func (d *QueryDocumentSet) UnmarshalJSON(data []byte) error {
	type Alias QueryDocumentSet
	var temp Alias
	err := json.Unmarshal(data, &temp)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &temp.ScalarFields)
	if err != nil {
		return err
	}
	reflectType := reflect.TypeOf(*d)
	for i := 0; i < reflectType.NumField(); i++ {
		field := reflectType.Field(i)
		tags := strings.Split(field.Tag.Get("json"), ",")
		if tags[0] == "-" {
			continue
		}
		delete(temp.ScalarFields, tags[0])
	}

	*d = QueryDocumentSet(temp)
	return nil
}

type GetDocumentSet struct {
	QueryDocumentSet
	// 这里复用list结构，文本前缀不再返回，而是返回文本内容
	TextPrefix string `json:"-"`
	Text       string `json:"text"`
}

type SearchDocument struct {
	Score       float64     `json:"score"`
	Chunk       Chunk       `json:"chunk"`
	DocumentSet DocumentSet `json:"documentSet"`
}

type Chunk struct {
	Text       string   `json:"text"`
	StartPos   int      `json:"startPos"`
	EndPos     int      `json:"endPos"`
	PreChunks  []string `json:"preChunks"`
	NextChunks []string `json:"nextChunks"`
}

type DocumentSet struct {
	Id              string                 `json:"id"`
	DocumentSetName string                 `json:"documentSetName"`
	ScalarFields    map[string]interface{} `json:"-"`
}

func (s DocumentSet) MarshalJSON() ([]byte, error) {
	type Alias DocumentSet
	res, err := json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(&s),
	})
	if err != nil {
		return nil, err
	}
	if len(s.ScalarFields) != 0 {
		field, err := json.Marshal(s.ScalarFields)
		if err != nil {
			return nil, err
		}
		if len(field) == 0 {
			return res, nil
		}
		// res = {}
		if len(res) == 2 {
			res = append(res[:1], field[1:]...)
		} else {
			res[len(res)-1] = ','
			res = append(res, field[1:]...)
		}
	}
	return res, nil
}

func (s *DocumentSet) UnmarshalJSON(data []byte) error {
	type Alias DocumentSet
	var temp Alias
	err := json.Unmarshal(data, &temp)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &temp.ScalarFields)
	if err != nil {
		return err
	}
	reflectType := reflect.TypeOf(*s)
	for i := 0; i < reflectType.NumField(); i++ {
		field := reflectType.Field(i)
		tags := strings.Split(field.Tag.Get("json"), ",")
		if tags[0] == "-" {
			continue
		}
		delete(temp.ScalarFields, tags[0])
	}

	*s = DocumentSet(temp)
	return nil
}
