Simple Web load tester

Used for testing concurrency primitives, context manipulation and error handling.


## Usage

go run . url1 url2 url3 ... --goroutines={max no. goroutines}


## Testing

go test -v -cover