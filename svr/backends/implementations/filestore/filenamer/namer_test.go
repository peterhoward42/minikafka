package filenamer

import (
	"testing"
    "strings"
    "fmt"

	"github.com/stretchr/testify/assert"
)

//-----------------------------------------------------------------------
// Mock NameChecker implementations.
//-----------------------------------------------------------------------
type AcceptAnyName struct{}

func (a AcceptAnyName) IsFilenameOk(name, context string) bool {
	return true
}

type AcceptSomeNames struct{}

func (a AcceptSomeNames) IsFilenameOk(name, context string) bool {
    disallow := []rune("0123456789")
    for  _, c := range disallow {
        cAsString := string(c)
        if strings.Contains(name, cAsString) {
            return false
        }
    }
    return true
}

//-----------------------------------------------------------------------
// Tests
//-----------------------------------------------------------------------

func TestNamesRightShape(t *testing.T) {
	fileNameChecker := AcceptAnyName{}
	name := NewMsgFilenameFor("topicA", fileNameChecker)
	expected := 8
	assert.Equal(t, expected, len(name))
}

func TestRejectionCallbackWorking(t *testing.T) {
	fileNameChecker := AcceptSomeNames{}
    // Ask for many names in sequence, making sure none offered contains our
    // disallowed runes.
    disallowed := []rune("0123456789")
    for i := 0; i < 1000; i++ {
        name := NewMsgFilenameFor("topicA", fileNameChecker)
        for  _, c := range disallowed {
            cAsString := string(c)
            if strings.Contains(name, cAsString) {
                msg := fmt.Sprintf("Name %s contains %s", name, cAsString)
                assert.FailNow(t, msg)
            }
        }
    }
}
