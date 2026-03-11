all:
	go build -tags 'noaugeas novirt'

test:
	go test -tags 'noaugeas novirt'
