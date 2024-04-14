package tbd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/raohwork/task"
)

func ExampleNew() {
	response, res, rej := New[*http.Response]()

	resolver := task.Task(func(ctx context.Context) error {
		req, err := http.NewRequestWithContext(ctx, "GET", "https://google.com", nil)
		if err != nil {
			return rej(err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return rej(err)
		}

		return res(resp)
	})

	saveToDB := func(ctx context.Context) error {
		resp, err := response.Get(ctx)
		if err != nil {
			return err
		}

		// parse the response and save it to db
		_ = resp
		return nil
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	err := task.Iter(resolver, saveToDB).Run(ctx)
	fmt.Print(err)
}
