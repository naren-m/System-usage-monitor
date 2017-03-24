docker run -p 9092:9092 -d \
          -v $PWD/kapacitor.conf:/etc/kapacitor/kapacitor.conf:ro \
          -v $PWD/tickscripts:/tickscripts \
          -v $PWD/templates:/templates \
          --name=kapacitor \
          -e KAPACITOR_INFLUXDB_0_URLS_0=http://influxdb:8086 \
          -h kapacitor \
          --net=influxdb \
          --rm \
          kapacitor
