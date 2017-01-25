docker run -d --name=telegraf \
      --net=influxdb \
      -v $PWD/telegraf.conf:/etc/telegraf/telegraf.conf:ro \
      telegraf
