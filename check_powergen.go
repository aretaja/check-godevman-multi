// Power generator status check
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/aretaja/godevman"
	"github.com/aretaja/icingahelper"
	"github.com/kr/pretty"
)

// Adds syncro check functionality to checkParams type
type checkPowerGen struct {
	subParams struct{ ctype, wVolt, cVolt, wCur, cCur, wPow, cPow, wFreq, cFreq, wBat, cBat, wFuel, cFuel, wTemp, cTemp string }
	checkParams
}

func (c *checkPowerGen) state() {
	c.initSubParams()

	// Initialize new check object
	check := icingahelper.NewCheck("GEN")

	// Exit if no host ip submitted
	if net.ParseIP(c.devParams.Ip) == nil {
		log.Println("error: valid host ip is required")
		os.Exit(check.RetVal())
	}

	md := c.initDevice()

	d, ok := md.(godevman.DevGenReader)
	if !ok {
		log.Println("error: power generator state check is not supported on this device type")
		os.Exit(check.RetVal())
	}

	switch c.subParams.ctype {
	case "common":
		res, err := c.getInfo(d, "Common")
		if err != nil {
			log.Printf("error: %v", err)
			os.Exit(check.RetVal())
		}
		c.common(check, res)
	case "electrical":
		res, err := c.getInfo(d, "Electrical")
		if err != nil {
			os.Exit(check.RetVal())
		}
		err = c.electrical(check, res)
		if err != nil {
			log.Printf("error: %v", err)
			os.Exit(check.RetVal())
		}
	case "engine":
		res, err := c.getInfo(d, "Engine")
		if err != nil {
			os.Exit(check.RetVal())
		}
		err = c.engine(check, res)
		if err != nil {
			log.Printf("error: %v", err)
			os.Exit(check.RetVal())
		}
	default:
		log.Printf("error: unknown check type - %s", c.subParams.ctype)
		os.Exit(check.RetVal())
	}

	fmt.Print(check.FinalMsg())
	os.Exit(check.RetVal())
}

