package utils

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func Unzip(zipPath string, dstPath string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		path := filepath.Join(dstPath, file.Name)

		if file.FileInfo().IsDir() {
			err = os.MkdirAll(path, 0755)
			if err != nil {
				return err
			}
			continue
		}

		err = os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			return err
		}
		in, err := file.Open()
		if err != nil {
			return err
		}

		out, err := os.Create(path)
		if err != nil {
			in.Close()
			return err
		}

		_, err = io.Copy(out, in)
		in.Close()
		out.Close()
		if err != nil {

			return err
		}

	}

	return nil
}
