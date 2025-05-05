package asyncapigendoc_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	asyncapigendoc "github.com/dnitsch/async-api-generator/cmd/async-api-gen-doc"
	"github.com/dnitsch/async-api-generator/internal/fshelper"
)

func Test_global_Analyis_runs_ok(t *testing.T) {

	t.Run("azblob source and output", func(t *testing.T) {
		t.Skip()

		cmd := asyncapigendoc.NewCmd(context.TODO())
		cmd.WithCommands()

		b := new(bytes.Buffer)

		cmd.Cmd.SetArgs([]string{"global-context", "-i",
			"azblob://stdevsandboxeuwdev/interim/current",
			"--output", "azblob://stdevsandboxeuwdev/processed"})

		cmd.Cmd.SetErr(b)
		cmd.Execute()
		out, err := io.ReadAll(b)
		if err != nil {
			t.Fatal(err)
		}
		if len(out) > 0 {
			t.Fatalf("expected error output buffer to be empty\ngot: %v\nwanted 0", string(out))
		}
	})

	t.Run("local source and local processed output", func(t *testing.T) {

		out, err := os.MkdirTemp("", "global")
		if err != nil {
			t.Fatal(err)
			return
		}
		defer os.RemoveAll(out)

		cmd := asyncapigendoc.NewCmd(context.TODO())
		cmd.WithCommands()

		baseDir := "test/interim-generated"

		b := new(bytes.Buffer)
		output := fmt.Sprintf("local://%s", out)
		cmd.Cmd.SetArgs([]string{"global-context", "-i",
			fmt.Sprintf("local://%s", fshelper.DebugDirHelper(t, baseDir, "cmd/async-api-gen-doc", "../../")),
			"--verbose",
			"--output", output})

		cmd.Cmd.SetErr(b)
		cmd.Execute()

		rb, err := io.ReadAll(b)
		if err != nil {
			t.Fatal(err)
		}

		if len(rb) > 0 {
			t.Fatalf("expected error output buffer to be empty\ngot: %v\nwanted 0", string(rb))
		}
		found := []string{}
		filepath.WalkDir(out, func(path string, d fs.DirEntry, err error) error {
			found = append(found, path)
			return nil
		})
		if len(found) < 2 {
			t.Fatal("no processed files were emitted")
		}
	})
}
