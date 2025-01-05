package cloud

type PhotoService interface {
	Connect() error
	ListPhotos() ([]Photo, error)
	DownloadPhoto(photo Photo) ([]byte, error)
	UploadPhoto(photo Photo, data []byte) error
}

type Photo struct {
	ID       string
	Name     string
	Path     string
	Metadata map[string]interface{}
}
