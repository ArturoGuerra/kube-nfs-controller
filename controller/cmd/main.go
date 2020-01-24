package main

import (
    "os"
    "github.com/golang/glog"
    "github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/controller"
    "k8s.io/apimachinery/pkg/util/wait"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/utils/exec"
)

const (
    provisioner = "arturoguerra/nfs"
)

var (
    server = os.Getenv("SERVER")
    mountPath = os.Getenv("MOUNTPATH")
    fakePath = os.Getenv("FAKEPATH")
)

type NfsProvisioner struct {
    runner exec.Interface
}

func NewNfsProvisioner() controller.Provisioner {
    return &NfsProvisioner{
        runner: exec.New(),
    }
}

func main() {
    config, err := rest.InClusterConfig()
    if err != nil {
        glog.Fatalf("Failed to create config: %v", err)
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        glog.Fatalf("Failed to create client: %v", err)
    }


    serverVersion, err := clientset.Discovery().ServerVersion()
    if err != nil {
        glog.Fatalf("Eror getting server version: %v", err)
    }

    nfsProvisioner := NewNfsProvisioner()

    // Provision Controller setup
    pc := controller.NewProvisionController(
        clientset,
        provisioner,
        nfsProvisioner,
        serverVersion.GitVersion,
    )

    pc.Run(wait.NeverStop)
}


