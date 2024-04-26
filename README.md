# SerialTerminalForWindowsTerminal
在开始这个项目之前，我发现Windows Terminal对串口设备的支持并不理想。

我试用了一段时间[Zhou-zhi-peng的SerialPortForWindowsTerminal](https://github.com/Zhou-zhi-peng/SerialPortForWindowsTerminal/)项目。

然而，这个项目存在着编码转换的问题，导致数据显示乱码，并且作者目前并没有进行后续支持。因此，我决定创建了这个项目。

## 功能进展
* [x] Hex接收发送(大写hex与原文同显)
* [x] 双向编码转换
* [x] 活动端口探测
* [x] 数据日志保存
* [x] Hex断帧设置
* [x] UDP数据转发
* [x] TCP数据转发
* [ ] 文件接收发送

## 运行示例

1. 参数帮助 `./COM`

    ![img1.png](image/img1.png)

2. 输入设备输出UTF8 终端输出GBK `./COM -p COM8 -b 115200 -o GBK`

    ![img2.png](image/img2.png)
3. 彩色终端输出

   ![img3.png](image/img3.png)

4. Hex接收 `./COM -p COM8 -b 115200 -i hex`
   
   ![img4.png](image/img4.png)
5. Hex发送 `./COM -p COM8 -b 115200`

   ![img5.png](image/img5.png)
