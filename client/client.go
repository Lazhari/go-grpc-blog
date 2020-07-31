package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/lazhari/blog-grpc/blogpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I am the sum client!")

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50051", opts)

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	createBlog(c)

}

func createBlog(c blogpb.BlogServiceClient) {
	fmt.Println("Create a new blog post")

	blog := &blogpb.Blog{
		AuthorId: "Lazhari",
		Title:    "Go lang",
		Content:  "The Go lang is amazing!!",
	}

	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})

	if err != nil {
		log.Fatalf("Error while creating blog: %v", err)
	}

	fmt.Println("Blog has been created:", res)

	blogID := res.GetBlog().GetId()

	// Reading a blog
	readBlog(c, "123213dff")
	readBlog(c, blogID)

	// Update Blog
	// Update the blog
	newBlog := &blogpb.Blog{
		Id:       blogID,
		AuthorId: "Lazhari Mohammed",
		Title:    "Go language is Amazing lang",
		Content:  "The Go lang is amazing!!",
	}

	updateBlog(c, newBlog)

	// Delete Blog
	deleteBlog(c, blogID)

	// Listing blogs
	listBlogs(c)
}

func readBlog(c blogpb.BlogServiceClient, id string) {

	res, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: id})

	if err != nil {
		log.Printf("[X] Error happened while reading: %v\n", err)
		return
	}

	fmt.Println("[+] Blog is: ", res)
}

func updateBlog(c blogpb.BlogServiceClient, newBlog *blogpb.Blog) {
	res, err := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: newBlog})

	if err != nil {
		log.Printf("Error while updating the blog: %v\n", err)
		return
	}
	fmt.Println("[+] the update blog is:", res)
}

func deleteBlog(c blogpb.BlogServiceClient, blogID string) {
	res, err := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: blogID})

	if err != nil {
		log.Printf("Error while deleting the blog: %v\n", err)
		return
	}
	fmt.Println("[+] The blog has been deleted: ", res.GetBlogId())
}

func listBlogs(c blogpb.BlogServiceClient) {
	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})

	if err != nil {
		log.Printf("Error while listing the blogs: %v", err)
		return
	}

	for {
		res, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Something happended: %v\n", err)
			return
		}

		fmt.Println(res.GetBlog())
	}

}
