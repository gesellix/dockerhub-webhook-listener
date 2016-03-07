package listener

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/network"
	"github.com/docker/go-connections/nat"
	"golang.org/x/net/context"
)

type Reloader struct{}

func (r *Reloader) Call(msg HubMessage) error {
	log.Println("reload in progress...")
	out, err := exec.Command("../reload.sh", msg.Repository.RepoName).Output()
	if err != nil {
		log.Println("reload error!")
		log.Println(err)
		return err
	}
	log.Println(string(out))

	log.Printf("reload done.")

	return nil
}

func (r *Reloader) Call2(msg HubMessage) error {
	log.Println("received message to reload ...")

	log.Printf("certPath %q, tls %v, host %v, api-version %v", os.Getenv("DOCKER_CERT_PATH"), os.Getenv("DOCKER_TLS_VERIFY"), os.Getenv("DOCKER_HOST"), os.Getenv("DOCKER_API_VERSION"))
	cli, err := client.NewEnvClient()
	//defaultHeaders := map[string]string{"User-Agent": "webhook-reloader"}
	//cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.22", nil, defaultHeaders)
	if err != nil {
		log.Print(err)
		return err
	}

	image := msg.Repository.RepoName
	tag := "latest"
	log.Printf("pull image %q with tag %q ...", image, tag)
	//ctx, cancel := context.WithTimeout(context.Background(), time.Duration(500) * time.Millisecond)
	//defer cancel()
	rc, err := cli.ImagePull(
		context.Background(),
		types.ImagePullOptions{
			ImageID: msg.Repository.RepoName,
			Tag:     tag}, nil)
	if err != nil {
		log.Print(err)
		return err
	}
	defer rc.Close()
	dec := json.NewDecoder(rc)
	for {
		var message jsonmessage.JSONMessage
		if err := dec.Decode(&message); err != nil {
			if err == io.EOF {
				break
			}
			log.Print(err)
			return err
		}
		log.Printf("%s", message)
	}

	containerName := "test"
	previousContainer := types.Container{}

	psOptions := types.ContainerListOptions{All: true}
	containers, err := cli.ContainerList(psOptions)
	if err != nil {
		log.Print(err)
		return err
	}
	for _, c := range containers {
		for _, name := range c.Names {
			log.Printf("%q/%q", c.ID, name)
			if name == fmt.Sprintf("/%s", containerName) {
				previousContainer = c
			}
		}
	}
	log.Printf("prev container %v", previousContainer)

	err = cli.ContainerStop(containerName, 10)
	if err != nil {
		log.Printf("stop %q: %v", containerName, err)
		//return err
	}

	rmOptions := types.ContainerRemoveOptions{ContainerID: containerName}
	err = cli.ContainerRemove(rmOptions)
	if err != nil {
		log.Printf("rm %q: %v", containerName, err)
		//return err
	}

	newContainerConfig := types.ContainerCreateConfig{}
	newContainerConfig.Config.Image = image

	port, err := nat.NewPort("8080", "http")
	exposedPorts := make(nat.PortSet)
	exposedPorts[port] = struct{}{}
	cli.ContainerCreate(
		&container.Config{
			Image:        image,
			ExposedPorts: exposedPorts},
		&container.HostConfig{},
		&network.NetworkingConfig{},
		containerName)

	log.Printf("done.")

	return nil
}
