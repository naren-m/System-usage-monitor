# Influxdata getting started with dockers


## Starting dockers

Use below scripts to launch dockers

    - start_influx_docker.sh
    - start_kapacitor_docker.sh
    - start_telegraf_docker.sh

## Trigger alert from Stream data using tick scripts

[Defining teplates](https://github.com/naren-m/influxdb_get_started/tree/master/tickscripts)

## Trigger alert from Stream data using templates

[Defining templates](https://github.com/naren-m/influxdb_get_started/tree/master/tickscripts/templates)

## Checking logs

```shell
docker logs -f influxdb
docker logs -f telegraf
docker logs -f kapacitor

```

## Creating Config file

InfluxDB  - docker run --rm influxdb incluxd config > influxdb.conf

Kapacitor - docker run --rm kapacitor kapacitord config > kapacitor.conf

Telegraf  - docker run --rm telegraf -sample-config -input-filter cpu:mem -output-filter influxdb > telegraf.conf

## References

1.[Kapacitor geting startedguide](https://docs.influxdata.com/kapacitor/v1.2/introduction/getting_started/)

2.[Kapacitpor API Documentation](https://docs.influxdata.com/kapacitor/v1.2/api/api)

3.[Kapacitor Templating Documentation](https://docs.influxdata.com/kapacitor/v1.2/examples/template_tasks/)