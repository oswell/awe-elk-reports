AWS Elk reports
===================

Pull Amazon AWS billing report data, and insert into elastic search.
The basic flow is as follows:

1.  Crawl S3, gathering a list of report files, and compare against the list of already
    processed files.

2.  Any unprocessed files are processed and indexed into elastic search.

Processed files, and associated data (size, hash, etc) are stored in a SQL store to
ensure we do not reprocess files when not necessary.

Configuration is handled via the TOML configuration, or more commonly, via environment variables.

```
Generated environment variables:
   CFG_BUCKET                   => S3 Bucket Name
   CFG_PREFIX                   => S3 report key prefix
   CFG_DBURL                    => MySQL connection string
   CFG_KAFKA_BROKERS            => List of Kafka brokers (comma delimited list)
   CFG_KAFKA_TLSCACERTIFICATE   => TLS CA certificate path
   CFG_KAFKA_TLSCERTIFICATE     => TLS certificiate path
   CFG_KAFKA_TLSKEY             => TLS key path
   CFG_KAFKA_TLSVERIFY          => Boolean, verify TLS certificates
   CFG_KAFKA_TOPIC              => Kafka topic to push to
   CFG_LOGLEVEL                 => Log level for output

AWS access is controlled via the environment variables:
   AWS_REGION
   AWS_ACCESS_KEY_ID
   AWS_SECRET_ACCESS_KEY
```
