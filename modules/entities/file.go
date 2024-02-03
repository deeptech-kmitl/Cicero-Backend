package entities

type ImageRes struct {
	Id       string `db:"id" json:"id"`
	Url      string `db:"url" json:"url"`
	Filename string `db:"filename" json:"filename"`
}
