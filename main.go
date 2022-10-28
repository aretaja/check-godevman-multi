// Copyright 2022 by Marko Punnar <marko[AT]aretaja.org>
// Use of this source code is governed by a Apache License 2.0 that can be found in the LICENSE file.

// check-godevman-multi is multipurpose plugin for Icinga2 compatible systems.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/aretaja/godevman"
	"github.com/kr/pretty"
)

// Version of release
const Version = "0.0.1"

// available checks
var checksInfo = map[string][]string{
	"power_gen": {"Power generator state checks.",
		"\tAlarms are based on provided or default arguments."},
	"sync_state": {"Syncronisation state check (Freq and Phase sync).",
		"\tCRITICAL - fsync signal not locked or psync not phase aligned.",
		"\tWARNING - sync source quality is bad.",
		"\tProvides long output and no performance data."},
}

// check options
type checkParams struct {
	execName  string
	checkName string
	subCheck  string
	subArgs   []string
	devParams godevman.Dparams
	dbg       bool
}

// Initialize CheckArgs using submitted command line options
func initParams() (checkParams, error) {
	H := flag.String("H", "", "<host ip>")
	V := flag.Int("V", 2, "[snmp version] (1|2|3)")
	u := flag.String("u", "public", "[username|community]")
	a := flag.String("a", "MD5", "[authentication protocol] (NoAuth|MD5|SHA)")
	A := flag.String("A", "", "[authentication protocol pass phrase]")
	l := flag.String("l", "authPriv", "[security level] (noAuthNoPriv|authNoPriv|authPriv)")
	x := flag.String("x", "DES", "[privacy protocol] (NoPriv|DES|AES|AES192|AES256|AES192C|AES256C)")
	X := flag.String("X", "", "[privacy protocol pass phrase]")
	d := flag.Bool("d", false, "Using this parameter will print out debug info")
	v := flag.Bool("v", false, "Using this parameter will display the version number and exit")
	usage := flag.Bool("usage", false, "Using this parameter will display general usage info and exit")

	flag.Parse()

	params := checkParams{
		devParams: godevman.Dparams{
			Ip: *H,
			SnmpCred: godevman.SnmpCred{
				Ver:      *V,
				User:     *u,
				Prot:     *a,
				Pass:     *A,
				Slevel:   *l,
				PrivProt: *x,
				PrivPass: *X,
			},
		},
		dbg: *d,
	}

	// Get executable name
	n, err := os.Executable()
	if err != nil {
		n = "check-godevman-multi"
	}
	params.execName = n

	// Show version
	if *v {
		fmt.Println("plugin version " + Version)
		os.Exit(3)
	}

	// Show usage info
	if *usage {
		fmt.Printf("Usage:\n\t%s <common args> <check_name> [check args]\n\n\tAvailable checks:\n", params.execName)
		keys := make([]string, 0, len(checksInfo))
		for k := range checksInfo {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			fmt.Printf("\t\t%s - %s\n", k, strings.Join(checksInfo[k], "\n\t\t"))
		}
		fmt.Print("\n\tTo get info of available common arguments:\n\t\tcheck-godevman-multi --help\n" +
			"\tTo get info of available check arguments:\n\t\tcheck-godevman-multi <common args> <check_name> --help\n")
		os.Exit(3)
	}

	// Retrieve the remaining cli arguments
	rargs := flag.Args()
	if len(rargs) == 0 {
		return params, fmt.Errorf("check name missing")
	}

	params.subCheck, params.subArgs = rargs[0], rargs[1:]
	// DEBUG
	if params.dbg {
		fmt.Printf("params: %# v\n", pretty.Formatter(params))
	}

	return params, nil
}

func main() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)

	p, err := initParams()
	if err != nil {
		log.Printf("error: initParams - %v\n", err)
		os.Exit(3)
	}

	switch p.subCheck {
	case "sync_state":
		p.checkName = "sync_state"
		c := checkSyncro{}
		c.checkParams = p
		c.state()
	case "power_gen":
		p.checkName = "power_gen"
		c := checkPowerGen{}
		c.checkParams = p
		c.state()
	default:
		log.Printf("error: unrecognized check name - %s\n", p.subCheck)
		os.Exit(3)
	}
}

// Initialize device
func (sd *checkParams) initDevice() any {
	p := sd.devParams
	device, err := godevman.NewDevice(&p)
	if err != nil {
		log.Printf("error: godevman.NewDevice: %v", err)
		os.Exit(3)
	}

	md := device.Morph()
	// DEBUG
	if sd.dbg {
		fmt.Printf("godevman morphed device: %# v\n", pretty.Formatter(md))
	}

	return md
}
