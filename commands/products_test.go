package commands_test

import (
	"errors"

	"github.com/pivotal-cf/jhanda"
	"github.com/pivotal-cf/om/api"
	"github.com/pivotal-cf/om/commands"
	"github.com/pivotal-cf/om/commands/fakes"
	"github.com/pivotal-cf/om/models"
	presenterfakes "github.com/pivotal-cf/om/presenters/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Products", func() {
	var (
		apService     *fakes.AvailableProductsService
		drService     *fakes.DiagnosticReportService
		fakePresenter *presenterfakes.FormattedPresenter

		command commands.Products
	)

	BeforeEach(func() {
		apService = &fakes.AvailableProductsService{}
		drService = &fakes.DiagnosticReportService{}
		fakePresenter = &presenterfakes.FormattedPresenter{}

		command = commands.NewProducts(apService, drService, fakePresenter)
	})

	Describe("Execute", func() {
		BeforeEach(func() {
			apService.ListAvailableProductsReturns(api.AvailableProductsOutput{
				ProductsList: []api.ProductInfo{
					api.ProductInfo{
						Name:    "first-product",
						Version: "1.2.3",
					},
					api.ProductInfo{
						Name:    "first-product",
						Version: "1.2.4",
					},
					api.ProductInfo{
						Name:    "second-product",
						Version: "4.5.6",
					},
				},
			}, nil)
		})

		It("lists available products, nothing staged or deployed yet", func() {

			drService.GetDiagnosticReportReturns(api.DiagnosticReport{
				StagedProducts:   []api.DiagnosticProduct{},
				DeployedProducts: []api.DiagnosticProduct{},
			}, nil)

			err := command.Execute([]string{})
			Expect(err).ToNot(HaveOccurred())

			Expect(fakePresenter.PresentProductsCallCount()).To(Equal(1))
			productVersions := fakePresenter.PresentProductsArgsForCall(0)
			Expect(productVersions).To(ConsistOf(
				models.ProductVersions{
					Name:              "first-product",
					AvailableVersions: "1.2.3, 1.2.4",
					StagedVersion:     "",
					DeployedVersion:   "",
				},
				models.ProductVersions{
					Name:              "second-product",
					AvailableVersions: "4.5.6",
					StagedVersion:     "",
					DeployedVersion:   "",
				},
			))
		})

		It("lists available and staged products, nothing deployed yet", func() {

			drService.GetDiagnosticReportReturns(api.DiagnosticReport{
				StagedProducts: []api.DiagnosticProduct{
					{
						Name:    "first-product",
						Version: "1.2.3",
					},
				},
				DeployedProducts: []api.DiagnosticProduct{},
			}, nil)

			err := command.Execute([]string{})
			Expect(err).ToNot(HaveOccurred())

			Expect(fakePresenter.PresentProductsCallCount()).To(Equal(1))
			productVersions := fakePresenter.PresentProductsArgsForCall(0)
			Expect(productVersions).To(ConsistOf(
				models.ProductVersions{
					Name:              "first-product",
					AvailableVersions: "1.2.3, 1.2.4",
					StagedVersion:     "1.2.3",
					DeployedVersion:   "",
				},
				models.ProductVersions{
					Name:              "second-product",
					AvailableVersions: "4.5.6",
					StagedVersion:     "",
					DeployedVersion:   "",
				},
			))
		})

		It("lists available, staged and deployed products", func() {

			drService.GetDiagnosticReportReturns(api.DiagnosticReport{
				StagedProducts: []api.DiagnosticProduct{
					{
						Name:    "first-product",
						Version: "1.2.3",
					},
				},
				DeployedProducts: []api.DiagnosticProduct{
					{
						Name:    "first-product",
						Version: "1.2.3",
					},
				},
			}, nil)

			err := command.Execute([]string{})
			Expect(err).ToNot(HaveOccurred())

			Expect(fakePresenter.PresentProductsCallCount()).To(Equal(1))
			productVersions := fakePresenter.PresentProductsArgsForCall(0)
			Expect(productVersions).To(ConsistOf(
				models.ProductVersions{
					Name:              "first-product",
					AvailableVersions: "1.2.3, 1.2.4",
					StagedVersion:     "1.2.3",
					DeployedVersion:   "1.2.3",
				},
				models.ProductVersions{
					Name:              "second-product",
					AvailableVersions: "4.5.6",
					StagedVersion:     "",
					DeployedVersion:   "",
				},
			))
		})

		When("there are no available products to list, but there are staged and deployed products", func() {
			It("prints a helpful message instead of a table", func() {
				command := commands.NewProducts(apService, drService, fakePresenter)

				apService.ListAvailableProductsReturns(api.AvailableProductsOutput{}, nil)
				drService.GetDiagnosticReportReturns(api.DiagnosticReport{
					StagedProducts: []api.DiagnosticProduct{
						{
							Name:    "p-bosh",
							Version: "2.8.3-build.217",
						},
					},
					DeployedProducts: []api.DiagnosticProduct{
						{
							Name:    "p-bosh",
							Version: "2.8.3-build.217",
						},
					},
				}, nil)

				err := command.Execute([]string{})
				Expect(err).ToNot(HaveOccurred())

				Expect(fakePresenter.PresentProductsCallCount()).To(Equal(1))
				productVersions := fakePresenter.PresentProductsArgsForCall(0)
				Expect(productVersions).To(ConsistOf(
					models.ProductVersions{
						Name:              "p-bosh",
						AvailableVersions: "",
						StagedVersion:     "2.8.3-build.217",
						DeployedVersion:   "2.8.3-build.217",
					},
				))

			})
		})

		When("the available products service fails", func() {
			It("returns the error", func() {
				command := commands.NewProducts(apService, drService, fakePresenter)

				apService.ListAvailableProductsReturns(api.AvailableProductsOutput{}, errors.New("blargh"))

				err := command.Execute([]string{})
				Expect(err).To(MatchError("blargh"))
			})
		})

		When("the diagnostic report service fails", func() {
			It("returns the error", func() {
				command := commands.NewProducts(apService, drService, fakePresenter)

				drService.GetDiagnosticReportReturns(api.DiagnosticReport{}, errors.New("ahalan"))

				err := command.Execute([]string{})
				Expect(err).To(MatchError("failed to retrieve diagnostic report ahalan"))
			})
		})
	})

	Describe("Usage", func() {
		It("returns usage information for the command", func() {
			command := commands.NewProducts(nil, nil, nil)
			Expect(command.Usage()).To(Equal(jhanda.Usage{
				Description:      "This authenticated command lists all available, staged, and deployed products.",
				ShortDescription: "list available, staged, and deployed products",
			}))
		})
	})
})
