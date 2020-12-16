test: 
	go test -v ./...

addlicense:
	# install with `go get github.com/google/addlicense`
	addlicense -c 'rnrch' -l apache -v .
