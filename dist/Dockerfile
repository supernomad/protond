# Copyright (c) 2017 Christian Saide <Supernomad>
# Licensed under the MPL-2.0, for details see https://github.com/Supernomad/protond/blob/master/LICENSE

FROM ubuntu

ENV DEBIAN_FRONTEND "noninteractive"

RUN apt-get update \
    && apt-get install -y --no-install-recommends --no-upgrade \
        tcpdump \
    && rm -rf /var/lib/apt

ENTRYPOINT ["/bin/protond"]
