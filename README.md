#  wechatgpt 微信机器人
### 记录：
* 尝试基于openai的ChatGPT搭建一个微信机器人，
* 在github网站上找到了一个相关项目：https://github.com/riba2534/openai-on-wechat ，阅读学习代码后尝试运行，发现获取ChatGPT的api token需要付费，于是查找免费使用api的方法。
* 在 https://github.com/chatanywhere/GPT_API_free 这个项目可以获取免费使用ChatGPT3.5 api的token，只是这样的话聊天时会缺少生成图片的功能。

### 待解决的问题
- [x] 回答超过限长1420不能回复问题：当ChatGPT生成的回答过长时（大于1420个字符），会无法发送到微信窗口。通过截取回答，分条发送，解决了这个问题。
- [ ] 程序不能后台静默运行，当关闭命令行窗口时程序也会停止。
- [ ] ChatGPT不能连网问题
- [ ] 尝试将程序搭载在手机上
