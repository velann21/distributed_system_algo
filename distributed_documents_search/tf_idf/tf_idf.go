package tf_idf

import (
	"github.com/coordination-service/distributed_documents_search/helpers"
	"math"
	"sync"
)

type TfIDF struct {

}

func (tf *TfIDF) CalculateTermFrequency(words []string, term string)float64{

	count := 0
	for _, v := range words{
		if v == term{
			count++
		}
	}
	tfreque := math.Trunc(float64(count)/float64(len(words))*1000000000) / 1000000000
	return tfreque
}



func (tf *TfIDF) CalculateOnlyTF(words []string, terms []string)map[string]float64{
	termsData := make(map[string]float64, 0)
	for _, term := range terms{
		termfreq := tf.CalculateTermFrequency(words, term)
		termsData[term] = termfreq
	}
	return termsData
}


func (tf *TfIDF) CreateDocumentData(words []string, terms []string)helpers.DocumentData{
	docData := make(map[string]float32,0)
	data := helpers.DocumentData{
		DocumentData:docData,
		Mutex:sync.Mutex{},
	}
	for _, term := range terms{
		_ = tf.CalculateTermFrequency(words, term)
		//data.PutTermFrequency(term, termfreq)
	}
	return data
}

func (tf *TfIDF) GetIDF(term string, docResult map[string]helpers.DocumentData)float64{
	nt := 0
	for k, _ := range docResult{
		docData := docResult[k]
		tf := docData.GetTermFrequency(term)
		if tf > 0.0{
			nt++
		}
	}
	return math.Log10(float64(len(docResult)/nt))
}

