# docker run -d --name=grafana --rm -p 3009:3009 --net=influxdb narenm/grafana
docker run -d -i --net=influxdb \
                 -p 3000:3000 \
                 -v $PWD/grafana:/var/lib/grafana \
                 --name grafana grafana/grafana
