package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
	"strings"
	"sort"

	"github.com/oswell/aws-elk-reports/report"
	"github.com/oswell/aws-elk-reports/config"

	"github.com/Sirupsen/logrus"
	"github.com/koding/multiconfig"
	// "github.com/jasonlvhit/gocron"
)

// Configuration
var configuration *config.Configuration

// findReports crawls the specified S3 bucket for compressed billing report files.
// Each file found will be read and processed if it hasn't been processed already.
func findReports() error {
	sess, err := session.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %s,", err)
	}

	svc := s3.New(sess, aws.NewConfig().WithRegion("us-west-2"))

	reportSuffix := "csv.zip"
	reports := report.Reports{}

	var token *string
	for {
		params := &s3.ListObjectsV2Input{
			Bucket:            aws.String(configuration.Bucket),
			MaxKeys:           aws.Int64(128),
			Prefix:            aws.String(configuration.Prefix),
			ContinuationToken: token,
		}
		resp, err := svc.ListObjectsV2(params)
		if err != nil {
			return fmt.Errorf("Error fetching S3 objects: %s", err)
		}

		for _, s3obj := range resp.Contents {
			if strings.HasSuffix(*s3obj.Key, reportSuffix) {
				reports = reports.AddReport(report.Report{
					FileName    : *s3obj.Key,
					FileSize    : *s3obj.Size,
					LastModified: *s3obj.LastModified,
					Config      : *configuration,
				})
			}
		}

		token = resp.NextContinuationToken
		if token == nil {
			break
		}
	}

	fmt.Printf("Found %d reports.\n", len(reports))
	sort.Sort(reports)
	for _, rep := range reports {
		err := rep.Process() ; if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	mc := multiconfig.NewWithPath("config.toml")
	configuration = new(config.Configuration)
	if err := mc.Load(configuration); err != nil {
		logrus.Errorf("Failed to load configuration, %s", err)
		os.Exit(-1)
	}
	mc.MustLoad(configuration)

	// cron := gocron.NewScheduler()
	// cron.Every(1).Hours().Do(fetchLatestBills)
	// <- cron.Start()

	findReports()

}
