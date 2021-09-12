package gen

import (
	"strings"

	"github.com/drykit-go/strcase"
)

const (
	providerStructSuffix    = "CheckerProvider"
	providerInterfaceSuffix = "CheckerProvider"
)

func structToInterfaceName(structName string) string {
	baseName := trimProviderStructSuffix(structName)
	interfaceName := appendProviderInterfaceSuffix(baseName)
	return exportName(interfaceName)
}

func structToVarName(structName string) string {
	baseName := trimProviderStructSuffix(structName)
	return exportName(baseName)
}

func trimProviderStructSuffix(structName string) string {
	return strings.TrimSuffix(structName, providerStructSuffix)
}

func appendProviderInterfaceSuffix(baseName string) string {
	return baseName + providerInterfaceSuffix
}

func exportName(unexported string) string {
	return strcase.Pascal(unexported)
}
