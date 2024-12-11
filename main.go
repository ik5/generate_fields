package main

import (
	"fmt"
	"strings"
	"time"
)

func printPayloads(info InfoToReport) {
	fmt.Println("Testing by Payload")

	jsonPayload, err := info.JSON()
	if err != nil {
		panic(err)
	}

	xmlPayload, err := info.XML(true)
	if err != nil {
		panic(err)
	}

	xmlPayload2, err := info.XML(false)
	if err != nil {
		panic(err)
	}

	urlPayload, err := info.EncodedValues()
	if err != nil {
		panic(err)
	}

	multipartPayload, err := info.Multipart()
	if err != nil {
		panic(err)
	}

	payloads := map[string][]byte{
		"json":           jsonPayload,
		"xml":            xmlPayload,
		"xml no header":  xmlPayload2,
		"encoded values": urlPayload,
		"multipart":      multipartPayload,
	}

	writer := strings.Builder{}
	for name, encoding := range payloads {
		writer.WriteString("--------------------------")
		writer.WriteString(" Payload: ")
		writer.WriteString(name)
		writer.WriteString("\n\t")
		writer.Write(encoding)
		writer.WriteString("\n")
		writer.WriteString("--------------------------\n\n")
	}

	fmt.Println(writer.String())
}

func printByContentType(info InfoToReport) {
	writer := strings.Builder{}

	fmt.Println("Testing by content-type")

	contentTypes := [...]string{JsonContentType, XmlContentType, MultipartContentType, UrlEncodedContentType}

	for _, contentType := range contentTypes {
		writer.WriteString("--------------------------")
		writer.WriteString(" contentType: ")
		writer.WriteString(contentType)
		writer.WriteString("\n")

		payload, err := info.ByContentType(contentType)
		if err != nil {
			writer.WriteString("error >>> ")
			writer.WriteString(err.Error())
			writer.WriteString("\n")
			writer.WriteString("--------------------------\n\n")
			continue
		}

		writer.Write(payload)
		writer.WriteString("\n")
		writer.WriteString("--------------------------\n\n")
	}

	fmt.Println(writer.String())
}

func main() {
	info := InfoToReport{
		ErrorCodes:    "000",
		Provider:      "Bla Bla Bla",
		PhoneNumber:   "1234567890",
		Identifier:    "foo-bar-baz",
		ReportingDate: time.Now().Local().Format("2006-01-02 03:04:05 Z07:00"),
	}

	fmt.Printf("info: \n\t%+#v\n\n", info)

	printPayloads(info)
	printByContentType(info)
}
