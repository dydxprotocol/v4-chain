# pprof support for Node.js

[![NPM Version][npm-image]][npm-url]
[![Build Status][build-image]][build-url]
[![Known Vulnerabilities][snyk-image]][snyk-url]

[pprof][pprof-url] support for Node.js.

## Prerequisites
1. Your application will need to be using Node.js 14 or greater.

2. The `pprof` module has a native component that is used to collect profiles 
with v8's CPU and Heap profilers. You may need to install additional
dependencies to build this module.
    * For Linux: `pprof` has prebuilt binaries available for Linux arm64/x64,
    Alpine Linux x64, macOS arm64/x64, windows x64 for Node 14/16/18/20.
    No additional dependencies are required.
    * For other environments: on environments that `pprof` does not have
    prebuilt binaries for, the module
    [`node-gyp`](https://www.npmjs.com/package/node-gyp) will be used to
    build binaries. See `node-gyp`'s
    [documentation](https://github.com/nodejs/node-gyp#installation)
    for information on dependencies required to build binaries with `node-gyp`.

3. The [`pprof`][pprof-url] CLI can be used to view profiles collected with
this module. Instructions for installing the `pprof` CLI can be found
[here][pprof-install-url].

## Basic Set-up

Install [`pprof`][npm-url] with `npm` or add to your `package.json`.
  ```sh
  # Install through npm while saving to the local 'package.json'
  npm install --save @datadog/pprof
  ```

## Using the Profiler

### Collect a Wall Time Profile

#### In code:
1. Update code to collect and save a profile:
    ```javascript
    const profile = await pprof.time.profile({
      durationMillis: 10000,    // time in milliseconds for which to 
                                // collect profile.
    });
    const buf = await pprof.encode(profile);
    fs.writeFile('wall.pb.gz', buf, (err) => {
      if (err) throw err;
    });
    ```

2. View the profile with command line [`pprof`][pprof-url]:
    ```sh
    pprof -http=: wall.pb.gz
    ```

#### Requiring from the command line

1. Start program from the command line:
    ```sh
    node --require @datadog/pprof app.js
    ```

2. A wall time profile for the job will be saved in 
`pprof-profile-${process.pid}.pb.gz`. View the profile with command line 
[`pprof`][pprof-url]:
    ```sh
    pprof -http=: pprof-profile-${process.pid}.pb.gz
    ```

### Collect a Heap Profile
1. Enable heap profiling at the start of the application:
    ```javascript
    // The average number of bytes between samples.
    const intervalBytes = 512 * 1024;

    // The maximum stack depth for samples collected.
    const stackDepth = 64;

    heap.start(intervalBytes, stackDepth); 
    ```
2. Collect heap profiles:
  
    * Collecting and saving a profile in profile.proto format:
        ```javascript
        const profile = await pprof.heap.profile();
        const buf = await pprof.encode(profile);
        fs.writeFile('heap.pb.gz', buf, (err) => {
          if (err) throw err;
        })
        ```

    * View the profile with command line [`pprof`][pprof-url].
        ```sh
        pprof -http=: heap.pb.gz
        ```
    
    * Collecting a heap profile with  V8 allocation profile format:
        ```javascript
          const profile = await pprof.heap.v8Profile();
        ``` 

[build-image]: https://github.com/Datadog/pprof-nodejs/actions/workflows/build.yml/badge.svg?branch=main
[build-url]: https://github.com/Datadog/pprof-nodejs/actions/workflows/build.yml
[coveralls-image]: https://coveralls.io/repos/google/pprof-nodejs/badge.svg?branch=main&service=github
[npm-image]: https://badge.fury.io/js/pprof.svg
[npm-url]: https://npmjs.org/package/pprof
[pprof-url]: https://github.com/google/pprof
[pprof-install-url]: https://github.com/google/pprof#building-pprof
[snyk-image]: https://snyk.io/test/github/google/pprof-nodejs/badge.svg
[snyk-url]: https://snyk.io/test/github/google/pprof-nodejs
