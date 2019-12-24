package docs_test

import (
	"github.com/onsi/gomega/gexec"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"

	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Documentation coverage", func() {
	Context("Commands doc", func() {
		It("Mentions every command we distribute", func() {
			commands := getCommandNames()
			taskDoc := readFile("../docs/README.md")

			var missing []string
			for _, command := range commands {
				if !strings.Contains(taskDoc, command) {
					missing = append(missing, command)
				}
			}

			Expect(missing).To(HaveLen(0), fmt.Sprintf("docs/README.md should have: \n%s\n run `go run docsgenerator/update-docs.go` to fix", strings.Join(missing, "\n")))
		})
	})
})

func readFile(docName string) (docContents string) {
	docPath, err := filepath.Abs(docName)
	Expect(err).ToNot(HaveOccurred())

	docContentsBytes, err := ioutil.ReadFile(docPath)
	docContents = string(docContentsBytes)
	Expect(err).ToNot(HaveOccurred())

	return docContents
}

func getCommandNames() []string {
	command := exec.Command(pathToMain, "--help")

	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).ToNot(HaveOccurred())

	session.Wait()

	output := strings.Split(string(session.Out.Contents()), "\n")

	var isCommand bool
	var commands []string
	for _, commandLine := range output {
		if strings.Contains(commandLine, "Commands:") && !isCommand {
			isCommand = true
			continue
		}

		if isCommand && commandLine != "" {
			splitCommandLine := strings.Fields(commandLine)
			commands = append(commands, splitCommandLine[0])
		}
	}

	return commands
}
