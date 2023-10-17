package src

type Endpoints struct {
	CreateImageFileEndpoint string
	CreateImagePostEndpoint string
	CreateCategoryEndpoint  string
	GetCategoryEndpoint     string
}
type HttpResponse[T any] struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Data    T      `json:"data"`
}
type HttpResponseError struct {
	Message   string
	Code      int
	Status    string
	Timestamp string
}

type SpringBootResponseError struct {
	Message string `json:"message"`
	Path    string `json:"path"`
	Status  int    `json:"status"`
	Error   string `json:"error"`
}

type Category struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	BgColor     string `json:"bgColor"`
	FgColor     string `json:"fgColor"`
}
type Image struct {
	Id          string   `json:"id"`
	Title       string   `json:"title"`
	ImageUrl    string   `json:"image_url"`
	Description string   `json:"description"`
	Category    []string `json:"category"`
	Tags        []string `json:"tags"`
}

type Tag struct {
	Name string `json:"name"`
}

type ImageFile struct {
	Id           string        `json:"id"`
	AccentColor  string        `json:"accentColor"`
	ColorPalette string        `json:"colorPalette"`
	Landscape    bool          `json:"landscape"`
	Variants     []FileVariant `json:"variants"`
}

type FileVariant struct {
	Id       int    `json:"id"`
	Original bool   `json:"original"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Size     int    `json:"size"`
	Title    string `json:"title"`
	Url      string `json:"url"`
	Md5      string `json:"md5"`
	Sha256   string `json:"sha256"`
}

type ImagePost struct {
	Id          string    `json:"id"`
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	Description string    `json:"content"`
	CreatedAt   string    `json:"createdAt"`
	UpdatedAt   string    `json:"updatedAt"`
	UserId      string    `json:"userId"`
	Public      bool      `json:"publik"`
	Ready       bool      `json:"ready"`
	Nsfw        bool      `json:"nsfw"`
	Thumbnail   string    `json:"thumbnail"`
	Image       ImageFile `json:"image"`
	Category    Category  `json:"category"`
	Tags        []Tag     `json:"tags"`
}

type ImagePostCreateModel struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Public   bool     `json:"public"`
	Nsfw     bool     `json:"nsfw"`
	Tags     []string `json:"tags"`
	Image    string   `json:"image"`
	Category string   `json:"category"`
}
type ImageSaveError struct {
	Message     string `json:"message"`
	OriginalErr error  `json:"error"`
	StatusCode  int    `json:"status_code"`
}

func (e *ImageSaveError) Error() string {
	return e.Message
}
