package asyncapigendoc_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	asyncapigendoc "github.com/dnitsch/async-api-generator/cmd/async-api-gen-doc"
	"github.com/dnitsch/async-api-generator/internal/fshelper"
)

func Test_single_repo_Analyis_runs_ok(t *testing.T) {
	baseDir := "test/foo.sample"

	searchParentDir := fmt.Sprintf("local://%s", fshelper.DebugDirHelper(t, baseDir, "cmd/async-api-gen-doc", "../../"))

	t.Run("local output and is-service set", func(t *testing.T) {
		cmd := asyncapigendoc.NewCmd(context.TODO())

		b := new(bytes.Buffer)
		cmd.Cmd.SetArgs([]string{"single-context", "--verbose", "--input", searchParentDir, "--is-service", "--bounded-ctx", "s2s", "--business-domain", "domain"})
		cmd.Cmd.SetErr(b)
		cmd.WithCommands()
		cmd.Execute()
		out, err := io.ReadAll(b)
		if err != nil {
			t.Fatal(err)
		}
		if len(out) > 0 {
			t.Fatalf("expected error output buffer to be empty\ngot: %v\nwanted 0", string(out))
		}
	})

	t.Run("dry-run service set", func(t *testing.T) {
		cmd := asyncapigendoc.NewCmd(context.TODO())
		cmd.WithCommands()

		b := new(bytes.Buffer)
		cmd.Cmd.SetArgs([]string{"single-context", "--verbose", "--is-service", "-i", searchParentDir, "--dry-run"})
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
	t.Run("local input and azblob output", func(t *testing.T) {
		// uncomment this for local testing only
		t.Skip()
		cmd := asyncapigendoc.NewCmd(context.TODO())
		cmd.WithCommands()

		b := new(bytes.Buffer)
		cmd.Cmd.SetArgs([]string{"single-context", "--verbose", "--is-service", "-i", searchParentDir, "--bounded-ctx", "s2s", "--business-domain", "domain", "--output", "azblob://stdevsandboxeuwdev/interim"})
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
}

func Test_Failures_on_incorrect_input(t *testing.T) {
	baseDir := "./test/negative"

	ttests := map[string]struct {
		flags         []string
		expErrMessage string
	}{
		"with wrong input protocol": {
			[]string{"single-context", "-i", fmt.Sprintf("file://%s", fshelper.DebugDirHelper(t, baseDir, "cmd/async-api-gen-doc", "../../")), "--verbose"},
			"unsupported protocol error",
		},
		"with wrong markers in found files": {
			[]string{"single-context", "-i", fmt.Sprintf("local://%s", fshelper.DebugDirHelper(t, baseDir, "cmd/async-api-gen-doc", "../../")), "--verbose", "--dry-run"},
			"GenDocBlox failed to parse all documents",
		},
	}
	for name, tt := range ttests {
		t.Run(name, func(t *testing.T) {
			b := new(bytes.Buffer)

			cmd := asyncapigendoc.NewCmd(context.TODO())
			cmd.WithCommands()
			cmd.Cmd.SetArgs(tt.flags)
			cmd.Cmd.SetErr(b)
			err := cmd.Execute()
			if err == nil {
				t.Fatal("should have failed with error")
			}
			if !strings.Contains(err.Error(), tt.expErrMessage) {
				t.Errorf("wrong error message got: %s\nwanted: %s", err.Error(), tt.expErrMessage)
			}
		})
	}
}
