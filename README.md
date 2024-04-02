## Instructions

1. Start coding on your features!
1. Once satisfied, create a tag `vX.Y.Z` and push it to trigger the GitHub Action that will create the release.
1. Navigate to the releases tab of your repository, you'll find a draft release based on your tag. Publish it to complete the release process.
1. Want to share your plugin with other developers? Consider adding it to the plugin catalog, by sending a PR editing the `plugins.json` file located in the public [Aspect CLI plugins registry](https://github.com/aspect-build/aspect-cli/tree/main/docs/plugins).

### Local development

1. Run `bazel build :dev` to build your plugin
1. In an existing codebase using aspect-cli, add the plugin like following:

```
# .aspect/cli/config.yaml
plugins:
  - name: <your-plugin-name>
    from: <relative path to the plugin>/bazel-bin/plugin
```

1. You're all set, just remember to build the `:dev` target whenever you change code in the plugin.

### Working with BEP events 

If you want to work with `BEPEventCallback(event *buildeventstream.BuildEvent) error` but you're not very familiar with the event structures, 
you can run the Bazel command you're looking to augment with the `--build_event_json_file=bep.json` and inspect the resulting `bep.json` to 
get a rough understanding about where you should look into.

See also the [Bazel user guide page about BEP](https://bazel.build/remote/bep).
