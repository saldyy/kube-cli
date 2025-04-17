package services

import "strings"

func CheckInstalledCharts(s string, ns string) bool {
	r, _ := RunSilent("helm", []string{"ls", "-n", ns})

	return strings.Contains(r, s)
}

func InstallOrUpdateChart(c string)
