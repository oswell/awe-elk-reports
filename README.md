AWS Elk reports
===================

Pull Amazon AWS billing report data, and insert into elastic search.
The basic flow is as follows:

1.  Crawl S3, gathering a list of report files, and compare against the list of already
    processed files.

2.  Any unprocessed files are processed and indexed into elastic search.

Processed files, and associated data (size, hash, etc) are stored in a SQL store to
ensure we do not reprocess files when not necessary.
