package asyncapigendoc

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dnitsch/async-api-generator/internal/fshelper"
	"github.com/dnitsch/async-api-generator/internal/generate"
	"github.com/dnitsch/async-api-generator/internal/storage"
	log "github.com/dnitsch/simplelog"
	"github.com/spf13/cobra"
)

type genDocContextFlags struct {
	businessDomain   string
	boundedCtxDomain string
	repoUrl          string
	repoLang         string
	isService        bool
	serviceId        string
}

func singleContextCmd(rootCmd *AsyncApiGenDocCmd) {
	f := &genDocContextFlags{}
	singleCtxCmd := &cobra.Command{
		Use:     "single-context",
		Aliases: []string{"sc", "single"},
		Short:   `Runs the gendoc against a single repo source`,
		Long:    `Runs the gendoc against a single repo source and emits the output to specified storage.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if rootCmd.rootFlags.verbose {
				rootCmd.logger = log.New(os.Stdout, log.DebugLvl)
			}

			conf, cleanUp, err := rootCmd.config(rootCmd.inputLocationStorageConfig, f)
			if err != nil {
				return err
			}

			defer cleanUp()

			files, err := fshelper.ListFiles(rootCmd.inputLocationStorageConfig.Destination)
			if err != nil {
				return err
			}

			gendoc := generate.New(conf, rootCmd.logger)

			gendoc.LoadInputsFromFiles(files)

			if err := gendoc.GenDocBlox(); err != nil {
				return err
			}

			if rootCmd.rootFlags.dryRun {
				rootCmd.logger.Debugf("--dry-run only not storing locally or remotely")
				return nil
			}

			// set out name for single repo analysis
			outName := fmt.Sprintf("current/%s.json", conf.SearchDirName)
			// select storage adapter
			sc, err := storage.ClientFactory(rootCmd.outputStorageConfig.Typ, rootCmd.outputStorageConfig.Destination)
			if err != nil {
				return err
			}

			storageUpldReq := &storage.StorageUploadRequest{
				ContainerName: rootCmd.outputStorageConfig.TopLevelFolder,
				BlobKey:       outName,
				Destination:   filepath.Join(rootCmd.outputStorageConfig.Destination, rootCmd.outputStorageConfig.TopLevelFolder, outName),
			}

			ctx, cancel := context.WithCancel(cmd.Context())
			defer cancel()

			return gendoc.CommitInterimState(ctx, sc, storageUpldReq)
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return rootCmd.setStorageLocation(rootCmd.rootFlags.inputLocation, rootCmd.rootFlags.outputLocation)
		},
	}

	singleCtxCmd.PersistentFlags().StringVarP(&f.businessDomain, "business-domain", "b", "", `businessDomain e.g. Warehouse Systems`)
	singleCtxCmd.PersistentFlags().StringVarP(&f.boundedCtxDomain, "bounded-ctx", "c", "", `boundedCtxDomain`)
	singleCtxCmd.PersistentFlags().StringVarP(&f.repoUrl, "repo", "r", "", `repoUrl`)
	singleCtxCmd.PersistentFlags().StringVarP(&f.repoLang, "lang", "", "C#", `Main Language used in repo`)
	singleCtxCmd.PersistentFlags().StringVarP(&f.serviceId, "service-id", "", "", `serviceId`)
	singleCtxCmd.PersistentFlags().BoolVarP(&f.isService, "is-service", "s", false, `whether the repo is a service repo`)
	rootCmd.Cmd.AddCommand(singleCtxCmd)
}
