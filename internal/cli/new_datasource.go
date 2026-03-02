package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chrispian/sigil/internal/config"
	"github.com/spf13/cobra"
)

// NewNewDataSourceCmd creates the "new datasource" subcommand.
func NewNewDataSourceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "datasource <alias>",
		Short: "Create a new datasource manifest",
		Long:  "Create a new datasource manifest from a starter template.",
		Args:  cobra.ExactArgs(1),
		RunE:  runNewDataSource,
	}
	cmd.Flags().String("capabilities", "search,filter,sort,paginate", "Comma-separated capabilities")
	cmd.Flags().Bool("force", false, "Overwrite existing file")
	return cmd
}

// dsManifest represents the datasource YAML output.
type dsManifest struct {
	Alias        string         `yaml:"alias"`
	Description  string         `yaml:"description"`
	Capabilities []string       `yaml:"capabilities"`
	Fields       []dsField      `yaml:"fields"`
	Relations    []interface{}  `yaml:"relations"`
	Endpoints    dsEndpoints    `yaml:"endpoints"`
	Defaults     dsDefaults     `yaml:"defaults"`
}

type dsField struct {
	Name    string `yaml:"name"`
	Type    string `yaml:"type"`
	Primary bool   `yaml:"primary,omitempty"`
	Hidden  bool   `yaml:"hidden,omitempty"`
}

type dsEndpoints struct {
	List   string `yaml:"list"`
	Create string `yaml:"create"`
	Read   string `yaml:"read"`
	Update string `yaml:"update"`
	Delete string `yaml:"delete"`
}

type dsDefaults struct {
	Sort     dsSort `yaml:"sort"`
	PageSize int    `yaml:"pageSize"`
}

type dsSort struct {
	Field     string `yaml:"field"`
	Direction string `yaml:"direction"`
}

func runNewDataSource(cmd *cobra.Command, args []string) error {
	alias := args[0]
	capsStr, _ := cmd.Flags().GetString("capabilities")
	force, _ := cmd.Flags().GetBool("force")

	caps := strings.Split(capsStr, ",")
	for i := range caps {
		caps[i] = strings.TrimSpace(caps[i])
	}

	outDir := filepath.Join(".sigil", "datasources")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("creating datasources directory: %w", err)
	}

	outPath := filepath.Join(outDir, alias+".yaml")
	if _, err := os.Stat(outPath); err == nil && !force {
		return fmt.Errorf("%s already exists (use --force to overwrite)", outPath)
	}

	// Pluralize and lowercase for endpoint paths
	resource := strings.ToLower(alias) + "s"

	manifest := dsManifest{
		Alias:        alias,
		Description:  "",
		Capabilities: caps,
		Fields: []dsField{
			{Name: "id", Type: "integer", Primary: true, Hidden: true},
		},
		Relations: []interface{}{},
		Endpoints: dsEndpoints{
			List:   "GET /api/" + resource,
			Create: "POST /api/" + resource,
			Read:   "GET /api/" + resource + "/{id}",
			Update: "PUT /api/" + resource + "/{id}",
			Delete: "DELETE /api/" + resource + "/{id}",
		},
		Defaults: dsDefaults{
			Sort:     dsSort{Field: "created_at", Direction: "desc"},
			PageSize: 25,
		},
	}

	if err := config.WriteFile(outPath, &manifest); err != nil {
		return fmt.Errorf("writing datasource manifest: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Created %s\n", outPath)
	return nil
}
