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
var missingList []string

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

func MergePDFs() {
	inFiles := []string{"C:/Users/admin/Desktop/ZOHO-GO/StudentContract/test/1.pdf",
		"C:/Users/admin/Desktop/ZOHO-GO/StudentContract/test/2.pdf",
		"C:/Users/admin/Desktop/ZOHO-GO/StudentContract/test/3.pdf",
		"C:/Users/admin/Desktop/ZOHO-GO/StudentContract/test/4.pdf",
		"C:/Users/admin/Desktop/ZOHO-GO/StudentContract/test/5.pdf",
		"C:/Users/admin/Desktop/ZOHO-GO/StudentContract/test/6.pdf",
		"C:/Users/admin/Desktop/ZOHO-GO/StudentContract/test/7.pdf"}
	err := api.MergeCreateFile(inFiles, "Test.pdf", nil)
	fmt.Print("Error happend is ", err)

}

func Extractor(wg *sync.WaitGroup) {

	for zipFile := range extactlist {

		r, err := zip.OpenReader(zipFile)
		if err != nil {
			fmt.Println("Error opening ZIP file:", err)
			return
		}

		// Create a directory to extract the files
		extractDir := "C:/Users/admin/Desktop/ZOHO-GO/StudentContract/" + strings.Split(zipFile, ".")[0]
		if err := os.MkdirAll(extractDir, 0755); err != nil {
			fmt.Println("Error creating extract directory:", err)
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
			path := extractDir + "/" + file.Name
			// fmt.Print("The extracted file path is", path)
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
			// fmt.Println("Extracted:", file.Name)

		}
		r.Close()
		filename := strings.Split(zipFile, ".")[0] + ".pdf"
		// MergePDFs(filename, fileList)
		tt := api.MergeCreateFile(fileList, filename, nil)
		if tt != nil {
			fmt.Print(tt)
			// missingList = append(missingList, extractDir)
		} else {
			fmt.Println("Extraction completed.")
		}

	}
	wg.Done()
	fmt.Print(missingList)
}
