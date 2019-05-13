# levee

## Usage

#### Create a limited reader
```go
now := time.Now()
// creating new bucket with the rate of of 8/s
localBucket := NewBucket(8/1, 8, now)
// creating new bucket with the rate of of 64/s
globalBucket := NewBucket(64/1, 64, now)
// wrap an io.Reader with both buckets and get a new reader
// that will assure no overflow above the requested bandwidth
limiter := LimitedReader(reader, localBucket, globalBucket)
```