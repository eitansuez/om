package commands

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/pivotal-cf/jhanda"
	"github.com/pivotal-cf/om/api"
	"github.com/pivotal-cf/om/extractor"
	"github.com/pivotal-cf/om/formcontent"
	"github.com/pivotal-cf/om/interpolate"
	"github.com/pivotal-cf/om/network"
	"github.com/pivotal-cf/om/presenters"
	"github.com/pivotal-cf/om/progress"
	"github.com/pivotal-cf/om/renderers"
	"github.com/pivotal/uilive"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type mainOptions struct {
	CACert               string `yaml:"ca-cert" long:"ca-cert" env:"OM_CA_CERT" description:"OpsManager CA certificate path or value"`
	ClientID             string `yaml:"client-id"             short:"c"  long:"client-id"             env:"OM_CLIENT_ID"                           description:"Client ID for the Ops Manager VM (not required for unauthenticated commands)"`
	ClientSecret         string `yaml:"client-secret"         short:"s"  long:"client-secret"         env:"OM_CLIENT_SECRET"                       description:"Client Secret for the Ops Manager VM (not required for unauthenticated commands)"`
	ConnectTimeout       int    `yaml:"connect-timeout"       short:"o"  long:"connect-timeout"       env:"OM_CONNECT_TIMEOUT"     default:"10"    description:"timeout in seconds to make TCP connections"`
	DecryptionPassphrase string `yaml:"decryption-passphrase" short:"d"  long:"decryption-passphrase" env:"OM_DECRYPTION_PASSPHRASE"             description:"Passphrase to decrypt the installation if the Ops Manager VM has been rebooted (optional for most commands)"`
	Env                  string `                             short:"e"  long:"env"                                                              description:"env file with login credentials"`
	Help                 bool   `                             short:"h"  long:"help"                                             default:"false" description:"prints this usage information"`
	Password             string `yaml:"password"              short:"p"  long:"password"              env:"OM_PASSWORD"                            description:"admin password for the Ops Manager VM (not required for unauthenticated commands)"`
	RequestTimeout       int    `yaml:"request-timeout"       short:"r"  long:"request-timeout"       env:"OM_REQUEST_TIMEOUT"     default:"1800"  description:"timeout in seconds for HTTP requests to Ops Manager"`
	SkipSSLValidation    bool   `yaml:"skip-ssl-validation"   short:"k"  long:"skip-ssl-validation"   env:"OM_SKIP_SSL_VALIDATION" default:"false" description:"skip ssl certificate validation during http requests"`
	Target               string `yaml:"target"                short:"t"  long:"target"                env:"OM_TARGET"                              description:"location of the Ops Manager VM"`
	Trace                bool   `yaml:"trace"                 short:"tr" long:"trace"                 env:"OM_TRACE"                               description:"prints HTTP requests and response payloads"`
	Username             string `yaml:"username"              short:"u"  long:"username"              env:"OM_USERNAME"                            description:"admin username for the Ops Manager VM (not required for unauthenticated commands)"`
	VarsEnv              string `                                                                     env:"OM_VARS_ENV"      experimental:"true" description:"load vars from environment variables by specifying a prefix (e.g.: 'MY' to load MY_var=value)"`
	Version              bool   `                             short:"v"  long:"version"                                          default:"false" description:"prints the om release version"`
}

type Main struct {
	stdout             io.Writer
	stderr             io.Writer
	stdin              io.Reader
	version            string
	applySleepDuration time.Duration
}

func NewMain(stdout, stderr io.Writer, stdin io.Reader, version string, applySleepDuration time.Duration) Main {
	return Main{
		stdout:             stdout,
		stderr:             stderr,
		stdin:              stdin,
		version:            version,
		applySleepDuration: applySleepDuration,
	}
}

