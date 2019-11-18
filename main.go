//
//	Created by Nick on 18-11-2019
//
//	Main.go
//	CSV to VCF converter, my first golang program

package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func getVCFDataFrom(name string, mobile *string) string {
	// Check cases where it might be possible that the mobile number is blank
	if *mobile == "NA" || *mobile == "" {
		*mobile = ""
		fmt.Printf("Ignored empty mobile number for: %v\n", name)
	} else if strings.Contains(*mobile, "/") { // Check cases where there might be more than 1 mobile number, delimited by a slash

		// Slice up the mobile and use the first elem
		fmt.Printf("Found 2 mobile numbers for %v, using the first one\n", name)
		*mobile = strings.Split(*mobile, "/")[0]
	}

	return fmt.Sprintf("BEGIN:VCARD\nVERSION:3.0\nREV:2019-11-18 12:00:07\nTEL;TYPE=cell,voice:%v\nN:;%v;;;;\nFN:%v\nEND:VCARD\n", *mobile, name, name)
}

//noinspection GoNilness
func main() {
	csvFile, err := os.Open("file.csv")
	defer csvFile.Close()

	// Check for errors
	if err != nil {
		log.Fatal(err)
	}

	// Create a new reader for CSV
	reader := csv.NewReader(bufio.NewReader(csvFile))

	// An array of strings that will hold our data, we will then write this using a buffered writer
	var contacts []string

	for {
		line, err := reader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		contacts = append(contacts, getVCFDataFrom(line[0], &line[1]))
	}

	// Create a file in the pwd
	file, err := os.Create("contact.vcf")
	defer file.Close()

	// Check error
	if err != nil {
		log.Fatal(err)
	}

	// Create a writer
	writer := bufio.NewWriter(file)

	// Loop and write the slice to the file
	for _, contact := range contacts {
		 _, _ = writer.WriteString(contact)
	}

	// Sync FS
	_ = writer.Flush()
	
	// Success
	fmt.Println("Successfully created the VCF file named: contacts.vcf")
}
