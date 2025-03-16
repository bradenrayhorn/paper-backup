package compress

import (
	"bytes"
	"compress/flate"
	"errors"
	"io"
)

const prefixCompressed = byte('1')
const prefixUncompressed = byte('2')

func Compress(data []byte) ([]byte, error) {
	// compress bytes
	var buf bytes.Buffer
	writer, err := flate.NewWriter(&buf, flate.BestCompression)
	if err != nil {
		return nil, err
	}

	_, err = writer.Write(data)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	compressedData := buf.Bytes()

	// is compression actually helpful?
	if len(compressedData) < len(data) {
		return append([]byte{prefixCompressed}, compressedData...), nil
	}

	return append([]byte{prefixUncompressed}, data...), nil
}

func Decompress(data []byte) ([]byte, error) {
	if len(data) < 1 {
		return nil, errors.New("data too short")
	}

	prefix := data[0]
	content := data[1:]

	switch prefix {
	case prefixUncompressed:
		return content, nil

	case prefixCompressed:
		reader := flate.NewReader(bytes.NewReader(content))
		defer func() { _ = reader.Close() }()

		var buf bytes.Buffer
		_, err := io.Copy(&buf, reader)
		if err != nil {
			return nil, err
		}

		return buf.Bytes(), nil

	default:
		return nil, errors.New("unknown compression prefix")
	}

}
