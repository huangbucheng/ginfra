package main

import (
	"context"
	"fmt"

	"github.com/spf13/pflag"

	"ginfra/plugin/k8sclient"
)

var (
	cfg       = pflag.StringP("config", "c", "", "kubernetes apiserver config file path.")
	cmd       = pflag.StringP("command", "C", "", "tool command:createjob|getjob.")
	namespace = pflag.StringP("namespace", "n", "", "kubernetes namespace.")
)

func main() {
	pflag.Parse()

	kclient := &k8sclient.KubeClient{}
	if len(*cfg) == 0 {
		kclient.InClusterConfig()
	}

	err := kclient.WithKubeConfig(*cfg)
	if err != nil {
		panic(err)
	}

	if *cmd == "createjob" {
		createjob(kclient)
	} else if *cmd == "getjob" {
		getjob(kclient)
	}
}

func createjob(kclient *k8sclient.KubeClient) {
	req := &k8sclient.JobRequest{
		Namespace:               *namespace,
		JobName:                 "demo-job",
		Image:                   "xxxx",
		CpuRequest:              "700m",
		MemoryRequest:           "512Mi",
		CpuLimit:                "700m",
		MemoryLimit:             "512Mi",
		TTLSecondsAfterFinished: 300,
		Mounts: []k8sclient.VolumeMount{{
			Name:      "code",
			MountPath: "/usr/local/service/runner/code",
			HostPath:  "/tmp/test",
		}},
		Envs: map[string]string{
			"RUNNER_DIR": "code",
		},
	}
	job, err := kclient.CreateFootballJob(context.TODO(), req)
	if err != nil {
		panic(err)
	}
	fmt.Printf("New job: %v\n", job)
}

func getjob(kclient *k8sclient.KubeClient) {
	job, err := kclient.GetJob(context.TODO(), *namespace, "demo-job")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Get job: %v\n", job)
	fmt.Printf("job status: %v\n", job.Status)
}
