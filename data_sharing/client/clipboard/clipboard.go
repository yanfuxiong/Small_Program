package clipboard

import (
	"context"
	"log"
	rtkCommon "rtk-cross-share/common"
	rtkGlobal "rtk-cross-share/global"
	rtkPlatform "rtk-cross-share/platform"
	rtkUtils "rtk-cross-share/utils"
	"time"
)

func WatchClipboardText(ctx context.Context, resultChan chan<- rtkCommon.ClipBoardData) {
	rtkPlatform.SetResetCBCallback(ResetOthersClipboard)
	rtkPlatform.WatchClipboardText(ctx, resultChan)
}

func ResetOthersClipboard(fmtType rtkCommon.ClipboardFmtType) {
	resetImage := func() {
		rtkGlobal.Handler.CopyImgHeader = rtkCommon.ImgHeader{
			Width:       0,
			Height:      0,
			Planes:      0,
			BitCount:    0,
			Compression: 0,
		}
	}
	resetFile := func() {
		rtkGlobal.Handler.CopyFilePath.Store("")
	}
	resetDataSize := func() {
		rtkGlobal.Handler.CopyDataSize = rtkCommon.FileSize{
			SizeHigh: 0,
			SizeLow:  0,
		}
	}
	switch fmtType {
	case rtkCommon.TEXT:
		resetImage()
		resetFile()
		resetDataSize()
	case rtkCommon.IMAGE:
		resetFile()
	case rtkCommon.FILE:
		resetImage()
	}
}

func WatchClipboardFiles(resultChan chan<- rtkCommon.ClipBoardData) {
	var lastContent []byte
	for {
		select {
		case <-time.After(100 * time.Millisecond):
			filePath := rtkGlobal.Handler.CopyFilePath.Load().(string)
			currentContent := []byte(filePath)
			if !rtkUtils.ContentEqual(lastContent, currentContent) {
				// TODO: check if currentContent empty
				ResetOthersClipboard(rtkCommon.FILE)
				hash, err := rtkUtils.CreateMD5Hash(currentContent)
				if err != nil {
					log.Fatalf("Failed to create hash: %v", err)
				}

				lastContent = currentContent
				if filePath != "" {
					resultChan <- rtkCommon.ClipBoardData{
						SourceID: rtkGlobal.NodeInfo.ID,
						FmtType:  rtkCommon.FILE,
						Hash:     hash.B58String(),
						Content:  currentContent,
					}
				}
			}
		}
	}
}

func imgHeaderEqual(a, b rtkCommon.ImgHeader) bool {
	if (a.Width != b.Width) || (a.Height != b.Height) || (a.Planes != b.Planes) ||
		(a.BitCount != b.BitCount) || (a.Compression != b.Compression) {
		return false
	}
	return true
}

func WatchClipboardImg(resultChan chan<- rtkCommon.ClipBoardData) {

	var lastHeader rtkCommon.ImgHeader
	var lastContent []byte

	isEmptyImgHeader := func(header rtkCommon.ImgHeader) bool {
		return header.Width == 0 && header.Height == 0 && header.Planes == 0 &&
			header.BitCount == 0 && header.Compression == 0
	}

	for {
		select {
		case <-time.After(100 * time.Millisecond):
			currentHeader := rtkCommon.ImgHeader(rtkGlobal.Handler.CopyImgHeader)
			currentContent := []byte(rtkGlobal.Handler.CopyImgData)
			if !imgHeaderEqual(lastHeader, currentHeader) && !rtkUtils.ContentEqual(lastContent, currentContent) {
				if !imgHeaderEqual(currentHeader, rtkCommon.ImgHeader{}) {
					ResetOthersClipboard(rtkCommon.IMAGE)
				}
				hash, err := rtkUtils.CreateMD5Hash(currentContent)
				if err != nil {
					log.Fatalf("Failed to create hash: %v", err)
				}
				lastHeader = currentHeader
				lastContent = currentContent

				if !isEmptyImgHeader(currentHeader) {
					log.Printf("WatchClipboardImg - got new Image  Wight:%d Height:%d, content len:[%d] \n\n", currentHeader.Width, currentHeader.Height, len(currentContent))
					resultChan <- rtkCommon.ClipBoardData{
						SourceID: rtkGlobal.NodeInfo.ID,
						FmtType:  rtkCommon.IMAGE,
						Hash:     hash.B58String(),
						Content:  currentContent,
					}
				}
			}
		}
	}
}
