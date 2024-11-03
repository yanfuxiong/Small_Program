#include "MSProgressBar.h"
#include "MSCommon.h"
#include "MSUtils.h"

#define WM_USER_UPDATE_PROGRESS (WM_USER + 1)
#define WM_USER_TRANS_ERR       (WM_USER + 2)
#define WM_USER_TRANS_TERM      (WM_USER + 3)
static HWND g_win = NULL;
static HWND g_progressBar = NULL;
static HWND g_msgLabel = NULL;
static HWND g_btn = NULL;

LRESULT CALLBACK WndProgressProc(HWND hwnd, UINT message, WPARAM wParam, LPARAM lParam) {
    switch (message) {
    case WM_CREATE:
        //DEBUG_LOG("[%s %d]", __func__, __LINE__);
        g_progressBar = CreateWindowEx(0, PROGRESS_CLASS, NULL, WS_CHILD | WS_VISIBLE | PBS_SMOOTH,
            10, 30, 260, 20, hwnd, NULL,
            ((LPCREATESTRUCT)lParam)->hInstance, NULL);
        if (g_progressBar == NULL) {
            break;
        }

        SendMessage(g_progressBar, PBM_SETRANGE, 0, MAKELPARAM(0, 100));
        SendMessage(g_progressBar, PBM_SETPOS, 0, 0);

        g_msgLabel = CreateWindowEx(0, "Warning", "", WS_CHILD | SS_CENTER,
            10, 30, 260, 20, hwnd, NULL,
            ((LPCREATESTRUCT)lParam)->hInstance, NULL);
        if (g_msgLabel == NULL) {
            break;
        }
        break;
    case WM_USER_UPDATE_PROGRESS:
        //DEBUG_LOG("[%s %d]", __func__, __LINE__);
        SendMessage(g_progressBar, PBM_SETPOS, wParam, 0);
        if (wParam >= 100)
        {
            SendMessage(hwnd, WM_CLOSE, 0, 0);
        }
        break;
    case WM_USER_TRANS_ERR:
    {
        ShowWindow(g_progressBar, SW_HIDE);
        ShowWindow(g_win, SW_HIDE);
        int result = MessageBox(hwnd, "File cannot be copied.\nPlease check if the file is currently in use.", "Warning",
            MB_OK | MB_ICONERROR);
        if (result == IDOK) {
            SendMessage(hwnd, WM_CLOSE, 0, 0);
        }
    }
        break;
    case WM_USER_TRANS_TERM:
    {
        // TODO: create cancel button and re-send the action to cgo
        // g_btn = CreateWindowEx(0, "BUTTON", "Cancel", WS_CHILD | WS_VISIBLE,
        //     10, 60, 80, 30, hwnd, (HMENU)1,
        //     ((LPCREATESTRUCT)lParam)->hInstance, NULL);
        ShowWindow(g_progressBar, SW_HIDE);
        ShowWindow(g_win, SW_HIDE);
        int result = MessageBox(hwnd, "Transmission failed: Network error\nPlease check your network connection and retry.", "Warning",
            MB_OK | MB_ICONERROR);
        if (result == IDOK) {
            SendMessage(hwnd, WM_CLOSE, 0, 0);
        }
    }
        break;
    case WM_DESTROY:
        PostQuitMessage(0);
        break;
    default:
        return DefWindowProc(hwnd, message, wParam, lParam);
    }
    return 0;
}

MSProgressBar::MSProgressBar(DWORD filesizeHigh, DWORD filesizeLow)
    : m_SetupReady(false)
{
    m_Thread = std::thread(&MSProgressBar::Init, this);
}

MSProgressBar::~MSProgressBar()
{
    m_Thread.join();
}

bool MSProgressBar::Init()
{
    HINSTANCE hInstance = GetModuleHandle(NULL);
    WNDCLASSEX wc = { sizeof(WNDCLASSEX), CS_HREDRAW | CS_VREDRAW, WndProgressProc, 0L, 0L,
                     hInstance, LoadIcon(NULL, IDI_APPLICATION), LoadCursor(NULL, IDC_ARROW),
                     (HBRUSH)(COLOR_WINDOW + 1), NULL, TEXT("ProgressWindow"), NULL };
    RegisterClassEx(&wc);

    int windowWidth = 350;
    int windowHeight = 100;

    int screenWidth = GetSystemMetrics(SM_CXSCREEN);
    int screenHeight = GetSystemMetrics(SM_CYSCREEN);

    int xPos = (screenWidth - windowWidth) / 2;
    int yPos = (screenHeight - windowHeight) / 2;

    g_win = CreateWindow(TEXT("ProgressWindow"), TEXT("Copy..."),
        WS_OVERLAPPED | WS_CAPTION, xPos, yPos,
        windowWidth, windowHeight, NULL, NULL, hInstance, NULL);

    if (g_win == NULL) return false;

    ShowWindow(g_win, SW_SHOW);
    SetForegroundWindow(g_win);
    SetWindowPos(g_win, HWND_TOPMOST, 0, 0, 0, 0, SWP_NOMOVE | SWP_NOSIZE);
    UpdateWindow(g_win);

    {
        std::lock_guard<std::mutex> lock(m_Mutex);
        m_SetupReady = true;
    }
    m_Cv.notify_one();

    MSG msg = { 0 };
    while (GetMessage(&msg, NULL, 0, 0)) {
        TranslateMessage(&msg);
        DispatchMessage(&msg);
    }

    return true;
}

void MSProgressBar::WaitSetupReady()
{
    std::unique_lock<std::mutex> lock(m_Mutex);
    m_Cv.wait(lock, [this] { return m_SetupReady; });
}

void MSProgressBar::SetProgress(int progress)
{
    if (g_win) {
        PostMessage(g_win, WM_USER_UPDATE_PROGRESS, progress, 0);
    }
    //SendMessage(hwndEdit, WM_PAINT, 0, (LPARAM)RGB(0, 255, 0));
}

void MSProgressBar::SetErrorMsg()
{
    if (g_win) {
        PostMessage(g_win, WM_USER_TRANS_ERR, 0, 0);
    }
}

void MSProgressBar::SetTransTerm()
{
    if (g_win) {
        PostMessage(g_win, WM_USER_TRANS_TERM, 0, 0);
    }
}
