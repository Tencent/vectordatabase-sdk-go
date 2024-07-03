package test

import (
	"fmt"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb"
	"testing"
	"time"
)

func TestHNSWIndex(t *testing.T) {
	cli.DropDatabase(ctx, database)
	db, err := cli.CreateDatabase(ctx, database)
	if err != nil {
		panic(err)
	}

	index := tcvectordb.Indexes{}
	index.VectorIndex = append(index.VectorIndex, tcvectordb.VectorIndex{
		FilterIndex: tcvectordb.FilterIndex{
			FieldName: "vector",
			FieldType: tcvectordb.Vector,
			IndexType: tcvectordb.HNSW,
		},
		Dimension:  3,
		MetricType: tcvectordb.COSINE,
		Params: &tcvectordb.HNSWParam{
			M:              16,
			EfConstruction: 200,
		},
	})
	index.FilterIndex = append(index.FilterIndex, tcvectordb.FilterIndex{FieldName: "id", FieldType: tcvectordb.String, IndexType: tcvectordb.PRIMARY})
	index.FilterIndex = append(index.FilterIndex, tcvectordb.FilterIndex{FieldName: "bookName", FieldType: tcvectordb.String, IndexType: tcvectordb.FILTER})
	index.FilterIndex = append(index.FilterIndex, tcvectordb.FilterIndex{FieldName: "page", FieldType: tcvectordb.Uint64, IndexType: tcvectordb.FILTER})

	db.WithTimeout(time.Second * 30)
	_, err = db.CreateCollection(ctx, collectionName, 3, 0, "test collection", index)
	if err != nil {
		panic(err)
	}
	collListRes, err := db.ListCollection(ctx)
	if err != nil {
		panic(err)
	}
	var col *tcvectordb.Collection
	for _, c := range collListRes.Collections {
		if c.CollectionName == collectionName {
			col = c
			break
		}
	}
	if col == nil {
		t.Fatal("Collection created failed")
	}

	if !compareVectorIndex(index.VectorIndex, col.Indexes.VectorIndex) {
		t.Fatal("index check failed")
	}

	if !compareFilterIndex(index.FilterIndex, col.Indexes.FilterIndex) {
		t.Fatal("index check failed")
	}
	cli.DropDatabase(ctx, database)
}

func TestIVFFlatIndex(t *testing.T) {
	cli.DropDatabase(ctx, database)
	db, err := cli.CreateDatabase(ctx, database)
	if err != nil {
		panic(err)
	}

	index := tcvectordb.Indexes{}
	index.VectorIndex = append(index.VectorIndex, tcvectordb.VectorIndex{
		FilterIndex: tcvectordb.FilterIndex{
			FieldName: "vector",
			FieldType: tcvectordb.Vector,
			IndexType: tcvectordb.IVF_FLAT,
		},
		Dimension:  3,
		MetricType: tcvectordb.IP,
		Params: &tcvectordb.IVFFLATParams{
			NList: 10,
		},
	})
	index.FilterIndex = append(index.FilterIndex, tcvectordb.FilterIndex{FieldName: "id", FieldType: tcvectordb.String, IndexType: tcvectordb.PRIMARY})
	index.FilterIndex = append(index.FilterIndex, tcvectordb.FilterIndex{FieldName: "bookName", FieldType: tcvectordb.String, IndexType: tcvectordb.FILTER})
	index.FilterIndex = append(index.FilterIndex, tcvectordb.FilterIndex{FieldName: "page", FieldType: tcvectordb.Uint64, IndexType: tcvectordb.FILTER})

	db.WithTimeout(time.Second * 30)
	_, err = db.CreateCollection(ctx, collectionName, 3, 0, "test collection", index)
	if err != nil {
		panic(err)
	}
	collListRes, err := db.ListCollection(ctx)
	if err != nil {
		panic(err)
	}
	var col *tcvectordb.Collection
	for _, c := range collListRes.Collections {
		if c.CollectionName == collectionName {
			col = c
			break
		}
	}
	if col == nil {
		t.Fatal("Collection created failed")
	}

	if !compareVectorIndex(index.VectorIndex, col.Indexes.VectorIndex) {
		t.Fatal("index check failed")
	}

	if !compareFilterIndex(index.FilterIndex, col.Indexes.FilterIndex) {
		t.Fatal("index check failed")
	}
	cli.DropDatabase(ctx, database)
}

func compareVectorIndex(a, b []tcvectordb.VectorIndex) bool {
	indexSet := make(map[string]struct{})
	for _, index := range a {
		indexSet[vectorToString(&index)] = struct{}{}
	}
	for _, index := range b {
		_, ok := indexSet[vectorToString(&index)]
		if !ok {
			return false
		}
		delete(indexSet, vectorToString(&index))
	}
	return len(indexSet) == 0
}

func compareFilterIndex(a, b []tcvectordb.FilterIndex) bool {
	indexSet := make(map[string]struct{})
	for _, index := range a {
		if index.IndexType == tcvectordb.PRIMARY {
			continue
		}
		indexSet[filterToString(&index)] = struct{}{}
	}
	for _, index := range b {
		if index.IndexType == tcvectordb.PRIMARY {
			continue
		}
		_, ok := indexSet[filterToString(&index)]
		if !ok {
			return false
		}
		delete(indexSet, filterToString(&index))
	}
	return len(indexSet) == 0
}

func filterToString(index *tcvectordb.FilterIndex) string {
	return fmt.Sprintf("%s%s%s%s", index.FieldName, index.FieldType, index.ElemType, index.IndexType)
}

func vectorToString(index *tcvectordb.VectorIndex) string {
	return fmt.Sprintf("%s%d%s%d%s", filterToString(&index.FilterIndex), index.Dimension, index.MetricType, index.IndexedCount, index.Params.Name())
}
