package common

type UpdateFile struct {
	Path            string `json:"path"`
	Name            string `json:"name"`
	Text            string `json:"text"`
	FileLastVersion string `json:"filelastversion"`
}
