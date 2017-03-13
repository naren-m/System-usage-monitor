# Trigger alert from Stream data using Templates

## Define Template

kapacitor define cpu_alert -type stream -tick cpu_alert.tick -dbrp telegraf.autogen
