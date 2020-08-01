package internal

import (
	"context"
	"fmt"

	"github.com/lazhari/blog-grpc/internal/helpers"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/lazhari/blog-grpc/blogpb"
	"github.com/lazhari/blog-grpc/internal/models"
)

// Server The gRPC server
type Server struct {
	Collection *mongo.Collection
}

func (s *Server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	fmt.Println("Create blog request")
	blog := req.GetBlog()

	blogItem, err := helpers.CreateBlogPost(s.Collection, blog)

	if err != nil {
		return nil, err
	}

	return &blogpb.CreateBlogResponse{
		Blog: blogItem,
	}, nil

}

func (s *Server) ReadBlog(ctx context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	fmt.Println("ReadBlog invoked!")
	blogID := req.GetBlogId()
	oid, err := primitive.ObjectIDFromHex(blogID)

	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID"),
		)
	}

	// Create an empty struct
	data := &models.BlogItem{}
	filter := bson.M{"_id": oid}

	res := s.Collection.FindOne(context.Background(), filter)

	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find blog with specified ID: %v", err),
		)
	}

	return &blogpb.ReadBlogResponse{
		Blog: helpers.DataToBlogPb(data),
	}, nil
}

func (s *Server) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {
	fmt.Println("UpdateBlog invoked")
	blog := req.GetBlog()
	oid, err := primitive.ObjectIDFromHex(blog.GetId())

	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID"),
		)
	}

	// Create an empty struct
	data := &models.BlogItem{}
	filter := bson.M{"_id": oid}

	res := s.Collection.FindOne(context.Background(), filter)

	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find blog with specified ID: %v", err),
		)
	}
	// Update the data
	data.AuthorID = blog.GetAuthorId()
	data.Content = blog.GetContent()
	data.Title = blog.GetTitle()

	_, updateErr := s.Collection.ReplaceOne(context.Background(), filter, data)

	if updateErr != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("We cannot perform an update: %v", updateErr),
		)
	}

	return &blogpb.UpdateBlogResponse{
		Blog: helpers.DataToBlogPb(data),
	}, nil
}

func (s *Server) DeleteBlog(ctx context.Context, req *blogpb.DeleteBlogRequest) (*blogpb.DeleteBlogResponse, error) {
	fmt.Println("DeleteBlog invoked!")
	oid, err := primitive.ObjectIDFromHex(req.GetBlogId())

	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID"),
		)
	}

	filter := bson.M{"_id": oid}

	res, err := s.Collection.DeleteOne(context.Background(), filter)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot delete the blog post in MongoDB : %v", err),
		)
	}

	if res.DeletedCount == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find the blog with"),
		)
	}

	return &blogpb.DeleteBlogResponse{
		BlogId: req.GetBlogId(),
	}, nil
}

func (s *Server) ListBlog(req *blogpb.ListBlogRequest, stream blogpb.BlogService_ListBlogServer) error {
	fmt.Println("ListBlog invoked")
	cursor, err := s.Collection.Find(context.Background(), bson.M{})

	if err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unknown internal error : %v", err),
		)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		data := &models.BlogItem{}
		err := cursor.Decode(data)

		if err != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("Unknown while decoding data from MongoDB: %v", err),
			)
		}
		stream.Send(&blogpb.ListBlogResponse{
			Blog: helpers.DataToBlogPb(data),
		})
	}

	if err := cursor.Err(); err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unknown internal error : %v", err),
		)
	}
	return nil
}
