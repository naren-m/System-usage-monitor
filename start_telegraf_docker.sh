docker run -d --name=telegraf \
      --net=influxdb \
      -v $PWD/telegraf.conf:/etc/telegraf/telegraf.conf:ro \
      -v /var/run/docker.sock:/var/run/docker.sock \
      telegraf
