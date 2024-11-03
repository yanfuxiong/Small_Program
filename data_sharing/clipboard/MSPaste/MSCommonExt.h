#ifndef __INCLUDED_MS_COMMON_EXT__
#define __INCLUDED_MS_COMMON_EXT__

#ifdef __cplusplus
extern "C" {
#endif

typedef enum EVENT_TYPE
{
    EVENT_TYPE_OPEN_FILE_ERR,
    EVENT_TYPE_RECV_TIMEOUT
} EVENT_TYPE;

typedef struct IMAGE_HEADER
{
    int width;
    int height;
    unsigned short planes;
    unsigned short bitCount;
    unsigned long compression;
} IMAGE_HEADER;

#ifdef __cplusplus
}
#endif

#endif //__INCLUDED_MS_COMMON_EXT__