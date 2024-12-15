package models

type UploadFile struct {
	Chunks     uint32 `json:"chunks" redis:"chunks"`
	ChunkIndex uint32 `json:"chunkIndex" redis:"chunkIndex"`
	FileSize   int64  `json:"fileSize" redis:"fileSize"`
	FileName   string `json:"fileName" redis:"fileName"`
	FilePath   string `json:"filePath" redis:"filePath"`
	UploadId   string `json:"uploadId" redis:"uploadId"`
}
