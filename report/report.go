package report

import (
    "io"
    "os"
    "fmt"
    "time"
    "strings"
    "strconv"
    "reflect"
    "io/ioutil"
    "archive/zip"
    "encoding/csv"
    "encoding/json"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/aws/aws-sdk-go/service/s3/s3manager"

    "github.com/oswell/aws-elk-reports/db"
    "github.com/oswell/aws-elk-reports/config"
)

type Report struct {
    FileName      string
    FileSize      int64
    LastModified  time.Time
    Config        config.Configuration
}

func (r Report) needsProcessing() (bool, error) {
    db := db.DB{ConnectionString: r.Config.DBUrl}
    process, err := db.ShouldProcess(&r.FileName, &r.FileSize)

    if err != nil {
        return true, fmt.Errorf("Error checking completion of report %s: %s", r.FileName, err)
    }
    return process, nil
}

func (r Report) Process() (error) {
    p, err := r.needsProcessing() ; if err != nil {
        fmt.Printf("Error: %v\n", err)
        return err
    }

    fmt.Printf("Processing report for %s? %v\n", r.FileName, p)
    r.Run()
    fmt.Printf("Finished processing report %s\n", r.FileName)

    db := db.DB{ConnectionString: r.Config.DBUrl}
    db.SaveReport(&r.FileName, &r.FileSize, &r.LastModified)

    return nil
}

// Run executes the main processing loop for each billing file
func (r Report) Run() error {
    err := r.handleFile()
	if err != nil {
		return err
	}

	return nil
}

func (r Report) handleFile() (error) {
	file, err := ioutil.TempFile("", "aws-bill") ; if err != nil {
		return err
	}

	defer os.Remove(file.Name())

    fmt.Printf("Downloading bill to %s\n", file.Name())

	downloader := s3manager.NewDownloader(session.New(&aws.Config{Region: aws.String("us-west-2")}))

	filename := strings.TrimPrefix(r.FileName, "/")
	input := &s3.GetObjectInput{
		Bucket: aws.String(r.Config.Bucket),
		Key:    aws.String(filename),
	}

	numBytes, err := downloader.Download(file, input) ; if err != nil {
		return err
	}

	err = r.uncompress(file.Name()); if err != nil {
		return err
	}

	fmt.Println("Downloaded file", file.Name(), numBytes, "bytes")

	return nil
}

// uncompress returns an uncompressed IO stream
// Returns the io.Reader, or an error if something went wrong.
func (r Report) uncompress(filename string) (error) {

    reader, err := zip.OpenReader(filename)
    if err != nil {
        return err
    }

	for _, file := range reader.File {
		fileReader, err := file.Open() ; if err != nil {
			return err
		}
		defer fileReader.Close()

		// 074822205801-aws-billing-detailed-line-items-with-resources-and-tags-2014-06.csv
		parts := strings.Split(file.Name, "-")
		index := fmt.Sprintf("%s.%s", parts[10], parts[11])

		r.parseFile(fileReader, index)
	}
	return nil
}

// ParseFile handles CSV parsing of the billing report
func (r Report) parseFile(reader io.Reader, indexDate string) error {
	rdr := csv.NewReader(reader)

	header, err := rdr.Read()
	if err != nil {
		return fmt.Errorf("Failed to read CSV record, %s", err.Error())
	}

    kafka := Kafka{ Config: r.Config.Kafka }
	for {
    	record, err := rdr.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			return fmt.Errorf("Failed to read CSV record, %s", err.Error())
		}

		hdrMap := r.headersToMap(header)
		jsonBytes, err := r.toJSON(hdrMap, record) ; if err != nil {
			return err
		}

        _ = kafka.Produce(*jsonBytes)
	}

	return nil
}

// processStruct handles processing for a single struct.
func (r Report) processStruct(headers map[string]int, record []string, obj interface{}) {
	recFields := reflect.ValueOf(obj).Elem()
	typeOfT := recFields.Type()

	for i := 0; i < recFields.NumField(); i++ {
		field := recFields.Field(i)


		if recKey, ok := headers[typeOfT.Field(i).Name] ; ok {
			value := record[recKey]

			switch field.Kind() {
			case reflect.String:
				field.SetString(value)

			case reflect.Float64:
				floatValue, err := strconv.ParseFloat(value, 64); if err != nil { }
				field.SetFloat(floatValue)

			case reflect.Bool:
				boolValue, err := strconv.ParseBool(value) ; if err != nil { }
				field.SetBool(boolValue)

			case reflect.Struct:
				form := "2006-01-02 15:04:00"
				t, _ := time.Parse(form, value)
				field.Set(reflect.ValueOf(t.Format(time.RFC3339)))
			}
		}
	}
}

func (r Report) headersToMap(headers []string) (map[string]int) {
	hdrMap := map[string]int{}
	for idx, hdr := range headers {
		hdrMap[hdr] = idx
	}
	return hdrMap
}

func (r Report) processTags(headers map[string]int, record []string, prefix string) (*map[string]string) {
	obj := map[string]string{}

	for key, idx := range headers {
		if strings.HasPrefix(key, "user:") {
			p := strings.Split(key, ":")
			keyName := strings.Replace(p[1], ".", "_", -1)
			obj[keyName] = record[idx]
		}
	}

	return &obj
}

// toJson converts the CSV record into a Map, then to JSON.
// Returns a string json representation of the record, or an error.
func (r Report) toJSON(headers map[string]int, record []string) (*[]byte, error) {

	var billing DetailedBilling
	r.processStruct(headers, record, &billing)

	billing.ResourceTags = *(r.processTags(headers, record, "resourceTags"))

	js, err :=  json.Marshal(&billing); if err != nil {
	     return nil, err
	}

	return &js, nil
}
