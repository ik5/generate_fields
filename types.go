package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"mime"
	"mime/multipart"
	"net/url"

	"github.com/google/go-querystring/query"
)

type InfoToReport struct {
	ErrorCodes    string `json:"error_codes" xml:"error_codes" url:"error_codes"`
	Provider      string `json:"provider" xml:"provider" url:"provider"`
	PhoneNumber   string `json:"phone_number" xml:"phone_number" url:"phone_number"`
	Identifier    string `json:"identifier" xml:"identifier" url:"identifier"`
	ReportingDate string `json:"reporting_date" xml:"reporting_date" url:"reporting_date"`
}

// JSON returns json based string from struct
func (info InfoToReport) JSON() ([]byte, error) {
	result, err := json.Marshal(info)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return result, nil
}

// XML returns XML base string from struct, and if withHeader is true, than <?xml ...?> header is returned without new
// line.
func (info InfoToReport) XML(withHeader bool) ([]byte, error) {
	result, err := xml.Marshal(info)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	header := []byte{}
	if withHeader {
		header = []byte(xml.Header[0 : len(xml.Header)-1])
	}

	return append(header, result...), nil
}

// Values returns url.Values data type out of struct.
func (info InfoToReport) Values() (url.Values, error) {
	result, err := query.Values(info)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return result, nil
}

// EncodedValues returns encoded EncodedValues string based on struct.
func (info InfoToReport) EncodedValues() ([]byte, error) {
	result, err := info.Values()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return []byte(result.Encode()), nil
}

// Multipart returns multipart string based on struct
func (info InfoToReport) Multipart() ([]byte, error) {
	urlValues, err := info.Values()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	buf := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(buf)

	for name, values := range urlValues {
		switch len(values) {
		case 0:
			err = writer.WriteField(name, "")
		case 1:
			err = writer.WriteField(name, values[0])
		default:
			for _, value := range values {
				err = writer.WriteField(name, value)
				if err != nil {
					return nil, fmt.Errorf("%w", err)
				}
			}

			continue
		}

		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
	}

	return buf.Bytes(), nil
}

// ByContentType returns a string payload based on supported content type.
func (info InfoToReport) ByContentType(contentType string) ([]byte, error) {
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	var result []byte

	switch mediaType {
	case JsonContentType:
		result, err = info.JSON()
	case XmlContentType:
		result, err = info.XML(true)
	case MultipartContentType:
		result, err = info.Multipart()
	case UrlEncodedContentType:
		result, err = info.EncodedValues()
	default:
		return nil, errors.New("unsupported content-type")
	}

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return result, nil
}
