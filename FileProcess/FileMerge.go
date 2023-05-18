package FileProcess

import (
	"archive/zip"
	"fmt"
	"io"
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func MergePDFs() {
	inFiles := []string{"C:/Users/admin/Desktop/test/CA_QUOTE_3000150934236.3 (1).pdf", "C:/Users/admin/Desktop/test/CA_QUOTE_3000151156765.1 (3).pdf"}
	err := api.MergeCreateFile(inFiles, "out.pdf", nil)
	fmt.Print(err)

}

func Extractor() {
	// Specify the path to the ZIP file you want to extract
	zipFile := "C:/Users/admin/Desktop/test1.zip"

	// Open the ZIP file for reading
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		fmt.Println("Error opening ZIP file:", err)
		return
	}
	defer r.Close()

	// Create a directory to extract the files
	extractDir := "C:/Users/admin/Desktop/test"
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		fmt.Println("Error creating extract directory:", err)
		return
	}

	// Extract each file in the ZIP archive
	for _, file := range r.File {
		// Open the file inside the ZIP archive
		rc, err := file.Open()
		if err != nil {
			fmt.Println("Error opening file inside ZIP:", err)
			return
		}
		defer rc.Close()

		// Create the corresponding file on disk
		path := extractDir + "/" + file.Name
		f, err := os.Create(path)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer f.Close()

		// Copy the contents of the file from the ZIP to the disk
		_, err = io.Copy(f, rc)
		if err != nil {
			fmt.Println("Error extracting file:", err)
			return
		}

		fmt.Println("Extracted:", file.Name)
	}

	fmt.Println("Extraction completed.")
}
