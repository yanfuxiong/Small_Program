#include <QTest>
#include "common_utils.h"

class TestQString : public QObject
{
    Q_OBJECT

private Q_SLOTS:
    void initTestCase() {}
    void cleanupTestCase() {}

    void test_utf16()
    {
        {
            QString str(R"(D:\jack_huang\Downloads\新增資料夾\測試.mp4)"); // 这里的形式是UTF8
            qInfo() << CommonUtils::toUtf16LE(str).toHex().toUpper().constData();
        }

        // 从 utf16 LE 转 utf8
        {
            QByteArray data;
            {
                QString str(R"(D:\jack_huang\Downloads\新增資料夾\測試.mp4)"); // 这里的形式是UTF8
                data = CommonUtils::toUtf16LE(str).toHex().toUpper();
            }

            QByteArray filePath = CommonUtils::toUtf8(QByteArray::fromHex(data));
            qInfo() << filePath.constData();
        }
    }

    void test_file_path_parse()
    {
        QVERIFY(CommonUtils::getFileNameByPath("D:/aa/bb/name.txt") == "name.txt");
        QVERIFY(CommonUtils::getFileNameByPath("D:\\aa\\测试\\測試.mp4") == "測試.mp4");
    }
};

QTEST_GUILESS_MAIN(TestQString)

#include "test_qstring.moc"

