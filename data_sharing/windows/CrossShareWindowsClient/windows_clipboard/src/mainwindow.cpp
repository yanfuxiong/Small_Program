#include "mainwindow.h"
#include "ui_mainwindow.h"

MainWindow::MainWindow(QWidget *parent)
    : QMainWindow(parent)
    , ui(new Ui::MainWindow)
    , m_testTimer(nullptr)
    , m_currentProgressVal(0)
{
    ui->setupUi(this);

    {
        connect(CommonSignals::getInstance(), &CommonSignals::dispatchMessage, this, &MainWindow::onDispatchMessage);
        connect(CommonSignals::getInstance(), &CommonSignals::logMessage, this, &MainWindow::onLogMessage);
        connect(CommonSignals::getInstance(), &CommonSignals::systemConfigChanged, this, &MainWindow::onSystemConfigChanged);
        //connect(CommonSignals::getInstance(), &CommonSignals::userAcceptFile, this, &MainWindow::onUserAcceptFile);
    }

    {
        auto deviceList = new DeviceInfoList(this);
        ui->deviceListWidget->addWidget(deviceList);
        ui->deviceListWidget->setCurrentWidget(deviceList);
    }

    QTimer::singleShot(0, this, [] {
        Q_EMIT CommonSignals::getInstance()->systemConfigChanged();
    });
}

MainWindow::~MainWindow()
{
    delete ui;
}

void MainWindow::onLogMessage(const QString &message)
{
    ui->log_browser->append(QString("[%1]: %2").arg(QDateTime::currentDateTime().toString("yyyy-MM-dd hh:mm:ss.zzz")).arg(message));
}

// 只用测试
void MainWindow::startTestTimer()
{
    if (m_testTimer) {
        if (m_testTimer->isActive()) {
            m_testTimer->stop();
            m_testTimer->deleteLater();
        }
    }

    m_testTimer = new QTimer(this);
    m_currentProgressVal = 0;
    m_testTimer->setInterval(500);
    connect(m_testTimer, &QTimer::timeout, this, [this] {
        m_currentProgressVal += 5;
        if (m_currentProgressVal > 100) {
            m_currentProgressVal = 100;
        }
        Q_EMIT CommonSignals::getInstance()->updateProgressInfo(m_currentProgressVal);
        if (m_currentProgressVal >= 100) {
            m_testTimer->stop();
        }
    });
    m_testTimer->start();
}

