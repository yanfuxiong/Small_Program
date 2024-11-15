package filedrop

import (
	rtkCommon "rtk-cross-share/common"
	rtkGlobal "rtk-cross-share/global"
)

func SendFileDropCmd(cmd rtkCommon.FileDropCmd, data interface{}) {
	switch cmd {
	case rtkCommon.FILE_DROP_REQUEST:
		rtkGlobal.Handler.CtxMutex.Lock()
		defer rtkGlobal.Handler.CtxMutex.Unlock()
		rtkGlobal.Handler.State.State = rtkCommon.FILE_DROP_INIT
		rtkGlobal.Handler.CopyDataSize.SizeHigh = data.(rtkCommon.FileInfo).FileSize_.SizeHigh
		rtkGlobal.Handler.CopyDataSize.SizeLow = data.(rtkCommon.FileInfo).FileSize_.SizeLow
		rtkGlobal.Handler.CopyFilePath.Store(data.(rtkCommon.FileInfo).FilePath)
	case rtkCommon.FILE_DROP_ACCEPT:
		rtkGlobal.Handler.CtxMutex.Lock()
		defer rtkGlobal.Handler.CtxMutex.Unlock()
		s, ok := data.(string)
		if ok {
			rtkGlobal.Handler.DstFilePath = s
		}
		rtkGlobal.Handler.State.State = rtkCommon.DEST_INIT
	case rtkCommon.FILE_DROP_CANCEL:
		// Do nothing
	}
}
