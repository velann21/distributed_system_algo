package tf_idf

import (
	"math"
)

type TfIDF struct {

}

func (tf *TfIDF) CalculateTermFrequency(words []string, term string)float64{
	count := 0
	for _, v := range words{
		if v == term{
			count++
		}
	}                                                           //Just to take the precision
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
