#include "mainwindow.h"
#include <QApplication>
#include <QProcess>
#include "common_utils.h"
#include "common_proxy_style.h"
#include "common_ui_process.h"
#include "process_message.h"

int main(int argc, char *argv[])
{
    qInstallMessageHandler(CommonUtils::commonMessageOutput);
    QApplication a(argc, argv);
    a.setStyleSheet(CommonUtils::getFileContent(":/resource/my.qss"));
    a.setStyle(new CustomProxyStyle);

    {
        g_globalRegister();
        CommonUiProcess::getInstance();
        ProcessMessage::getInstance();
    }

    MainWindow w;
    w.setWindowTitle("CrossShare Client");
    w.show();
    qInfo() << qApp->applicationFilePath().toUtf8().constData();

    for (const auto &path : g_getPipeServerExePathList()) {
        qInfo() << path;
    }

    // 启动 NamedPipe Server
    QTimer::singleShot(10, qApp, [] {
        bool exists = false;
        for (const auto &serverExePath : g_getPipeServerExePathList()) {
            if (CommonUtils::processIsRunning(serverExePath)) {
                qInfo() << "process is running:" << serverExePath;
                return;
            }
        }
        for (const auto &serverExePath : g_getPipeServerExePathList()) {
            if (QFile::exists(serverExePath)) {
                exists = true;
                QProcess process;
                process.startDetached(serverExePath);
                break;
            }
        }
        if (exists == false) {
            Q_EMIT CommonSignals::getInstance()->showWarningMessageBox("warning",
                                                                       QString("startup failed [%1]").arg(g_getPipeServerExePathList().front()));
        }
    });
    return a.exec();
}
