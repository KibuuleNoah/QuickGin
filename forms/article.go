package forms

type CreateArticleForm struct {
	Title    string `json:"title" validate:"required,min=3,max=200"`
	Body     string `json:"body" validate:"required,min=10"`
	AuthorID string `json:"author_id" validate:"required,uuid"`
}

type UpdateArticleForm struct {
	ID        string  `json:"-" validate:"required,uuid"`
	Title     *string `json:"title" validate:"omitempty,min=3,max=200"`
	Body      *string `json:"body" validate:"omitempty,min=10"`
	Published *bool   `json:"published"`
}

var CreateArticleFormMessages = ValidationMessages{
	"Title": {
		"required": "Please enter the article title",
		"min":      "Title should be between 3 to 100 characters",
		"max":      "Title should be between 3 to 100 characters",
	},
	"Body": {
		"required": "Please enter the article content",
		"min":      "Content should be between 3 to 1000 characters",
		"max":      "Content should be between 3 to 1000 characters",
	},
	"AuthorID": {
		"required": "Author ID is required",
		"uuid":     "Author ID must be a valid UUID",
	},
}

var UpdateArticleFormMessages = ValidationMessages{
	"Title": {
		"min": "Title must be at least 3 characters",
		"max": "Title must not exceed 200 characters",
	},
	"Body": {
		"min": "Body must be at least 10 characters",
	},
}
