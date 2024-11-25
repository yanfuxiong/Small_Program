#include "device_info.h"
#include "ui_device_info.h"
#include "common_signals.h"
#include <QFileInfo>
#include <QMessageBox>
#include <QFileDialog>

DeviceInfo::DeviceInfo(QWidget *parent) :
    QWidget(parent),
    ui(new Ui::DeviceInfo)
{
    ui->setupUi(this);
    ui->progressBar->setValue(0);

    QHBoxLayout *pHBoxLayout = new QHBoxLayout;
    pHBoxLayout->setMargin(0);
    pHBoxLayout->setSpacing(10);
    ui->progressBar->setLayout(pHBoxLayout);

    {
        m_imageLabel = new QLabel;
        m_imageLabel->setObjectName("ImageLabel");
        m_imageLabel->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Expanding);
        m_imageLabel->setFixedSize(55, 55);
        m_imageLabel->setStyleSheet("border-image:url(:/resource/image.svg);");

        pHBoxLayout->addSpacing(5);
        pHBoxLayout->addWidget(m_imageLabel);
    }

    {
        m_deviceInfoLabel = new QLabel;
        m_deviceInfoLabel->setObjectName("DeviceInfo");
        m_deviceInfoLabel->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Expanding);
        m_deviceInfoLabel->setFixedWidth(140);

        m_deviceInfoLabel->setAlignment(Qt::AlignLeft | Qt::AlignVCenter);
        m_deviceInfoLabel->setText("设备1\n192.168.0.1");

        pHBoxLayout->addWidget(m_deviceInfoLabel);
    }

    pHBoxLayout->addStretch();

    {
        m_fileIcon = new QPushButton;
        m_fileIcon->setObjectName("openFile");
        m_fileIcon->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Expanding);
        m_fileIcon->setFixedSize(50, 50);
        pHBoxLayout->addWidget(m_fileIcon);
        connect(m_fileIcon, &QPushButton::clicked, this, &DeviceInfo::onOpenFile);
    }

    pHBoxLayout->addSpacing(5);

    // 信号连接
    {
        connect(CommonSignals::getInstance(), &CommonSignals::updateProgressInfoWithID, this, &DeviceInfo::onUpdateProgressInfoWithID);
        connect(CommonSignals::getInstance(), &CommonSignals::updateControlStatus, this, &DeviceInfo::onUpdateControlStatus);
    }
}

DeviceInfo::~DeviceInfo()
{
    delete ui;
}

void DeviceInfo::setDeviceInfo(const nlohmann::json &deviceInfoJson)
{
    Q_ASSERT(deviceInfoJson.is_object());
    Q_ASSERT(deviceInfoJson.contains("clientName"));
    Q_ASSERT(deviceInfoJson.contains("clientID"));
    m_deviceInfoJson = deviceInfoJson;

    {
        QString descInfo;
        descInfo += deviceInfoJson["clientName"].get<std::string>().c_str();
        descInfo += "\n";
        descInfo += deviceInfoJson["ip"].get<std::string>().c_str();
        m_deviceInfoLabel->setText(descInfo);
    }
}

QByteArray DeviceInfo::getClientID() const
{
    Q_ASSERT(m_deviceInfoJson.empty() == false);
    try {
        return m_deviceInfoJson.at("clientID").get<std::string>().c_str();
    } catch (const std::exception &e) {
        qWarning() << e.what();
        return {};
    }
}

void DeviceInfo::onOpenFile()
{
    UpdateClientStatusMsgPtr ptrClient = nullptr;
    for (const auto &data : g_getGlobalData()->m_clientVec) {
        if (data->clientID == getClientID()) {
            ptrClient = data;
            break;
        }
    }

    if (ptrClient == nullptr) {
        return;
    }

//    if (g_getGlobalData()->selectedFileName.isEmpty()) {
//        Q_EMIT CommonSignals::getInstance()->showWarningMessageBox("warning", "No file selected, unable to send......");
//        return;
//    }

    {
        g_getGlobalData()->selectedFileName.clear(); // 先清空

        QString fileName = QFileDialog::getOpenFileName(this, tr("Select File"),
                                                        CommonUtils::desktopDirectoryPath(),
                                                        ("Files (*.*)"));
        if (fileName.isEmpty()) {
            return;
        }
        qInfo() << "[FILE]:" << fileName;
        Q_EMIT CommonSignals::getInstance()->logMessage(QString("[SELECT]: %1").arg(fileName));
        g_getGlobalData()->selectedFileName = fileName; // 保存选中的文件名
    }

    {
        SendFileRequestMsg message;
        message.ip = ptrClient->ip;
        message.port = ptrClient->port;
        message.clientID = ptrClient->clientID;
        message.fileSize = static_cast<uint64_t>(QFileInfo(g_getGlobalData()->selectedFileName).size());
        message.timeStamp = QDateTime::currentDateTime().toUTC().toMSecsSinceEpoch();
        message.fileName = g_getGlobalData()->selectedFileName;

        QByteArray data = SendFileRequestMsg::toByteArray(message);
        Q_EMIT CommonSignals::getInstance()->logMessage(QString("开始发送: %1").arg(message.fileName));

        Q_EMIT CommonSignals::getInstance()->sendDataToServer(data);

        Q_EMIT CommonSignals::getInstance()->updateControlStatus(false);
    }
}

void DeviceInfo::onUpdateProgressInfoWithID(int currentVal, const QByteArray &clientID)
{
    if (clientID != getClientID()) {
        ui->progressBar->setValue(0);
        return;
    }
    ui->progressBar->setValue(currentVal);
    if (currentVal >= 100) {
        Q_EMIT CommonSignals::getInstance()->updateControlStatus(true);
    }
}

void DeviceInfo::onUpdateControlStatus(bool status)
{
    if (status) {
        m_fileIcon->setEnabled(true);
        m_fileIcon->show();
    } else {
        m_fileIcon->setEnabled(false);
    }
}

