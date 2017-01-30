Influxdata getting started with dockers


**Starting dockers**

Use below scripts to launch dockers

    - start_influx_docker.sh
    - start_kapacitor_docker.sh
    - start_telegraf_docker.sh


**Defining Kapacitor**

kapacitor define cpu_alert -type stream -tick cpu_alert.tick -dbrp telegraf.autogen


**Enable Kapacitor**

Now that we know it’s working, let’s change it back to a more reasonable
threshold. Are you happy with the threshold? If so, let’s enable the task so it
can start processing the live data stream with:
```
kapacitor enable cpu_alert
```
Show Kapacitor task       
```
kapacitor show cpu_alert
```
Running task on Kapacitor to bring CPU usage to above 70
```
while true; do i=0; done
```

For detail explanation please look at the [Kapacitor geting started
guide](https://docs.influxdata.com/kapacitor/v1.2/introduction/getting_started/)

**Testing the Tick script**

However nothing is going to happen until we enable the task. Before we enable
the task, we should test it first so we do not spam ourselves with alerts.
Record the current data stream for a bit so we can use it to test our task with

```
kapacitor record stream -task cpu_alert -duration 20s
```

Since we defined the task with a database and retention policy pair, the
recording knows to only record data from that database and retention policy. Now
grab that ID that was returned and let’s put it in a bash variable for easy use
later (your ID will be different):

```
rid=cd158f21-02e6-405c-8527-261ae6f26153
kapacitor list recordings $rid
```
You should see some output like:
```
ID                                      Type    Status    Size      Date
cd158f21-02e6-405c-8527-261ae6f26153    stream  finished  1.6 MB    04 May 16 11:44 MDT
```

OK, we have a snapshot of data recorded from the stream, so we can now replay
that data to our task. The replay action replays data only to a specific task.
This way we can test the task in complete isolation:
```
kapacitor replay -recording $rid -task cpu_alert
```

Check logs as below:
```
cat /tmp/alerts.log
```

**Checking logs**
```
docker logs -f influxdb
docker logs -f telegraf
docker logs -f kapacitor

```

**Creating Config file**

InfluxDB  - docker run --rm influxdb incluxd config > influxdb.conf

Kapacitor - docker run --rm kapacitor kapacitord config > kapacitor.conf

Telegraf  - docker run --rm telegraf -sample-config -input-filter cpu:mem -output-filter influxdb > telegraf.conf
