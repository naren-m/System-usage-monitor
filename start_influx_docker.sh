docker run -d -p 8083:8083 -p 8086:8086 \
          -v $PWD/influxdb:/var/lib/influxdb \
          --name=influxdb \
          --net=robot \
          --rm \
          influxdb
