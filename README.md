# NordVPN GTK [WIP]

NordVPN GTK is a GTK+ Linux client for NordVPN built using [OpenNord](https://github.com/adamdb5/opennord).

Current, NordVPN GTK is in development and not all features are complete.

TODO:
 - Partial config saving
 - Complete whitelist implementation
 - Refactor to use hierachial structure (stop passing a reference of the App to all functions)
 - Consider re-writing this and OpenNord in a more performant language like C++ or Rust. (Unfortunately, there are no C bindings for gRPC)


## Screeshots
These are some screenshots of the current application. The UI layout may change in the future.

Connect Tab:

![img](/docs/images/ConnectTab.png)


Session Tab:

![img](/docs/images/SessionTab.png)


Configure Tab:

![img](/docs/images/ConfigureTab.png)


Whitelist Tab:

![img](/docs/images/WhitelistTab.png)


Account Tab:

![img](/docs/images/AccountTab.png)


About Tab:

![img](/docs/images/AboutTab.png)


## Contributing
If you run into any issues or see something that you would like to improve, please feel free to create an issue or raise
a pull request.

## License
This library is licensed under the MIT license.
