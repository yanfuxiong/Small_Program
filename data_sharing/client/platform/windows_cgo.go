//go:build windows
// +build windows

package platform

/*
#cgo LDFLAGS: -L../../clipboard/libs -lclipboard -Wl,-Bstatic
#cgo LDFLAGS: -lgdi32 -lole32 -luuid -lstdc++
#include <stdlib.h>
#include <wchar.h>
#include "../../clipboard/MSPaste/MSCommonExt.h"

typedef void (*ClipboardCopyFileCallback)(wchar_t*, unsigned long, unsigned long);
typedef void (*ClipboardPasteFileCallback)(char*);
typedef void (*FileDropCmdCallback)(unsigned long, wchar_t*);
typedef void (*ClipboardCopyImgCallback)(IMAGE_HEADER, unsigned char*, unsigned long);

extern void SetClipboardCopyFileCallback(ClipboardCopyFileCallback callback);
extern void SetClipboardPasteFileCallback(ClipboardPasteFileCallback callback);
extern void SetFileDropCmdCallback(FileDropCmdCallback callback);
extern void SetClipboardCopyImgCallback(ClipboardCopyImgCallback callback);

extern void StartClipboardMonitor();
extern void StopClipboardMonitor();

extern void SetupDstPasteFile(wchar_t* desc, wchar_t* fileName, unsigned long fileSizeHigh, unsigned long fileSizeLow);
extern void SetupFileDrop(wchar_t* desc, wchar_t* fileName, unsigned long fileSizeHigh, unsigned long fileSizeLow);
extern void SetupDstPasteImage(wchar_t* desc, IMAGE_HEADER imgHeader, unsigned long dataSize);
extern void DataTransfer(unsigned char* data, unsigned int size);
extern void EventHandle(EVENT_TYPE event);

void clipboardCopyFileCallback(wchar_t* content, unsigned long, unsigned long);
void clipboardPasteFileCallback(char* content);
void fileDropCmdCallback(unsigned long, wchar_t*);
void clipboardCopyImgCallback(IMAGE_HEADER, unsigned char*, unsigned long);
*/
import "C"
import (
	"context"
	"fmt"
	"log"
	"os"
	rtkCommon "rtk-cross-share/common"
	rtkFileDrop "rtk-cross-share/filedrop"
	rtkGlobal "rtk-cross-share/global"
	rtkUtils "rtk-cross-share/utils"
	"unsafe"

	"github.com/libp2p/go-libp2p/core/crypto"
	"golang.design/x/clipboard"
)

type Callback interface {
	CallbackMethod(result string)
	CallbackMethodImage(content []byte)
	LogMessageCallback(msg string)
	EventCallback(event int)
	CallbackMethodFileConfirm(platform string, fileSize int64)
}

var CallbackInstance Callback = nil

func SetCallback(cb Callback) {
	CallbackInstance = cb
}

type CallbackFunc func(rtkCommon.ClipboardFmtType)

var CallbackInstanceResetCB CallbackFunc = nil

func SetResetCBCallback(cb CallbackFunc) {
	CallbackInstanceResetCB = cb
}

func WatchClipboardText(ctx context.Context, resultChan chan<- rtkCommon.ClipBoardData) {
	log.Printf("windows WatchClipboardText...")
	changeText := clipboard.Watch(ctx, clipboard.FmtText)

	for {
		select {
		case <-ctx.Done():
			return
		case contentText := <-changeText:
			if string(contentText) == "" || len(contentText) == 0 {
				continue
			}
			log.Printf("windows watchClipboardText - got new message: [%s]", string(contentText))
			if CallbackInstanceResetCB != nil {
				CallbackInstanceResetCB(rtkCommon.TEXT)
			}

			hash, err := rtkUtils.CreateMD5Hash(contentText)
			if err != nil {
				log.Fatalf("Failed to create hash: %v", err)
			}
			resultChan <- rtkCommon.ClipBoardData{
				SourceID: rtkGlobal.NodeInfo.ID,
				Content:  contentText,
				FmtType:  rtkCommon.TEXT,
				Hash:     hash.B58String(),
			}
		}
	}
}

