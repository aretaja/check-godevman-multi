# check-godevman-multi
Icinga2 plugin designed to do multiple checks using godevman (not publicly available yet) module

## Usage
```
$ check-godevman-multi --usage
Usage:
        /usr/local/lib/icinga2/libexec/check-godevman-multi <common args> <check_name> [check args]

        Available checks:
                power_gen - Power generator state checks.
                        Alarms are based on provided or default arguments.
                sync_state - Syncronisation state check (Freq and Phase sync).
                        CRITICAL - fsync signal not locked or psync not phase aligned.
                        WARNING - sync source quality is bad.
                        Provides long output and no performance data.

        To get info of available common arguments:
                check-godevman-multi --help
        To get info of available check arguments:
                check-godevman-multi <common args> <check_name> --help
```
```
$ check-godevman-multi --help
Usage of check-godevman-multi:
  -A string
        [authentication protocol pass phrase]
  -H string
        <host ip>
  -V int
        [snmp version] (1|2|3) (default 2)
  -X string
        [privacy protocol pass phrase]
  -a string
        [authentication protocol] (NoAuth|MD5|SHA) (default "MD5")
  -d    Using this parameter will print out debug info
  -l string
        [security level] (noAuthNoPriv|authNoPriv|authPriv) (default "authPriv")
  -u string
        [username|community] (default "public")
  -usage
        Using this parameter will display general usage info and exit
  -v    Using this parameter will display the version number and exit
  -x string
        [privacy protocol] (NoPriv|DES|AES|AES192|AES256|AES192C|AES256C) (default "DES")
```
```
$ check-godevman-multi power_gen --help
Usage of power_gen:
  -cb string
        [critical level for battery voltage] (V*10). ctype - engine (default "120:155")
  -cc string
        [critical level for gen. current] (A). ctype - electrical (default "27")
  -cf string
        [critical level for gen. freq.] (Hz). ctype - electrical (default "46:54")
  -cl string
        [critical level for fuel level] (%). ctype - engine (default "10:100")
  -cp string
        [critical level for gen. power] (kW). ctype - electrical (default "15")
  -ct string
        [critical level for coolant temp] (°C). ctype - engine (default "104")
  -cv string
        [critical level for mains and gen. voltage] (V). ctype - electrical (default "210:250")
  -info
        About check
  -t string
        <check type>
                electrical - check electrical parameters
                engine - check engine parameters
                common - check common status
    
  -wb string
        [warning level for battery voltage] (V*10). ctype - engine (default "130:145")
  -wc string
        [warning level for gen. current] (A). ctype - electrical (default "24")
  -wf string
        [warning level for gen. freq.] (Hz). ctype - electrical (default "48:52")
  -wl string
        [warning level for fuel level] (%). ctype - engine (default "20:100")
  -wp string
        [warning level for gen. power] (kW). ctype - electrical (default "13")
  -wt string
        [warning level for coolant temp] (°C). ctype - engine (default "98")
  -wv string
        [warning level for mains and gen. voltage] (V). ctype - electrical (default "215:245")
```
```
check-godevman-multi sync_state --help
Usage of sync_state:
  -info
        About check
```
