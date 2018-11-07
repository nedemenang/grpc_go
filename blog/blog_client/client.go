package main

import (
	"io"
	"context"
	"github.com/nedemenang/grpc_go/blog/blogpb"
	"log"
	"google.golang.org/grpc"
	"fmt"
)

func main() {

	fmt.Println("Blog client")

	opts := grpc.WithInsecure()

	
	cc, err := grpc.Dial("localhost:50051", opts)

	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	fmt.Println("Creating the blog")
	blog := &blogpb.Blog{
		AuthorId: "Nnamso",
		Title: "My First Blog",
		Content: "Content of the first blog",
	}

	createBlogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected error: %v\n", err)
	}
	fmt.Printf("Blog has been created: %v\n", createBlogRes)
	blogId := createBlogRes.GetBlog().GetId()

	fmt.Println("Reading the Blog")

	_, err2 := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "lsljdl"})
	if err2 != nil {
		fmt.Printf("Error happend while reading: %v\n", err2)
	}

	readBlogReq := &blogpb.ReadBlogRequest{BlogId: blogId}
	readBlogRes, readBlogErr := c.ReadBlog(context.Background(), readBlogReq)
	if readBlogErr != nil {
		fmt.Printf("Error happend while reading: %v\n", readBlogErr)
	}

	fmt.Printf("Blog was read: %v\n", readBlogRes)

	// update blog

	newBlog := &blogpb.Blog{
		Id: blogId,
		AuthorId: "Changed Author",
		Title: "My First Blog(edited)",
		Content: "Content of the first blog with additions",
	}
	updateRes, updateErr := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog : newBlog})

	if updateErr != nil {
		fmt.Printf("Error happend while updating: %v\n", updateErr)
	} 

	fmt.Printf("Blog was read: %v", updateRes)

	deleteRes, deleteErr := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: blogId}) 
	if deleteErr != nil {
		fmt.Printf("Error happend while deleting: %v\n", deleteErr)
	} 
	fmt.Printf("Blog was deleted: %v\n", deleteRes)

	// list blogs
	

	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{}) 

	if err != nil {
		log.Fatalf("Error while calling list blog RPC: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}
		log.Println(res.GetBlog())
	}

}