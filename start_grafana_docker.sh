# docker run -d --name=grafana --rm -p 3000:3000 --net=robot narenm/grafana
docker run -d -i --net=influx -p 3000:3000 -v $PWD/grafana:/var/lib/grafana --name grafana grafana/grafana

