export GOOSE_DRIVER=sqlite3
export GOOSE_DBSTRING=./lidardb.sqlite3
export GOOSE_MIGRATION_DIR=./migrations

status:
	goose status

up:
	goose up
	
down:
	goose down



app:
	go clean
	go build -ldflags="-w -s"


clean:
	go clean 