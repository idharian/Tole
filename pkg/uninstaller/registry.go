package uninstaller

import (
	"sort"
	"strings"

	"golang.org/x/sys/windows/registry"
)

type InstalledApp struct {
	ID              string
	DisplayName     string
	DisplayVersion  string
	Publisher       string
	InstallLocation string
	UninstallString string
	EstimatedSizeKB uint64
}

// GetInstalledApps scans the Windows Registry for installed applications.
func GetInstalledApps() ([]*InstalledApp, error) {
	paths := []struct {
		root registry.Key
		path string
	}{
		{registry.LOCAL_MACHINE, `Software\Microsoft\Windows\CurrentVersion\Uninstall`},
		{registry.LOCAL_MACHINE, `Software\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall`},
		{registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Uninstall`},
	}

	appMap := make(map[string]*InstalledApp)

	for _, p := range paths {
		k, err := registry.OpenKey(p.root, p.path, registry.READ)
		if err != nil {
			continue
		}

		subkeys, err := k.ReadSubKeyNames(-1)
		k.Close()
		if err != nil {
			continue
		}

		for _, subkey := range subkeys {
			sub, err := registry.OpenKey(p.root, p.path+`\`+subkey, registry.READ)
			if err != nil {
				continue
			}

			name, _, _ := sub.GetStringValue("DisplayName")
			name = strings.TrimSpace(name)
			if name == "" {
				sub.Close()
				continue
			}

			// Ignore system updates or components if SystemComponent = 1
			sysComp, _, _ := sub.GetIntegerValue("SystemComponent")
			if sysComp == 1 {
				sub.Close()
				continue
			}

			uninstallStr, _, _ := sub.GetStringValue("UninstallString")
			if uninstallStr == "" {
				uninstallStr, _, _ = sub.GetStringValue("QuietUninstallString")
			}

			version, _, _ := sub.GetStringValue("DisplayVersion")
			publisher, _, _ := sub.GetStringValue("Publisher")
			installLoc, _, _ := sub.GetStringValue("InstallLocation")
			sizeKB, _, _ := sub.GetIntegerValue("EstimatedSize")

			sub.Close()

			// Distinct apps: same name different arch/version kept separate.
			key := strings.ToLower(name) + "|" + strings.ToLower(version) + "|" + strings.ToLower(publisher) + "|" + strings.ToLower(subkey)
			if _, exists := appMap[key]; !exists {
				appMap[key] = &InstalledApp{
					ID:              subkey,
					DisplayName:     name,
					DisplayVersion:  version,
					Publisher:       publisher,
					InstallLocation: installLoc,
					UninstallString: uninstallStr,
					EstimatedSizeKB: sizeKB,
				}
			}
		}
	}

	apps := make([]*InstalledApp, 0, len(appMap))
	for _, a := range appMap {
		apps = append(apps, a)
	}

	sort.Slice(apps, func(i, j int) bool {
		return strings.ToLower(apps[i].DisplayName) < strings.ToLower(apps[j].DisplayName)
	})

	return apps, nil
}
