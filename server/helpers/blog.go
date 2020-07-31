package helpers

import (
	"github.com/lazhari/blog-grpc/blogpb"
	"github.com/lazhari/blog-grpc/server/models"
)

// DataToBlogPb transform document data to Blog protocol buffers
func DataToBlogPb(data *models.BlogItem) *blogpb.Blog {
	return &blogpb.Blog{
		Id:       data.ID.Hex(),
		AuthorId: data.AuthorID,
		Content:  data.Content,
		Title:    data.Title,
	}
}
