package aws

import (
	"fmt"
	"os/exec"
)

// openURL opens the specified URL in the default browser of the user.
// https://gist.github.com/sevkin/9798d67b2cb9d07cb05f89f14ba682f8
// https://stackoverflow.com/questions/39320371/how-start-web-server-to-open-page-in-browser-in-golang
func openURL(url string) error {
	return exec.Command("open", url).Start()
}

func sendNotification(title string, text string) error {
	message := fmt.Sprintf("display notification \"%s\" with title \"%s\"", text, title)
	return exec.Command("osascript", "-e", message).Start()
}
