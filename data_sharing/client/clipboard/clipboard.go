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
	rtkPlatform.SetResetCBCallback(ResetClipboard)
	rtkPlatform.WatchClipboardText(ctx, resultChan)
}

func ResetClipboard(resetTypeVal rtkCommon.ClipboardResetType) {
	resetImage := func() {
		rtkGlobal.Handler.CopyImgHeader = rtkCommon.ImgHeader{
			Width:       0,
			Height:      0,
			Planes:      0,
			BitCount:    0,
			Compression: 0,
		}
		rtkGlobal.Handler.CopyImgData = nil
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

	if (resetTypeVal & rtkCommon.CLIPBOARD_RESET_TYPE_TEXT) != 0 {
		resetImage()
		resetFile()
		resetDataSize()
	}
	if (resetTypeVal & rtkCommon.CLIPBOARD_RESET_TYPE_IMAGE) != 0 {
		resetImage()
	}
	if (resetTypeVal & rtkCommon.CLIPBOARD_RESET_TYPE_FILE) != 0 {
		resetFile()
	}
}

func WatchClipboardFiles(ctx context.Context, resultChan chan<- rtkCommon.ClipBoardData, ipAddr string) {
	var lastContent []byte
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(100 * time.Millisecond):
			filePath, ok := rtkGlobal.Handler.CopyFilePath.Load().(string)
			if ok && filePath == "" {
				continue
			}

			currentContent := []byte(filePath)
			if !rtkUtils.ContentEqual(lastContent, currentContent) {
				ResetClipboard(rtkCommon.CLIPBOARD_RESET_TYPE_IMAGE)

				lastContent = currentContent
				if rtkGlobal.Handler.IsFileDropMap[ipAddr] { // file drop  cannot  trigger clipboard change
					rtkGlobal.Handler.IsFileDropMap[ipAddr] = false
					continue
				}
				hash, err := rtkUtils.CreateMD5Hash(currentContent)
				if err != nil {
					log.Fatalf("Failed to create hash: %v", err)
				}

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

func WatchClipboardImg(ctx context.Context, resultChan chan<- rtkCommon.ClipBoardData) {
	var lastContent []byte

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(100 * time.Millisecond):
			currentHeader := rtkCommon.ImgHeader(rtkGlobal.Handler.CopyImgHeader)
			currentContent := []byte(rtkGlobal.Handler.CopyImgData)
			if !rtkUtils.ContentEqual(lastContent, currentContent) {
				hash, err := rtkUtils.CreateMD5Hash(currentContent)
				if err != nil {
					log.Fatalf("Failed to create hash: %v", err)
				}
				lastContent = currentContent

				if len(currentContent) > 0 {
					ResetClipboard(rtkCommon.CLIPBOARD_RESET_TYPE_FILE)

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
