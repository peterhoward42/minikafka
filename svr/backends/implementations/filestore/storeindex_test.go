package filestore

import (
	"testing"
    "time"
)

func TestSomething(t *testing.T) {
    makeIndexProgrammatically()
}

//--------------------------------------------------------------------------------
// Auxilliary code.
//--------------------------------------------------------------------------------

func makeIndexProgrammatically() *StoreIndex {

    idx := StoreIndex{}

    idx["foo_topic"] = []FileMeta{}
    foo1Meta := FileMeta{
        "foo1", 
        MsgMeta{1, time.Now().Add(-9 * 24 * time.Hour),},
        MsgMeta{10, time.Now().Add(-8 * 24 * time.Hour),},
    }
    foo2Meta := FileMeta{
        "foo2", 
        MsgMeta{11, time.Now().Add(-7 * 24 * time.Hour),},
        MsgMeta{20, time.Now().Add(-6 * 24 * time.Hour),},
    }
    idx["foo_topic"] = append(idx["foo_topic"], foo1Meta)
    idx["foo_topic"] = append(idx["foo_topic"], foo2Meta)

    idx["bar_topic"] = []FileMeta{}
    bar1Meta := FileMeta{
        "bar1", 
        MsgMeta{1, time.Now().Add(-5 * 24 * time.Hour),},
        MsgMeta{10, time.Now().Add(-4 * 24 * time.Hour),},
    }
    idx["bar_topic"] = append(idx["bar_topic"], bar1Meta)

    return &idx
}