func wcharToString(wstr *C.wchar_t) string {
	var goStr string
	for ptr := wstr; *ptr != 0; ptr = (*C.wchar_t)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + unsafe.Sizeof(*ptr))) {
		goStr += string(rune(*ptr))
	}
	return goStr
}

func stringToWChar(str string) *C.wchar_t {
	utf16Str := make([]uint16, len(str)+1)
	for i, r := range str {
		utf16Str[i] = uint16(r)
	}
	utf16Str[len(str)] = 0
	return (*C.wchar_t)(unsafe.Pointer(&utf16Str[0]))
}

//export clipboardCopyFileCallback
func clipboardCopyFileCallback(cContent *C.wchar_t, cFileSizeHigh C.ulong, cFileSizeLow C.ulong) {
	content := wcharToString(cContent)
	fileSizeHigh := uint32(cFileSizeHigh)
	fileSizeLow := uint32(cFileSizeLow)
	/*rtkGlobal.Handler.CopyFilePath.Store(content)
	rtkGlobal.Handler.CopyDataSize.SizeHigh = fileSizeHigh
	rtkGlobal.Handler.CopyDataSize.SizeLow = fileSizeLow*/
	//rtkGlobal.Handler.AppointIpAddr = ipAddr
	var fileInfo = rtkCommon.FileInfo{
		FileSize_: rtkCommon.FileSize{
			SizeHigh: fileSizeHigh,
			SizeLow:  fileSizeLow,
		},
		FilePath: content,
	}
	rtkFileDrop.SendFileDropCmd(rtkCommon.FILE_DROP_REQUEST, fileInfo)
	fmt.Println("windows Clipboard file content:", rtkGlobal.Handler.CopyFilePath, "fileSize high:", fileSizeHigh, "low:", fileSizeLow)
}

// For DEBUG
func GoClipboardPasteFileCallback(content string) {
	cContent := C.CString(content)
	defer C.free(unsafe.Pointer(cContent))
	clipboardPasteFileCallback(cContent)
}

//export clipboardPasteFileCallback
func clipboardPasteFileCallback(cContent *C.char) {
	content := C.GoString(cContent)
	rtkGlobal.Handler.CtxMutex.Lock()
	defer rtkGlobal.Handler.CtxMutex.Unlock()
	rtkGlobal.Handler.State.State = rtkCommon.DEST_INIT
	fmt.Println("Paste Clipboard file content:", content)
}

//export fileDropCmdCallback
func fileDropCmdCallback(cCmd C.ulong, cContent *C.wchar_t) {
	fmt.Println("File Drop CMD:", cCmd)
	value := uint64(cCmd)

	// REF: FILE_DROP_CMD
	// rtkGlobal.Handler.CtxMutex is already protected in rtkPlatform.GoSetupFileDrop
	switch value {
	case 0: // FILE_DROP_REQUEST
		{
		}
	case 1: // FILE_DROP_ACCEPT
		{
			path := wcharToString(cContent)
			rtkGlobal.Handler.DstFilePath = path
			rtkGlobal.Handler.State.State = rtkCommon.DEST_INIT
			fmt.Println("path:", path)
		}
	case 2: // FILE_DROP_CANCEL
		{
		}
	}
}

//export clipboardCopyImgCallback
func clipboardCopyImgCallback(cHeader C.IMAGE_HEADER, cData *C.uchar, cDataSize C.ulong) {
	imgHeader := rtkCommon.ImgHeader{
		Width:       int32(cHeader.width),
		Height:      int32(cHeader.height),
		Planes:      uint16(cHeader.planes),
		BitCount:    uint16(cHeader.bitCount),
		Compression: uint32(cHeader.compression),
	}
	data := C.GoBytes(unsafe.Pointer(cData), C.int(cDataSize))
	dataSize := uint32(cDataSize)
	rtkGlobal.Handler.CopyImgHeader = imgHeader
	rtkGlobal.Handler.CopyImgData = data
	rtkGlobal.Handler.CopyDataSize.SizeHigh = 0
	rtkGlobal.Handler.CopyDataSize.SizeLow = dataSize
	fmt.Printf("Clipboard image content, width[%d] height[%d] data size:[%d] \n", imgHeader.Width, imgHeader.Height, dataSize)
}

