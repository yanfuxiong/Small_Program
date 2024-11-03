//go:build windows
// +build windows

package windows

import (
	"fmt"
	rtkCmd "rtk-cross-share/cmd"
	"syscall"
	"unsafe"
)

var (
	user32          = syscall.NewLazyDLL("user32.dll")
	procFindWindow  = user32.NewProc("FindWindowW")
	procSendMessage = user32.NewProc("SendMessageW")
	hwndCPlusPlus   uintptr
)

const (
	WM_COPYDATA       = 0x004A
	GO_WNDOWS_TITLE   = "GO_FILE_TRANSFER"
	CPP_WINDOWS_TITLE = "CPP_FILE_TRANSFER"
)

type COPYDATASTRUCT struct {
	DwData uintptr
	CbData uint32
	LpData unsafe.Pointer
}

func main() {
	var wc syscall.WNDCLASSEX
	wc.Size = uint32(unsafe.Sizeof(wc))
	wc.WndProc = syscall.NewCallback(WindowProc)
	wc.Instance = syscall.GetModuleHandle(nil)
	wc.ClassName = syscall.StringToUTF16Ptr("GoWindowClass")

	if _, err := syscall.RegisterClassEx(&wc); err != nil {
		fmt.Println("Failed to register window class:", err)
		return
	}

	hwnd, _, _ := syscall.CreateWindowEx(0, wc.ClassName, syscall.StringToUTF16Ptr(GO_WINDOWS_TITLE),
		0, 0, 0, 0, 0, 0, 0, wc.Instance, 0)
	if hwnd == 0 {
		fmt.Println("Failed to create window")
		return
	}

	hwndCPlusPlus, _, _ = procFindWindow.Call(0, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(CPP_WINDOWS_TITLE))))

	var msg syscall.MSG
	for {
		syscall.GetMessage(&msg, 0, 0, 0)
		syscall.TranslateMessage(&msg)
		syscall.DispatchMessage(&msg)
	}
}

func NotifyFileTransCommand() {

}

func WindowProc(hwnd uintptr, msg uint32, wparam, lparam uintptr) uintptr {
	if msg == WM_COPYDATA {
		cds := (*COPYDATASTRUCT)(unsafe.Pointer(lparam))
		receivedData := syscall.UTF16ToString((*[1 << 10]uint16)(cds.LpData)[:cds.CbData/2])
		fmt.Println("Received from windows ", receivedData)
		rtkCmd.SendFileTransCmd(receivedData)

		/*
			if hwndCPlusPlus != 0 {
				response :=
				cds := COPYDATASTRUCT{
					DwData: 1,
					CbData: uint32(len(response) * 2),
					LpData: unsafe.Pointer(syscall.StringToUTF16Ptr(response)),
				}
				procSendMessage.Call(hwndCPlusPlus, WM_COPYDATA, hwnd, uintptr(unsafe.Pointer(&cds)))
			}
		*/
		return 0
	}
	return syscall.DefWindowProc(hwnd, msg, wparam, lparam)
}
