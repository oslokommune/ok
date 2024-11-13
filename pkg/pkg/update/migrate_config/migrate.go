package migrate_config

import (
	"crypto/sha256"
	"fmt"
	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/oslokommune/ok/pkg/pkg/update/migrate_config/add_apex_domain"
	"io"
	"os"
)

func MigratePackageConfig(packagesToUpdate []common.Package) error {
	for _, pkg := range packagesToUpdate {
		for _, varFile := range pkg.VarFiles {
			fileHash, err := getFileHash(varFile)
			if err != nil {
				return fmt.Errorf("getting file hash: %w", err)
			}

			err = migrate(varFile, pkg)
			if err != nil {
				err = tryToGracefullyHandleError(varFile, fileHash, err)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func migrate(varFile string, pkg common.Package) error {
	// NOTE: Be careful with the order of the functions here. In general:
	// - Always append function calls to new updates at the end of this function.
	// - Do not change the order of the functions.
	//
	// This is to ensure that previously executed migrations/updates do not get messed up somehow, because of
	// dependencies between them. Of course, if you know what you are doing, go ahead.

	err := add_apex_domain.AddApexDomainSupport(varFile, pkg)
	if err != nil {
		return err
	}

	return nil
}

func getFileHash(filePath string) (string, error) {
	// https://pkg.go.dev/crypto/sha256#New
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err = io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func tryToGracefullyHandleError(varFile string, oldHash string, cause error) error {
	fileHash, err := getFileHash(varFile)
	if err != nil {
		return fmt.Errorf("getting file hash from file %s: %w", varFile, err)
	}

	if oldHash != fileHash {
		return err
	}

	fmt.Printf("WARNING: Auto migrating package config failed. "+
		"However, as the config file has not changed, we're ignoring this error. Config file: %s. Error:%s\n", varFile, cause)

	return nil
}
