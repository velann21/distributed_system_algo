package helpers

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type FileOperations struct {
	data  *os.File
	part  []byte
	err   error
	count int
	Chunksize int
}


func (fo *FileOperations) CreateDocumentToTokens(dirs []string, documents map[string][]string){
	for _, v := range dirs{
		_, size := fo.openFile("/Users/singaravelannandakumar/github/src/github.com/coordination-service/distributed_documents_search/docs/"+v)
		words, err := fo.Tokens(size)
		if err != nil{
			continue
		}
		documents[v] = words
	}
}

func (fo *FileOperations)  openFile(filename string) (byteCount int, buffer *bytes.Buffer) {

	fo.data, fo.err = os.Open(filename)
	if fo.err != nil {
		log.Fatal(fo.err)
	}
	defer fo.data.Close()

	reader := bufio.NewReader(fo.data)
	buffer = bytes.NewBuffer(make([]byte, 0))
	fo.part = make([]byte, fo.Chunksize)

	for {
		if fo.count, fo.err = reader.Read(fo.part); fo.err != nil {
			break
		}
		buffer.Write(fo.part[:fo.count])
	}
	if fo.err != io.EOF {
		log.Fatal("Error Reading ", filename, ": ", fo.err)
	} else {
		fo.err = nil
	}

	byteCount = buffer.Len()
	return
}


func (fo *FileOperations)  IOReadDir(root string) ([]string, error) {
	var files []string
	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return files, err
	}
	for _, file := range fileInfo {
		files = append(files, file.Name())
	}
	return files, nil
}

func (fo *FileOperations) Tokens(buffer *bytes.Buffer)([]string,error){
	words := []string{}
	scanner := bufio.NewScanner(buffer)
	for scanner.Scan() {
		tokens := fo.BreakLineIntoToken(scanner.Text())
		for _, v := range tokens{
			words = append(words, v)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return words, nil

}

func (fo *FileOperations) BreakLineIntoToken(line string)[]string{
	tokens := fo.SplitAny(line, "(\\.)+|(,)+|( )+|(-)+|(\\?)+|(!)+|(;)+|(:)+([)+(])+(***)")
	return tokens
}

func (fo *FileOperations) SplitAny(s string, seps string) []string {
	splitter := func(r rune) bool {
		return strings.ContainsRune(seps, r)
	}
	return strings.FieldsFunc(s, splitter)
}
