package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unsafe"
)

//func main() {
//	fo := helpers.FileOperations{Chunksize:1024}
//
//	documents := map[string][]string{}
//	dirs, err := fo.IOReadDir("/Users/singaravelannandakumar/github/src/github.com/coordination-service/distributed_documents_search/docs")
//	if err != nil{
//		fmt.Println(err)
//	}
//	fo.CreateDocumentToTokens(dirs, documents)
//	tfidf := tf_idf.TfIDF{}
//	terms := "eBook Swift"
//
//	docScores := []map[string]map[string]float64{}
//	for k, v := range documents{
//		docScore := tfidf.CalculateOnlyTF(v, strings.Split(terms, " "))
//		score := map[string]map[string]float64{}
//		score[k] = docScore
//		docScores = append(docScores, score)
//
//	}
//
//	finalScore := map[string]float64{}
//	for _, termsScores := range docScores{
//		for valueKey, valueMap := range termsScores {
//			cummulativeScore := 0.0
//			for _, v := range valueMap{
//				cummulativeScore += v
//			}
//			finalScore[valueKey] = cummulativeScore
//		}
//	}
//
//
//}

func main(){
	// integer for convert
	num := 1.2344
	fmt.Println("Original number:", num)

	// integer to byte array
	byteArr := float64ToByte(num)
	fmt.Println("Byte Array", byteArr)

	buf := bytes.NewReader(byteArr)
	err := binary.Read(buf, binary.LittleEndian, &flty)

	//// byte array to integer again
	//numAgain := ByteArrayToInt(byteArr)
	//fmt.Println("Converted number:", numAgain)
}

func FloatToByteArray(num float64)[]byte{
	size := int(unsafe.Sizeof(num))
	arr := make([]byte, size)
	for i := 0 ; i < size ; i++ {
		byt := *(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&num)) + uintptr(i)))
		arr[i] = byt
	}
	return arr
}

func IntToByteArray(num int64) []byte {
	size := int(unsafe.Sizeof(num))
	arr := make([]byte, size)
	for i := 0 ; i < size ; i++ {
		byt := *(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&num)) + uintptr(i)))
		arr[i] = byt
	}
	return arr
}

func float64ToByte(f float64) []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, f)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return buf.Bytes()
}
func ByteArrayToInt(arr []byte) int64{
	val := int64(0)
	size := len(arr)
	for i := 0 ; i < size ; i++ {
		*(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&val)) + uintptr(i))) = arr[i]
	}
	return val
}