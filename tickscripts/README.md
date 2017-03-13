# Trigger alert from Stream data using Tickscript

## Define Tickscript

kapacitor define cpu_alert -type stream -tick cpu_alert.tick -dbrp telegraf.autogen

## Enable Tickscript

Now that we know it’s working, let’s change it back to a more reasonable
threshold. Are you happy with the threshold? If so, let’s enable the task so it
can start processing the live data stream with:

```s
kapacitor enable cpu_alert
```

Show Kapacitor task

```s
kapacitor show cpu_alert
```

Running task on Kapacitor to bring CPU usage to above 70

```shell
while true; do i=0; done
```

For detail explanation please look at the [Kapacitor geting started
guide](https://docs.influxdata.com/kapacitor/v1.2/introduction/getting_started/)

## Testing the Tick script

However nothing is going to happen until we enable the task. Before we enable
the task, we should test it first so we do not spam ourselves with alerts.
Record the current data stream for a bit so we can use it to test our task with

```shell
kapacitor record stream -task cpu_alert -duration 20s
```

Since we defined the task with a database and retention policy pair, the
recording knows to only record data from that database and retention policy. Now
grab that ID that was returned and let’s put it in a bash variable for easy use
later (your ID will be different):

```s
rid=cd158f21-02e6-405c-8527-261ae6f26153
kapacitor list recordings $rid
```

You should see some output like:

```s
ID                                      Type    Status    Size      Date
cd158f21-02e6-405c-8527-261ae6f26153    stream  finished  1.6 MB    04 May 16 11:44 MDT
```

OK, we have a snapshot of data recorded from the stream, so we can now replay
that data to our task. The replay action replays data only to a specific task.
This way we can test the task in complete isolation:

```s
kapacitor replay -recording $rid -task cpu_alert
```

Check logs as below:

```shell
cat /tmp/alerts.log
```