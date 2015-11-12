package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"time"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a binary.",
	Long:  `run starts a binary if there's an actually new Consul event.`,
	Run:   startRun,
}

func startRun(cmd *cobra.Command, args []string) {
	var oldEvent int64
	checkFlags()
	start := time.Now()

	if Exec != "" {
		stdin := readStdin()
		if stdin != "" {
			EventName, lTime := decodeStdin(stdin)
			lTimeString := strconv.FormatInt(int64(lTime), 10)
			ConsulKey := createKey(EventName)

			c, _ := Connect()
			ConsulData := Get(c, ConsulKey)
			if ConsulData != "" {
				oldEvent, _ = strconv.ParseInt(ConsulData, 10, 64)
			}

			if ConsulData == "" || oldEvent < lTime {
				Set(c, ConsulKey, lTimeString)
				runCommand(Exec)
				RunTime(start, "complete", fmt.Sprintf("exec='%s' ltime='%d'", Exec, lTime))
			} else {
				RunTime(start, "duplicate", fmt.Sprintf("exec='%s' ltime='%d'", Exec, lTime))
			}

		} else {
			RunTime(start, "blank", fmt.Sprintf("exec='%s'", Exec))
		}
	}

}

func checkFlags() {
	if Exec == "" {
		fmt.Println("Need a command to exec with '-e'")
		os.Exit(0)
	}
}

func init() {
	RootCmd.AddCommand(runCmd)
}
