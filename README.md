# 这是一个接入discord的小机器人，有定时提醒，投骰子，问天气, 接入google ai等功能，
# This is a small bot integrated with Discord, featuring functions such as scheduled reminders, dice rolling, and weather inquiries.


# 先进入discord developer Portal中心；
1、创建应用；能得到client_id;
2、在bot那块重置token, 得到token;
3、打开Message Content Intent权限；
4、邀请机器人加入；
https://discord.com/api/oauth2/authorize?client_id=[你的Client ID]&permissions=274877974592&scope=bot

# 在OpenWeathMap注册
1、地址为：https://openweathermap.org/api
2、填入相关信息：得到api-key 

# 将token 和 api-key填入.env文件中。

# 运行init.sh

# 编译
go build -o discord-bot main.go

# 将编译后的discord-bot 和 .env文件上传到正式服务器

# sudo vim /etc/systemd/system/kairos.service
# 将kairos.service文件中的内容复制

sudo systemctl daemon-reload
sudo systemctl enable kairos
sudo systemctl start kairos

# 申请google ai 
https://aistudio.google.com/api-keys
