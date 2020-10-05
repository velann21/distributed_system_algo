package main

import (
	"fmt"
	"github.com/coordination-service/distributed_documents_search/helpers"
	"github.com/coordination-service/distributed_documents_search/tf_idf"
	"strings"
)

func main() {
	fo := helpers.FileOperations{Chunksize:1024}

	documents := map[string][]string{}
	dirs, err := fo.IOReadDir("/Users/singaravelannandakumar/github/src/github.com/coordination-service/distributed_documents_search/docs")
	if err != nil{
		fmt.Println(err)
	}
	fo.CreateDocumentToTokens(dirs, documents)
	tfidf := tf_idf.TfIDF{}
	terms := "eBook Swift"

	docScores := []map[string]map[string]float64{}
	for k, v := range documents{
		docScore := tfidf.CalculateOnlyTF(v, strings.Split(terms, " "))
		score := map[string]map[string]float64{}
		score[k] = docScore
		docScores = append(docScores, score)

	}

	finalScore := map[string]float64{}
	for _, termsScores := range docScores{
		for valueKey, valueMap := range termsScores {
			cummulativeScore := 0.0
			for _, v := range valueMap{
				cummulativeScore += v
			}
			finalScore[valueKey] = cummulativeScore
		}
	}


}
