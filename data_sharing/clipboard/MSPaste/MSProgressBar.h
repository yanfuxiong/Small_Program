#ifndef __INCLUDED_MS_PROGRESS_BAR__
#define __INCLUDED_MS_PROGRESS_BAR__

#include <windows.h>
#include <commctrl.h>
#include <stdint.h>
#include <thread>
#include <mutex>
#include <condition_variable>

LRESULT CALLBACK WndProc(HWND, UINT, WPARAM, LPARAM);

class MSProgressBar
{
public:
    MSProgressBar(DWORD filesizeHigh, DWORD filesizeLow);
    ~MSProgressBar();
    void WaitSetupReady();
    static void SetProgress(int progress);
    static void SetErrorMsg();
    static void SetTransTerm();

private:
    bool Init();
    std::thread m_Thread;
    std::mutex  m_Mutex;
    std::condition_variable m_Cv;
    bool        m_SetupReady;
};

#endif  //__INCLUDED_MS_PROGRESS_BAR__