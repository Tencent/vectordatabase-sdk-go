// Copyright (C) 2023 Tencent Cloud.
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the vectordb-sdk-java), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is furnished
// to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
// INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
// PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
// SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package test

import (
	"log"
	"strconv"
	"sync"
	"testing"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb"
	"github.com/tencent/vectordatabase-sdk-go/tcvectordb/api/collection_view"
)

func TestParallel(t *testing.T) {
	parallel := 10

	counts := 600
	ch := make(chan int, counts)
	for i := 0; i < counts; i++ {
		ch <- i
	}

	var wg sync.WaitGroup
	wg.Add(parallel)

	for i := 0; i < parallel; i++ {
		parrllelId := i
		go func(i int) {
			defer wg.Done()

			index := tcvectordb.Indexes{
				FilterIndex: []tcvectordb.FilterIndex{
					{
						FieldName: "author_name",
						FieldType: tcvectordb.String,
						IndexType: tcvectordb.FILTER,
					},
				},
			}

			enableWordsEmbedding := true
			appendTitleToChunk := true
			appendKeywordsToChunk := false

			for {
				if len(ch) == 0 {
					break
				}
				dbNum := <-ch

				db, err := cli.CreateAIDatabase(ctx, aiDatabase+strconv.Itoa(dbNum))
				printErr(err)
				t.Logf("create database success, %s", db.DatabaseName)

				name := collectionViewName + strconv.Itoa(dbNum)
				coll, err := db.CreateCollectionView(ctx, name, tcvectordb.CreateCollectionViewParams{
					Description: "test ai collectionView",
					Indexes:     index,
					Embedding: &collection_view.DocumentEmbedding{
						Language:             string(tcvectordb.LanguageChinese),
						EnableWordsEmbedding: &enableWordsEmbedding,
					},
					SplitterPreprocess: &collection_view.SplitterPreprocess{
						AppendTitleToChunk:    &appendTitleToChunk,
						AppendKeywordsToChunk: &appendKeywordsToChunk,
					},
					ExpectedFileNum: 204800,
					AverageFileSize: 10240,
				})
				if err != nil {
					log.Printf("CreateCollectionView fail: %v: %v: %v", db, name, err.Error())
				} else {
					log.Printf("CreateCollectionView success: %v: %v", coll.DatabaseName, name)
				}
			}
		}(parrllelId)
	}

	wg.Wait()
	close(ch)
}

func TestListParallel(t *testing.T) {
	parallel := 100

	counts := 286
	ch := make(chan int, counts)
	for i := 0; i < counts; i++ {
		ch <- i
	}

	var wg sync.WaitGroup
	wg.Add(parallel)

	for i := 0; i < parallel; i++ {
		parrllelId := i
		go func(i int) {
			defer wg.Done()
			for {
				if len(ch) == 0 {
					break
				}
				dbNum := <-ch
				db := cli.AIDatabase(aiDatabase + strconv.Itoa(dbNum))

				_, err := db.ListCollectionViews(ctx)
				if err != nil {
					log.Printf("ListCollectionViews fail: %v: %v", db, err.Error())
				} else {
					log.Printf("ListCollectionViews success: %v", db)
				}
			}
		}(parrllelId)
	}

	wg.Wait()
	close(ch)
}
