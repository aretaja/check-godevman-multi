module github.com/aretaja/check-godevman-multi

go 1.19

// For local development
// replace github.com/aretaja/godevman => ../godevman

require github.com/aretaja/godevman v0.0.1-devel.3

require google.golang.org/genproto/googleapis/rpc v0.0.0-20231016165738-49dd2c1f3d0b // indirect

require (
	github.com/aretaja/icingahelper v1.1.1
	github.com/aretaja/snmphelper v1.1.3 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/goexpect v0.0.0-20210430020637-ab937bf7fd6f // indirect
	github.com/google/goterm v0.0.0-20200907032337-555d40f16ae2 // indirect
	github.com/gosnmp/gosnmp v1.36.1 // indirect
	github.com/kr/pretty v0.3.1
	github.com/kr/text v0.2.0 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/praserx/ipconv v1.2.1 // indirect
	github.com/rogpeppe/go-internal v1.11.0 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	google.golang.org/grpc v1.59.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
)
