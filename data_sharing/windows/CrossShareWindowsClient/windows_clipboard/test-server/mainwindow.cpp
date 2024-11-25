#include "mainwindow.h"
#include "ui_mainwindow.h"
#include "common_signals.h"
#include <QTimer>
#include <QEventLoop>

MainWindow::MainWindow(QWidget *parent) :
    QMainWindow(parent),
    ui(new Ui::MainWindow)
{
    ui->setupUi(this);

    {
        connect(&m_server, &NamedPipeServer::recvData, this, &MainWindow::onRecvClientData);
        connect(CommonSignals::getInstance(), &CommonSignals::logMessage, this, &MainWindow::onLogMessage);
    }

    QTimer::singleShot(0, this, [this] {
        m_server.startServer("CrossSharePipe");
        Q_EMIT CommonSignals::getInstance()->logMessage("----------------start server");
    });
}

MainWindow::~MainWindow()
{
    delete ui;
}

void MainWindow::on_add_client_clicked()
{
    Q_EMIT CommonSignals::getInstance()->addTestClient();
}

void MainWindow::onLogMessage(const QString &message)
{
    ui->log_browser->append(QString("[%1]: %2").arg(QDateTime::currentDateTime().toString("yyyy-MM-dd hh:mm:ss.zzz")).arg(message));
}

void MainWindow::onRecvClientData(const QByteArray &data)
{
    Q_EMIT CommonSignals::getInstance()->logMessage(data.toHex().toUpper().constData());
    m_buffer.append(data);

    uint8_t typeValue = 0;
    uint8_t code = 0;
    // 解析消息
    while (g_getCodeFromByteArray(QByteArray(m_buffer.peek(), m_buffer.readableBytes()), typeValue, code)) {
        switch (code) {
        case 1: {
            GetConnStatusRequestMsg message;
            if (GetConnStatusRequestMsg::fromByteArray(QByteArray(m_buffer.peek(), m_buffer.readableBytes()), message)) {
                m_buffer.retrieve(message.getMessageLength());
                {
                    nlohmann::json infoJson;
                    infoJson["message"] = "GetConnStatusRequestMsg";
                    infoJson["type"] = message.headerInfo.type;
                    infoJson["code"] = message.headerInfo.code;
                    infoJson["contentLength"] = message.headerInfo.contentLength;
                    Q_EMIT CommonSignals::getInstance()->logMessage(infoJson.dump(4).c_str());
                }

                GetConnStatusResponseMsg responseMessage;
                responseMessage.statusCode = 1;

                Q_EMIT CommonSignals::getInstance()->sendDataForTestServer(GetConnStatusResponseMsg::toByteArray(responseMessage));
            }
            break;
        }
        case 2: {
            Q_ASSERT(false);
            break;
        }
        case 3: {
            Q_ASSERT(false);
            break;
        }
        case 4: {
            if (typeValue == PipeMessageType::Response) {
                SendFileResponseMsg message;
                if (SendFileResponseMsg::fromByteArray(QByteArray(m_buffer.peek(), m_buffer.readableBytes()), message)) {
                    m_buffer.retrieve(message.getMessageLength());

                    {
                        nlohmann::json infoJson;
                        infoJson["message"] = "SendFileResponseMsg";
                        infoJson["fileName"] = message.fileName.toStdString();
                        infoJson["clientID"] = message.clientID.toStdString();
                        infoJson["ip"] = message.ip.toStdString();
                        infoJson["port"] = message.port;
                        infoJson["statusCode"] = message.statusCode;
                        infoJson["desc"] = (message.statusCode == 0 ? "拒绝" : "接受");
                        Q_EMIT CommonSignals::getInstance()->logMessage(infoJson.dump(4).c_str());
                    }

                    if (message.statusCode == 0) {
                        break; // 拒绝了就不需要处理进度信息了, 直接 break
                    }

                    {
                        auto pProgressMsg = new UpdateProgressMsg;
                        pProgressMsg->ip = message.ip;
                        pProgressMsg->port = message.port;
                        pProgressMsg->clientID = message.clientID;
                        pProgressMsg->fileSize = 0; // 先初始化为0
                        pProgressMsg->timeStamp = message.timeStamp;
                        pProgressMsg->fileName = message.fileName;

                        m_cacheProgressMsgPtr.reset(pProgressMsg);
                    }
                }
            } else {
                Q_ASSERT(typeValue == PipeMessageType::Request);
                SendFileRequestMsg message;
                if (SendFileRequestMsg::fromByteArray(QByteArray(m_buffer.peek(), m_buffer.readableBytes()), message)) {
                    m_buffer.retrieve(message.getMessageLength());
                    {
                        nlohmann::json infoJson;
                        infoJson["message"] = "SendFileRequestMsg";
                        infoJson["fileName"] = message.fileName.toStdString();
                        infoJson["clientID"] = message.clientID.toStdString();
                        infoJson["ip"] = message.ip.toStdString();
                        infoJson["port"] = message.port;
                        Q_EMIT CommonSignals::getInstance()->logMessage(infoJson.dump(4).c_str());
                    }

                    {
                        auto pProgressMsg = new UpdateProgressMsg;
                        pProgressMsg->ip = message.ip;
                        pProgressMsg->port = message.port;
                        pProgressMsg->clientID = message.clientID;
                        pProgressMsg->fileSize = 0; // 先初始化为0
                        pProgressMsg->timeStamp = message.timeStamp;
                        pProgressMsg->fileName = message.fileName;

                        m_cacheProgressMsgPtr.reset(pProgressMsg);
                    }
                }
            }
            if (m_timer) {
                m_timer->stop();
                m_timer->deleteLater();
            }
            m_timer = new QTimer;
            m_timer->setTimerType(Qt::TimerType::PreciseTimer);
            m_timer->setInterval(60);
            m_cacheProgressVal = 0;
            connect(m_timer, &QTimer::timeout, this, &MainWindow::sendProgressData);
            m_timer->start();
            break;
        }
        case 5: {
            Q_ASSERT(false);
            break;
        }
        default: {
            Q_ASSERT(false);
            break;
        }
        }
    }
}