func (m Main) Execute(args ...string) error {
	var global mainOptions

	stderr := log.New(m.stderr, "", 0)
	stdout := log.New(m.stdout, "", 0)

	args, err := jhanda.Parse(&global, args)
	if err != nil {
		return err
	}

	err = setEnvFileProperties(&global)
	if err != nil {
		return err
	}

	globalFlagsUsage, err := jhanda.PrintUsage(global)
	if err != nil {
		return err
	}

	var command string
	if len(args) > 0 {
		command, args = args[0], args[1:]
	}

	if global.Version {
		command = "version"
	}

	if global.Help {
		command = "help"
	}

	if command == "" {
		command = "help"
	}

	requestTimeout := time.Duration(global.RequestTimeout) * time.Second
	connectTimeout := time.Duration(global.ConnectTimeout) * time.Second

	var unauthenticatedClient, authedClient, authedCookieClient, unauthenticatedProgressClient, authedProgressClient httpClient
	unauthenticatedClient, _ = network.NewUnauthenticatedClient(global.Target, global.SkipSSLValidation, global.CACert, connectTimeout, requestTimeout)
	if err != nil {
		return err
	}

	authedClient, err = network.NewOAuthClient(global.Target, global.Username, global.Password, global.ClientID, global.ClientSecret, global.SkipSSLValidation, global.CACert, connectTimeout, requestTimeout)

	if err != nil {
		return err
	}

	if global.DecryptionPassphrase != "" {
		authedClient = network.NewDecryptClient(authedClient, unauthenticatedClient, global.DecryptionPassphrase, m.stderr)
	}

	authedCookieClient, err = network.NewOAuthClient(global.Target, global.Username, global.Password, global.ClientID, global.ClientSecret, global.SkipSSLValidation, "", connectTimeout, requestTimeout)
	if err != nil {
		return err
	}

	liveWriter := uilive.New()
	liveWriter.Out = m.stderr
	unauthenticatedProgressClient = network.NewProgressClient(unauthenticatedClient, progress.NewBar(), liveWriter)
	authedProgressClient = network.NewProgressClient(authedClient, progress.NewBar(), liveWriter)

	if global.Trace {
		unauthenticatedClient = network.NewTraceClient(unauthenticatedClient, m.stderr)
		unauthenticatedProgressClient = network.NewTraceClient(unauthenticatedProgressClient, m.stderr)
		authedClient = network.NewTraceClient(authedClient, m.stderr)
		authedCookieClient = network.NewTraceClient(authedCookieClient, m.stderr)
		authedProgressClient = network.NewTraceClient(authedProgressClient, m.stderr)
	}

	api := api.New(api.ApiInput{
		Client:                 authedClient,
		UnauthedClient:         unauthenticatedClient,
		ProgressClient:         authedProgressClient,
		UnauthedProgressClient: unauthenticatedProgressClient,
		Logger:                 stderr,
	})

	logWriter := NewLogWriter(m.stdout)
	tableWriter := tablewriter.NewWriter(m.stdout)

	form := formcontent.NewForm()

	metadataExtractor := extractor.MetadataExtractor{}

	presenter := presenters.NewPresenter(presenters.NewTablePresenter(tableWriter), presenters.NewJSONPresenter(m.stdout))
	envRendererFactory := renderers.NewFactory(renderers.NewEnvGetter())

	commandSet := jhanda.CommandSet{}
	commandSet["activate-certificate-authority"] = NewActivateCertificateAuthority(api, stdout)
	commandSet["apply-changes"] = NewApplyChanges(api, api, logWriter, stdout, m.applySleepDuration)
	commandSet["assign-multi-stemcell"] = NewAssignMultiStemcell(api, stdout)
	commandSet["assign-stemcell"] = NewAssignStemcell(api, stdout)
	commandSet["available-products"] = NewAvailableProducts(api, presenter, stdout)
	commandSet["bosh-env"] = NewBoshEnvironment(api, stdout, global.Target, envRendererFactory)
	commandSet["certificate-authorities"] = NewCertificateAuthorities(api, presenter)
	commandSet["certificate-authority"] = NewCertificateAuthority(api, presenter, stdout)
	commandSet["config-template"] = NewConfigTemplate(DefaultProvider())
	commandSet["configure-authentication"] = NewConfigureAuthentication(os.Environ, api, stdout)
	commandSet["configure-director"] = NewConfigureDirector(os.Environ, api, stdout)
	commandSet["configure-ldap-authentication"] = NewConfigureLDAPAuthentication(os.Environ, api, stdout)
	commandSet["configure-product"] = NewConfigureProduct(os.Environ, api, global.Target, stdout)
	commandSet["configure-saml-authentication"] = NewConfigureSAMLAuthentication(os.Environ, api, stdout)
	commandSet["create-certificate-authority"] = NewCreateCertificateAuthority(api, presenter)
	commandSet["create-vm-extension"] = NewCreateVMExtension(os.Environ, api, stdout)
	commandSet["credential-references"] = NewCredentialReferences(api, presenter, stdout)
	commandSet["credentials"] = NewCredentials(api, presenter, stdout)
	commandSet["curl"] = NewCurl(api, stdout, stderr)
	commandSet["delete-certificate-authority"] = NewDeleteCertificateAuthority(api, stdout)
	commandSet["delete-installation"] = NewDeleteInstallation(api, logWriter, stdout, os.Stdin, m.applySleepDuration)
	commandSet["delete-product"] = NewDeleteProduct(api)
	commandSet["delete-ssl-certificate"] = NewDeleteSSLCertificate(api, stdout)
	commandSet["delete-unused-products"] = NewDeleteUnusedProducts(api, stdout)
	commandSet["deployed-manifest"] = NewDeployedManifest(api, stdout)
	commandSet["deployed-products"] = NewDeployedProducts(presenter, api)
	commandSet["diagnostic-report"] = NewDiagnosticReport(presenter, api)
	commandSet["disable-director-verifiers"] = NewDisableDirectorVerifiers(presenter, api, stdout)
	commandSet["disable-product-verifiers"] = NewDisableProductVerifiers(presenter, api, stdout)
	commandSet["download-product"] = NewDownloadProduct(os.Environ, stdout, stderr, m.stderr)
	commandSet["errands"] = NewErrands(presenter, api)
	commandSet["expiring-certificates"] = NewExpiringCertificates(api, stdout)
	commandSet["export-installation"] = NewExportInstallation(api, stderr)
	commandSet["generate-certificate"] = NewGenerateCertificate(api, stdout)
	commandSet["generate-certificate-authority"] = NewGenerateCertificateAuthority(api, presenter)
	commandSet["help"] = NewHelp(m.stdout, globalFlagsUsage, commandSet)
	commandSet["import-installation"] = NewImportInstallation(form, api, global.DecryptionPassphrase, stdout)
	commandSet["installation-log"] = NewInstallationLog(api, stdout)
	commandSet["installations"] = NewInstallations(api, presenter)
	commandSet["interpolate"] = NewInterpolate(os.Environ, stdout, os.Stdin)
	commandSet["pending-changes"] = NewPendingChanges(presenter, api)
	commandSet["pre-deploy-check"] = NewPreDeployCheck(presenter, api, stdout)
	commandSet["regenerate-certificates"] = NewRegenerateCertificates(api, stdout)
	commandSet["revert-staged-changes"] = NewRevertStagedChanges(api, stdout)
	commandSet["ssl-certificate"] = NewSSLCertificate(api, presenter)
	commandSet["stage-product"] = NewStageProduct(api, stdout)
	commandSet["staged-config"] = NewStagedConfig(api, stdout)
	commandSet["staged-director-config"] = NewStagedDirectorConfig(api, stdout, stderr)
	commandSet["staged-manifest"] = NewStagedManifest(api, stdout)
	commandSet["staged-products"] = NewStagedProducts(presenter, api)
	commandSet["product-metadata"] = NewProductMetadata(stdout)
	commandSet["tile-metadata"] = NewDeprecatedProductMetadata(stdout)
	commandSet["unstage-product"] = NewUnstageProduct(api, stdout)
	commandSet["update-ssl-certificate"] = NewUpdateSSLCertificate(api, stdout)
	commandSet["upload-product"] = NewUploadProduct(form, metadataExtractor, api, stdout)
	commandSet["upload-stemcell"] = NewUploadStemcell(form, api, stdout)
	commandSet["version"] = NewVersion(m.version, m.stdout)

	err = commandSet.Execute(command, args)
	if err != nil {
		return err
	}

	return nil
}

