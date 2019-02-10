package main

type vulnerableDockerAPI struct {
	Endpoint      string
	Host          string
	DockerVersion string
	Port          uint16
	SocketError   error

	Info       dockerInfo
	Containers []dockerContainer
	Images     []string
}

type dockerContainer struct {
	Image  string
	Mounts string
	ID     string
}

type dockerInfo struct {
	ContainersStopped int
	ContainersRunning int
	Images            int
	OS                string
}
