package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	protobuf "github.com/wilkermichael/go-proto-grpc/internal/protobuf_v1"
	"google.golang.org/grpc"
)

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "missing subcommand: list or add")
		os.Exit(1)
	}

	conn, err := grpc.Dial(":8888", grpc.WithInsecure())
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not connect to backend: %v\n", err)
		os.Exit(1)
	}
	client := protobuf.NewTasksClient(conn)

	switch cmd := flag.Arg(0); cmd {
	case "list":
		err = list(context.Background(), client)
	case "add":
		err = add(context.Background(), client, strings.Join(flag.Args()[1:], " "))
	default:
		err = fmt.Errorf("unknown subcommand %s", cmd)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func add(ctx context.Context, client protobuf.TasksClient, text string) error {
	_, err := client.Add(ctx, &protobuf.Text{Text: text})
	if err != nil {
		return fmt.Errorf("could not add task in the backend: %v", err)
	}

	fmt.Println("task added successfully")
	return nil
}

func list(ctx context.Context, client protobuf.TasksClient) error {
	l, err := client.List(ctx, &protobuf.Void{})
	if err != nil {
		return fmt.Errorf("could not fetch tasks: %v", err)
	}
	for _, t := range l.Tasks {
		if t.Done {
			fmt.Printf("👍")
		} else {
			fmt.Printf("😱")
		}
		fmt.Printf(" %s\n", t.Text)
	}
	return nil
}
