// Syncronization status check
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sort"
	"strings"

	"github.com/aretaja/godevman"
	"github.com/aretaja/icingahelper"
	"github.com/kr/pretty"
)

// Adds syncro check functionality to checkParams type
type checkSyncro struct {
	checkParams
}

func (c *checkSyncro) state() {
	c.initSubParams()

	// Initialize new check object
	check := icingahelper.NewCheck("SYNC")

	// Exit if no host ip submitted
	if net.ParseIP(c.devParams.Ip) == nil {
		log.Println("error: valid host ip is required")
		os.Exit(check.RetVal())
	}

	md := c.initDevice()

	fd, ok := md.(godevman.DevFreqSyncReader)
	if !ok {
		log.Println("error: freq sync state check is not supported on this device type")
		os.Exit(check.RetVal())
	}

	pd, ok := md.(godevman.DevPhaseSyncReader)
	if !ok {
		log.Println("error: phase sync state check is not supported on this device type")
		os.Exit(check.RetVal())
	}

	// Get freq sync information from device
	resf, err := fd.FreqSyncInfo()
	if err != nil {
		log.Printf("error: FreqSyncInfo: %v", err)
		l := check.RetVal()
		if strings.HasSuffix(err.Error(), "not configured") {
			l = 0
		}
		os.Exit(l)
	}
	// DEBUG
	if c.dbg {
		fmt.Printf("freq sync info: %# v\n", pretty.Formatter(resf))
	}

	// Get phase sync information from device
	resp, err := pd.PhaseSyncInfo()
	if err != nil {
		log.Printf("error: PhaseSyncInfo: %v", err)
		l := check.RetVal()
		if strings.HasSuffix(err.Error(), "not configured") {
			l = 0
		}
		os.Exit(l)
	}
	// DEBUG
	if c.dbg {
		fmt.Printf("phase sync info: %# v\n", pretty.Formatter(resp))
	}

	check.SetRetVal(0)
	// Freq sync data to icingahelper
	fl := ""
	if resf.SrcsQaLevel != nil {
		p := resf.SrcsQaLevel
		keys := make([]string, 0, len(p))
		for k := range p {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		m := []string{"Configured frequency sync sources:"}
		for _, k := range keys {
			m = append(m, fmt.Sprintf(" %s: %s", k, p[k]))
		}
		fl = strings.Join(m, "\n")
	}

	if resf.ClockMode.IsSet {
		val := resf.ClockMode.Value
		if val != "locked" {
			check.SetRetVal(2)
		}
		check.AddMsg(check.RetVal(), fmt.Sprintf("Fsync Mode: %s", val), fmt.Sprintf("%s\nFreq sync ", fl))
	} else {
		level := check.RetVal()
		if level != 2 {
			level = 3
		}
		check.SetRetVal(level)
		check.AddMsg(3, "Fsync Mode: Na", fmt.Sprintf("%s\nFreq sync ", fl))
	}

	if resf.ClockQaLevel.IsSet {
		val := resf.ClockQaLevel.Value
		l := 0
		if val != "PRC" {
			l = 1
			if check.RetVal() != 2 {
				check.SetRetVal(1)
			}
		}
		check.AddMsg(l, fmt.Sprintf("Fsync Qa: %s", val), "")
	} else {
		level := check.RetVal()
		if level != 2 {
			level = 3
		}
		check.SetRetVal(level)
		check.AddMsg(3, "Fsync Qa: Na", "")
	}

	// Phase sync data to icingahelper
	pm := []string{}

	if resp.SrcsState != nil {
		p := resp.SrcsState
		keys := make([]string, 0, len(p))
		for k := range p {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		pm = append(pm, "Configured phase sync sources:")
		for _, k := range keys {
			pm = append(pm, fmt.Sprintf(" %s: %s", k, p[k]))
		}
	}

	if resp.HopsToGm.IsSet {
		pm = append(pm, fmt.Sprintf("Hops to GM: %d", resp.HopsToGm.Value))
	} else {
		pm = append(pm, "Hops to GM: Na")
	}

	pl := strings.Join(pm, "\n")

	if resp.State.IsSet {
		val := resp.State.Value
		l := 0
		if val != "phaseAligned" {
			check.SetRetVal(2)
			l = 2
		}
		check.AddMsg(l, fmt.Sprintf("PTP Mode: %s", val), fmt.Sprintf("%s\nPhase sync ", pl))
	} else {
		level := check.RetVal()
		if level == 0 {
			level = 3
		}
		check.SetRetVal(level)
		check.AddMsg(3, "PTP Mode: Na", fmt.Sprintf("%s\nPhase sync ", pl))
	}

	if resp.ParentGmClass.IsSet {
		val := resp.ParentGmClass.Value
		l := 0
		if val != "prtcLock(6)" {
			l = 2
			if val != "holdover(7)" {
				if check.RetVal() != 2 {
					check.SetRetVal(1)
					l = 1
				}
			} else {
				check.SetRetVal(2)
			}
		}
		check.AddMsg(l, fmt.Sprintf("PTP GM Class: %s", val), "")
	} else {
		level := check.RetVal()
		if level == 0 {
			level = 3
		}
		check.SetRetVal(level)
		check.AddMsg(3, "PTP GM Class: Na", "")
	}

	if resp.ParentGmIdent.IsSet {
		check.AddMsg(0, fmt.Sprintf("GrandMaster: %s", resp.ParentGmIdent.Value), "")
	} else {

		check.AddMsg(0, fmt.Sprintf("GrandMaster: %s", resp.ParentGmIdent.Value), "")
	}

	fmt.Print(check.FinalMsg())
	os.Exit(check.RetVal())
}

func (c *checkSyncro) initSubParams() {
	flag := flag.NewFlagSet("sync_state", flag.ExitOnError)
	var info = flag.Bool("info", false, "About check")
	flag.Parse(c.subArgs)

	// Show info about check
	if *info {
		if i, ok := checksInfo[c.checkName]; ok {
			fmt.Printf("%s: %s\n", c.checkName, strings.Join(i, "\n"))
		} else {
			fmt.Printf("%s: no additional information\n", c.checkName)
		}
		os.Exit(3)
	}
}