func setEnvFileProperties(global *mainOptions) error {
	if global.Env == "" {
		return nil
	}

	var opts mainOptions
	_, err := os.Open(global.Env)
	if err != nil {
		return fmt.Errorf("env file does not exist: %s", err)
	}

	contents, err := interpolate.Execute(interpolate.Options{
		TemplateFile:  global.Env,
		EnvironFunc:   os.Environ,
		VarsEnvs:      []string{global.VarsEnv},
		ExpectAllKeys: false,
	})
	if err != nil {
		return err
	}

	err = yaml.UnmarshalStrict(contents, &opts)
	if err != nil {
		return fmt.Errorf("could not parse env file: %s", err)
	}

	if global.ClientID == "" {
		global.ClientID = opts.ClientID
	}
	if global.ClientSecret == "" {
		global.ClientSecret = opts.ClientSecret
	}
	if global.Password == "" {
		global.Password = opts.Password
	}
	if global.ConnectTimeout == 10 && opts.ConnectTimeout != 0 {
		global.ConnectTimeout = opts.ConnectTimeout
	}
	if global.RequestTimeout == 1800 && opts.RequestTimeout != 0 {
		global.RequestTimeout = opts.RequestTimeout
	}
	if global.SkipSSLValidation == false {
		global.SkipSSLValidation = opts.SkipSSLValidation
	}
	if global.Target == "" {
		global.Target = opts.Target
	}
	if global.Trace == false {
		global.Trace = opts.Trace
	}
	if global.Username == "" {
		global.Username = opts.Username
	}
	if global.DecryptionPassphrase == "" {
		global.DecryptionPassphrase = opts.DecryptionPassphrase
	}
	if global.CACert == "" {
		global.CACert = opts.CACert
	}

	err = checkForVars(global)
	if err != nil {
		return fmt.Errorf("found problem in --env file: %s", err)
	}

	return nil
}

