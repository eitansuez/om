package commands

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pivotal-cf/jhanda"
	"github.com/pivotal-cf/om/api"
	"github.com/pivotal-cf/om/models"
	"github.com/pivotal-cf/om/presenters"
)

type Products struct {
	apService availableProductsService
	drService diagnosticReportService
	presenter presenters.FormattedPresenter
}

func NewProducts(
	apService availableProductsService,
	drService diagnosticReportService,
	presenter presenters.FormattedPresenter) Products {
	return Products{
		apService: apService,
		drService: drService,
		presenter: presenter,
	}
}

func (p Products) Execute(args []string) error {

	output, err := p.apService.ListAvailableProducts()
	if err != nil {
		return err
	}

	var availableProducts []models.Product
	for _, product := range output.ProductsList {
		availableProducts = append(availableProducts, models.Product{
			Name:    product.Name,
			Version: product.Version,
		})
	}

	diagnosticReport, err := p.drService.GetDiagnosticReport()
	if err != nil {
		return fmt.Errorf("failed to retrieve diagnostic report %s", err)
	}

	stagedProducts := diagnosticReport.StagedProducts
	deployedProducts := diagnosticReport.DeployedProducts

	p.presenter.SetFormat("table")
	allProductVersions := combinedProductVersions(availableProducts, stagedProducts, deployedProducts)
	p.presenter.PresentProducts(allProductVersions)

	return nil
}

func combinedProductVersions(availableProducts []models.Product,
	stagedProducts []api.DiagnosticProduct,
	deployedProducts []api.DiagnosticProduct) []models.ProductVersions {

	availMap := map[string][]string{}

	for _, p := range availableProducts {
		availMap[p.Name] = append(availMap[p.Name], p.Version)
	}

	stagedMap := map[string]string{}
	for _, p := range stagedProducts {
		stagedMap[p.Name] = p.Version
	}

	deployedMap := map[string]string{}
	for _, p := range deployedProducts {
		deployedMap[p.Name] = p.Version
	}

	allProducts := map[string]bool{}

	for name := range availMap {
		allProducts[name] = true
	}
	for name := range stagedMap {
		allProducts[name] = true
	}
	for name := range deployedMap {
		allProducts[name] = true
	}

	productVersions := []models.ProductVersions{}
	for name := range allProducts {

		availableVersions := ""
		if len(availMap[name]) > 0 {
			availableVersions = strings.Join(availMap[name], ", ")
		}

		productVersions = append(productVersions, models.ProductVersions{
			Name:              name,
			AvailableVersions: availableVersions,
			StagedVersion:     stagedMap[name],
			DeployedVersion:   deployedMap[name],
		})
	}

	sort.Slice(productVersions, func(i, j int) bool {
		return strings.ToLower(productVersions[i].Name) < strings.ToLower(productVersions[j].Name)
	})

	return productVersions
}

func (p Products) Usage() jhanda.Usage {
	return jhanda.Usage{
		Description:      "This authenticated command lists all available, staged, and deployed products.",
		ShortDescription: "list available, staged, and deployed products",
	}
}
