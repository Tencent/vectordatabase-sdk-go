package ai_document_set

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api"
)

// Document document struct for document api
type QueryDocumentSet struct {
	DocumentSetId      string                      `json:"documentSetId"`
	DocumentSetName    string                      `json:"documentSetName"`
	Text               *string                     `json:"text,omitempty"`
	TextPrefix         *string                     `json:"textPrefix,omitempty"`
	DocumentSetInfo    *DocumentSetInfo            `json:"documentSetInfo,omitempty"`
	ScalarFields       map[string]interface{}      `json:"-"`
	SplitterPreprocess *DocumentSplitterPreprocess `json:"splitterPreprocess,omitempty"`
	ParsingProcess     *api.ParsingProcess         `json:"parsingProcess,omitempty"`
}

type DocumentSetInfo struct {
	TextLength      *uint64 `json:"textLength,omitempty"`
	ByteLength      *uint64 `json:"byteLength,omitempty"`
	IndexedProgress *uint64 `json:"indexedProgress,omitempty"`
	IndexedStatus   *string `json:"indexedStatus,omitempty"` // Ready | New | Loading | Failure
	CreateTime      *string `json:"createTime,omitempty"`
	LastUpdateTime  *string `json:"lastUpdateTime,omitempty"`
	IndexedErrorMsg *string `json:"indexedErrorMsg,omitempty"`
	Keywords        *string `json:"keywords,omitempty"`
}

// [DocumentSplitterPreprocess] holds the parameters for splitting document chunks.
//
// Fields:
//   - AppendTitleToChunk:  (Optional) Whether to append the title to the splitting chunks (default to false).
//   - AppendKeywordsToChunk:  (Optional) Whether to append the keywords to the splitting chunks,
//     which are extracted from the full text (default to true).
//   - ChunkSplitter: The regular expression used to configure the document splitting method.
//     For example, "\n{2,}" can be used in QA contents.
type DocumentSplitterPreprocess struct {
	AppendTitleToChunk    *bool   `json:"appendTitleToChunk,omitempty"`
	AppendKeywordsToChunk *bool   `json:"appendKeywordsToChunk,omitempty"`
	ChunkSplitter         *string `json:"chunkSplitter,omitempty"`
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
	ds := json.NewDecoder(bytes.NewReader(data))
	ds.UseNumber()
	err = ds.Decode(&temp.ScalarFields)
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

type SearchDocument struct {
	Score       float64           `json:"score"`
	Data        SearchData        `json:"data"`
	DocumentSet SearchDocumentSet `json:"documentSet"`
}

type SearchData struct {
	Text                     string   `json:"text"`
	StartPos                 int      `json:"startPos"`
	EndPos                   int      `json:"endPos"`
	Pre                      []string `json:"pre"`
	Next                     []string `json:"next"`
	ParagraphTitle           string   `json:"paragraphTitle"`
	AllParentParagraphTitles []string `json:"allParentParagraphTitles"`
}

type SearchDocumentSet struct {
	DocumentSetId   string                 `json:"documentSetId"`
	DocumentSetName string                 `json:"documentSetName"`
	ScalarFields    map[string]interface{} `json:"-"`
}

func (s SearchDocumentSet) MarshalJSON() ([]byte, error) {
	type Alias SearchDocumentSet
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

func (s *SearchDocumentSet) UnmarshalJSON(data []byte) error {
	type Alias SearchDocumentSet
	var temp Alias
	err := json.Unmarshal(data, &temp)
	if err != nil {
		return err
	}
	ds := json.NewDecoder(bytes.NewReader(data))
	ds.UseNumber()
	err = ds.Decode(&temp.ScalarFields)
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

	*s = SearchDocumentSet(temp)
	return nil
}
