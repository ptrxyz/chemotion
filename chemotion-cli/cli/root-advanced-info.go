package cli

import (
	"fmt"
	"math"
	"runtime"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
)

// helper function that is also used by infoInstanceRootCmd
func systemInfo() (info string) {
	// CPU
	info += fmt.Sprintln("- CPU Cores:", runtime.NumCPU())
	// Disk space
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		var disk syscall.Statfs_t
		if err := syscall.Statfs(workDir.String(), &disk); err == nil {
			info += fmt.Sprintf("- Disk space:\n  - %7.1fGi (total) %7.1fGi (free)\n", float64(disk.Blocks*uint64(disk.Bsize))/math.Pow(2, 30), float64(disk.Bavail*uint64(disk.Bsize))/math.Pow(2, 30))
		} else {
			zboth.Warn().Err(err).Msgf("Failed to retrieve information about disk space.")
		}
	} else {
		zboth.Warn().Err(fmt.Errorf("running on %s", runtime.GOOS)).Msgf("Cannot retrieve disk space information for this operating system.")
		// TODO-maybe-v2 write implementation for windows
	}
	// Memory
	if runtime.GOOS == "linux" {
		var mem syscall.Sysinfo_t
		if err := syscall.Sysinfo(&mem); err == nil {
			info += fmt.Sprintf("- Memory:\n  - %7.1fGi (total) %7.1fGi (free)\n", float64(mem.Totalram)/math.Pow(2, 30), float64(mem.Freeram)/math.Pow(2, 30))
		} else {
			zboth.Warn().Err(err).Msgf("Failed to retrieve information about memory.")
		}
	} else {
		zboth.Warn().Err(fmt.Errorf("running on %s", runtime.GOOS)).Msgf("Cannot retrieve memory information for this operating system.")
	}
	info += fmt.Sprintln("Used software versions:")
	printVersionOf := []string{"docker", "ruby", "passenger", "node", "npm"}
	for _, software := range printVersionOf {
		info += fmt.Sprintf("- %s: %s\n", strings.ToTitle(software), findVersion(software))
	}
	return
}

// Show host machine information to the user
// See also, chemotion instance info
var infoAdvancedRootCmd = &cobra.Command{
	Use:   "info",
	Args:  cobra.NoArgs,
	Short: "get information about the system",
	Run: func(cmd *cobra.Command, args []string) {
		logWhere()
		confirmInstalled()
		info := systemInfo()
		if currentState.quiet {
			if err := workDir.Join("system.info").WriteFile([]byte(info)); err != nil {
				zboth.Debug().Msgf(info)
			}
		} else {
			if currentState.debug {
				zboth.Debug().Msgf(info)
			} else {
				fmt.Println("This is what we know about the host machine:")
				fmt.Println(info)
			}
		}
	},
}

func init() {
	advancedRootCmd.AddCommand(infoAdvancedRootCmd)
}
