.PHONY: test


test:
	docker-compose exec app go test ./pkg/service/test -v