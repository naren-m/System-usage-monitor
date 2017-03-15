# Trigger alert from Stream data using Templates

## Define Template

## Reference

[Defining template examples](https://docs.influxdata.com/kapacitor/v1.0/examples/template_tasks/)

[Kapacitor template APIS](https://docs.influxdata.com/kapacitor/v1.0/api/api/#templates)


Steps to define and enable alerts using templates

Defining Template
```shell
kapacitor define-template basic_template -tick basic_template.tick -type stream
kapacitor show-template basic_template
```

CPU

```shell
kapacitor define cpu_alert -template basic_template -vars cpu_vars.json -dbrp telegraf.autogen
kapacitor show cpu_alert
kapacitor enable cpu_alert
```


Memory

```shell
kapacitor define mem_alert -template basic_template -vars mem_vars.json -dbrp telegraf.autogen
kapacitor show mem_alert
```