## Proto

This directory defines all protos for `v4`. We follow the Cosmos-SDK convention of using a tool called
[buf](https://github.com/bufbuild/buf) to manage proto dependencies. You can think of `buf` as being like `npm` for
protocol buffers. See the `buf` [documentation](https://docs.buf.build/how-to/iterate-on-modules#update-dependencies)
for further details.

### First time setup
Install `buf` locally:
```shell
brew install buf
```

### Update protos
After updating the hashes/tags in `buf.yaml`, update the protos and regenerate the code (from `v4/` directory):
```shell
cd proto
buf mod update
buf build
cd ..
make proto-all
```


