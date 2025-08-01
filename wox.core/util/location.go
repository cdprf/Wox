package util

import (
	"fmt"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/mitchellh/go-homedir"
)

var locationInstance *Location
var locationOnce sync.Once

type Location struct {
	// wox data directory is the directory that contains all wox data, including logs, hosts, etc.
	woxDataDirectory string

	// user data directory is the directory that contains all user data, including plugins, settings, etc.
	// user may change the user data directory to another location, E.g. icloud, google drive, etc.
	userDataDirectory string

	userDataDirectoryShortcutPath string // A file named .wox.location that contains the user data directory path
}

func GetLocation() *Location {
	locationOnce.Do(func() {
		locationInstance = &Location{}
	})
	return locationInstance
}

func (l *Location) Init() error {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home dir: %w", err)
	}

	// check if wox data directory exists, if not, create it
	l.woxDataDirectory = path.Join(dirname, ".wox")
	if directoryErr := l.EnsureDirectoryExist(l.woxDataDirectory); directoryErr != nil {
		return directoryErr
	}

	l.userDataDirectoryShortcutPath = path.Join(l.woxDataDirectory, ".userdata.location")
	if _, statErr := os.Stat(l.userDataDirectoryShortcutPath); os.IsNotExist(statErr) {
		// shortcut file does not exist, create and write default data directory path to it
		file, createErr := os.Create(l.userDataDirectoryShortcutPath)
		if createErr != nil {
			return fmt.Errorf("failed to create shortcut file: %w", createErr)
		}
		defer file.Close()

		// write data directory path to file
		_, writeErr := file.WriteString(path.Join(l.woxDataDirectory, "wox-user"))
		if writeErr != nil {
			return fmt.Errorf("failed to write user data directory path to shortcut file: %w", writeErr)
		}
	}

	// read data directory path from shortcut file
	file, openErr := os.Open(l.userDataDirectoryShortcutPath)
	if openErr != nil {
		return fmt.Errorf("failed to open shortcut file: %w", openErr)
	}
	defer file.Close()

	// read data directory path from file
	readFile, readFileErr := os.ReadFile(l.userDataDirectoryShortcutPath)
	if readFileErr != nil {
		return fmt.Errorf("failed to read shortcut file: %w", readFileErr)
	}
	userDataDirectory, _ := homedir.Expand(string(readFile))
	userDataDirectory = strings.ReplaceAll(userDataDirectory, "\n", "")
	l.userDataDirectory = userDataDirectory

	if directoryErr := l.EnsureDirectoryExist(l.userDataDirectory); directoryErr != nil {
		return directoryErr
	}
	if directoryErr := l.EnsureDirectoryExist(l.GetLogDirectory()); directoryErr != nil {
		return directoryErr
	}
	if directoryErr := l.EnsureDirectoryExist(l.GetLogHostsDirectory()); directoryErr != nil {
		return directoryErr
	}
	if directoryErr := l.EnsureDirectoryExist(l.GetLogPluginDirectory()); directoryErr != nil {
		return directoryErr
	}
	if directoryErr := l.EnsureDirectoryExist(l.GetPluginDirectory()); directoryErr != nil {
		return directoryErr
	}
	if directoryErr := l.EnsureDirectoryExist(l.GetThemeDirectory()); directoryErr != nil {
		return directoryErr
	}
	if directoryErr := l.EnsureDirectoryExist(l.GetHostDirectory()); directoryErr != nil {
		return directoryErr
	}
	if directoryErr := l.EnsureDirectoryExist(l.GetUpdatesDirectory()); directoryErr != nil {
		return directoryErr
	}
	if directoryErr := l.EnsureDirectoryExist(l.GetPluginSettingDirectory()); directoryErr != nil {
		return directoryErr
	}
	if directoryErr := l.EnsureDirectoryExist(l.GetUIDirectory()); directoryErr != nil {
		return directoryErr
	}
	if directoryErr := l.EnsureDirectoryExist(l.GetOthersDirectory()); directoryErr != nil {
		return directoryErr
	}
	if directoryErr := l.EnsureDirectoryExist(l.GetScriptPluginTemplatesDirectory()); directoryErr != nil {
		return directoryErr
	}
	if directoryErr := l.EnsureDirectoryExist(l.GetUserScriptPluginsDirectory()); directoryErr != nil {
		return directoryErr
	}
	if directoryErr := l.EnsureDirectoryExist(l.GetCacheDirectory()); directoryErr != nil {
		return directoryErr
	}
	if directoryErr := l.EnsureDirectoryExist(l.GetImageCacheDirectory()); directoryErr != nil {
		return directoryErr
	}
	if directoryErr := l.EnsureDirectoryExist(l.GetBackupDirectory()); directoryErr != nil {
		return directoryErr
	}

	return nil
}

