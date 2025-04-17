package services

import (
	"fmt"
	"os"
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
	{
		HelmRepo: "https://jkroepke.github.io/helm-charts/",
		Name:     "jkroepke",
	},
}

func InitCluster() {
	ResumeCluster()

	UpdateHelmCharts()

	installIstio()
	installCertManager()
	installEKSPodIdentityWebhook()
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

	RunWithLiveOutput("minikube", args)
}

func UpdateCluster() {
	UpdateHelmCharts()
	installIstio()
	installCertManager()
	installEKSPodIdentityWebhook()
}

func DestroyCluster() {
	args := []string{
		"delete",
		"-p",
		PROFILE_NAME,
	}

	RunWithLiveOutput("minikube", args)
}

func UpdateHelmCharts() {
	fmt.Printf("Updating Helm repositories...\n")
	for _, d := range dependencies {
		args := []string{
			"repo", "add", d.Name, d.HelmRepo,
		}
		RunWithLiveOutput("helm", args)
	}

	RunWithLiveOutput("helm", []string{"repo", "update"})
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
	RunWithLiveOutput("helm", args)

	// Install Istio discovery
	args = []string{
		"install",
		"istiod",
		"istio/istiod",
		"-n",
		"istio-system",
	}
	RunWithLiveOutput("helm", args)

	// Install Istio Ingress gateway
	args = []string{
		"install",
		"istio-ingress",
		"istio/gateway",
		"-n",
		"istio-system",
	}
	RunWithLiveOutput("helm", args)
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
	RunWithLiveOutput("helm", args)
	fmt.Printf("Finish install Cert Manager...\n")
}

func installEKSPodIdentityWebhook() {
	// Install base components
	fmt.Printf("Installing EKS Pod Identity Webhook...\n")
	args := []string{
		"install",
		"amazon-eks-pod-identity-webhook",
		"jkroepke/amazon-eks-pod-identity-webhook",
	}
	RunWithLiveOutput("helm", args)
	fmt.Printf("Finish install EKS Pod Identity Webhook...\n")
}
