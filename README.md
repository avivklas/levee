# levee

## Usage

#### Create a bucket limiter
```go
now := time.Now()
// creating new bucket with the rate of of 8/s
localBucket := NewBucket(8/1, 8, now)
```

#### Create another bucket for global limit
```go
// creating new bucket with the rate of of 64/s
globalBucket := NewBucket(64/1, 64, now)
```

#### Create the limited reader
```go
// wrap an io.Reader with both buckets and get a new reader
// that will assure no overflow above the requested bandwidth
limiter := LimitedReader(reader, localBucket, globalBucket)
```