package aws

import (
	"os/exec"
)

// openURL opens the specified URL in the default browser of the user.
// https://gist.github.com/sevkin/9798d67b2cb9d07cb05f89f14ba682f8
// https://stackoverflow.com/questions/39320371/how-start-web-server-to-open-page-in-browser-in-golang
func openURL(url string) error {
	args := []string{"/c", "start", url}
	// args[0] is used for 'start' command argument, to prevent issues with URLs starting with a quote
	args = append(args[:1], append([]string{""}, args[1:]...)...)
	return exec.Command("cmd", args...).Start()
}

func sendMacNotification(title string, text string) error {
	// Not implemented for this platform
	return nil
}
