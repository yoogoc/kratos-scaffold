package conf

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/yoogoc/kratos-scaffold/pkg/util"
)

const serviceMessage = `  message Service {
    string address = 1;
    google.protobuf.Duration timeout = 2;
  }
`

func EnsureServiceInConf(confDir, targetModel, namespace string, isSingle bool) error {
	protoPath := path.Join(confDir, "conf.proto")
	content, err := os.ReadFile(protoPath)
	if err != nil {
		return fmt.Errorf("read conf.proto: %w", err)
	}

	src := string(content)
	modified := false

	if !hasServiceMessage(src) {
		src = addServiceMessage(src)
		modified = true
	}

	fieldName := strcase.ToSnake(targetModel)
	if !hasServiceField(src, fieldName) {
		src = addServiceField(src, fieldName)
		modified = true
	}

	if !hasDurationImport(src) {
		src = addDurationImport(src)
		modified = true
	}

	if !modified {
		return nil
	}

	if err := os.WriteFile(protoPath, []byte(src), 0o644); err != nil {
		return fmt.Errorf("write conf.proto: %w", err)
	}

	fmt.Println("compiling conf.proto...")
	return runProtoc(confDir, namespace, isSingle)
}

func hasServiceMessage(src string) bool {
	return strings.Contains(src, "message Service")
}

func hasDurationImport(src string) bool {
	return strings.Contains(src, "google/protobuf/duration.proto")
}

func addDurationImport(src string) string {
	importLine := `import "google/protobuf/duration.proto";`
	if idx := strings.Index(src, "option go_package"); idx != -1 {
		lineEnd := strings.Index(src[idx:], "\n")
		if lineEnd != -1 {
			insertAt := idx + lineEnd + 1
			return src[:insertAt] + "\n" + importLine + "\n" + src[insertAt:]
		}
	}
	return importLine + "\n" + src
}

func addServiceMessage(src string) string {
	re := regexp.MustCompile(`(message\s+Data\s*\{)`)
	return re.ReplaceAllString(src, "${1}\n"+serviceMessage)
}

func hasServiceField(src, fieldName string) bool {
	re := regexp.MustCompile(`Service\s+` + regexp.QuoteMeta(fieldName) + `\s*=`)
	return re.MatchString(src)
}

func addServiceField(src, fieldName string) string {
	nextNum := findNextDataFieldNumber(src)
	field := fmt.Sprintf("  Service %s = %d;\n", fieldName, nextNum)

	dataBlockEnd := findDataBlockClosingBrace(src)
	if dataBlockEnd < 0 {
		return src
	}
	return src[:dataBlockEnd] + field + src[dataBlockEnd:]
}

func findNextDataFieldNumber(src string) int {
	dataStart := strings.Index(src, "message Data")
	if dataStart < 0 {
		return 1
	}
	block := src[dataStart:]
	braceEnd := findMatchingBrace(block)
	if braceEnd < 0 {
		return 1
	}
	dataBlock := block[:braceEnd]

	re := regexp.MustCompile(`=\s*(\d+)\s*;`)
	matches := re.FindAllStringSubmatch(dataBlock, -1)
	max := 0
	for _, m := range matches {
		if n, err := strconv.Atoi(m[1]); err == nil && n > max {
			max = n
		}
	}
	return max + 1
}

func findDataBlockClosingBrace(src string) int {
	dataStart := strings.Index(src, "message Data")
	if dataStart < 0 {
		return -1
	}
	block := src[dataStart:]
	pos := findMatchingBrace(block)
	if pos < 0 {
		return -1
	}
	return dataStart + pos
}

func findMatchingBrace(s string) int {
	start := strings.Index(s, "{")
	if start < 0 {
		return -1
	}
	depth := 0
	for i := start; i < len(s); i++ {
		switch s[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

func runProtoc(confDir, namespace string, isSingle bool) error {
	if isSingle {
		cmd := exec.Command("protoc",
			"--proto_path=./internal/conf",
			"--proto_path=./third_party",
			"--go_out=paths=source_relative:./internal/conf",
			"conf.proto",
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
	ns := namespace
	if ns == "" {
		return fmt.Errorf("namespace is required for mono-repo")
	}
	return util.Exec("make", "config-"+ns)
}
