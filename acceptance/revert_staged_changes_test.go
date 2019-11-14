package acceptance

import (
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/ghttp"
	"github.com/pivotal-cf/om/commands"
	"net/http"
	"time"

	. "github.com/onsi/gomega"
)

var _ = FDescribe("revert-staged-changes command", func() {
	var (
		server *ghttp.Server
		main   commands.Main
		stdout *gbytes.Buffer
		stderr *gbytes.Buffer
	)

	BeforeEach(func() {
		server = createTLSServer()
		stdout = gbytes.NewBuffer()
		stderr = gbytes.NewBuffer()
		main = commands.NewMain(
			stdout,
			stderr,
			nil,
			"",
			time.Second,
		)
	})

	It("reverts the staged changes on the Ops Manager", func() {
		ensureHandler := &ensureHandler{}
		server.AppendHandlers(
			ensureHandler.Ensure(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("DELETE", "/api/v0/staged"),
					ghttp.RespondWith(http.StatusOK, ""),
				),
			)...,
		)

		err := main.Execute(
			"--target", server.URL(),
			"--username", "some-username",
			"--password", "some-password",
			"--skip-ssl-validation",
			"revert-staged-changes",
		)

		Expect(err).ToNot(HaveOccurred())
		Expect(ensureHandler.Handlers()).To(HaveLen(0))
		Eventually(stdout).Should(gbytes.Say("Changes Reverted."))
	})

	When("there are no changes to revert", func() {
		It("does nothing", func() {
			ensureHandler := &ensureHandler{}
			server.AppendHandlers(
				ensureHandler.Ensure(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("DELETE", "/api/v0/staged"),
						ghttp.RespondWith(http.StatusNotModified, ""),
					),
				)...,
			)

			err := main.Execute(
				"--target", server.URL(),
				"--username", "some-username",
				"--password", "some-password",
				"--skip-ssl-validation",
				"revert-staged-changes",
			)

			Expect(err).ToNot(HaveOccurred())

			Expect(ensureHandler.Handlers()).To(HaveLen(0))
			Eventually(stdout).Should(gbytes.Say("No changes to revert."))
		})
	})

	When("the revert is forbidden", func() {
		It("errors", func() {
			ensureHandler := &ensureHandler{}
			server.AppendHandlers(
				ensureHandler.Ensure(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("DELETE", "/api/v0/staged"),
						ghttp.RespondWith(http.StatusForbidden, ""),
					),
				)...,
			)

			err := main.Execute(
				"--target", server.URL(),
				"--username", "some-username",
				"--password", "some-password",
				"--skip-ssl-validation",
				"revert-staged-changes",
			)

			Expect(err).To(MatchError(ContainSubstring("revert staged changes command failed: request failed: unexpected response from /api/v0/staged")))
			Expect(ensureHandler.Handlers()).To(HaveLen(0))
		})
	})
})
