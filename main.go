package main

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"encoding/base64"
	"fmt"
	"github.com/klauspost/compress/zstd"
	"io/ioutil"
	"os"
)

func main() {
	sourceFile, err := os.Open("data.json")
	if err != nil {
		fmt.Println(err)
	}
	defer sourceFile.Close()

	fullContent, err := ioutil.ReadAll(sourceFile)
	if err != nil {
		fmt.Println(err)
	}

	err = doGzip(fullContent)
	if err != nil {
		fmt.Println(err)
	}

	err = doZstd(fullContent)
	if err != nil {
		fmt.Println(err)
	}

	err = doFlate(fullContent)
	if err != nil {
		fmt.Println(err)
	}

	err = doZlib(fullContent)
	if err != nil {
		fmt.Println(err)
	}

	doFullTurnAround(fullContent)

}

func doGzip(buff []byte) error {
	var b bytes.Buffer
	destWriter := gzip.NewWriter(&b)

	_, err := destWriter.Write(buff)
	if err != nil {
		return err
	}
	destWriter.Close()

	fmt.Printf("gzip size=%d, comp=%f \n", b.Len(), (float32(b.Len())*100)/float32(len(buff)))
	return nil
}

func doZstd(buff []byte) error {
	var b bytes.Buffer
	destWriter, err := zstd.NewWriter(&b)
	if err != nil {
		return err
	}

	_, err = destWriter.Write(buff)
	if err != nil {
		return err
	}
	destWriter.Close()

	fmt.Printf("zstd size=%d, comp=%f \n", b.Len(), (float32(b.Len())*100)/float32(len(buff)))
	return nil
}

func doFlate(buff []byte) error {
	var b bytes.Buffer
	destWriter, err := flate.NewWriter(&b, flate.BestCompression)
	if err != nil {
		return err
	}

	_, err = destWriter.Write(buff)
	if err != nil {
		return err
	}
	destWriter.Close()

	fmt.Printf("flate size=%d, comp=%f \n", b.Len(), (float32(b.Len())*100)/float32(len(buff)))
	return nil
}

func doZlib(buff []byte) error {
	var b bytes.Buffer
	destWriter, err := zlib.NewWriterLevel(&b, zlib.BestCompression)
	if err != nil {
		return err
	}

	_, err = destWriter.Write(buff)
	if err != nil {
		return err
	}
	destWriter.Close()

	fmt.Printf("zlib size=%d, comp=%f \n", b.Len(), (float32(b.Len())*100)/float32(len(buff)))
	return nil
}

func doFullTurnAround(buff []byte) error {
	// Compress
	var b bytes.Buffer
	destWriter, err := zstd.NewWriter(&b)
	if err != nil {
		return err
	}

	_, err = destWriter.Write(buff)
	if err != nil {
		return err
	}
	destWriter.Close()

	// Encode
	b64Encoded := base64.StdEncoding.EncodeToString(b.Bytes())

	// Decode
	b64Decoded, err := base64.StdEncoding.DecodeString(b64Encoded)
	if err != nil {
		return err
	}

	// Uncompress
	sourceReader, err := zstd.NewReader(bytes.NewReader(b64Decoded))
	if err != nil {
		return err
	}
	defer sourceReader.Close()

	uncompressed, err := ioutil.ReadAll(sourceReader)
	if err != nil {
		return err
	}

	fmt.Printf("document size=%d, compresed size=%d, encoded size=%d, decoded size=%d, uncompressed size=%d \n", len(buff), b.Len(), len(b64Encoded), len(b64Decoded), len(uncompressed))
	return nil
}