func (c *checkPowerGen) initSubParams() {
	flag := flag.NewFlagSet("power_gen", flag.ExitOnError)
	var t = flag.String("t", "", "<check type>\n"+
		"\telectrical - check electrical parameters\n"+
		"\tengine - check engine parameters\n"+
		"\tcommon - check common status\n",
	)
	var wv = flag.String("wv", "215:245", "[warning level for mains and gen. voltage] (V). ctype - electrical")
	var cv = flag.String("cv", "210:250", "[critical level for mains and gen. voltage] (V). ctype - electrical")
	var wc = flag.String("wc", "24", "[warning level for gen. current] (A). ctype - electrical")
	var cc = flag.String("cc", "27", "[critical level for gen. current] (A). ctype - electrical")
	var wp = flag.String("wp", "13", "[warning level for gen. power] (kW). ctype - electrical")
	var cp = flag.String("cp", "15", "[critical level for gen. power] (kW). ctype - electrical")
	var wf = flag.String("wf", "48:52", "[warning level for gen. freq.] (Hz). ctype - electrical")
	var cf = flag.String("cf", "46:54", "[critical level for gen. freq.] (Hz). ctype - electrical")
	var wb = flag.String("wb", "130:145", "[warning level for battery voltage] (V*10). ctype - engine")
	var cb = flag.String("cb", "120:155", "[critical level for battery voltage] (V*10). ctype - engine")
	var wl = flag.String("wl", "20:100", "[warning level for fuel level] (%). ctype - engine")
	var cl = flag.String("cl", "10:100", "[critical level for fuel level] (%). ctype - engine")
	var wt = flag.String("wt", "98", "[warning level for coolant temp] (°C). ctype - engine")
	var ct = flag.String("ct", "104", "[critical level for coolant temp] (°C). ctype - engine")
	var info = flag.Bool("info", false, "About check")

	flag.Parse(c.subArgs)

	c.subParams.ctype = *t
	c.subParams.wVolt = *wv
	c.subParams.cVolt = *cv
	c.subParams.wCur = *wc
	c.subParams.cCur = *cc
	c.subParams.wPow = *wp
	c.subParams.cPow = *cp
	c.subParams.wFreq = *wf
	c.subParams.cFreq = *cf
	c.subParams.wBat = *wb
	c.subParams.cBat = *cb
	c.subParams.wFuel = *wl
	c.subParams.cFuel = *cl
	c.subParams.wTemp = *wt
	c.subParams.cTemp = *ct
	// DEBUG
	if c.dbg {
		fmt.Printf("powergen params: %# v\n", pretty.Formatter(c))
	}

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

func (c *checkPowerGen) getInfo(d godevman.DevGenReader, t string) (godevman.GenInfo, error) {
	res, err := d.GeneratorInfo([]string{t})
	if err != nil {
		return res, err
	}
	// DEBUG
	if c.dbg {
		fmt.Printf("powergen %s info: %# v\n", t, pretty.Formatter(res))
	}
	return res, err
}

func (c *checkPowerGen) common(check *icingahelper.IcingaCheck, i godevman.GenInfo) {
	check.SetRetVal(0)
	if i.GenMode.IsSet {
		val := i.GenMode.Value
		if val != "Auto" {
			check.SetRetVal(2)
		}
		check.AddMsg(check.RetVal(), fmt.Sprintf("Mode: %s", val), "")
	} else {
		level := check.RetVal()
		if level != 2 {
			level = 3
		}
		check.SetRetVal(level)
		check.AddMsg(level, "Mode: Na", "")
	}

	if i.BreakerState.IsSet {
		val := i.BreakerState.Value
		if val != "MainsOper" {
			check.SetRetVal(2)
		}
		check.AddMsg(check.RetVal(), fmt.Sprintf("Breaker: %s", val), "")
	} else {
		level := check.RetVal()
		if level != 2 {
			level = 3
		}
		check.SetRetVal(level)
		check.AddMsg(level, "Breaker: Na", "")
	}

	if i.EngineState.IsSet {
		val := i.EngineState.Value
		if val != "Ready" {
			check.SetRetVal(2)
		}
		check.AddMsg(check.RetVal(), fmt.Sprintf("Engine: %s", val), "")
	} else {
		level := check.RetVal()
		if level != 2 {
			level = 3
		}
		check.SetRetVal(level)
		check.AddMsg(level, "Engine: Na", "")
	}
}

func (c *checkPowerGen) electrical(check *icingahelper.IcingaCheck, i godevman.GenInfo) error {
	data := map[string]godevman.SensorVal{
		"Mains Voltage L1": i.MainsVoltL1,
		"Mains Voltage L2": i.MainsVoltL2,
		"Mains Voltage L3": i.MainsVoltL3,
		"Gen Voltage L1":   i.GenVoltL1,
		"Gen Voltage L2":   i.GenVoltL2,
		"Gen Voltage L3":   i.GenVoltL3,
		"Gen Current L1":   i.GenCurrentL1,
		"Gen Current L2":   i.GenCurrentL2,
		"Gen Current L3":   i.GenCurrentL3,
		"Gen Power":        i.GenPower,
		"Gen Frequency":    i.GenFreq,
	}

	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		switch data[k].Unit {
		case "V":
			if data[k].IsSet {
				val := data[k].Value
				level := 0
				if strings.Contains(k, "Mains") || val != 0 {
					l, err := check.AlarmLevel(int64(val), c.subParams.wVolt, c.subParams.cVolt)
					if err != nil {
						return fmt.Errorf("voltage alarm level error: %v", err)
					}
					level = l
				}

				check.AddMsg(level, fmt.Sprintf("%s: %dV", k, val), "")
				check.AddPerfData(fmt.Sprintf("'%s'", k), strconv.Itoa(int(val)), "", c.subParams.wVolt, c.subParams.wVolt, "0", "")
			} else {
				check.AddMsg(3, fmt.Sprintf("%s: Na", k), "")
			}
		case "A":
			if data[k].IsSet {
				val := data[k].Value
				level, err := check.AlarmLevel(int64(val), c.subParams.wCur, c.subParams.cCur)
				if err != nil {
					return fmt.Errorf("current alarm level error: %v", err)
				}

				check.AddMsg(level, fmt.Sprintf("%s: %dA", k, val), "")
				check.AddPerfData(fmt.Sprintf("'%s'", k), strconv.Itoa(int(val)), "", c.subParams.wCur, c.subParams.cCur, "0", "")
			} else {
				check.AddMsg(3, fmt.Sprintf("%s: Na", k), "")
			}
		case "kW":
			if data[k].IsSet {
				val := data[k].Value
				level, err := check.AlarmLevel(int64(val), c.subParams.wPow, c.subParams.cPow)
				if err != nil {
					return fmt.Errorf("power alarm level error: %v", err)
				}

				check.AddMsg(level, fmt.Sprintf("%s: %dkW", k, val), "")
				check.AddPerfData(fmt.Sprintf("'%s'", k), strconv.Itoa(int(val)), "", c.subParams.wPow, c.subParams.cPow, "0", "")
			} else {
				check.AddMsg(3, fmt.Sprintf("%s: Na", k), "")
			}
		case "Hz":
			if data[k].IsSet {
				val := data[k].Value
				level := 0
				if val != 0 {
					l, err := check.AlarmLevel(int64(val), c.subParams.wFreq, c.subParams.cFreq)
					if err != nil {
						return fmt.Errorf("power alarm level error: %v", err)
					}
					level = l
				}

				rVal := float64(val) / float64(data[k].Divisor)
				check.AddMsg(level, fmt.Sprintf("%s: %.1fHz", k, rVal), "")
				check.AddPerfData(fmt.Sprintf("'%s'", k), strconv.Itoa(int(val)), "", c.subParams.wFreq, c.subParams.cFreq, "0", "")
			} else {
				check.AddMsg(3, fmt.Sprintf("%s: Na", k), "")
			}
		default:
			return fmt.Errorf("unexpected results from godevman")
		}
	}

	return nil
}

