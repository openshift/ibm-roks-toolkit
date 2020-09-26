package release

import (
	"github.com/pkg/errors"

	"fmt"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"os"
	"strings"

	"github.com/openshift/oc/pkg/cli/admin/release"
)

// ReleaseInfo includes image references and versions for a given release
type ReleaseInfo struct {
	Images   map[string]string
	Versions map[string]string
}

func GetReleaseInfo(image string, originReleasePrefix string, pullSecretFile string) (*ReleaseInfo, error) {
	streams := genericclioptions.IOStreams{
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}
	options := release.NewInfoOptions(streams)
	options.SecurityOptions.RegistryConfig = pullSecretFile
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

	return &ReleaseInfo{
		Images:   images,
		Versions: versions,
	}, nil
}
