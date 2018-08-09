package main

// The following import seems to be necessary in order for me to access the
// Ethos Kubernetes integration cluster in Azure
import _ "k8s.io/client-go/plugin/pkg/client/auth"
