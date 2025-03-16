package services

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
}

func InitCluster() {
	args := []string{
		"start",
		"-p",
		PROFILE_NAME,
		"--addons",
		"metallb",
		"--cpus=4",
		"--memory=4000mb",
		"--kubernetes-version=v1.31.0",
	}

	RunCommand("minikube", CommandOptions{Args: args, WithOutput: true})

  UpdateHelmCharts()
  installIstio()
}

func DestroyCluster() {
	args := []string{
		"destroy",
		"-p",
		PROFILE_NAME,
	}

	RunCommand("minikube", CommandOptions{Args: args, WithOutput: true})
}

func UpdateHelmCharts() {
	for _, d := range dependencies {
		args := []string{
			"repo", "add", d.Name, d.HelmRepo,
		}
		RunCommand("helm", CommandOptions{Args: args})
	}

	RunCommand("helm", CommandOptions{Args: []string{"repo", "update"}})
}

func installIstio() {
	// Install base components
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

}
