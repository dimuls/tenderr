package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-nlp/dmmclust"
	"github.com/xtgo/set"
	"github.com/xuri/excelize/v2"
)

func makeCorpus(a []string) map[string]int {
	retVal := make(map[string]int)
	var id int
	for _, s := range a {
		for _, f := range strings.Fields(s) {
			if _, ok := retVal[f]; !ok {
				retVal[f] = id
				id++
			}
		}
	}
	return retVal
}

func makeDocuments(a []string, c map[string]int, allowRepeat bool) []dmmclust.Document {
	retVal := make([]dmmclust.Document, 0, len(a))
	for _, s := range a {
		var ts []int
		for _, f := range strings.Fields(s) {
			id := c[f]
			ts = append(ts, id)
		}
		if !allowRepeat {
			ts = set.Ints(ts) // this uniquifies the sentence
		}
		retVal = append(retVal, dmmclust.TokenSet(ts))
	}
	return retVal
}

type logLine struct {
	text  string
	count int
}

type logLines []logLine

func (ls logLines) Len() int {
	return len(ls)
}
func (ls logLines) Less(i, j int) bool {
	return ls[i].count > ls[j].count
}
func (ls logLines) Swap(i, j int) {
	ls[i], ls[j] = ls[j], ls[i]
}

func main() {

	f, err := excelize.OpenFile("logs.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("Логи")
	if err != nil {
		fmt.Println(err)
		return
	}

	var data []string

	for i, row := range rows {
		if i == 0 {
			continue
		}
		if row[2] == "Unable to get integration token" || row[2] == "The wait operation timed out" {
			continue
		}
		data = append(data, row[2])
		if len(data) >= 10000 {
			break
		}
	}

	fmt.Println("data loaded")

	corp := makeCorpus(data)
	fmt.Println("corpus made")

	r := rand.New(rand.NewSource(1337))
	conf := dmmclust.Config{
		K:          10,                   // maximum 10 clusters expected
		Vocabulary: len(corp),            // simple example: the vocab is the same as the corpus size
		Iter:       1000,                 // iterate 100 times
		Alpha:      0.0001,               // smaller probability of joining an empty group
		Beta:       0.1,                  // higher probability of joining groups like me
		Score:      dmmclust.Algorithm4,  // use Algorithm3 to score
		Sampler:    dmmclust.NewGibbs(r), // use Gibbs to sample
	}
	var clustered []dmmclust.Cluster
	// Using Algorithm4, where repeat words are allowed
	docs := makeDocuments(data, corp, true)

	fmt.Println("docs made")

	if clustered, err = dmmclust.FindClusters(docs, conf); err != nil {
		fmt.Println(err)
	}

	stats := map[int]map[string]int{}

	fmt.Println("\nClusters (Algorithm4):")
	for i, clust := range clustered {
		if stats[clust.ID()] == nil {
			stats[clust.ID()] = map[string]int{}
		}
		stats[clust.ID()][data[i]]++
	}

	result := map[int]logLines{}

	for cid, c := range stats {
		for line, count := range c {
			result[cid] = append(result[cid], logLine{
				text:  line,
				count: count,
			})
		}
	}

	for _, clust := range clustered {
		sort.Sort(result[clust.ID()])
	}

	spew.Dump(result)
}
