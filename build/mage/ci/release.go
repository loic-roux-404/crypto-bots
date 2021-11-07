package ci

import (
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Release namespace
type Release mg.Namespace

const (
	semverExe = "semantic-release"
)

var semverFlags = []string{
	"allow-initial-development-versions",
	"download-plugins",
}

// Semantic run release generation
func (Release) Semantic(prerelease bool, noCi bool) error {
	if prerelease {
		semverFlags = append(semverFlags, "prerelease")
	}

	if noCi {
		semverFlags = append(semverFlags, "no-ci")
	}

	finalFlags := strings.Split(strings.Join(semverFlags, " --"), " ")

	return sh.Run("semantic-release", finalFlags...)
}
