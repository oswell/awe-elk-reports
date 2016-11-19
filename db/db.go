package db

import (
    "fmt"
    "time"
    "strings"

    "github.com/jinzhu/gorm"

    // Required for use with the GORM library.
    _ "github.com/go-sql-driver/mysql"
)

// DB is used for handling all DB related functionality
type DB struct {
    ConnectionString string
    conn *gorm.DB
}

// Report describes a single report entry
type Report struct {
    gorm.Model

    FileName      *string
    Size          *int64
    LastModified  *time.Time
}

const (
    //parseTimeSuffix is a required suffix for the MySQL connection string.
    parseTimeSuffix = "?parseTime=true"
)

func (db *DB) getConnection() (*gorm.DB, error) {

    if db.conn != nil {
        return db.conn, nil
    }

    // The connection string *must* contain "?parseTime=true" or it won't work.  Rather than forcing the
    // user to ensure they have it, we'll just add if it's missing.
    if strings.Index(db.ConnectionString, parseTimeSuffix) == -1 {
        db.ConnectionString = fmt.Sprintf("%s%s", db.ConnectionString, parseTimeSuffix)
    }

    var err error
    db.conn, err = gorm.Open("mysql", db.ConnectionString) ; if err != nil {
      return nil, fmt.Errorf("failed to connect database, %s", err)
    }

    db.conn.LogMode(false)
    db.conn.AutoMigrate(&Report{})

    return db.conn, nil
}

// ShouldProcess returns true if either the report doesn't exist in the database or
// the current file size is larger than what is in the database.
func (db *DB) ShouldProcess(filename *string, size *int64) (bool, error) {
    report, err := db.FindReport(filename) ; if err != nil {
        return false, err
    }

    if report == nil || *report.Size < *size {
        return true, nil
    }

    return false, nil
}

// FindReport returns a Report object from MySQL if a matching item is found
// for the given filename.
func (db *DB) FindReport(filename *string) (*Report, error) {
    conn, err := db.getConnection() ; if err != nil {
        return nil, err
    }

    report := Report{}
    conn.Where(&Report{FileName: filename}).First(&report)

    if report.ID == 0 {
        return nil, nil
    }

    return &report, nil
}

// SaveReport saves the current report to the MySQL database
// Returns an error if something goes wrong.
func (db *DB) SaveReport(fileName *string, size *int64, lastModified *time.Time) (error) {
    conn, err := db.getConnection() ; if err != nil {
        return err
    }

    rm := Report{FileName: fileName, Size: size, LastModified: lastModified}
    conn.Create(&rm)

    return nil
}