func (l *Location) EnsureDirectoryExist(directory string) error {
	if _, statErr := os.Stat(directory); os.IsNotExist(statErr) {
		mkdirErr := os.MkdirAll(directory, os.ModePerm)
		if mkdirErr != nil {
			return fmt.Errorf("failed to create directory [%s]: %w", directory, mkdirErr)
		}
	}

	return nil
}

func (l *Location) GetLogDirectory() string {
	return path.Join(l.woxDataDirectory, "log")
}

func (l *Location) GetWoxDataDirectory() string {
	return l.woxDataDirectory
}

func (l *Location) GetLogPluginDirectory() string {
	return path.Join(l.GetLogDirectory(), "plugins")
}

func (l *Location) GetLogHostsDirectory() string {
	return path.Join(l.GetLogDirectory(), "hosts")
}

func (l *Location) GetPluginDirectory() string {
	return path.Join(l.userDataDirectory, "plugins")
}

func (l *Location) GetUserScriptPluginsDirectory() string {
	return path.Join(l.GetPluginDirectory(), "scripts")
}

func (l *Location) GetThemeDirectory() string {
	return path.Join(l.userDataDirectory, "themes")
}

func (l *Location) GetPluginSettingDirectory() string {
	return path.Join(l.userDataDirectory, "settings")
}

func (l *Location) GetUserDataDirectory() string {
	return l.userDataDirectory
}

func (l *Location) GetWoxSettingPath() string {
	return path.Join(l.GetPluginSettingDirectory(), "wox.json")
}

func (l *Location) GetWoxAppDataPath() string {
	return path.Join(l.GetPluginSettingDirectory(), "wox.data.json")
}

func (l *Location) GetHostDirectory() string {
	return path.Join(l.woxDataDirectory, "hosts")
}

func (l *Location) GetUpdatesDirectory() string {
	return path.Join(l.woxDataDirectory, "updates")
}

func (l *Location) GetUIDirectory() string {
	return path.Join(l.woxDataDirectory, "ui")
}

func (l *Location) GetOthersDirectory() string {
	return path.Join(l.woxDataDirectory, "others")
}

func (l *Location) GetScriptPluginTemplatesDirectory() string {
	return path.Join(l.woxDataDirectory, "script_plugin_templates")
}

func (l *Location) GetCacheDirectory() string {
	return path.Join(l.woxDataDirectory, "cache")
}

func (l *Location) GetImageCacheDirectory() string {
	return path.Join(l.GetCacheDirectory(), "images")
}

func (l *Location) GetBackupDirectory() string {
	return path.Join(l.woxDataDirectory, "backup")
}

func (l *Location) GetUIAppPath() string {
	if IsWindows() {
		return path.Join(l.GetUIDirectory(), "flutter", "wox", "wox-ui.exe")
	}
	if IsLinux() {
		return path.Join(l.GetUIDirectory(), "flutter", "wox", "wox-ui")
	}
	if IsMacOS() {
		return path.Join(l.GetUIDirectory(), "flutter", "wox-ui.app", "Contents", "MacOS", "wox-ui")
	}
	return ""
}

func (l *Location) GetAppLockPath() string {
	return path.Join(l.GetWoxDataDirectory(), "wox.lock")
}

func (l *Location) UpdateUserDataDirectory(newDirectory string) {
	l.userDataDirectory = newDirectory
}

// Get the path to the shortcut file that stores the user data directory path
func (l *Location) GetUserDataDirectoryShortcutPath() string {
	return l.userDataDirectoryShortcutPath
}
