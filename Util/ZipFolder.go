/*
	UTILITY PACKAGE TO ZIP THE GIVEN FOLDER TO THE GIVEN FILE NAME(with location).
*/

package Util

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Zipit(source, target string) error {

	// create target dir
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	// check file information
	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	// check if file is a directory/folder
	// if yes, then set the base addr as the folder name
	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	// go through every file in the dir
	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// get header of the file
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		// if it is a folder then only add "/", else add relative path of the file.
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		// create header for the created zip file.
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

	return err
}
