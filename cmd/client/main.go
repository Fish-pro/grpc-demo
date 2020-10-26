package main

import (
	"context"
	v1 "github.com/Fish-pro/grpc-demo/api/proto/v1"
	"github.com/Fish-pro/grpc-demo/helper"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"log"
	"time"
)

const apiVsersion = "v1"

func main() {
	cred := helper.GetClientCred()
	conn, err := grpc.Dial(":8081", grpc.WithTransportCredentials(cred))
	if err != nil {
		log.Fatal("failed to connect server:", err)
	}
	defer conn.Close()

	c := v1.NewToDoServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t := time.Now().In(time.UTC)
	reminder, _ := ptypes.TimestampProto(t)
	pfx := t.Format(time.RFC3339Nano)

	req1 := v1.CreateRequest{
		Api: apiVsersion,
		ToDo: &v1.ToDo{
			Title:       "title (" + pfx + ")",
			Description: "description (" + pfx + ")",
			Reminder:    reminder,
		},
	}
	res1, err := c.Create(ctx, &req1)
	if err != nil {
		log.Fatal("create error:", err)
	}
	log.Printf("create result:%v", res1)
	id := res1.Id

	req2 := v1.ReadRequest{Id: id, Api: apiVsersion}
	res2, err := c.Read(ctx, &req2)
	if err != nil {
		log.Fatal("Read failed:", err)
	}
	log.Printf("Read result: %v", res2)

	req3 := v1.UpdateRequest{
		Api: apiVsersion,
		ToDo: &v1.ToDo{
			Id:          res2.ToDo.Id,
			Title:       res2.ToDo.Title,
			Description: res2.ToDo.Description + "updated",
			Reminder:    res2.ToDo.Reminder,
		},
	}
	res3, err := c.Update(ctx, &req3)
	if err != nil {
		log.Fatal("update error:", err)
	}
	log.Printf("update result: %v", res3)

	req4 := v1.ReadAllRequest{
		Api: apiVsersion,
	}
	res4, err := c.ReadAll(ctx, &req4)
	if err != nil {
		log.Fatal("read all error:", err)
	}
	log.Printf("read all result: %v", res4)

	req5 := v1.DeleteRequest{Api: apiVsersion, Id: id}
	res5, err := c.Delete(ctx, &req5)
	if err != nil {
		log.Fatal("delete error:", err)
	}
	log.Printf("delete result: %v", res5)
}