func (c *checkPowerGen) engine(check *icingahelper.IcingaCheck, i godevman.GenInfo) error {
	data := map[string]godevman.SensorVal{
		"Running Hours":       i.RunHours,
		"Fuel level":          i.FuelLevel,
		"Fuel Consumption":    i.FuelConsum,
		"Battery Voltage":     i.BatteryVolt,
		"Coolant Temperature": i.CoolantTemp,
	}

	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		switch k {
		case "Battery Voltage":
			if data[k].IsSet {
				val := data[k].Value
				level, err := check.AlarmLevel(int64(val), c.subParams.wBat, c.subParams.cBat)
				if err != nil {
					return fmt.Errorf("voltage alarm level error: %v", err)
				}

				rVal := float64(val) / float64(data[k].Divisor)
				check.AddMsg(level, fmt.Sprintf("%s: %.1fV", k, rVal), "")
				check.AddPerfData(fmt.Sprintf("'%s'", k), strconv.Itoa(int(val)), "", c.subParams.wBat, c.subParams.cBat, "0", "")
			} else {
				check.AddMsg(3, fmt.Sprintf("%s: Na", k), "")
			}
		case "Coolant Temperature":
			if data[k].IsSet {
				val := data[k].Value
				level, err := check.AlarmLevel(int64(val), c.subParams.wTemp, c.subParams.cTemp)
				if err != nil {
					return fmt.Errorf("coolant alarm level error: %v", err)
				}

				check.AddMsg(level, fmt.Sprintf("%s: %d%s", k, val, data[k].Unit), "")
				check.AddPerfData(fmt.Sprintf("'%s'", k), strconv.Itoa(int(val)), "", c.subParams.wTemp, c.subParams.cTemp, "0", "")
			} else {
				check.AddMsg(3, fmt.Sprintf("%s: Na", k), "")
			}
		case "Fuel level":
			if data[k].IsSet {
				val := data[k].Value
				level, err := check.AlarmLevel(int64(val), c.subParams.wFuel, c.subParams.cFuel)
				if err != nil {
					return fmt.Errorf("fuel alarm level error: %v", err)
				}

				check.AddMsg(level, fmt.Sprintf("%s: %d%s", k, val, data[k].Unit), "")
				check.AddPerfData(fmt.Sprintf("'%s'", k), strconv.Itoa(int(val)), data[k].Unit, c.subParams.wFuel, c.subParams.cFuel, "0", "")
			} else {
				check.AddMsg(3, fmt.Sprintf("%s: Na", k), "")
			}
		default:
			if data[k].IsSet {
				val := data[k].Value

				rVal := float64(val)
				if data[k].Divisor != 0 {
					rVal = rVal / float64(data[k].Divisor)
				}

				check.AddMsg(0, fmt.Sprintf("%s: %.1f%s", k, rVal, data[k].Unit), "")
				check.AddPerfData(fmt.Sprintf("'%s'", k), strconv.Itoa(int(val)), "", "", "", "0", "")
			} else {
				check.AddMsg(3, fmt.Sprintf("%s: Na", k), "")
			}
		}
	}

	name := "Number of Starts"
	if i.NumStarts.IsSet {
		val := i.NumStarts.Value
		check.AddMsg(0, fmt.Sprintf("%s: %d", name, val), "")
		check.AddPerfData(fmt.Sprintf("'%s'", name), strconv.Itoa(int(val)), "", "", "", "0", "")
	} else {
		level := check.RetVal()
		if level != 2 {
			level = 3
		}
		check.AddMsg(level, fmt.Sprintf("%s: Na", name), "")
	}

	return nil
}
