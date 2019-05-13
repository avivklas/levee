package levee

import (
	"io"
	"math"
	"time"
)

// levee implements io.Reader by wrapping another reader and limits
// its bandwidth by rate limits that are defined by the bucket and the global bucket
type levee struct {
	reader       io.Reader
	bucket       *Bucket
	globalBucket *Bucket
}

// LimitedReader returns a levee reader that limits the bandwidth of reads
// of the underlying reader by rate limited buckets - local and global
func LimitedReader(reader io.Reader, bucket *Bucket, globalBucket *Bucket) io.Reader {
	return &levee{
		reader:       reader,
		bucket:       bucket,
		globalBucket: globalBucket,
	}
}

// Read implementation of io.Reader
func (r *levee) Read(buf []byte) (int, error) {
	n, err := r.reader.Read(buf)
	if n <= 0 {
		return n, err
	}
	now := time.Now()
	bucketWait := r.bucket.take(now, int64(n))
	globalBucketWait := r.globalBucket.take(now, int64(n))
	wait := math.Max(float64(bucketWait), float64(globalBucketWait))
	time.Sleep(time.Duration(wait))
	return n, err
}
