package dto

type UploadResp struct {
	FileName     string `json:"file_name"`
	FileUrl      string `json:"file_url"`
	Checksum     string `json:"checksum"`
	HashFunction string `json:"hash_function"`
}
