package btesting

import (
	"fmt"
	"strings"
	"testing"

	bolt "go.etcd.io/bbolt"
)

var _ bolt.Journal = (*StringJournal)(nil)

type StringJournal struct {
	strings.Builder
	paused bool
}

func (j *StringJournal) Pause() {
	j.paused = true
	j.WriteString("// snip\n")
}

func (j *StringJournal) Resume() {
	j.paused = false
}

func (j *StringJournal) Verify(t *testing.T, expected ...string) {
	t.Helper()
	joined := strings.Join(expected, "\n") + "\n"
	actual := j.String()
	if joined != actual {
		t.Fatalf("unexpected write operations:\n%s", actual)
	}
	j.Reset()
}

func (j *StringJournal) WriteTxStarted(id uint64) {
	if j.paused {
		return
	}
	fmt.Fprintf(j, "WriteTxStarted(%v)\n", id)
}

func (j *StringJournal) BucketCreated(bucket *bolt.Bucket) {
	if j.paused {
		return
	}
	j.WriteString("BucketCreated(")
	j.writeFullBucketName(bucket)
	j.WriteString(")\n")
}

func (j *StringJournal) BucketDeleted(bucket *bolt.Bucket) {
	if j.paused {
		return
	}
	j.WriteString("BucketDeleted(")
	j.writeFullBucketName(bucket)
	j.WriteString(")\n")
}

func (j *StringJournal) BucketMoved(oldParent *bolt.Bucket, moved *bolt.Bucket) {
	if j.paused {
		return
	}
	j.WriteString("BucketMoved(")
	j.writeFullBucketName(oldParent)
	j.WriteString(", ")
	j.writeFullBucketName(moved)
	j.WriteString(")\n")
}

func (j *StringJournal) KeyDeleted(bucket *bolt.Bucket, key []byte) {
	if j.paused {
		return
	}
	j.WriteString("KeyDeleted(")
	j.writeFullBucketName(bucket)
	fmt.Fprintf(j, ", %q)\n", string(key))
}

func (j *StringJournal) KeyUpdated(bucket *bolt.Bucket, key []byte, value []byte) {
	if j.paused {
		return
	}
	j.WriteString("KeyUpdated(")
	j.writeFullBucketName(bucket)
	fmt.Fprintf(j, ", %q, %q)\n", string(key), string(value))
}

func (j *StringJournal) SequenceUpdated(bucket *bolt.Bucket, value uint64) {
	if j.paused {
		return
	}
	j.WriteString("SequenceUpdated(")
	j.writeFullBucketName(bucket)
	fmt.Fprintf(j, ", %d)\n", value)
}

func (j *StringJournal) WriteTxCommitted() error {
	if j.paused {
		return nil
	}
	j.WriteString("WriteTxCommitted()\n")
	return nil
}

func (j *StringJournal) WriteTxRolledBack() {
	if j.paused {
		return
	}
	j.WriteString("WriteTxRolledBack()\n")
}

func (j *StringJournal) writeFullBucketName(b *bolt.Bucket) {
	if p := b.Parent(); p != nil {
		j.writeFullBucketName(p)
		j.WriteString("/")
	}
	j.Write(b.Name())
}
