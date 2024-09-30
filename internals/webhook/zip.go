package webhook

import (
    "archive/zip"
    "fmt"
    "io"
    "os"
    "path/filepath"
)

func extractZip(zipFilePath, folder string) error {
	r, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %v", err)
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(folder, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if filepath.Ext(fpath) == ".zip" {
			outFile, err := os.Create(fpath)
			if err != nil {
				return fmt.Errorf("failed to create nested zip file: %v", err)
			}

			rc, err := f.Open()
			if err != nil {
				return err
			}

			_, err = io.Copy(outFile, rc)
			outFile.Close()
			rc.Close()

			if err != nil {
				return fmt.Errorf("failed to save nested zip file: %v", err)
			}

			err = extractZip(fpath, folder)
			if err != nil {
				return fmt.Errorf("failed to extract nested zip file: %v", err)
			}
		} else {
			if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}

			rc, err := f.Open()
			if err != nil {
				return err
			}

			_, err = io.Copy(outFile, rc)
			outFile.Close()
			rc.Close()

			if err != nil {
				return err
			}
		}
	}

	fmt.Printf("Extracted log files from %s\n", zipFilePath)
	return nil
}
