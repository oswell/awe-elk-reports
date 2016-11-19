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
   CONFIGURATION_BUCKET                   => S3 Bucket Name
   CONFIGURATION_PREFIX                   => S3 report key prefix
   CONFIGURATION_DBURL                    => MySQL connection string
   CONFIGURATION_KAFKA_BROKERS            => List of Kafka brokers
   CONFIGURATION_KAFKA_TLSCACERTIFICATE   => TLS CA certificate path
   CONFIGURATION_KAFKA_TLSCERTIFICATE     => TLS certificiate path
   CONFIGURATION_KAFKA_TLSKEY             => TLS key path
   CONFIGURATION_KAFKA_TLSVERIFY          => Boolean, verify TLS certificates
   CONFIGURATION_KAFKA_TOPIC              => Kafka topic to push to

AWS access is controlled via the environment variables:
   AWS_REGION
   AWS_ACCESS_KEY_ID
   AWS_SECRET_ACCESS_KEY
```
