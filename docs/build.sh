go fmt *.go
go fmt handlers/*.go
go fmt model/*.go
go fmt srcp/*.go

go clean
go build

cd test
go fmt *.go
go test
cd ..
