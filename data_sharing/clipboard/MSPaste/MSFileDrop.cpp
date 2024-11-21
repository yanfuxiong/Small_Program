#include "MSFileDrop.h"
#include "MSProgressBar.h"
#include "MSUtils.h"
#include <iostream>
#include <shobjidl.h>
#include <windows.h>

MSFileDrop::MSFileDrop(FileDropCmdCallback& cmdCb):
    mFileList({}),
    mCmdCb(cmdCb),
    mCurProgressBar(NULL),
    mCurProgress(0)
{

}

MSFileDrop::~MSFileDrop()
{

}

bool MSFileDrop::SetupDropFilePath(std::vector<FILE_INFO>& fileList)
{
    mFileList = fileList;
    SetupDialog();
    return true;
}

bool MSFileDrop::UpdateProgressBar(unsigned long progress)
{
    // FIXME
    int64_t fileSize = (int64_t)mFileList[0].fileSizeHigh << 32 | (int64_t)mFileList[0].fileSizeLow;
    mCurProgress += progress;
    MSUtils::PrintProgress(mCurProgress, (uint32_t)fileSize);
    return true;
}

void MSFileDrop::DeinitProgressBar()
{
    if (mCurProgressBar) {
        delete mCurProgressBar;
        mCurProgressBar = NULL;
    }
    mCurProgress = 0;
}

void MSFileDrop::SetupDialog()
{
    int response = MessageBox(
        NULL,
        TEXT("Do you accept a remote file drop?"),
        TEXT(""),
        MB_OKCANCEL
    );

    if (response == IDOK) {
        std::wstring userSelectedPath;
        PopSelectPathDialog(userSelectedPath);
        int length = userSelectedPath.size() + mFileList[0].fileName.size();
        wchar_t* wideCPath = new wchar_t[length]();
        std::wcscat(wideCPath, userSelectedPath.c_str());
        const wchar_t* back_slash = L"\\";
        std::wcscat(wideCPath, back_slash);
        std::wcscat(wideCPath, mFileList[0].fileName.c_str());
        mCmdCb((unsigned long)FILE_DROP_ACCEPT, wideCPath);
        // FIXME
        if (!mCurProgressBar) {
            mCurProgressBar = new MSProgressBar(mFileList[0].fileSizeHigh, mFileList[0].fileSizeLow);
        }
        delete[] wideCPath;
    } else if (response == IDCANCEL) {
        mCmdCb((unsigned long)FILE_DROP_CANCEL, NULL);
    }
}

void MSFileDrop::PopSelectPathDialog(std::wstring &userSelectedPath)
{
    HRESULT hr = CoInitialize(NULL);
    if (SUCCEEDED(hr)) {
        IFileDialog *pFileDialog;
        hr = CoCreateInstance(CLSID_FileOpenDialog, NULL, CLSCTX_INPROC_SERVER, IID_PPV_ARGS(&pFileDialog));

        if (SUCCEEDED(hr)) {
            DWORD options;
            pFileDialog->GetOptions(&options);
            pFileDialog->SetOptions(options | FOS_PICKFOLDERS);

            hr = pFileDialog->Show(NULL);
            if (SUCCEEDED(hr)) {
                IShellItem *pItem;
                hr = pFileDialog->GetResult(&pItem);
                if (SUCCEEDED(hr)) {
                    PWSTR pszFilePath;
                    hr = pItem->GetDisplayName(SIGDN_FILESYSPATH, &pszFilePath);
                    if (SUCCEEDED(hr)) {
                        userSelectedPath = pszFilePath;
                        // std::wcout << L"Selected folder: " << pszFilePath << std::endl;
                        CoTaskMemFree(pszFilePath);
                    }
                    pItem->Release();
                }
            }
            pFileDialog->Release();
        }
        CoUninitialize();
    }
}