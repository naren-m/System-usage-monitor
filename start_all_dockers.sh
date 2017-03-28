./start_influx_docker.sh
echo "Sleeping for 1 sec"
sleep 1
./start_kapacitor_docker.sh
./start_telegraf_docker.sh
./start_grafana_docker.sh
