package debug

import (
	"bufio"
	"fmt"
	"os"
	rtkCommon "rtk-cross-share/common"
	rtkFileDrop "rtk-cross-share/filedrop"
	rtkGlobal "rtk-cross-share/global"
	rtkPlatform "rtk-cross-share/platform"
	rtkUtils "rtk-cross-share/utils"
	"strings"
)

type TestCase struct {
	FileName     string
	FileSizeHigh uint32
	FileSizeLow  uint32
}

func DebugCmdLine() {
	test_case_1 := TestCase{
		FileName:     "D:\\share\\PC1.txt",
		FileSizeHigh: 0,
		FileSizeLow:  13902,
	}

	test_case_2 := TestCase{
		FileName:     "D:\\share\\library.zip",
		FileSizeHigh: 0,
		FileSizeLow:  6452991,
	}

	test_case_3 := TestCase{
		FileName:     "E:\\CODE\\png\\3.png",
		FileSizeHigh: 0,
		FileSizeLow:  1508,
	}

	test_case_4 := TestCase{
		FileName:     "/Users/hp/myGolang/test.png",
		FileSizeHigh: 0,
		FileSizeLow:  109939,
	}

	test_case_5 := TestCase{
		FileName:     "/Users/hp/myGolang/test.mp4",
		FileSizeHigh: 0,
		FileSizeLow:  8986659,
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter text to debug:")
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("You entered:", line)
		if strings.Contains(line, "COPY_TEST_1") {
			rtkGlobal.Handler.CopyFilePath.Store(test_case_1.FileName)
			rtkGlobal.Handler.CopyDataSize.SizeHigh = test_case_1.FileSizeHigh
			rtkGlobal.Handler.CopyDataSize.SizeLow = test_case_1.FileSizeLow
			//rtkGlobal.Handler.AppointIpAddr = "192.168.153.32:6666"
			fmt.Println("Clipboard file content:", rtkGlobal.Handler.CopyFilePath)
		} else if strings.Contains(line, "COPY_TEST_2") {
			rtkGlobal.Handler.CopyFilePath.Store(test_case_2.FileName)
			rtkGlobal.Handler.CopyDataSize.SizeHigh = test_case_2.FileSizeHigh
			rtkGlobal.Handler.CopyDataSize.SizeLow = test_case_2.FileSizeLow
			//rtkGlobal.Handler.AppointIpAddr = "192.168.153.32:6666"
			fmt.Println("Clipboard file content:", rtkGlobal.Handler.CopyFilePath)
		} else if strings.Contains(line, "COPY_TEST_3") {
			rtkGlobal.Handler.CopyFilePath.Store(test_case_3.FileName)
			rtkGlobal.Handler.CopyDataSize.SizeHigh = test_case_3.FileSizeHigh
			rtkGlobal.Handler.CopyDataSize.SizeLow = test_case_3.FileSizeLow
			//rtkGlobal.Handler.AppointIpAddr = "192.168.153.237:6622"
			fmt.Println("Clipboard file content:", rtkGlobal.Handler.CopyFilePath)
		} else if strings.Contains(line, "PASTE_FILE") {
			rtkPlatform.GoClipboardPasteFileCallback("123")
		} else if strings.Contains(line, "FILE_DROP_TEST_1") {
			rtkGlobal.Handler.CopyFilePath.Store(test_case_4.FileName)
			rtkGlobal.Handler.CopyDataSize.SizeHigh = test_case_4.FileSizeHigh
			rtkGlobal.Handler.CopyDataSize.SizeLow = test_case_4.FileSizeLow
			var fileInfo = rtkCommon.FileInfo{
				FileSize_: rtkCommon.FileSize{
					SizeHigh: test_case_4.FileSizeHigh,
					SizeLow:  test_case_4.FileSizeLow,
				},
				FilePath: test_case_4.FileName,
			}
			rtkFileDrop.SendFileDropCmd(rtkCommon.FILE_DROP_REQUEST, fileInfo)
		} else if strings.Contains(line, "FILE_DROP_TEST_2") {
			rtkGlobal.Handler.CopyFilePath.Store(test_case_5.FileName)
			rtkGlobal.Handler.CopyDataSize.SizeHigh = test_case_5.FileSizeHigh
			rtkGlobal.Handler.CopyDataSize.SizeLow = test_case_5.FileSizeLow
			var fileInfo = rtkCommon.FileInfo{
				FileSize_: rtkCommon.FileSize{
					SizeHigh: test_case_5.FileSizeHigh,
					SizeLow:  test_case_5.FileSizeLow,
				},
				FilePath: test_case_5.FileName,
			}
			rtkFileDrop.SendFileDropCmd(rtkCommon.FILE_DROP_REQUEST, fileInfo)
		} else if strings.Contains(line, "GetClientList") {
			clientList := rtkUtils.GetClientList()
			fmt.Println("GetClientList : ", clientList)
		}
	}
}