// export SetupDstPasteFile
func GoSetupDstPasteFile(desc string, fileName string, fileSizeHigh uint32, fileSizeLow uint32) {
	cDesc := stringToWChar(desc)
	cFileName := stringToWChar(fileName)
	cFileSizeHigh := C.ulong(fileSizeHigh)
	cFileSizeLow := C.ulong(fileSizeLow)
	C.SetupDstPasteFile(cDesc, cFileName, cFileSizeHigh, cFileSizeLow)
}

// export SetupFileDrop
func GoSetupFileDrop(desc string, fileName string, fileSizeHigh uint32, fileSizeLow uint32) {
	cDesc := stringToWChar(desc)
	cFileName := stringToWChar(fileName)
	cFileSizeHigh := C.ulong(fileSizeHigh)
	cFileSizeLow := C.ulong(fileSizeLow)
	C.SetupFileDrop(cDesc, cFileName, cFileSizeHigh, cFileSizeLow)
}

// export SetupDstPasteImage
func GoSetupDstPasteImage(desc string, content []byte, imgHeader rtkCommon.ImgHeader, dataSize uint32) {
	cDesc := stringToWChar(desc)
	cImgHeader := C.IMAGE_HEADER{
		width:       C.int(imgHeader.Width),
		height:      C.int(imgHeader.Height),
		planes:      C.ushort(imgHeader.Planes),
		bitCount:    C.ushort(imgHeader.BitCount),
		compression: C.ulong(imgHeader.Compression),
	}
	cDataSize := C.ulong(dataSize)
	C.SetupDstPasteImage(cDesc, cImgHeader, cDataSize)
}

// export DataTransfer
func GoDataTransfer(data []byte) {
	cData := (*C.uchar)(unsafe.Pointer(&data[0]))
	cSize := C.uint(len(data))
	C.DataTransfer(cData, cSize)
}

// export EventHandle
func GoEventHandle(eventType int) {
	C.EventHandle(C.EVENT_TYPE(eventType))
}

func SetupCallbackSettings() {
	C.SetClipboardCopyFileCallback((C.ClipboardCopyFileCallback)(unsafe.Pointer(C.clipboardCopyFileCallback)))
	C.SetClipboardPasteFileCallback((C.ClipboardPasteFileCallback)(unsafe.Pointer(C.clipboardPasteFileCallback)))
	C.SetFileDropCmdCallback((C.FileDropCmdCallback)(unsafe.Pointer(C.fileDropCmdCallback)))
	C.SetClipboardCopyImgCallback((C.ClipboardCopyImgCallback)(unsafe.Pointer(C.clipboardCopyImgCallback)))
	C.StartClipboardMonitor()
}

func GoSetupDstPasteText(content []byte) {
	log.Println("GoSetupDstPasteText :", string(content))
	clipboard.Write(clipboard.FmtText, content)
}

func ReceiveFileConfirm(fileSize int64) {

	log.Println("ReceiveFileConfirm:", fileSize)
	//CallbackInstance.CallbackMethodFileConfirm("android", fileSize)
}
func ReceiveCopyDataDone(fileType rtkCommon.ClipboardFmtType, fileSize int64) {
	log.Println("ReceiveFileDone FmtType:", fileType, " size: ", fileSize)
}

func FoundPeer() {
	//CallbackInstance.CallbackMethodFoundPeer()
	log.Println("CallbackMethodFoundPeer")
}

func GenKey() crypto.PrivKey {
	privKeyFile := ".priv.pem"
	return rtkUtils.GenKey(privKeyFile)
}

func IsHost() bool {
	return rtkUtils.FileExists(".HostID")
}

func GetHostID() string {
	file, err := os.Open(".HostID")
	if err != nil {
		log.Println(err)
		return rtkGlobal.HOST_ID
	}
	defer file.Close()

	data := make([]byte, 1024)
	_, err = file.Read(data)
	if err != nil {
		log.Println(err)
		return rtkGlobal.HOST_ID
	}

	return string(data)
}

func GetIDPath() string {
	return ".ID"
}

func GetHostIDPath() string {
	return ".HostID"
}
