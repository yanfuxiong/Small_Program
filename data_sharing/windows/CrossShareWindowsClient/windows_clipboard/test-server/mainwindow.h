#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include <QMainWindow>
#include "server.h"
#include "common_utils.h"
#include "common_signals.h"

namespace Ui {
class MainWindow;
}

class MainWindow : public QMainWindow
{
    Q_OBJECT

public:
    explicit MainWindow(QWidget *parent = nullptr);
    ~MainWindow();

    // 这里只是模拟
    void sendProgressData();

private slots:
    void on_add_client_clicked();
    void onLogMessage(const QString &message);
    void onRecvClientData(const QByteArray &data);

    void on_send_file_clicked();

private:
    Ui::MainWindow *ui;
    NamedPipeServer m_server;
    Buffer m_buffer;
    QPointer<QTimer> m_timer;
    int m_cacheProgressVal; // 进度值, 只用于模拟
    std::unique_ptr<UpdateProgressMsg> m_cacheProgressMsgPtr;
};

#endif // MAINWINDOW_H
