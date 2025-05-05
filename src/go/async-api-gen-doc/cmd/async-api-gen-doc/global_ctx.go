package asyncapigendoc

import (
	"context"
	"os"

	"github.com/dnitsch/async-api-generator/internal/fshelper"
	"github.com/dnitsch/async-api-generator/internal/generate"
	"github.com/dnitsch/async-api-generator/internal/storage"
	log "github.com/dnitsch/simplelog"
	"github.com/spf13/cobra"
)

func globalCtxCmd(rootCmd *AsyncApiGenDocCmd) {

	cmd := &cobra.Command{
		Use:     "global-context",
		Aliases: []string{"gc", "global"},
		Short:   `Runs the gendoc against a directory containing processed GenDocBlox.`,
		Long: `Runs the gendoc against a directory containing processed GenDocBlox. Builds a hierarchical tree with the generated interim states across multiple contexts. 
Source must be specified [see output] option for examples and structure`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if rootCmd.rootFlags.verbose {
				rootCmd.logger = log.New(os.Stdout, log.DebugLvl)
			}

			conf, cleanUp, err := rootCmd.config(rootCmd.inputLocationStorageConfig, &genDocContextFlags{})
			if err != nil {
				return err
			}
			defer cleanUp()
			rootCmd.logger.Debugf("interim output: %s", conf.InterimOutputDir)
			rootCmd.logger.Debugf("download output: %s", conf.DownloadDir)

			ctx, cancel := context.WithCancel(cmd.Context())
			defer cancel()

			if err := fetchPrep(ctx, conf, rootCmd.inputLocationStorageConfig); err != nil {
				return err
			}

			files, err := fshelper.ListFiles(conf.DownloadDir)
			if err != nil {
				return err
			}

			g := generate.New(conf, rootCmd.logger)

			g.LoadInputsFromFiles(files)

			if err := g.ConvertProcessed(); err != nil {
				return err
			}

			if err := g.BuildContextTree(); err != nil {
				return err
			}

			if err := g.AsyncAPIFromProcessedTree(); err != nil {
				return err
			}
			return uploadPrep(ctx, g, rootCmd.outputStorageConfig)
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return rootCmd.setStorageLocation(rootCmd.rootFlags.inputLocation, rootCmd.rootFlags.outputLocation)
		},
	}
	rootCmd.Cmd.AddCommand(cmd)
}

// fetchPrep
func fetchPrep(ctx context.Context, conf *generate.Config, storageConf *storage.Conf) error {
	// storage adapter for source
	sc, err := storage.ClientFactory(storageConf.Typ, storageConf.Destination)
	if err != nil {
		return err
	}

	fetchReq := &storage.StorageFetchRequest{
		Destination:   storageConf.Destination,
		ContainerName: storageConf.TopLevelFolder,
		EmitPath:      conf.DownloadDir,
	}

	if err := sc.Fetch(ctx, fetchReq); err != nil {
		return err
	}
	return nil
}

// uploadPrep
func uploadPrep(ctx context.Context, g *generate.Generate, conf *storage.Conf) error {
	// storage adapter for output
	storageClient, err := storage.ClientFactory(conf.Typ, conf.Destination)
	if err != nil {
		return err
	}

	uploadReq := &storage.StorageUploadRequest{
		ContainerName: conf.TopLevelFolder,
		Destination:   conf.Destination,
		BlobKey:       ""} //blobKey i.e. las portion of the path are handled by the committer
	return g.CommitProcessedState(ctx, storageClient, uploadReq)
}
