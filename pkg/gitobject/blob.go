package gitobject

type GitBlob struct {
	format   []byte
	BlobData []byte
}

func NewGitBlob(data []byte) *GitBlob {
	return &GitBlob{
		format:   []byte("blob"),
		BlobData: data,
	}
}

func (g *GitBlob) Serialize() []byte {
	return g.BlobData
}

func (g *GitBlob) Deserialize(data []byte) {
	g.BlobData = data
}

func (g *GitBlob) Init() {
	// do nothing
}

func (g *GitBlob) GetFormat() []byte {
	return g.format
}
