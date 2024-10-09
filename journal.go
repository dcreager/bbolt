package bbolt

// Journal lets you inspect all of the writes that are applied to a [DB].
//
// This interface includes “event” methods for each update that can be made to a
// database within a write transaction.  There are no events for reads.
//
// Note that we only call an event method when that update actually needed to be
// performed.  For instance, if you use [Bucket.CreateBucketIfNotExists], we
// will only call [BucketCreated] if the bucket did not already exist, and
// actually needed to be created.
type Journal interface {
	// WriteTxStarted is called when a new write transaction is opened.
	WriteTxStarted(id uint64)

	// BucketCreated is called when a new bucket is created.
	BucketCreated(bucket *Bucket)

	// BucketDeleted is called when a bucket is deleted.
	BucketDeleted(bucket *Bucket)

	// BucketMoved is called when a bucket is moved from one parent to another.
	BucketMoved(oldParent *Bucket, moved *Bucket)

	// KeyDeleted is called when a key is deleted from a bucket.  The key is
	// only valid for the duration of the call.
	KeyDeleted(bucket *Bucket, key []byte)

	// KeyUpdated is called when a new value is assigned to a key in a bucket.
	// The key and value are only valid for the duration of the call.
	KeyUpdated(bucket *Bucket, key []byte, value []byte)

	// SequenceUpdated is called when the sequence number for a bucket is
	// updated.
	SequenceUpdated(bucket *Bucket, value uint64)

	// WriteTxCommitted is called when the current write transaction is about to
	// be committed.  If this method returns an error, the transaction is rolled
	// back, and that error is returned from [Tx.Commit].
	WriteTxCommitted() error

	// WriteTxRolledBack is called after the current write transaction is rolled
	// back.
	WriteTxRolledBack()
}
