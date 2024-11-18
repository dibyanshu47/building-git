package gitobject

type GitObject interface {
	Serialize() []byte

	Deserialize(data []byte)

	Init()

	GetFormat() []byte
}
