package commands_test

import (
	"errors"

	"github.com/pivotal-cf/om/api"
	"github.com/pivotal-cf/om/commands"
	"github.com/pivotal-cf/om/commands/fakes"
	"github.com/pivotal-cf/om/models"
	presenterfakes "github.com/pivotal-cf/om/presenters/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AvailableProducts", func() {
	var (
		apService     *fakes.AvailableProductsService
		fakePresenter *presenterfakes.FormattedPresenter
		logger        *fakes.Logger

		command *commands.AvailableProducts
	)

	BeforeEach(func() {
		apService = &fakes.AvailableProductsService{}
		fakePresenter = &presenterfakes.FormattedPresenter{}
		logger = &fakes.Logger{}

		command = commands.NewAvailableProducts(apService, fakePresenter, logger)
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
						Name:    "second-product",
						Version: "4.5.6",
					},
				},
			}, nil)
		})

		It("lists the available products", func() {
			err := executeCommand(command, []string{}, nil)
			Expect(err).NotTo(HaveOccurred())

			Expect(fakePresenter.PresentAvailableProductsCallCount()).To(Equal(1))
			products := fakePresenter.PresentAvailableProductsArgsForCall(0)
			Expect(products).To(ConsistOf(
				models.Product{
					Name:    "first-product",
					Version: "1.2.3",
				},
				models.Product{
					Name:    "second-product",
					Version: "4.5.6",
				},
			))
		})

		Context("when the json flag is provided", func() {
			It("sets the format to json on the presenter", func() {
				err := executeCommand(command, []string{"--format", "json"}, nil)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakePresenter.SetFormatCallCount()).To(Equal(1))
				Expect(fakePresenter.SetFormatArgsForCall(0)).To(Equal("json"))
			})
		})

		Context("when there are no products to list", func() {
			It("prints a helpful message instead of a table", func() {
				command := commands.NewAvailableProducts(apService, fakePresenter, logger)

				apService.ListAvailableProductsReturns(api.AvailableProductsOutput{}, nil)

				err := executeCommand(command, []string{}, nil)
				Expect(err).NotTo(HaveOccurred())

				Expect(logger.PrintfArgsForCall(0)).To(Equal("no available products found"))
				Expect(fakePresenter.PresentAvailableProductsCallCount()).To(Equal(0))
			})
		})

		Context("when the service fails to return the list", func() {
			It("returns the error", func() {
				command := commands.NewAvailableProducts(apService, fakePresenter, logger)

				apService.ListAvailableProductsReturns(api.AvailableProductsOutput{}, errors.New("blargh"))

				err := executeCommand(command, []string{}, nil)
				Expect(err).To(MatchError("blargh"))
			})
		})

		Context("when an unknown flag is passed", func() {
			It("returns an error", func() {
				err := executeCommand(command, []string{"--unknown-flag"}, nil)
				Expect(err).To(MatchError("unknown flag `unknown-flag'"))
			})
		})
	})
})
