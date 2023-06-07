package FileProcess

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

var extactlist = make(chan string)

func LoopFile() {
	os.Chdir("StudentContract")
	// currentDirectory, err := os.Getwd()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	filepath.Walk("C:/Users/admin/Desktop/ZOHO-GO/StudentContract/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf(err.Error())
		}

		extactlist <- info.Name()

		return nil
	})
	close(extactlist)
}
func CreateWorkerPool(noOfWorkers int) {
	var wg sync.WaitGroup
	for i := 0; i < noOfWorkers; i++ {
		wg.Add(1)
		go Extractor(&wg)
	}
	wg.Wait()

}
func MergePDFs(outputfile string, inFiles []string) {
	// inFiles := []string{"C:/Users/admin/Desktop/test/CA_QUOTE_3000150934236.3 (1).pdf", "C:/Users/admin/Desktop/test/CA_QUOTE_3000151156765.1 (3).pdf"}
	err := api.MergeCreateFile(inFiles, outputfile, nil)
	fmt.Print(err)

}

func Extractor(wg *sync.WaitGroup) {
	// os.Chdir("C:/Users/admin/Desktop/ZOHO-GO/")
	// Specify the path to the ZIP file you want to extract
	// zipFile := "C:/Users/admin/Desktop/test1.zip"
	for zipFile := range extactlist {
		// Open the ZIP file for reading
		// fmt.Print("the zip file path is ", zipFile)
		r, err := zip.OpenReader(zipFile)
		if err != nil {
			fmt.Println("Error opening ZIP file:", err)
			currentDirectory, _ := os.Getwd()
			fmt.Println("The current directory is ", currentDirectory)
			return
		}

		// Create a directory to extract the files
		// extractDir := "C:/Users/admin/Desktop/test"
		if err := os.MkdirAll("C:/Users/admin/Desktop/ZOHO-GO/StudentContract/test", 0755); err != nil {
			fmt.Println("Error creating extract directory:", err)
			currentDirectory, _ := os.Getwd()
			fmt.Println("The current directory is ", currentDirectory)
			return
		}
		var fileList []string
		// Extract each file in the ZIP archive
		for _, file := range r.File {
			// Open the file inside the ZIP archive
			rc, err := file.Open()
			if err != nil {
				fmt.Println("Error opening file inside ZIP:", err)
				return
			}

			// Create the corresponding file on disk
			path := file.Name
			f, err := os.Create(path)
			if err != nil {
				fmt.Println("Error creating file:", err)
				return
			}
			fileList = append(fileList, path)

			// Copy the contents of the file from the ZIP to the disk
			_, err = io.Copy(f, rc)
			if err != nil {
				fmt.Println("Error extracting file:", err)
				return
			}
			rc.Close()
			f.Close()
			fmt.Println("Extracted:", file.Name)

		}
		r.Close()
		filename := strings.Split(zipFile, ".")[0] + ".pdf"
		MergePDFs(filename, fileList)

		fmt.Println("Extraction completed.")
	}
	wg.Done()

}
