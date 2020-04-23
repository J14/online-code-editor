package run

import (
	"context"
	"log"
	"io/ioutil"
	"bytes"

	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

func Run(code string) string {
	ctx := context.Background()

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	_, err = cli.ImagePull(ctx, "python", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "python",
		Cmd: []string{"python", "-c", code},
		Tty: true,
	}, nil, nil, "")

	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	result, err := ioutil.ReadAll(out)
	if err != nil {
		log.Fatal(err)
	}

	buf := bytes.NewBuffer(result)
	return buf.String()
}