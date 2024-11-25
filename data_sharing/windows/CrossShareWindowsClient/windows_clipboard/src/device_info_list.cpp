#include "device_info_list.h"
#include "ui_device_info_list.h"

DeviceInfoList::DeviceInfoList(QWidget *parent) :
    QWidget(parent),
    ui(new Ui::DeviceInfoList)
{
    ui->setupUi(this);

    {
        connect(CommonSignals::getInstance(), &CommonSignals::updateClientList, this, &DeviceInfoList::onUpdateClientList);
        connect(CommonSignals::getInstance(), &CommonSignals::pipeDisconnected, this, [this] {
            QTimer::singleShot(50, this, [this] {
                onUpdateClientList();
            });
        });
    }

    QTimer::singleShot(0, this, [this] {
        onUpdateClientList();
    });
}

DeviceInfoList::~DeviceInfoList()
{
    delete ui;
}

void DeviceInfoList::onUpdateClientList()
{
    while (ui->content->count() > 0) {
        auto widget = ui->content->widget(0);
        ui->content->removeWidget(widget);
        widget->deleteLater();
    }

    const auto &clientVec = g_getGlobalData()->m_clientVec;
    if (clientVec.empty()) {
        auto imageWidget = new QLabel(this);
        imageWidget->setStyleSheet("border-image:url(:/resource/background.jpg);");
        ui->content->addWidget(imageWidget);
        ui->content->setCurrentIndex(0);
        return;
    }

    auto area = new QScrollArea;
    {
        area->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Expanding);
        area->setHorizontalScrollBarPolicy(Qt::ScrollBarPolicy::ScrollBarAlwaysOff);
        area->setVerticalScrollBarPolicy(Qt::ScrollBarPolicy::ScrollBarAsNeeded);
        area->setWidgetResizable(true);
    }

    auto backgoundWidget = new QWidget;
    auto flowLayout = new QVBoxLayout;
    flowLayout->setSpacing(8);
    flowLayout->setMargin(5);
    backgoundWidget->setLayout(flowLayout);
    area->setWidget(backgoundWidget); // 放入滚动区域
    ui->content->addWidget(area);

    for (const auto &data : clientVec) {
        auto widget = new DeviceInfo;
        nlohmann::json deviceInfoJson;
        deviceInfoJson["clientName"] = data->clientName.toStdString();
        deviceInfoJson["clientID"] = data->clientID.toStdString();
        deviceInfoJson["ip"] = data->ip.toStdString();
        widget->setDeviceInfo(deviceInfoJson);
        flowLayout->addWidget(widget);
    }
    flowLayout->addStretch();
    ui->content->setCurrentIndex(0);
}
