package filenamer

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

//-----------------------------------------------------------------------
// Mock PreviouslyUsed implementations.
//-----------------------------------------------------------------------
type AlwaysFalse struct{}

func (checker AlwaysFalse) PreviouslyUsed(name, context string) bool {
	return false
}

type SayNamesWithNumbersPreviouslyUsed struct{}

func (checker SayNamesWithNumbersPreviouslyUsed) PreviouslyUsed(
	name, context string) bool {
	disallow := []rune("0123456789")
	for _, c := range disallow {
		cAsString := string(c)
		if strings.Contains(name, cAsString) {
			return true
		}
	}
	return false
}

//-----------------------------------------------------------------------
// Tests
//-----------------------------------------------------------------------

func TestNamesRightShape(t *testing.T) {
	previouslyUsedChecker := AlwaysFalse{}
	name := NewMsgFilenameFor("topicA", previouslyUsedChecker)
	expected := 8
	assert.Equal(t, expected, len(name))
}

func TestRejectionCallbackWorking(t *testing.T) {
	previouslyUsedChecker := SayNamesWithNumbersPreviouslyUsed{}
	// Ask for many names in sequence, making sure none offered contains our
	// disallowed runes.
	disallowed := []rune("0123456789")
	for i := 0; i < 1000; i++ {
		name := NewMsgFilenameFor("topicA", previouslyUsedChecker)
		for _, c := range disallowed {
			cAsString := string(c)
			if strings.Contains(name, cAsString) {
				msg := fmt.Sprintf("Name %s contains %s", name, cAsString)
				assert.FailNow(t, msg)
			}
		}
	}
}
