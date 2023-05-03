package common

type UpdateFile struct {
	FilePath    string `json:"path"`
	FileName    string `json:"name"`
	FileText    string `json:"text"`
	FileVersion string `json:"filelast_version"`
}
