package moby

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

func setupInspect(t *testing.T, label Image) types.ImageInspect {
	var inspect types.ImageInspect
	var config container.Config

	labelJSON, err := json.Marshal(label)
	if err != nil {
		t.Error(err)
	}
	config.Labels = map[string]string{"org.mobyproject.config": string(labelJSON)}

	inspect.Config = &config

	return inspect
}

func TestOverrides(t *testing.T) {
	var yamlCaps = []string{"CAP_SYS_ADMIN"}

	var yaml = Image{
		Name:         "test",
		Image:        "testimage",
		Capabilities: &yamlCaps,
	}

	var labelCaps = []string{"CAP_SYS_CHROOT"}

	var label = Image{
		Capabilities: &labelCaps,
		Cwd:          "/label/directory",
	}

	inspect := setupInspect(t, label)

	oci, err := ConfigInspectToOCI(yaml, inspect)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(oci.Process.Capabilities.Bounding, yamlCaps) {
		t.Error("Expected yaml capabilities to override but got", oci.Process.Capabilities.Bounding)
	}
	if oci.Process.Cwd != label.Cwd {
		t.Error("Expected label Cwd to be applied, got", oci.Process.Cwd)
	}
}

func TestInvalidCap(t *testing.T) {
	yaml := Image{
		Name:  "test",
		Image: "testimage",
	}

	labelCaps := []string{"NOT_A_CAP"}
	var label = Image{
		Capabilities: &labelCaps,
	}

	inspect := setupInspect(t, label)

	_, err := ConfigInspectToOCI(yaml, inspect)
	if err == nil {
		t.Error("expected error, got valid OCI config")
	}
}