package release

import (
	"regexp"
	"runtime"

	"github.com/pkg/errors"

	"fmt"
	"os"
	"strings"

	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/openshift/oc/pkg/cli/admin/release"
)

// Info includes image references and versions for a given release
type Info struct {
	Images   map[string]string
	Versions map[string]string
}

func GetReleaseInfo(image string, originReleasePrefix string, pullSecretFile string) (*Info, error) {
	streams := genericclioptions.IOStreams{
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}
	options := release.NewInfoOptions(streams)
	options.SecurityOptions.RegistryConfig = pullSecretFile
	options.FilterOptions.DefaultOSFilter = true
	pattern := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
	options.FilterOptions.FilterByOS = pattern
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	options.FilterOptions.OSFilter = re
	info, err := options.LoadReleaseInfo(image, false)
	if err != nil {
		return nil, err
	}
	if info.References == nil {
		return nil, errors.New("release image does not contain image references")
	}

	var newImagePrefix string
	if !strings.Contains(image, originReleasePrefix) {
		newImagePrefix = strings.Replace(image, ":", "-", -1)
	}
	images := make(map[string]string)
	for _, tag := range info.References.Spec.Tags {
		name := tag.From.Name
		if len(newImagePrefix) > 0 {
			name = fmt.Sprintf("%s@%s", newImagePrefix, strings.Split(tag.From.Name, "@")[1])
		}
		images[tag.Name] = name
	}

	versions := make(map[string]string)
	if info.Metadata != nil {
		versions["release"] = info.Metadata.Version
	}
	for component, version := range info.ComponentVersions {
		versions[component] = version.Version
	}

	return &Info{
		Images:   images,
		Versions: versions,
	}, nil
}
