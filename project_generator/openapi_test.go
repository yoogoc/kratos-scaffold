package project_generator

import (
	"strings"
	"testing"
)

func TestOpenAPITemplateUsesV3Generator(t *testing.T) {
	if !strings.Contains(bufGenOpenapiYaml, "protoc-gen-openapi") {
		t.Fatal("OpenAPI template must use the OpenAPI v3 generator")
	}
	if strings.Contains(bufGenOpenapiYaml, "protoc-gen-openapiv2") {
		t.Fatal("OpenAPI template still contains the v2 generator")
	}
	if strings.Contains(bufGenOpenapiYaml, "allow_merge") {
		t.Fatal("OpenAPI template still contains v2-only options")
	}
}
