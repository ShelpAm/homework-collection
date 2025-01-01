package main

import (
	"io"
	"os"
)

type OnProgressUpdate func(progress float64)

type ProgressWriter struct {
	Writer   io.Writer
	Total    int64
	Written  int64
	OnUpdate OnProgressUpdate
}

func (pw *ProgressWriter) Write(p []byte) (n int, err error) {
	n, err = pw.Writer.Write(p)
	pw.Written += int64(n)

	if pw != nil {
		progress := float64(pw.Written) / float64(pw.Total)
		pw.OnUpdate(progress)
	}

	return n, err
}

func writeToFileWithProgress(src SizeableReader, dstPath string, callback OnProgressUpdate) error {
	// Open the destination file
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Create a progress writer with the provided callback
	pw := &ProgressWriter{
		Writer:   dstFile,  // The actual writer (destination file)
		Total:    src.Size, // The size of the source file
		OnUpdate: callback, // User-defined callback function
	}

	// Copy the file with progress tracking
	_, err = io.Copy(pw, &src)
	if err != nil {
		return err
	}

	return nil
}
