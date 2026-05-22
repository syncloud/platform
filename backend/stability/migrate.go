package stability

import (
	"errors"
	"os"

	"go.uber.org/zap"
)

func MigrateEventLog(oldPath, newPath string, logger *zap.Logger) {
	if _, err := os.Stat(newPath); err == nil {
		return
	} else if !errors.Is(err, os.ErrNotExist) {
		logger.Warn("stability: stat new event log failed", zap.Error(err))
		return
	}
	if _, err := os.Stat(oldPath); err != nil {
		return
	}
	if err := os.Rename(oldPath, newPath); err != nil {
		logger.Warn("stability: migrate event log failed", zap.String("from", oldPath), zap.String("to", newPath), zap.Error(err))
		return
	}
	logger.Info("stability: migrated event log", zap.String("from", oldPath), zap.String("to", newPath))
}