void MainWindow::onDispatchMessage(const QVariant &data)
{
    if (data.canConvert<GetConnStatusResponseMsgPtr>() == true) {
        GetConnStatusResponseMsgPtr ptr_msg = data.value<GetConnStatusResponseMsgPtr>();
        if (ptr_msg->statusCode == 1) {
            Q_EMIT CommonSignals::getInstance()->showInfoMessageBox("connectionStatus", "connected to server.");
        } else {
            Q_EMIT CommonSignals::getInstance()->showWarningMessageBox("connectionStatus", "The server is in a disconnected state.");
        }
        return;
    }

    if (data.canConvert<UpdateClientStatusMsgPtr>() == true) {
        UpdateClientStatusMsgPtr ptr_msg = data.value<UpdateClientStatusMsgPtr>();

        auto &clientVec = g_getGlobalData()->m_clientVec;
        bool exists = false;
        for (auto itr = clientVec.begin(); itr != clientVec.end(); ++itr) {
            if ((*itr)->clientID == ptr_msg->clientID) {
                exists = true;
                if ((*itr)->status == 0) {
                    clientVec.erase(itr);
                } else {
                    *itr = ptr_msg;
                }
                break;
            }
        }
        if (exists == false) {
            clientVec.push_back(ptr_msg);
        }
        Q_EMIT CommonSignals::getInstance()->updateClientList();
        return;
    }

    if (data.canConvert<SendFileRequestMsgPtr>() == true) {
        SendFileRequestMsgPtr ptr_msg = data.value<SendFileRequestMsgPtr>();
        AcceptFileDialog dialog;
        {
//            nlohmann::json infoJson;
//            infoJson["ip"] = ptr_msg->ip.toStdString();
//            infoJson["port"] = ptr_msg->port;
//            infoJson["fileSize"] = ptr_msg->fileSize;
//            infoJson["timeStamp"] = ptr_msg->timeStamp;
//            infoJson["fileName"] = ptr_msg->fileName.toStdString();
//            infoJson["clientID"] = ptr_msg->clientID.toStdString();

            QString infoText;
            for (const auto &deviceInfo : g_getGlobalData()->m_clientVec) {
                if (deviceInfo->clientID == ptr_msg->clientID) {
                    infoText += QString("From: %1\n").arg(deviceInfo->clientName);
                    break;
                }
            }
            infoText += QString("IP: %1\n").arg(ptr_msg->ip);
            infoText += QString("FileName: %1\n").arg(ptr_msg->fileName);
            infoText += QString("FileSize: %1\n").arg(ptr_msg->fileSize);

            dialog.setFileInfo(infoText);
        }
        dialog.setWindowTitle("是否接受檔案?");
        if (dialog.exec() == QDialog::Accepted) {
            QString dir = QFileDialog::getExistingDirectory(this, "Open Directory",
                                                            CommonUtils::desktopDirectoryPath(),
                                                            QFileDialog::ShowDirsOnly | QFileDialog::DontResolveSymlinks);
            if (dir.isEmpty()) {
                return;
            }

            QString newFileName = dir + "/" + CommonUtils::getFileNameByPath(ptr_msg->fileName);

            SendFileResponseMsg responseMessage;
            responseMessage.statusCode = 1; // 接受
            responseMessage.ip = ptr_msg->ip;
            responseMessage.port = ptr_msg->port;
            responseMessage.clientID = ptr_msg->clientID;
            responseMessage.fileSize = ptr_msg->fileSize;
            responseMessage.timeStamp = ptr_msg->timeStamp;
            responseMessage.fileName = newFileName;

            QByteArray data = SendFileResponseMsg::toByteArray(responseMessage);
            Q_EMIT CommonSignals::getInstance()->sendDataToServer(data);

            // 屏蔽用户的点击操作
            Q_EMIT CommonSignals::getInstance()->updateControlStatus(false);
        } else {
            SendFileResponseMsg responseMessage;
            responseMessage.statusCode = 0; // 拒绝
            responseMessage.ip = ptr_msg->ip;
            responseMessage.port = ptr_msg->port;
            responseMessage.clientID = ptr_msg->clientID;
            responseMessage.fileSize = ptr_msg->fileSize;
            responseMessage.timeStamp = ptr_msg->timeStamp;
            responseMessage.fileName = ptr_msg->fileName;

            QByteArray data = SendFileResponseMsg::toByteArray(responseMessage);
            Q_EMIT CommonSignals::getInstance()->sendDataToServer(data);
        }
        return;
    }

    if (data.canConvert<UpdateProgressMsgPtr>() == true) {
        UpdateProgressMsgPtr ptr_msg = data.value<UpdateProgressMsgPtr>();
        if (1)
        {
            nlohmann::json infoJson;
            infoJson["ip"] = ptr_msg->ip.toStdString();
            infoJson["port"] = ptr_msg->port;
            infoJson["clientID"] = ptr_msg->clientID.toStdString();
            infoJson["fileSize"] = ptr_msg->fileSize;
            infoJson["sentSize"] = ptr_msg->sentSize;
            infoJson["timeStamp"] = QDateTime::fromMSecsSinceEpoch(ptr_msg->timeStamp).toString("yyyy-MM-dd hh:mm:ss.zzz").toStdString();
            infoJson["fileName"] = ptr_msg->fileName.toStdString();
            Q_EMIT CommonSignals::getInstance()->logMessage(infoJson.dump(4).c_str());
        }

//        if (m_progressDialog == nullptr) {
//            m_progressDialog = new ProgressBarDialog;
//            m_progressDialog->setWindowTitle(QString("send file [%1]").arg(ptr_msg->fileName));
//            m_progressDialog->setModal(true);
//            m_progressDialog->show();
//        } else {
//            uint64_t totalSize = ptr_msg->fileSize;
//            uint64_t sentSize = ptr_msg->sentSize;
//            int progressVal = static_cast<int>((sentSize / double(totalSize)) * 100);
//            if (sentSize >= totalSize) {
//                progressVal = 100;
//            }
//            Q_EMIT CommonSignals::getInstance()->updateProgressInfo(progressVal);
//        }
        {
            uint64_t totalSize = ptr_msg->fileSize;
            uint64_t sentSize = ptr_msg->sentSize;
            int progressVal = static_cast<int>((sentSize / double(totalSize)) * 100);
            if (sentSize >= totalSize) {
                progressVal = 100;
            }
            Q_EMIT CommonSignals::getInstance()->updateProgressInfoWithID(progressVal, ptr_msg->clientID);
        }
        return;
    }
}

void MainWindow::onSystemConfigChanged()
{
    const auto &config = g_getGlobalData()->systemConfig;
    if (config.displayLogSwitch) {
        ui->bottom_box->show();
    } else {
        ui->bottom_box->hide();
    }
}

void MainWindow::on_settings_btn_clicked()
{
    auto &config = g_getGlobalData()->systemConfig;
    config.displayLogSwitch = !config.displayLogSwitch;
    Q_EMIT CommonSignals::getInstance()->systemConfigChanged();
}

void MainWindow::closeEvent(QCloseEvent *event)
{
    CommonUtils::killServer();
    QMainWindow::closeEvent(event);
}

//void MainWindow::on_conn_status_clicked()
//{
//    if (g_getGlobalData()->namedPipeConnected == false) {
//        Q_EMIT CommonSignals::getInstance()->showWarningMessageBox("connectionStatus", "The server is in a disconnected state.");
//        return;
//    }
//    GetConnStatusRequestMsg message;
//    QByteArray data = GetConnStatusRequestMsg::toByteArray(message);
//    Q_EMIT CommonSignals::getInstance()->sendDataToServer(data);
//}
