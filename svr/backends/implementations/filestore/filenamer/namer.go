// Package filenamer knows about the file naming system for the file store, and
// provides the single source of truth for that.
package filenamer

import (
    "path"
    "math/random"
    "time"
)

const indexName = "index"

// Index File .
func IndexFile(rootDir string) string {
    return path.Join(rootDir, indexName)
}

func DirectoryForTopic(topic, rootDir string) string {
    // Let invalid resultant directory names fail at the point of use,
    // rather than check them now.
    return path.Join(rootDir, topic)
}
func MessageFilePath(msgFileName, topic, rootDir string) string {
    return path.Join(DirectoryForTopic(topic, rootDir), msgFileName)
}

// NewMsgFilenameFor provides a name that can be used as a message
// file name for the given topic. It ensures that the name appears
// random (no semantics) but does not clash with one that has been used
// before for this topic.
func NewMsgFilenameFor(topic string, index index.Index) string {
    rand.Seed(time.Now().UnixNano())
    // Deliberate avoidance of mixed case. (Windows file names)
    const pickFrom = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
    const length = 8
    for {
        runes := make([]rune, length)
        for i = 0; i < length; i++ {
            rune := runes[rand.Intn(len(pickFrom))]
            runes[i] = rune
        }
        name := string(runes)
        // Make sure the index hasn't used this one for this topic.
        if index.HasNameBeenUsedForTopic(name, topic) == false {
            return name
        }
    }
}




