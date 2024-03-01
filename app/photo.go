package app

import "mime/multipart"

type Photo struct {
	Title   string                `form:"title"`
	Caption string                `form:"caption"`
	Photo   *multipart.FileHeader `form:"photo" binding:"-"`
}
