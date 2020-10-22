package v1

import (
	"context"
	"database/sql"
	"fmt"
	v1 "github.com/Fish-pro/grpc-demo/api/proto/v1"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

const apiVersion = "v1"

type ToDoServiceServer struct {
	db *sql.DB
}

func NewToDoServiceServer(db *sql.DB) v1.ToDoServiceServer {
	return &ToDoServiceServer{db: db}
}

func (s *ToDoServiceServer) checkAPI(api string) error {
	if apiVersion != api {
		return status.Error(codes.Unimplemented, fmt.Sprintf("unsupported API version:service implements API version '%s',but given '%s'", apiVersion, api))
	}
	return nil
}

func (s *ToDoServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("failed to connect to database:%s", err.Error()))
	}
	return c, nil
}

func (s *ToDoServiceServer) Create(ctx context.Context, req *v1.CreateRequest) (*v1.CreateResponse, error) {
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	reminder, err := ptypes.Timestamp(req.ToDo.Reminder)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("params error:%s", err.Error()))
	}

	res, err := c.ExecContext(ctx, "INSERT INTO ToDo(`Title`,`Description`,`Reminder`) VALUES(?,?,?)", req.ToDo.Title, req.ToDo.Description, reminder)
	if err != nil {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("create error:%s", err.Error()))
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("get new id error:%s", err.Error()))
	}
	return &v1.CreateResponse{Api: apiVersion, Id: id}, nil
}

func (s *ToDoServiceServer) Read(ctx context.Context, req *v1.ReadRequest) (*v1.ReadResponse, error) {
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	rows, err := c.QueryContext(ctx, "SELECT `ID`,`Title`,`Description`,`Reminder` FROM ToDo WHERE `ID`=?", req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("query error:%s", err.Error()))
	}
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, fmt.Sprintf("query error:%s", err.Error()))
		}
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ID = %d not found", req.Id))
	}
	var td v1.ToDo
	var reminder time.Time
	if err := rows.Scan(&td.Id, &td.Title, &td.Description, &reminder); err != nil {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("query data error:%s", err.Error()))
	}
	if rows.Next() {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("query many data ID=%d", req.Id))
	}
	return &v1.ReadResponse{Api: apiVersion, ToDo: &td}, nil
}

func (s *ToDoServiceServer) Update(ctx context.Context, req *v1.UpdateRequest) (*v1.UpdateResponse, error) {
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	reminder, err := ptypes.Timestamp(req.ToDo.Reminder)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid reminder:%s", err.Error()))
	}

	res, err := c.ExecContext(
		ctx,
		"UPDATE ToDo Set`Title`=?,`Description`=?,`Reminder`=? where `ID`=?",
		req.ToDo.Title,
		req.ToDo.Description,
		reminder,
		req.ToDo.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("update error:%s", err.Error()))
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("failed to affect in update:%s", err.Error()))
	}
	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ID = %d not found", req.ToDo.Id))
	}
	return &v1.UpdateResponse{Api: apiVersion, Updated: rows}, nil
}

func (s *ToDoServiceServer) Delete(ctx context.Context, req *v1.DeleteRequest) (*v1.DeleteResponse, error) {
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	res, err := c.ExecContext(ctx, "DELETE FROM ToDo where `ID`=?", req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("delete error:%s", err.Error()))
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("delete error in affect:%s", err.Error()))
	}
	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ID = %d not found", req.Id))
	}
	return &v1.DeleteResponse{Api: apiVersion, Deleted: rows}, nil
}

func (s *ToDoServiceServer) ReadAll(ctx context.Context, req *v1.ReadAllRequest) (*v1.ReadAllResponse, error) {
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	rows, err := c.QueryContext(ctx, "SELECT `ID`,`Title`,`Description`,`Reminder` FROM ToDo")
	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("query error:%s", err.Error()))
	}
	defer rows.Close()

	var remainder time.Time
	list := []*v1.ToDo{}
	for rows.Next() {
		td := new(v1.ToDo)
		if err := rows.Scan(&td.Id, &td.Title, &td.Description, &remainder); err != nil {
			return nil, status.Error(codes.Unknown, fmt.Sprintf("query error:%s", err.Error()))
		}
		td.Reminder, err = ptypes.TimestampProto(remainder)
		if err != nil {
			return nil, status.Error(codes.Unknown, fmt.Sprintf("reminder error:%s", err.Error()))
		}
		list = append(list, td)
	}
	if err := rows.Err(); err != nil {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("get data error:%s", err.Error()))
	}
	return &v1.ReadAllResponse{Api: apiVersion, ToDos: list}, nil

}
