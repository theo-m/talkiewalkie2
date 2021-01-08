nukeDb:
	dropdb talkiewalkie
	createdb talkiewalkie

newMigration:
	migrate create -ext sql -dir migrations -seq XXX

migrate:
	migrate -source file://migrations -database postgres://theo:@localhost:5432/talkiewalkie?sslmode=disable up

envoy:
	envoy -c envoy.yaml

