# Analyze Postmark DMARC Raw Reports

Get raw DMARC reports for your domain from [Postmark's API](https://dmarc.postmarkapp.com/api/) and analyze the records.


## Prerequisites

Assuming you have Postmark as `mailto` in your `_dmarc` DNS record you'll of course need an API token in order to interact with their API. Also don't forget to designate the number of workers, meaning the number of concurrent requests to the API. You can do that by exporting `TOKEN` and `WORKERS` environment variables or designate them in an `.env` file.

## Usage

Get all raw reports' metadata for a given date range (i.e. from 2024-11-18 to 2024-11-25).
```
go run reports/reports.go 2024-11-18 2024-11-25
```

This will create a `reports.json` file.

Get all the records from all of the reports.
```
go run records/records.go
```

This will create a `records.json` file.

Extract DKIM and SPF pass/fail statistics from the records.
```
go run explore.go
```

This will create a `stats.json` file


## License

[![License: MIT](https://img.shields.io/github/license/vlatan/postmark-reports?label=License)](/LICENSE "License: MIT")