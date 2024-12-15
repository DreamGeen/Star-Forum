package models

type PreUpload struct {
	FileName string `form:"fileName"`
	Chunks   uint32 `form:"chunks"`
}
