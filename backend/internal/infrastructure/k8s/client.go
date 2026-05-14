package k8s

import (
    "os"
    "path/filepath"

    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
)

type Client struct {
    Clientset kubernetes.Interface
}

func NewKubernetesClient() (*Client, error) {
    config, err := rest.InClusterConfig()
    if err != nil {
        kubeconfig := filepath.Join(homeDir(), ".kube", "config")
        config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
        if err != nil {
            return nil, err
        }
    }
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        return nil, err
    }
    return &Client{Clientset: clientset}, nil
}

func homeDir() string {
    if h, err := os.UserHomeDir(); err == nil {
        return h
    }
    return "."
}
