package ai_document

import (
	"encoding/json"
	"reflect"
	"strings"
)

// Document document struct for document api
type QueryDocument struct {
	Id             string                 `json:"id"`
	ContentLength  uint64                 `json:"_content_length"`
	FileName       string                 `json:"_file_name"`
	FileKeywards   string                 `json:"_file_keywards"`
	Indexed        uint64                 `json:"_indexed"`
	IndexedStatus  uint64                 `json:"_indexed_status"`
	CreateTime     uint64                 `json:"_create_time"`
	LastUpdateTime uint64                 `json:"_last_update_time"`
	TextLength     uint64                 `json:"_text_length"`
	FileMetadata   string                 `json:"_file_metadata"`
	ScalarFields   map[string]interface{} `json:"-"`
}

func (d QueryDocument) MarshalJSON() ([]byte, error) {
	type Alias QueryDocument
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

func (d *QueryDocument) UnmarshalJSON(data []byte) error {
	type Alias QueryDocument
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

	*d = QueryDocument(temp)
	return nil
}

type SearchDocument struct {
	Score      float64
	Chunk      Chunk
	SourceFile SourceFile
}

type Chunk struct {
	Text       string   `json:"text"`
	StartPos   int      `json:"startPos"`
	EndPos     int      `json:"endPos"`
	PreChunks  []string `json:"preChunks"`
	NextChunks []string `json:"nextChunks"`
}

type SourceFile struct {
	Id           string                 `json:"id"`
	FileName     string                 `json:"_file_name"`
	FileMetadata string                 `json:"_file_metadata"`
	CreateTime   string                 `json:"_create_time"`
	ScalarFields map[string]interface{} `json:"-"`
}

func (s SourceFile) MarshalJSON() ([]byte, error) {
	type Alias SourceFile
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

func (s *SourceFile) UnmarshalJSON(data []byte) error {
	type Alias SourceFile
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

	*s = SourceFile(temp)
	return nil
}