func checkForVars(opts *mainOptions) error {
	var errBuffer []string

	interpolateRegex := regexp.MustCompile(`\(\(.*\)\)`)

	if interpolateRegex.MatchString(opts.DecryptionPassphrase) {
		errBuffer = append(errBuffer, "* use OM_DECRYPTION_PASSPHRASE environment variable for the decryption-passphrase value")
	}

	if interpolateRegex.MatchString(opts.ClientID) {
		errBuffer = append(errBuffer, "* use OM_CLIENT_ID environment variable for the client-id value")
	}

	if interpolateRegex.MatchString(opts.ClientSecret) {
		errBuffer = append(errBuffer, "* use OM_CLIENT_SECRET environment variable for the client-secret value")
	}

	if interpolateRegex.MatchString(opts.Password) {
		errBuffer = append(errBuffer, "* use OM_PASSWORD environment variable for the password value")
	}

	if interpolateRegex.MatchString(opts.Target) {
		errBuffer = append(errBuffer, "* use OM_TARGET environment variable for the target value")
	}

	if interpolateRegex.MatchString(opts.Username) {
		errBuffer = append(errBuffer, "* use OM_USERNAME environment variable for the username value")
	}

	if len(errBuffer) > 0 {
		errBuffer = append([]string{"env file contains YAML placeholders. Pleases provide them via interpolation or environment variables."}, errBuffer...)
		errBuffer = append(errBuffer, "Or, to enable interpolation of env.yml with variables from env-vars,")
		errBuffer = append(errBuffer, "set the OM_VARS_ENV env var and put export the needed vars.")

		return fmt.Errorf(strings.Join(errBuffer, "\n"))
	}

	return nil
}
