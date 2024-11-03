package main

import (
	"fmt"
	rtkBuildConfig "rtk-cross-share/buildConfig"
	rtkCmd "rtk-cross-share/cmd"
	rtkMdns "rtk-cross-share/mdns"
)

func main() {
	fmt.Println("========================")
	fmt.Println("Version: ", rtkBuildConfig.Version)
	fmt.Println("Build Date: ", rtkBuildConfig.BuildDate)
	fmt.Printf("========================\n\n")

	rtkMdns.MdnsCfg = rtkMdns.ParseFlags()
	rtkCmd.Run()
}
