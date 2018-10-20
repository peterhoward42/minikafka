// package indexing is centred around its Index type, which is an in-memory data
// structure that keeps track of which message filenames have been used for each
// topic, and for each, which range of message numbers and creation times they
// hold. The Index type also provides methods whereby an instance can be
// gob-serialized and deserialized, and then an additional pair of methods to
// save and retrieve this serialized representation to disk. The Index holds
// message file names as file basenames and has no knowledge about where these
// files are.
package indexing

// The types' fields are exported so they can be automatically gob-encoded
// without bothering with structure tags.

// Index is the top level index object which organises the information it holds
// by topic.
type Index struct {
	MessageFileLists map[string]*MessageFileList
}

// NewIndex creates and initialized an Index.
func NewIndex() *Index {
	return &Index{map[string]*MessageFileList{}}
}

// GetMessageFileListFor provides access to the MesageFileList for the
// given topic. It copes gracefully with the topic being hithertoo unknown.
func (index *Index) GetMessageFileListFor(topic string) *MessageFileList {
	_, ok := index.MessageFileLists[topic]
	if ok == false {
		index.MessageFileLists[topic] = NewMessageFileList()
	}
	return index.MessageFileLists[topic]
}

// NextMessageNumberFor provides the next-available message number for a topic.
// It copes gracefully with the two special cases of the index having no
// record of that topic, or it having never yet contained any messages.
func (index Index) NextMessageNumberFor(topic string) int32 {
	messageFileList, ok := index.MessageFileLists[topic]
	if ok == false {
		return 1
	}
	nTopicFiles := len(messageFileList.Names)
	if nTopicFiles == 0 {
		return 1
	}
	// Consult the meta data for the newest file.
	newestName := messageFileList.Names[nTopicFiles-1]
	newestFileMeta := messageFileList.Meta[newestName]
	return newestFileMeta.Newest.MsgNum + 1
}

// CurrentMsgFileNameFor provides the name of the file that is currently being
// used to store incoming messages for a topic. It copes gracefully with there
// not being one - by returning an empty string.
func (index Index) CurrentMsgFileNameFor(topic string) string {
	msgFileList, ok := index.MessageFileLists[topic]
	if ok == false {
		return ""
	}
	if len(msgFileList.Names) == 0 {
		return ""
	}
	n := len(msgFileList.Names)
	return msgFileList.Names[n-1]
}

// PreviouslyUsed indicates if the given file base name has been used at any
// time previously as a message file for the given topic.
func (index Index) PreviouslyUsed(name string, topic string) bool {
	msgFileList, ok := index.MessageFileLists[topic]
	if ok == false {
		return false
	}
	if len(msgFileList.Names) == 0 {
		return false
	}
	for _, existingName := range msgFileList.Names {
		if existingName == name {
			return true
		}
	}
	return false
}

//-----------------------------------------------------------------------
