IMPLEMENTATION NOTES

Added an "archiving", only in campaigns as I will simply get the graph and move this data to an archived folder in "s3"

Decided on a pg vector column that is manually maintained for each row in each table of data that is "searchable". To search across all I can do a union across all or allow a search across just one object.

Developer Install notes

1) Install docker
2) run `docker-compose -up -d`
3) Install sql-migrate: `go get -u github.com/rubenv/sql-migrate/...`
4) run `make migrate-up`

Review to ensure the oms schema exists with the tables.