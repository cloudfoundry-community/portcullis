portcullis: main.go api/*.go broker/*.go broker/bindparser/*.go store/*.go store/postgres/*.go config/*.go
	go build -ldflags='-X "github.com/cloudfoundry-community/portcullis/config.Version=(development build): $(shell /bin/date '+%Y-%m-%d %H:%M:%S')"' .

ARTIFACTS := artifacts/portcullis-{{.OS}}-{{.Arch}}
LDFLAGS := -X "github.com/cloudfoundry-community/portcullis/config.Version=$(VERSION)"
release:
	@echo "Checking that VERSION was defined in the calling environment"
	@test -n "$(VERSION)"
	@echo "OK.  VERSION=$(VERSION)"
	
	@echo "Checking that TARGETS was defined in the calling environment"
	@test -n "$(TARGETS)"
	@echo "OK.  TARGETS='$(TARGETS)'"
	rm -rf artifacts
	gox -osarch="$(TARGETS)" -ldflags='$(LDFLAGS)' --output="$(ARTIFACTS)/portcullis" .
	
	cd artifacts && for x in portcullis-*; do tar -czvf $$x.tar.gz $$x; rm -r $$x;  done

clean:
	rm -rf artifacts
	rm -f ./portcullis

test:
	go test ./api ./broker ./broker/bindparser ./config ./store

coverage: 
	ginkgo -cover ./api ./broker ./broker/bindparser ./config ./store
