package _struct

type Books struct {
	NovelName  string
	NovelID    string
	IsFinish   bool
	MarkCount  string
	NovelCover string
	AuthorName string
	CharCount  string
	SignStatus string
}
type MyBookInfoJsonPro struct {
	Book       Books
	NewBooks   map[string]string
	ConfigPath string
	CoverPath  string
	//BackupsPath  string
	DownloadList []string
}
