export APP_PORT=3000

build:
	go build

local: export APP_MODE=local
local: export DB_HOST=localhost
local: export DB_NAME=weight_tracker
local: export DB_SCHEMA=public
local: export DB_USER=sam
local: export DB_PASSWORD="1030Jaco"
local: export DB_SSLMODE=disable
local:
	./weight_tracker
