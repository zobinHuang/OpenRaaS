FROM ubuntu:20.04

RUN mkdir -p /winecontainer
WORKDIR /winecontainer

# 写入国内源
RUN echo "\
deb http://mirrors.aliyun.com/ubuntu/ focal main restricted\n\
deb http://mirrors.aliyun.com/ubuntu/ focal-updates main restricted\n\
deb http://mirrors.aliyun.com/ubuntu/ focal universe\n\
deb http://mirrors.aliyun.com/ubuntu/ focal-updates universe\n\
deb http://mirrors.aliyun.com/ubuntu/ focal multiverse\n\
deb http://mirrors.aliyun.com/ubuntu/ focal-updates multiverse\n\
deb http://mirrors.aliyun.com/ubuntu/ focal-backports main restricted universe multiverse\n\
deb http://mirrors.aliyun.com/ubuntu/ focal-security main restricted\n\
deb http://mirrors.aliyun.com/ubuntu/ focal-security universe\n\
deb http://mirrors.aliyun.com/ubuntu/ focal-security multiverse" > /etc/apt/sources.list

RUN dpkg --add-architecture i386
RUN apt update && apt-get update -y
RUN apt-get install -y --no-install-recommends --assume-yes wget software-properties-common gpg-agent supervisor xvfb mingw-w64 ffmpeg cabextract aptitude vim pulseaudio
RUN wget -qO- https://dl.winehq.org/wine-builds/winehq.key | apt-key add -

RUN echo "\
deb http://mirrors.aliyun.com/ubuntu/ focal main restricted\n\
deb http://mirrors.aliyun.com/ubuntu/ focal-updates main restricted\n\
deb http://mirrors.aliyun.com/ubuntu/ focal universe\n\
deb http://mirrors.aliyun.com/ubuntu/ focal-updates universe\n\
deb http://mirrors.aliyun.com/ubuntu/ focal multiverse\n\
deb http://mirrors.aliyun.com/ubuntu/ focal-updates multiverse\n\
deb http://mirrors.aliyun.com/ubuntu/ focal-backports main restricted universe multiverse\n\
deb http://mirrors.aliyun.com/ubuntu/ focal-security main restricted\n\
deb http://mirrors.aliyun.com/ubuntu/ focal-security universe\n\
deb http://mirrors.aliyun.com/ubuntu/ focal-security multiverse\n\
deb http://dl.winehq.org/wine-builds/ubuntu/ focal main" > /etc/apt/sources.list
RUN apt update -y

RUN apt install -y --install-recommends winehq-stable winbind

RUN wget -nv -O /usr/bin/winetricks https://raw.githubusercontent.com/Winetricks/winetricks/master/src/winetricks \
    && chmod +x /usr/bin/winetricks

# Silence all the "fixme: blah blah blah" messages from wine
ENV WINEDEBUG fixme-all
RUN winetricks d3dx9_43
# uncomment it for lutris game
#RUN winetricks --force -q dotnet48

COPY ./containerfiles/default.pa /etc/pulse/
COPY ./containerfiles/supervisord.conf /etc/supervisor/conf.d/
COPY ./containerfiles/download_gecko_and_mono.sh /winecontainer/download_gecko_and_mono.sh
COPY ./containerfiles/syncinput.cpp /winecontainer/syncinput.cpp

# Download gecko and mono installers
RUN chmod +x /winecontainer/download_gecko_and_mono.sh
RUN /winecontainer/download_gecko_and_mono.sh "$(dpkg -s wine-stable | grep "^Version:\s" | awk '{print $2}' | sed -E 's/~.*$//')"

# Compile syncinput.exe
RUN x86_64-w64-mingw32-g++ /winecontainer/syncinput.cpp -o /winecontainer/syncinput.exe -lws2_32 -lpthread -static

RUN wget https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz
RUN apt-get install xz-utils
RUN tar -xvf ./ffmpeg-release-amd64-static.tar.xz
RUN mv ./ffmpeg-6.0-amd64-static/ffmpeg /usr/bin/ffmpeg_audio

ENTRYPOINT [ "supervisord" ]