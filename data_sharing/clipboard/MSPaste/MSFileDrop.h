#ifndef __INCLUDED_MS_FILE_DROP_IMPL__
#define __INCLUDED_MS_FILE_DROP_IMPL__

#include "MSCommon.h"
#include <windef.h>
#include <vector>
#include <mutex>

class MSProgressBar;
class MSFileDrop
{
public:
    MSFileDrop(FileDropCmdCallback& cmdCb);
    ~MSFileDrop();
    bool SetupDropFilePath(std::vector<FILE_INFO>& fileList);
    bool UpdateProgressBar(unsigned long progress);
    void DeinitProgressBar();

private:
    void SetupDialog();
    void PopSelectPathDialog(std::wstring &userSelectedPath);

    std::vector<FILE_INFO> mFileList;
    FileDropCmdCallback& mCmdCb;
    MSProgressBar* mCurProgressBar;
    unsigned int mCurProgress;
};

#endif //__INCLUDED_MS_FILE_DROP_IMPL__