package cloud

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// BackupManager handles backup operations
type BackupManager struct {
	backupDir string
}

// NewBackupManager creates a new backup manager
func NewBackupManager(backupDir string) *BackupManager {
	return &BackupManager{backupDir: backupDir}
}

// Initialize creates the backup directory
func (bm *BackupManager) Initialize() error {
	return os.MkdirAll(bm.backupDir, 0755)
}

// CreateBackup creates a backup of an account's data
func (bm *BackupManager) CreateBackup(accountID, containerID string) (string, error) {
	// Create backup filename
	timestamp := time.Now().Format("20060102_150405")
	backupFile := filepath.Join(bm.backupDir, fmt.Sprintf("%s_%s.tar.gz", accountID, timestamp))

	// Create backup file
	file, err := os.Create(backupFile)
	if err != nil {
		return "", fmt.Errorf("create backup file: %w", err)
	}
	defer file.Close()

	// Create gzip writer
	gw := gzip.NewWriter(file)
	defer gw.Close()

	// Create tar writer
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// Backup data directory
	dataDir := fmt.Sprintf("/var/lib/hypernexus/%s", containerID)
	if err := bm.addDirectoryToTar(tw, dataDir, ""); err != nil {
		return "", fmt.Errorf("backup data: %w", err)
	}

	// Backup configuration
	configDir := fmt.Sprintf("/var/lib/hypernexus/%s/config", containerID)
	if err := bm.addDirectoryToTar(tw, configDir, "config"); err != nil {
		log.Printf("[Backup] Warning: config directory not found for %s", containerID)
	}

	log.Printf("[Backup] Created backup for account %s: %s", accountID, backupFile)
	return backupFile, nil
}

// RestoreBackup restores a backup to an account's container
func (bm *BackupManager) RestoreBackup(backupFile, containerID string) error {
	// Open backup file
	file, err := os.Open(backupFile)
	if err != nil {
		return fmt.Errorf("open backup file: %w", err)
	}
	defer file.Close()

	// Create gzip reader
	gr, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("create gzip reader: %w", err)
	}
	defer gr.Close()

	// Create tar reader
	tr := tar.NewReader(gr)

	// Extract files
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read tar: %w", err)
		}

		// Create target path
		targetPath := filepath.Join("/var/lib/hypernexus", containerID, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return fmt.Errorf("create directory: %w", err)
			}
		case tar.TypeReg:
			// Create directory if needed
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return fmt.Errorf("create directory: %w", err)
			}

			// Create file
			file, err := os.Create(targetPath)
			if err != nil {
				return fmt.Errorf("create file: %w", err)
			}

			// Copy content
			if _, err := io.Copy(file, tr); err != nil {
				file.Close()
				return fmt.Errorf("copy file: %w", err)
			}
			file.Close()
		}
	}

	log.Printf("[Backup] Restored backup %s to container %s", backupFile, containerID)
	return nil
}

// ListBackups lists all backups for an account
func (bm *BackupManager) ListBackups(accountID string) ([]string, error) {
	pattern := filepath.Join(bm.backupDir, fmt.Sprintf("%s_*.tar.gz", accountID))
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("glob backups: %w", err)
	}
	return matches, nil
}

// DeleteBackup deletes a backup file
func (bm *BackupManager) DeleteBackup(backupFile string) error {
	return os.Remove(backupFile)
}

// addDirectoryToTar adds a directory to a tar archive
func (bm *BackupManager) addDirectoryToTar(tw *tar.Writer, dirPath, prefix string) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create header
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return fmt.Errorf("create header: %w", err)
		}

		// Update name
		relPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			return fmt.Errorf("relative path: %w", err)
		}
		if prefix != "" {
			header.Name = filepath.Join(prefix, relPath)
		} else {
			header.Name = relPath
		}

		// Write header
		if err := tw.WriteHeader(header); err != nil {
			return fmt.Errorf("write header: %w", err)
		}

		// If it's a file, write content
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("open file: %w", err)
			}
			defer file.Close()

			if _, err := io.Copy(tw, file); err != nil {
				return fmt.Errorf("copy file: %w", err)
			}
		}

		return nil
	})
}

// ScheduleBackup schedules automatic backups
func (bm *BackupManager) ScheduleBackup(accountID, containerID string, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			backupFile, err := bm.CreateBackup(accountID, containerID)
			if err != nil {
				log.Printf("[Backup] Error creating backup for %s: %v", accountID, err)
				continue
			}
			log.Printf("[Backup] Scheduled backup created: %s", backupFile)
		}
	}()
}
