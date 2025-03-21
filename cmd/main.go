package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ibcode15/BetterFs/internal/Tokenizer"
)

type path = string
type TermName = string
type TermCount = uint32
type TermMap = map[TermName]TermCount
type DocMap = map[path]Document

var print = fmt.Printf

type Index struct {
	documents DocMap
}

type Document struct {
	TermFreq TermMap
}

type Optional[T any] struct {
	Data    T
	hasData bool
}

func existsAsFolder(path string) bool {
	stat, err := os.Stat(path)
	if err == nil && stat != nil {
		if !stat.IsDir() {
			print("cannot find %+v as dir", path)
			return false
		}
		return true
	}

	return false
}

func readFile(fileName string) Optional[string] {
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		return Optional[string]{hasData: false}
	}
	return Optional[string]{hasData: true, Data: string(fileContent)}
}

func checkForSuffers(path string, suffixes []string) Optional[string] {
	for _, s := range suffixes {
		if strings.HasSuffix(path, s) {
			return Optional[string]{Data: s, hasData: true}
		}
	}
	return Optional[string]{hasData: false}

}

func indexPath(path string, termfreq *TermMap, absoluteRoot string, is_file bool) {
	relaviePath := strings.TrimPrefix(path, absoluteRoot)

	tokenizer := Tokenizer.CreateTokenizer(relaviePath)
	for token := range tokenizer.IterateTokens() {
		_, ok := (*termfreq)[token]
		if token == "\\" {
			continue
		}
		if ok {
			(*termfreq)[token] += 1
		} else {
			(*termfreq)[token] = 1
		}
	}

	fileData := Optional[string]{hasData: false}
	suffex := ""
	if is_file {
		suffix := checkForSuffers(path, []string{"txt"})
		if suffix.hasData {
			fileData = readFile(path)
			suffex = suffix.Data
		}
	}
	if fileData.hasData {
		if suffex == "txt" {
			tokenizer = Tokenizer.CreateTokenizer(fileData.Data)
			for token := range tokenizer.IterateTokens() {
				_, ok := (*termfreq)[token]
				if ok {
					(*termfreq)[token] += 1
				} else {
					(*termfreq)[token] = 1
				}
			}
		}

	}
}

func indexing(path string, index *Index, absoluteRoot string) {
	stat, _ := os.Stat(path)

	if stat == nil {
		return
	}

	if stat.IsDir() {
		entries, _ := os.ReadDir(path)

		for _, e := range entries {
			fileinfo, _ := e.Info()
			if fileinfo == nil {
				continue
			}
			absolute_path := filepath.Join(path, e.Name())
			indexing(absolute_path, index, absoluteRoot)
		}

		if path != absoluteRoot {
			index.documents[path] = Document{TermFreq: make(TermMap)}
			TermFreq := index.documents[path].TermFreq
			indexPath(path, &TermFreq, absoluteRoot, false)
		}

	} else { //is file

		index.documents[path] = Document{TermFreq: make(TermMap)}
		TermFreq := index.documents[path].TermFreq
		indexPath(path, &TermFreq, absoluteRoot, true)
	}

}
func IndexRoot(rootDir string, index *Index) {

	exists := existsAsFolder(rootDir)

	if !exists {
		return
	}
	print("found folder")
	indexing(rootDir, index, rootDir)
}

// original version
func CalulateTF_OG(term string, termfreq *TermMap) float32 {
	sum_of_terms := uint32(0)
	for _, v := range *termfreq {
		sum_of_terms += v
	}
	termfreq_in_map, has_term := (*termfreq)[term]
	if !has_term {
		termfreq_in_map = 1
	}

	return float32(termfreq_in_map) / float32(sum_of_terms)
}

// new version made for filesystem as smaller paths will always get ranked higher
func CalulateTF(term string, termfreq *TermMap) float32 {
	sum_of_terms := uint32(0)
	for _, v := range *termfreq {
		sum_of_terms += v
	}
	termfreq_in_map, has_term := (*termfreq)[term]
	if !has_term {
		termfreq_in_map = 1
	} else if len(term) > 1 {
		termfreq_in_map = termfreq_in_map * 60

	}

	return float32(termfreq_in_map) / float32(sum_of_terms)
}

func CalulateIDF(term string, docs *DocMap) float32 {
	number_of_documents := float32(len((*docs)))
	count_of_term_in_documents := 0

	for _, item := range *docs {
		_, has_term := item.TermFreq[term]
		if has_term {
			count_of_term_in_documents += 1
		}
	}
	if count_of_term_in_documents == 0 {
		count_of_term_in_documents = 1
	}
	return float32(math.Log10(float64(number_of_documents / float32(count_of_term_in_documents))))

}

type Pair struct {
	filename string
	rank     float32
}

func ComputeResult(term string, index *Index) []Pair {
	tokenizer := (Tokenizer.CreateTokenizer(term))
	tokens := tokenizer.ToArray()
	documents := &((*index).documents)

	result := []Pair{}
	for filename, doc := range *documents {

		rank := float32(0)
		for _, token := range tokens {

			TF := CalulateTF(token, &doc.TermFreq)
			IDF := CalulateIDF(token, documents)
			// print("token = %+v IDF = %+v TF = %+v filename = %+v\n", token, IDF, TF, filename)
			rank += TF * IDF
		}
		result = append(result, Pair{filename, rank})
		// print("%+v %+v\n", rank, filename)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].rank > result[j].rank
	})
	for _, p := range result {
		fmt.Printf("%+v %+v\n", p.filename, p.rank)
	}
	return result

}

func main() {

}
