package generator

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/klauspost/compress/zstd"
)

type compositeReadCloser struct {
	reader io.Reader
	closer func() error
}

func (reader compositeReadCloser) Read(buffer []byte) (int, error) {
	return reader.reader.Read(buffer)
}

func (reader compositeReadCloser) Close() error {
	return reader.closer()
}

func openWordIndexReader(path string) (io.ReadCloser, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open word index %s: %w", path, err)
	}

	lowerPath := strings.ToLower(path)
	if strings.HasSuffix(lowerPath, ".zst") {
		zstdReader, err := zstd.NewReader(file)
		if err != nil {
			_ = file.Close()
			return nil, fmt.Errorf("open zstd word index %s: %w", path, err)
		}

		return compositeReadCloser{
			reader: zstdReader,
			closer: func() error {
				zstdReader.Close()
				return file.Close()
			},
		}, nil
	}

	if !strings.HasSuffix(lowerPath, ".gz") {
		return file, nil
	}

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		_ = file.Close()
		return nil, fmt.Errorf("open gzip word index %s: %w", path, err)
	}

	return compositeReadCloser{
		reader: gzipReader,
		closer: func() error {
			gzipErr := gzipReader.Close()
			fileErr := file.Close()
			if gzipErr != nil {
				return gzipErr
			}
			return fileErr
		},
	}, nil
}
