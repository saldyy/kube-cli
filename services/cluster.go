package services

import (
	"fmt"
	"os"
	"strings"
)

const PROFILE_NAME = "saldyy"

type Dependencies struct {
	HelmRepo string
	Name     string
}

var dependencies = []Dependencies{
	{
		HelmRepo: "https://istio-release.storage.googleapis.com/charts",
		Name:     "istio",
	},
	{
		HelmRepo: "https://kubernetes-sigs.github.io/metrics-server/",
		Name:     "metric-server",
	},
	{
		HelmRepo: "https://charts.jetstack.io",
		Name:     "jetstack",
	},
}

func InitCluster() {
	ResumeCluster()

	UpdateHelmCharts()

	installIstio()
	installCertManager()
}

func ResumeCluster() {
	homeDir := os.Getenv("HOME")

	args := []string{
		"start",
		"-p",
		PROFILE_NAME,
		"--cpus=4",
		"--memory=4000mb",
		"--kubernetes-version=v1.31.0",
		"--mount",
		fmt.Sprintf("--mount-string=%s/.kube/local-cluster:/var/lib/minikube/certs/keys/irsa-key", homeDir),
		"--extra-config=apiserver.service-account-key-file=/var/lib/minikube/certs/keys/irsa-key/sa-signer-pkcs8.pub",
		"--extra-config=apiserver.service-account-signing-key-file=/var/lib/minikube/certs/keys/irsa-key/sa-signer.key",
		"--extra-config=apiserver.service-account-issuer=https://s3.ap-southeast-1.amazonaws.com/local-cluster-identity",
		"--extra-config=apiserver.api-audiences=sts.amazonaws.com",
	}

	fmt.Printf("Command: %s\n", strings.Join(args, " "))
	RunCommand("minikube", CommandOptions{Args: args, WithOutput: true})
}

func UpdateCluster() {
	UpdateHelmCharts()
	installIstio()
	installCertManager()
}

func DestroyCluster() {
	args := []string{
		"delete",
		"-p",
		PROFILE_NAME,
	}

	RunCommand("minikube", CommandOptions{Args: args, WithOutput: true})
}

func UpdateHelmCharts() {
	fmt.Printf("Updating Helm repositories...\n")
	for _, d := range dependencies {
		args := []string{
			"repo", "add", d.Name, d.HelmRepo,
		}
		RunCommand("helm", CommandOptions{Args: args})
	}

	RunCommand("helm", CommandOptions{Args: []string{"repo", "update"}})
	fmt.Printf("Finished update Helm repositories...\n")
}

func installIstio() {
	// Install base components
	fmt.Printf("Installing Istio...\n")
	args := []string{
		"install",
		"istio-base",
		"istio/base",
		"-n",
		"istio-system",
		"--create-namespace",
	}
	RunCommand("helm", CommandOptions{Args: args})

	// Install Istio discovery
	args = []string{
		"install",
		"istiod",
		"istio/istiod",
		"-n",
		"istio-system",
	}
	RunCommand("helm", CommandOptions{Args: args})

	// Install Istio Ingress gateway
	args = []string{
		"install",
		"istio-ingress",
		"istio/gateway",
		"-n",
		"istio-system",
	}
	RunCommand("helm", CommandOptions{Args: args})
	fmt.Printf("Finish install Istio...\n")
}

func installCertManager() {
	// Install base components
	fmt.Printf("Installing Cert Manager...\n")
	args := []string{
		"install",
		"cert-manager",
		"jetstack/cert-manager",
		"-n",
		"cert-manager",
		"--create-namespace",
		"-f",
		"./manifests/cert-manager.yaml",
	}
	RunCommand("helm", CommandOptions{Args: args})
	fmt.Printf("Finish install Cert Manager...\n")
}
