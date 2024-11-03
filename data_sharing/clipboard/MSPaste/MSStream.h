#ifndef __INCLUDED_MS_STREAM__
#define __INCLUDED_MS_STREAM__

#include "MSCommon.h"
#include "MSITransData.h"
#include <objidl.h>
#include <mutex>

class MSStream : public IStream, public MSITransData
{
public:
    MSStream(ClipboardPasteFileCallback& pasteCb, DWORD fileSizeHigh, DWORD fileSizeLow, IStreamObserver* pObserver);
    ~MSStream();
    // IUnknown
    STDMETHODIMP QueryInterface(REFIID riid, void **ppv);
    STDMETHODIMP_(ULONG) AddRef();
    STDMETHODIMP_(ULONG) Release();
    // IStream
    HRESULT STDMETHODCALLTYPE Read(void* pv, ULONG cb, ULONG* pcbRead);
    HRESULT STDMETHODCALLTYPE Write(const void* pv, ULONG cb, ULONG* pcbWritten);
    HRESULT STDMETHODCALLTYPE Seek(LARGE_INTEGER dlibMove, DWORD dwOrigin, ULARGE_INTEGER* plibNewPosition);
    HRESULT STDMETHODCALLTYPE SetSize(ULARGE_INTEGER libNewSize);
    HRESULT STDMETHODCALLTYPE CopyTo(IStream* pstm, ULARGE_INTEGER cb, ULARGE_INTEGER* pcbRead, ULARGE_INTEGER* pcbWritten);
    HRESULT STDMETHODCALLTYPE Commit(DWORD grfCommitFlags);
    HRESULT STDMETHODCALLTYPE Revert();
    HRESULT STDMETHODCALLTYPE LockRegion(ULARGE_INTEGER libOffset, ULARGE_INTEGER cb, DWORD dwLockType);
    HRESULT STDMETHODCALLTYPE UnlockRegion(ULARGE_INTEGER libOffset, ULARGE_INTEGER cb, DWORD dwLockType);
    HRESULT STDMETHODCALLTYPE Stat(STATSTG* pstatstg, DWORD grfStatFlag);
    HRESULT STDMETHODCALLTYPE Clone(IStream** ppstm);

    void StartDownload() override;
    void WriteFile(BYTE* data, unsigned int size) override;
    void Cancel() override;

private:
    HRESULT InitFile(DWORD fileSizeHigh, DWORD fileSizeLow);

    ULONG m_cRef;
    std::mutex mMutexRead;
    std::mutex mMutexWrite;
    HANDLE mData;
    BYTE* m_pBuffer;
    DWORD mFileSize;
    DWORD mProgressSize;
    ULONG mPosition;
    bool mWaitDataFlag;
    bool mCancel;

    ClipboardPasteFileCallback& mPasteCb;
    IStreamObserver* mpObserver;
};

#endif //__INCLUDED_MS_STREAM__