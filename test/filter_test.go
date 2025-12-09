package test

import (
	"fmt"
	"testing"

	"github.com/tencent/vectordatabase-sdk-go/tcvectordb"
)

func Test_FilterUint64OneCondition(t *testing.T) {
	filter := tcvectordb.NewFilter("a=100")
	fmt.Println(filter.Cond())
}

func Test_FilterStringOneCondition(t *testing.T) {
	filter := tcvectordb.NewFilter("author=\"jerry\"")
	fmt.Println(filter.Cond())
}

func Test_FilterStringInList(t *testing.T) {
	filter := tcvectordb.NewFilter(tcvectordb.In("stringkey", []string{"v1", "v2", "v3"}))
	fmt.Println(filter.Cond())
}

func Test_FilterStringNotInList(t *testing.T) {
	filter := tcvectordb.NewFilter(tcvectordb.NotIn("stringkey", []string{"v1", "v2", "v3"}))
	fmt.Println(filter.Cond())
}

func Test_FilterArrayInclude(t *testing.T) {
	filter := tcvectordb.NewFilter(tcvectordb.Include("arraykey", []string{"v1", "v2", "v3"}))
	fmt.Println(filter.Cond())
}

func Test_FilterArrayExclude(t *testing.T) {
	filter := tcvectordb.NewFilter(tcvectordb.Exclude("arraykey", []string{"v1", "v2", "v3"}))
	fmt.Println(filter.Cond())
}

func Test_FilterArrayIncludeAll(t *testing.T) {
	filter := tcvectordb.NewFilter(tcvectordb.IncludeAll("arraykey", []string{"v1", "v2", "v3"}))
	fmt.Println(filter.Cond())
}

func Test_FilterMultiConditions(t *testing.T) {
	filter := tcvectordb.NewFilter("author=\"jerry\"").And("a=1").Or("r=\"or\"").OrNot("rn=2").AndNot("an=\"andNot\"")
	fmt.Println(filter.Cond())
}

func Test_FilerAndCondtions(t *testing.T) {
	filter := tcvectordb.NewFilter(`author="jerry" or a=1`).And("b=2 or c=3")
	// expect: (author="jerry" or a=1) and (b=2 or c=3)
	// not expect: author="jerry" or a=1 and (b=2 or c=3)
	fmt.Println(filter.Cond())
}
