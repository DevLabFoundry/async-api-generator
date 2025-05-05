package asyncapigendoc

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dnitsch/async-api-generator/internal/generate"
	"github.com/dnitsch/async-api-generator/internal/parser"
	"github.com/dnitsch/async-api-generator/internal/storage"
	log "github.com/dnitsch/simplelog"
	"github.com/spf13/cobra"
)

var (
	Version  string = "0.0.1"
	Revision string = "1111aaaa"
)

type flags struct {
	verbose        bool
	dryRun         bool
	outputLocation string
	inputLocation  string
}

type AsyncApiGenDocCmd struct {
	ctx                        context.Context
	Cmd                        *cobra.Command
	logger                     log.Loggeriface
	rootFlags                  *flags
	outputStorageConfig        *storage.Conf
	inputLocationStorageConfig *storage.Conf
}

func NewCmd(ctx context.Context) *AsyncApiGenDocCmd {
	f := &flags{}
	aagd := &AsyncApiGenDocCmd{
		ctx:       ctx,
		logger:    log.New(os.Stderr, log.ErrorLvl),
		rootFlags: f,
	}
	aagd.Cmd = &cobra.Command{
		Use:     "gendoc",
		Aliases: []string{"aadg", "generator"},
		Short:   "Generator for AsyncAPI documents",
		Long: `Generator for AsyncAPI documents, functions by performing lexical analysis on source files in a given base directory. 
	These can then be further fed into other generator tools, e.g. client/server generators`,
		Example:      "",
		SilenceUsage: true,
		Version:      fmt.Sprintf("%s-%s", Version, Revision),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return aagd.setStorageLocation(f.inputLocation, f.outputLocation)
		},
	}
	aagd.Cmd.PersistentFlags().StringVarP(&f.outputLocation, "output", "o", "local://$HOME/.gendoc", `Output type and destination, currently only supports [local://, azblob://]. if dry-run is set then this is ignored`)
	aagd.Cmd.PersistentFlags().StringVarP(&f.inputLocation, "input", "i", "local://.", `Path to start the search in, Must include the protocol - see output for options`)
	aagd.Cmd.PersistentFlags().BoolVarP(&f.verbose, "verbose", "v", false, "Verbose output")
	aagd.Cmd.PersistentFlags().BoolVarP(&f.dryRun, "dry-run", "", false, "Dry run only runs in validate mode and does not emit anything")
	return aagd
}

func (c *AsyncApiGenDocCmd) WithCommands() {
	for _, fn := range []func(*AsyncApiGenDocCmd){singleContextCmd, globalCtxCmd} {
		fn(c)
	}
}

func (c *AsyncApiGenDocCmd) Execute() error {
	return c.Cmd.ExecuteContext(c.ctx)
}

// config bootstraps pflags into useable config
func (c *AsyncApiGenDocCmd) config(outConf *storage.Conf, sf *genDocContextFlags) (*generate.Config, func(), error) {
	dirName := filepath.Base(outConf.Destination)

	conf := &generate.Config{
		ParserConfig:  parser.Config{ServiceRepoUrl: sf.repoUrl, BusinessDomain: sf.businessDomain, BoundedDomain: sf.boundedCtxDomain, ServiceLanguage: sf.repoLang},
		SearchDirName: dirName,
		Output:        outConf,
	}
	if sf.isService {
		// use the current search dir name as the serviceId
		// this allows certain objects to __not__ have parentId or id specified
		conf.ParserConfig.ServiceId = dirName
	}

	if !c.rootFlags.dryRun {
		// create interim local dirs for interim state or interim download storage
		interim, err := os.MkdirTemp("", ".gendoc-interim-*")
		if err != nil {
			return nil, nil, err
		}
		download, err := os.MkdirTemp("", ".gendoc-download-*")
		if err != nil {
			return nil, nil, err
		}
		conf.InterimOutputDir = interim
		conf.DownloadDir = download
	}

	return conf, func() {
		_ = os.RemoveAll(conf.InterimOutputDir)
		_ = os.RemoveAll(conf.DownloadDir)
	}, nil
}

// setStorageLocation sets the input/output locations
func (c *AsyncApiGenDocCmd) setStorageLocation(input, output string) error {
	inStoreConf, err := storage.ParseStorageOutputConfig(input)
	if err != nil {
		return err
	}
	outStoreConf, err := storage.ParseStorageOutputConfig(output)
	if err != nil {
		return err
	}
	c.inputLocationStorageConfig = inStoreConf
	c.outputStorageConfig = outStoreConf
	return nil
}
