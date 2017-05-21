docker run -d -p 8083:8083 -p 8086:8086 \
          -v $PWD/influxdb.conf:/etc/influxdb/influxdb.conf:ro \
          -v $PWD/influxdb:/var/lib/influxdb \
          --name=influxdb \
          --net=influxdb \
          -e INFLUXDB_ADMIN_ENABLED=true \
          influxdb -config /etc/influxdb/influxdb.conf
