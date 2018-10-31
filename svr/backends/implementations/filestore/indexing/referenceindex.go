package indexing

import (
	"time"
)

// MakeReferenceIndex provides a repeatable Index for testing purposes.
// the creation times returned are in oldest to newest sequence:
// topicA/file1/msg1
// topicA/file1/msg2
// topicA/file1/msg3
// topicA/file2/msg4
// topicA/file2/msg5
// topicA/file2/msg6
// topicB/file1/msg1
// ...
func MakeReferenceIndex() (index *Index, creationTimes []time.Time) {

	idx := NewIndex()

	ctimes := []time.Time{}
	// Use two topics.
	for _, topic := range []string{"topicA", "topicB"} {
		idx.RegisterTopic(topic)
		msgFileList := idx.GetMessageFileListFor(topic)
		// Register two files in each topic.
		for _, fileName := range []string{"file1", "file2"} {
			msgFileList.RegisterNewFile(fileName)
			fileMeta := msgFileList.Meta[fileName]
			// Register 3 messages in each file.
			for i := 0; i < 3; i++ {
				time.Sleep(time.Duration(100 * time.Millisecond))
				msgNumber := idx.GetAndIncrementMessageNumberFor(topic)
				msgSize := int64(1024)
				now := time.Now()
				ctimes = append(ctimes, now)
				fileMeta.RegisterNewMessage(msgNumber, msgSize)
			}
		}
	}
	return idx, ctimes
}
