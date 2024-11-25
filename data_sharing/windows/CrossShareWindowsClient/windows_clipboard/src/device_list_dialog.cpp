#include "device_list_dialog.h"
#include "ui_device_list_dialog.h"
#include "device_info.h"
#include "common_signals.h"
#include <QTimer>

DeviceListDialog::DeviceListDialog(QWidget *parent) :
    QDialog(parent),
    ui(new Ui::DeviceListDialog)
{
    ui->setupUi(this);
    {
        connect(CommonSignals::getInstance(), &CommonSignals::sendDataToServer, this, &DeviceListDialog::onSendDataToServer);
        connect(CommonSignals::getInstance(), &CommonSignals::updateClientList, this, &DeviceListDialog::onUpdateClientList);
        connect(CommonSignals::getInstance(), &CommonSignals::pipeDisconnected, this, [this] {
            QTimer::singleShot(50, this, [this] {
               onUpdateClientList();
            });
        });
    }

    {
        //m_flowLayout = new FlowLayout(this, 5, 8, 8);
        //ui->content->setLayout(m_flowLayout);
    }

    QTimer::singleShot(0, this, [this] {
        onUpdateClientList();
    });
}

DeviceListDialog::~DeviceListDialog()
{
    delete ui;
}

void DeviceListDialog::onSendDataToServer(const QByteArray &data)
{
    Q_UNUSED(data)
    // 此时用户已经点击了发送, dialog可以关闭了
    accept();
}

void DeviceListDialog::onUpdateClientList()
{
    while (ui->content->count() > 0) {
        auto widget = ui->content->widget(0);
        ui->content->removeWidget(widget);
        widget->deleteLater();
    }

    auto backgoundWidget = new QWidget;
    //auto flowLayout = new FlowLayout(5, 8, 8);
    auto flowLayout = new QVBoxLayout;
    flowLayout->setSpacing(8);
    flowLayout->setMargin(5);
    backgoundWidget->setLayout(flowLayout);
    ui->content->addWidget(backgoundWidget);

    const auto &clientVec = g_getGlobalData()->m_clientVec;
    for (const auto &data : clientVec) {
        auto widget = new DeviceInfo;
        nlohmann::json deviceInfoJson;
        deviceInfoJson["clientName"] = data->clientName.toStdString();
        deviceInfoJson["clientID"] = data->clientID.toStdString();
        flowLayout->addWidget(widget);
    }
    flowLayout->addStretch();
    ui->content->setCurrentIndex(0);
}
