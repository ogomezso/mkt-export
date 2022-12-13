# Marketplace *Export Tool

## Use

Simple Client for one-off export if existing resources on Confluent Cluster to Santander's Marketplace.

There are 2 main commands:

**export** : Take excel file with the agreed structure as inpunt and export data for event and subscription json payloads matching with the marketplace events API. The json files will be copied to `events` and ``
**register** : Take the absolute path to the folder containing the output of **export** command and do the request needed to register both events and subscriptions.

## Configuration

For **export** command you must provide a `config.yaml` configuration file as following:

```yaml
input: <absolute path to input excel file>
output: <absolute path to root output folder>
```

For **register** command the configuration needed is:

```yaml
input : <absolute path of the json files>
output : . <not used but required>
marketplace:
  appkey: <app key used for login on marketplace>
  mktplaceurl: <root-path of marketplace>
  credentials:
    bearer: <bearer used for authorization on marketplace>
```


