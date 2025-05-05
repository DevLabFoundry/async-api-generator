package asyncapigendoc_test

import (
	"bytes"
	"context"
	"io"
	"testing"

	asyncapigendoc "github.com/dnitsch/async-api-generator/cmd/async-api-gen-doc"
	"github.com/dnitsch/async-api-generator/internal/fshelper"
)

func Test_root_ok(t *testing.T) {
	baseDir := "test/foo.sample"
	b := new(bytes.Buffer)

	cmd := asyncapigendoc.NewCmd(context.TODO())
	cmd.WithCommands()
	fshelper.DebugDirHelper(t, baseDir, "cmd/async-api-gen-doc", "../../")

	cmd.Cmd.SetArgs([]string{"--version"})
	cmd.Cmd.SetErr(b)
	cmd.Execute()
	out, err := io.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) > 0 {
		t.Fatalf("expected error output buffer to be empty\ngot: %v\nwanted 0", string(out))
	}
}
