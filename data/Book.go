package data

type BookFile struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Link      string `json:"link"`
	Extension string `json:"ext"`
}

type JSONReadWriter interface {
	WriteData([]BookFile) error
	ReadData() ([]byte, error)

	GetBooks() ([]BookFile, error)
}
