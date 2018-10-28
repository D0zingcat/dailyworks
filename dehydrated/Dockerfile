FROM stretch 
MAINTAINER lee tang <d0zingcat@outlook.com>

RUN apt-get update && apt-get install  \
    bash \
    curl \
    openssl \
    git \
    vim \
    less \ 
    git clone https://github.com/lukas2511/dehydrated.git

WORKDIR /dehydrated

ENTRYPOINT ["./dehydrated"]
CMD ['-h']
