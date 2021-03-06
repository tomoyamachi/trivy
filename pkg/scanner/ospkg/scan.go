package ospkg

import (
	"github.com/knqyf263/fanal/analyzer"
	fos "github.com/knqyf263/fanal/analyzer/os"
	_ "github.com/knqyf263/fanal/analyzer/os/alpine"
	_ "github.com/knqyf263/fanal/analyzer/os/debianbase"
	_ "github.com/knqyf263/fanal/analyzer/os/redhatbase"
	_ "github.com/knqyf263/fanal/analyzer/pkg/apk"
	_ "github.com/knqyf263/fanal/analyzer/pkg/dpkg"
	"github.com/knqyf263/fanal/extractor"
	"github.com/knqyf263/trivy/pkg/log"
	"github.com/knqyf263/trivy/pkg/scanner/ospkg/alpine"
	"github.com/knqyf263/trivy/pkg/scanner/ospkg/debian"
	"github.com/knqyf263/trivy/pkg/scanner/ospkg/redhat"
	"github.com/knqyf263/trivy/pkg/scanner/ospkg/ubuntu"
	"github.com/knqyf263/trivy/pkg/types"
	"golang.org/x/xerrors"
)

type Scanner interface {
	Detect(string, []analyzer.Package) ([]types.Vulnerability, error)
}

func Scan(files extractor.FileMap) (string, string, []types.Vulnerability, error) {
	os, err := analyzer.GetOS(files)
	if err != nil {
		return "", "", nil, xerrors.Errorf("failed to analyze OS: %w", err)
	}
	log.Logger.Debugf("OS family: %s, OS version: %s", os.Family, os.Name)

	var s Scanner
	switch os.Family {
	case fos.Alpine:
		s = alpine.NewScanner()
	case fos.Debian:
		s = debian.NewScanner()
	case fos.Ubuntu:
		s = ubuntu.NewScanner()
	case fos.RedHat, fos.CentOS:
		s = redhat.NewScanner()
	default:
		return "", "", nil, xerrors.New("unsupported os")
	}
	pkgs, err := analyzer.GetPackages(files)
	if err != nil {
		return "", "", nil, xerrors.Errorf("failed to analyze OS packages: %w", err)
	}
	log.Logger.Debugf("the number of packages: %d", len(pkgs))

	vulns, err := s.Detect(os.Name, pkgs)
	if err != nil {
		return "", "", nil, xerrors.Errorf("failed to detect vulnerabilities: %w", err)
	}

	return os.Family, os.Name, vulns, nil
}
