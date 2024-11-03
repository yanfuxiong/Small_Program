#ifndef __INCLUDED_MS_COMMON__
#define __INCLUDED_MS_COMMON__

#include <string>
#include "MSCommonExt.h"

#define DEBUG 1

#if DEBUG
#define DEBUG_LOG(format, ...) printf(format "\n", ##__VA_ARGS__)
#else
#define DEBUG_LOG(format, ...) printf("not debug\n")
#endif

class IStream;

enum PASTE_TYPE
{
    PASTE_TYPE_FILE,
    PASTE_TYPE_DIB,
    PASTE_TYPE_UNKNOWN = -1,
};

struct FILE_INFO
{
    std::wstring desc; // TODO: P2P ID, or other info
    std::wstring fileName;
    unsigned long fileSizeHigh;
    unsigned long fileSizeLow;
    IStream* pStream = nullptr;
};

struct IMAGE_INFO
{
    std::wstring desc; // TODO: P2P ID, or other info
    IMAGE_HEADER picHeader;
    unsigned long dataSize;
};

typedef void (*ClipboardPasteFileCallback)(char*);

class IStreamObserver {
public:
    virtual void OnStreamEOF() = 0;
};

#endif //__INCLUDED_MS_COMMON__