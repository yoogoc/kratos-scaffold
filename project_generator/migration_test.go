package project_generator

import (
	"strings"
	"testing"
)

func TestMigrationCommandOnlyAllowsUpAndDown(t *testing.T) {
	for _, want := range []string{
		`ValidArgs:           []string{"up", "down"}`,
		"cobra.OnlyValidArgs",
		"goose.RunContext(cmd.Context(), args[0], db, dir)",
	} {
		if !strings.Contains(cmdMigrationTmpl, want) {
			t.Fatalf("migration command template does not contain %q", want)
		}
	}
	for _, unwanted := range []string{`case "create"`, `case "fix"`} {
		if strings.Contains(cmdMigrationTmpl, unwanted) {
			t.Fatalf("migration command template still contains %q", unwanted)
		}
	}
}

func TestMigrationRecipeUsesGoToolGoose(t *testing.T) {
	for name, recipe := range map[string]string{
		"single": justfileSingleContent,
		"mono":   justfileMonoContent,
	} {
		if !strings.Contains(recipe, "go tool goose -dir") {
			t.Errorf("%s migration recipe does not use go tool goose", name)
		}
		if strings.Contains(recipe, "cmd migrate create") {
			t.Errorf("%s migration recipe still delegates create to cmd", name)
		}
	}
}
