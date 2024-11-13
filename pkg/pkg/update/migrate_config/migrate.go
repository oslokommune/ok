package migrate_config

import (
	"crypto/sha256"
	"fmt"
	"github.com/oslokommune/ok/pkg/pkg/common"
	"github.com/oslokommune/ok/pkg/pkg/update/migrate_config/add_apex_domain"
	"github.com/oslokommune/ok/pkg/pkg/update/migrate_config/metadata"
	"io"
	"log/slog"
	"os"
)

func MigratePackageConfig(packagesToUpdate []common.Package) error {
	for _, pkg := range packagesToUpdate {
		for _, varFile := range pkg.VarFiles {
			fileHash, err := getFileHash(varFile)
			if err != nil {
				return fmt.Errorf("getting file hash: %w", err)
			}

			err = updateVarFile(varFile)
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

func updateVarFile(varFile string) error {
	slog.Debug("updating var file", slog.String("varFile", varFile))

	firstLine, err := readFirstLine(varFile)
	if err != nil {
		return fmt.Errorf("reading first line from %s: %w", varFile, err)
	}

	varFileMetadata, err := metadata.ParseMetadata(firstLine)
	if err != nil {
		return fmt.Errorf("getting metadata from var file %s: %w", varFile, err)
	}

	err = update(varFile, varFileMetadata)
	if err != nil {
		return fmt.Errorf("updating varFile %s: %w", varFile, err)
	}

	return nil
}

func update(varFile string, metadata metadata.VarFileMetadata) error {
	err := add_apex_domain.AddApexDomainSupport(varFile, metadata)
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
