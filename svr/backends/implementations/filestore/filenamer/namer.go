// Package filenamer knows about the file naming system for the file store, and
// provides the single source of truth for that.
package filenamer

import (
	"math/rand"
	"path"
	"time"
)

const indexName = "index"

// IndexFile .
func IndexFile(rootDir string) string {
	return path.Join(rootDir, indexName)
}

// DirectoryForTopic .
func DirectoryForTopic(topic, rootDir string) string {
	// Let invalid resultant directory names fail at the point of use,
	// rather than check them now.
	return path.Join(rootDir, topic)
}

// MessageFilePath .
func MessageFilePath(msgFileName, topic, rootDir string) string {
	return path.Join(DirectoryForTopic(topic, rootDir), msgFileName)
}

// NewMsgFilenameFor provides a name that can be used as a message
// file name for the given topic. It ensures that the name appears
// random (no semantics) and also satisfies the NameChecker (interface) passed
// in.
func NewMsgFilenameFor(topic string, nameChecker NameChecker) string {
	// Deliberate avoidance of mixed case. (Windows file names)
	pickFrom := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	nChoices := len(pickFrom)
	const requiredLength = 8
	rand.Seed(time.Now().UnixNano())
	for {
		runes := make([]rune, requiredLength)
		for i := 0; i < requiredLength; i++ {
			randomIndex := rand.Intn(nChoices - 1)
			randomRune := pickFrom[randomIndex]
			runes[i] = randomRune
		}
		name := string(runes)
		// Is the name checker ok with this name?
		if nameChecker.IsFilenameOk(name, topic) == true {
			return name
		}
		// If not keep trying.
	}
}

// NameChecker is a thing that will check whether a filename is ok from
// its perspective in the given context.
type NameChecker interface {
	IsFilenameOk(name string, context string) bool
}
