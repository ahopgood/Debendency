# Debendency
A tool to output debian package dependencies

## Parameters
Output from `go run main.go -help` at this point in time:
```
Usage of /tmp/go-build1711827394/b001/exe/main:
  -d    output dependencies as a diagram
  -o string
        output directory to save installer files to
  -p string
        .deb package name to calculate dependencies for
  -s    output dependencies as salt code

```

## To-Do 
* Bug fixes
  * Support case where file already has been downloaded and we don't create a model as we cannot parse the parameters from the download output
    * dpkg -I might be needed to list these details if the download doesn't happen
* Create command line flags
  * Mandatory vs Optional, can we support these
* Test command line flags
* Add testing command line flags to only run if on linux - integration tests
* Add integration tests
* Logging
  * Add standard out logger
  * Add standard error logger
  * Add verbose flag to set logging level to DEBUG
* Add puml output
    * Move model parsing to puml.go
    * Update puml.go to add the starting dependency to the title of the graph
    * Add puml tests
* Add Makefile with goals
  * Build
  * Test
  * Integration test
  * Linting?
  * Formatting
  * Coverage
  * Go generate for fakes
* Supplementary Features
  * apt-get update before first download
  * Download cache
    * Default location of cache `~/debendency/cache`
    * Flag to specify location for the cache
    * Flag to delete download cache after use?
  * Hide already installed packages - useful if you only need to know the dependencies that still need to be installed on a system
  * Salt output
    * Salt output tests
  * Puppet output