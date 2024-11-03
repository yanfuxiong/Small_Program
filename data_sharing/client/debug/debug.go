package debug

import (
	"bufio"
	"fmt"
	"os"
	rtkGlobal "rtk-cross-share/global"
	rtkPlatform "rtk-cross-share/platform"
	"strings"
)

func DebugCmdLine() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter text to debug:")
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("You entered:", line)
		if strings.Contains(line, "COPY_TEST_1") {
			rtkGlobal.Handler.CopyFilePath.Store("/Users/hp/myGolang/test.png")
			rtkGlobal.Handler.CopyDataSize.SizeHigh = 0
			rtkGlobal.Handler.CopyDataSize.SizeLow = 109939
			fmt.Println("Clipboard file content:", rtkGlobal.Handler.CopyFilePath)
		} else if strings.Contains(line, "COPY_TEST_2") {
			rtkGlobal.Handler.CopyFilePath.Store("/Users/hp/myGolang/test2.png")
			rtkGlobal.Handler.CopyDataSize.SizeHigh = 0
			rtkGlobal.Handler.CopyDataSize.SizeLow = 3241468
			fmt.Println("Clipboard file content:", rtkGlobal.Handler.CopyFilePath)
		} else if strings.Contains(line, "COPY_TEST_3") {
			rtkGlobal.Handler.CopyFilePath.Store("D:\\xyf\\client2.exe")
			rtkGlobal.Handler.CopyDataSize.SizeHigh = 0
			rtkGlobal.Handler.CopyDataSize.SizeLow = 33904128
			rtkGlobal.Handler.IpAddr = "192.168.153.237:6622"
			fmt.Println("Clipboard file content:", rtkGlobal.Handler.CopyFilePath)
		} else if strings.Contains(line, "COPY_TEST_4") {
			rtkGlobal.Handler.CopyFilePath.Store("D:\\xyf\\p2p.log")
			rtkGlobal.Handler.CopyDataSize.SizeHigh = 0
			rtkGlobal.Handler.CopyDataSize.SizeLow = 706
			rtkGlobal.Handler.IpAddr = "192.168.153.237:6622"
			fmt.Println("Clipboard file content:", rtkGlobal.Handler.CopyFilePath)
		} else if strings.Contains(line, "PASTE_FILE") {
			rtkPlatform.GoClipboardPasteFileCallback("123")
		}
	}
}
