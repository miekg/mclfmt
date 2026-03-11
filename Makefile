all:
	go build -tags 'noaugeas novirt'

test:
	go test -v -tags 'noaugeas novirt'