void MainWindow::sendProgressData()
{
    if (m_cacheProgressMsgPtr == nullptr) {
        return;
    }

    if (m_cacheProgressVal >= 100) {
        m_timer->stop();
        m_timer->deleteLater();
        m_timer = nullptr;
        m_cacheProgressMsgPtr.reset(nullptr);
        return;
    }

    UpdateProgressMsg message;
    message.ip = m_cacheProgressMsgPtr->ip;
    message.port = m_cacheProgressMsgPtr->port;
    message.clientID = m_cacheProgressMsgPtr->clientID;
    message.fileSize = QFileInfo(m_cacheProgressMsgPtr->fileName).size();
    auto delta = message.fileSize / 100.0;
    if (delta < 1.0) {
        delta = 1.0;
    }
    uint64_t currentSentFileSize = static_cast<uint64_t>((++m_cacheProgressVal) * delta);
    if (m_cacheProgressVal >= 100) {
        currentSentFileSize = message.fileSize;
    }
    message.sentSize = currentSentFileSize;
    message.timeStamp = QDateTime::currentDateTime().toUTC().toMSecsSinceEpoch();
    message.fileName = m_cacheProgressMsgPtr->fileName;

    Q_EMIT CommonSignals::getInstance()->sendDataForTestServer(UpdateProgressMsg::toByteArray(message));
}


void MainWindow::on_send_file_clicked()
{
    SendFileRequestMsg msg;
    msg.ip = "192.168.30.1";
    msg.port = 12345;
    msg.clientID = "QmQ7obXFx1XMFr6hCYXtovn9zREFqSXEtH5hdtpBDLjrAz";
    //msg.fileSize = static_cast<uint64_t>(QFileInfo(__FILE__).size());
    msg.fileSize = 60727169;
    msg.timeStamp = QDateTime::currentDateTime().toUTC().toMSecsSinceEpoch();
    msg.fileName = R"(C:\Users\TS\Desktop\test_data.mp4)";

    {
        nlohmann::json infoJson;
        infoJson["message"] = "SendFileRequestMsg";
        infoJson["ip"] = msg.ip.toStdString();
        infoJson["port"] = msg.port;
        infoJson["clientID"] = msg.clientID.toStdString();
        infoJson["fileSize"] = msg.fileSize;
        infoJson["timeStamp"] = msg.timeStamp;
        infoJson["fileName"] = msg.fileName.toUtf8().constData();
        Q_EMIT CommonSignals::getInstance()->logMessage(infoJson.dump(4).c_str());
    }

    Q_EMIT CommonSignals::getInstance()->sendDataForTestServer(SendFileRequestMsg::toByteArray(msg));
}
