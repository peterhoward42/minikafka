package indexing

import (
    "time"
)

// MakeReferenceIndex provides a repeatable Index for testing purposes.
func MakeReferenceIndex() *Index {

	idx := NewIndex()

	msgNum := int32(0)
	minutes := 1
	for _, topic := range []string{"topicA", "topicB"} {
		msgFileList := idx.GetMessageFileListFor(topic)
		for _, fileName := range []string{"file1", "file2"} {
			msgFileList.RegisterNewFile(fileName)

			fileMeta := msgFileList.Meta[fileName]

			fileMeta.Oldest.MsgNum = msgNum + 1
			fileMeta.Oldest.Created = nowMinusNMinutes(minutes)

			fileMeta.Newest.MsgNum = msgNum + 5
			fileMeta.Newest.Created = nowMinusNMinutes(minutes + 5)

			msgNum += 10
			minutes += 15
		}
	}
	return idx
}

func nowMinusNMinutes(minutes int) time.Time {
	now := time.Now()
	duration := time.Duration(minutes) * time.Minute
	return now.Add(-duration)
}
