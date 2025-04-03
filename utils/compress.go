package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func ZipFolder(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	sourceInfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	var baseDir string
	if sourceInfo.IsDir() {
		baseDir = filepath.Base(source)
	}

	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			relPath, err := filepath.Rel(source, path)
			if err != nil {
				return err
			}
			if relPath == "." {
				return nil
			}
			header.Name = filepath.Join(baseDir, relPath)
		}

		header.Method = zip.Deflate

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Name = strings.ReplaceAll(header.Name, "\\", "/")
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		return err
	})
}

func UnzipFile(source, target string) error {
	reader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		err := extractFile(file, target)
		if err != nil {
			return err
		}
	}

	return nil
}

func extractFile(file *zip.File, target string) error {
	path := filepath.Join(target, file.Name)

	if !strings.HasPrefix(path, filepath.Clean(target)+string(os.PathSeparator)) {
		return fmt.Errorf("illegal path: %s", path)
	}

	if file.FileInfo().IsDir() {
		return os.MkdirAll(path, file.Mode())
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	dstFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	srcFile, err := file.Open()
	if err != nil {
		return err
	}
	defer srcFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
