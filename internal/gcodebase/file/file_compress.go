package file

import (
	"archive/zip"
	"fmt"
	"path/filepath"
	"os"
	"io"
)

func ZipCompress(paths []string, zipFile string) error {
	var err error
	if !Exist(filepath.Dir(zipFile)) {
		if err = MkDirAll(filepath.Dir(zipFile)); err != nil {
			return err
		}
	}
	zf,err := os.Create(zipFile)
	if err != nil {
		return err
	}
	defer zf.Close()
	zw := zip.NewWriter(zf)
	defer zw.Close()
	for _,f := range paths {
		if Exist(f) {
			return fmt.Errorf("%s NOT exist", f)
		}
		err = zipCompress(f, "", zw)
		if err != nil {
			return err
		}
	}
	return nil
}

func zipCompress(filePath string, prefix string, zw *zip.Writer) error {
	var err error
	file,err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	if fileInfo.IsDir() {
		prefix = filepath.Join(prefix, fileInfo.Name())
		subFilesInfo,err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range subFilesInfo {
			err = zipCompress(filepath.Join(filePath, fi.Name()), prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(fileInfo)
		header.Name = filepath.Join(prefix, header.Name)
		if err != nil {
			return err
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		if err != nil {
			return err
		}
	}
	return nil
}

func ZipDecompress(zipFile, targetDir string) error {
	var err error
	reader,err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()
	if Exist(targetDir) && !IsDir(targetDir) {
		return fmt.Errorf("%s Specified Must be a Dir", targetDir)
	}
	for _,f := range reader.File {
		if f.FileInfo().IsDir() {
			if err = MkDirAll(filepath.Join(targetDir, f.Name)); err != nil {
				return err
			}
			continue
		}
		srcFile,err := f.Open()
		if err != nil {
			return err
		}
		defer srcFile.Close()
		targetFilePath := filepath.Join(targetDir, f.Name)
		if !Exist(filepath.Dir(targetFilePath)) {
			err = MkDirAll(filepath.Dir(targetFilePath))
			if err != nil {
				return err
			}
		}
		targetFile,err := os.Create(targetFilePath)
		if err != nil {
			return err
		}
		defer targetFile.Close()
		io.Copy(targetFile, srcFile)
	}
	return nil
}