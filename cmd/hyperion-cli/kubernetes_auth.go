package main

// The following import seems to be necessary in order for me to access the
// Ethos Kubernetes integration cluster in Azure
// https://wiki.corp.adobe.com/display/CoreServicesTeam/Ethos+Kubernetes+Integration+Cluster
import _ "k8s.io/client-go/plugin/pkg/client/auth"
