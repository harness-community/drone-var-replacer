A plugin to replace and override variables with values in a template file.

# Usage

The plugin looks for the template file where the variable is defined.

- PLUGIN_TEMPLATE_FILE_PATH

Below is an example that uses this plugin.

```yaml
- step:
    type: Plugin
    name: var-replacer
    identifier: var-replacer
    spec:
      connectorRef: harnessdocker
      image: plugins/var-replacer:latest
      settings:
        template_file_path: template.yaml
        pythonversion: python3.9
        httpmethod: get
```

Here, `pythonversion` and `httpmethod` are the variables defined in the template.yaml file.

# Building

Build the plugin binary:

```text
scripts/build.sh
```

Build the plugin image:

```text
docker build -t {{DockerRepository}} -f docker/Dockerfile .
```
