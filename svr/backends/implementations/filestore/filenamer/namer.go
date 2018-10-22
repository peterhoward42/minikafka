// Package filenamer knows about the file naming system for the file store, and
// provides the single source of truth for that.
package filenamer

import (
	"math/rand"
	"path"
	"time"
)

const indexName = "index"

// IndexFile provides the full path of the index file.
func IndexFile(rootDir string) string {
	return path.Join(rootDir, indexName)
}

// DirectoryForTopic provides the directory that should be used for the
// given topic.
func DirectoryForTopic(topic, rootDir string) string {
	// Let invalid resultant directory names fail at the point of use,
	// rather than check them now.
	return path.Join(rootDir, topic)
}

// MessageFilePath provides the full path of where a message file with a given
// basename can be found for a given topic.
func MessageFilePath(msgFileName, topic, rootDir string) string {
	return path.Join(DirectoryForTopic(topic, rootDir), msgFileName)
}

// NewMsgFilenameFor provides a file base name that can be used as a message
// file name for the given topic. It ensures that the name appears
// random (no implied semantics) and also satisfies the NotPreviouslyUsed
// interface passed in. (Specified as an interface to allow the use of a mock
// during testing).
func NewMsgFilenameFor(
	topic string, previouslyUsedChecker PreviouslyUsedChecker) string {
	// Names are randomly generated, 8 characters long, from a specified
	// palette

	// Deliberate avoidance of mixed case for Windows suitability.
	pickFrom := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	nChoices := len(pickFrom)
	const requiredLength = 8
	rand.Seed(time.Now().UnixNano())
	// Keep making them up for as long as the previously-used checker
	// rejects them.
	for {
		runes := make([]rune, requiredLength)
		for i := 0; i < requiredLength; i++ {
			randomIndex := rand.Intn(nChoices - 1)
			randomRune := pickFrom[randomIndex]
			runes[i] = randomRune
		}
		name := string(runes)
		// Has this name been used before?
		if previouslyUsedChecker.PreviouslyUsed(name, topic) == false {
			return name
		}
	}
}

// PreviouslyUsedChecker is a thing that will check whether a name has
// already been used before in a given context.
type PreviouslyUsedChecker interface {
	PreviouslyUsed(name string, context string) bool
}
