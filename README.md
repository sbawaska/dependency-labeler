# deplab

## Introduction
Deplab adds metadata about a container image's dependencies as a label to the container image.

## Dependencies
Docker is required to be installed and available on your path, test this by running `docker version`.
API version 1.39 or higher is required.

## Usage
Download the latest `deplab` binary from the releases page.
To run the tool run the following command:
```bash
./deplab --image <image name> --git <path to git repo>
```

* `<image name>` is the name of the image that you want to add the metadata to.
* `<path to git repo>` is a path to a directory under git version control.

This returns the sha256 of the new image with added metadata.
Currently this will add the label `io.pivotal.metadata` along with the necessary metadata.

To visualise the metadata this command can be run

```bash
docker inspect $(./deplab --image <image-name> --git <path to git repo>) \
  | jq -r '.[0].Config.Labels."io.pivotal.metadata"' \ 
  | jq .
```

### Multiple git repositories

You can specify as many git repositories as required by passing more than one
git flag into the command.
```bash
./deplab --image <image name> --git <path to git repo> --git <path to another git repo>
```

The output will look like:
```json
{
  "dependencies": [
    {
      "type": "package",
      "source": {
        "type": "git",
        "version": {
          "commit":  "d2c3ccdffd3c5a014891e40a3ed8ba020d00eefd"
         },
        "metadata": {
          "url": "https://github.com/pivotal/deplab.git",
          "refs": ["0.5.0"]
        }
      }
    },
    {
      "type": "package",
      "source": {
        "type": "git",
        "version": {
          "commit":  "d2a3ccdffd3c5a014891e40a3ed8ba020d00eefd"
         },
        "metadata": {
          "url": "https://github.com/pivotal/anotherdeplab.git",
          "refs": ["0.3.0"]
        }
      }
    }
  ]
}
```

### Usage with tarball

Alternatively, deplab can be used with an image stored locally in tar format.

```bash
./deplab --image-tar <path to image.tar> --git <path to git repo>
```

* `<path to image.tar>` is the path to the tarball.
* `<path to git repo>` is, as above, a path to a directory under git version control (it can be specified multiple times).

This returns the sha256 of the new image with added metadata.
Currently this will add the label `io.pivotal.metadata` along with the necessary metadata.

### Usage in Concourse

Please see [CONCOURSE.md](CONCOURSE.md) for information about using Deplab as a task in your
Concourse pipeline.

## Optional Arguments

### Metadata file
Deplab can output the metadata to a file providing the path with the argument `--metadata-file` or `-m` 

```bash
./deplab -i <image name> -g <path to git repo> --metadata-file <metadata file>
```

If the file path cannot be created, deplab will return the newly labelled image, and return an error for the writing of the metadata file. 

If a file exists at the given path, the file will be overwritten.

### dpkg file
Deplab can output the debian package list portion of the metadata to a file with the argument `--dpkg-file` or `-d`

```bash
./deplab -i <image name> -g <path to git repo> --dpkg-file <dpkg file>
```

If the file path cannot be created, deplab will return the newly labelled image, and return an error for the writing of the dpkg file. 

If a file exists at the given path, the file will be overwritten.

This file is approximately similar to the file which will be output by running `dpkg -l`, with the addition of an extra header which provides an ID for this list.

### Output as image tarball

Deplab can output the image in tar format.

```bash
./deplab -i <image name> -g <path to git repo> --output-tar ./path/to/image.tar
```

If the file path cannot be created deplab will process the image and store it in Docker, but will also return an error for the writing of the tar. 

If a file exists at the given path, the file will be overwritten.

### Tag
Deplab can add a tag to the output image

```bash
./deplab -i <image name> -g <path to git repo> --tag <tag name>
```
The SHA256 will be returned as normal.

You can inspect the tag has been added to the image:
```bash
docker inspect $(./deplab --image <image-name> --git <path to git repo> --tag <tag name>) \
 | jq '.[0].RepoTags'
```
 
## Data

##### debian package list

The `debian_package_list` requires `dpkg` to be a package on the image being instrumented on. If not present, the dependency of type `debian_package_list` will be omitted.

`version` contains the _sha256_ of the `json` content of the metadata. Successive run of deplab on containers with the same `packages` and `apt_sources` are going to generate the same digest.

The debian package list is generated with the following format.

```json
{
  "dependencies": [
    {
      "type": "debian_package_list",
      "version": {
        "sha256": "a56...42b"
      },
      "source": {
        "type": "inline",
        "version": null,
        "metadata": {
          "packages": [...],
          "apt_sources": [...]
        }
      }
    }
  ]
}
```



Example of a package item in field `packages` 

```json
{
  "package": "zlib1g",
  "version": "1:1.2.11.dfsg-0ubuntu2",
  "architecture": "amd64",
  "source": {
    "package": "zlib",
    "version": "1:1.2.11.dfsg-0ubuntu2",
    "upstreamVersion": "1.2.11.dfsg"
  }
}
```

Example of `apt_sources` content

```json
[
  "deb http://archive.ubuntu.com/ubuntu/ bionic main restricted",
  "deb http://archive.ubuntu.com/ubuntu/ bionic-updates main restricted",
  "deb http://archive.ubuntu.com/ubuntu/ bionic universe",
  "deb http://archive.ubuntu.com/ubuntu/ bionic-updates universe",
  "deb http://archive.ubuntu.com/ubuntu/ bionic multiverse",
  "deb http://archive.ubuntu.com/ubuntu/ bionic-updates multiverse",
  "deb http://archive.ubuntu.com/ubuntu/ bionic-backports main restricted universe multiverse",
  "deb http://security.ubuntu.com/ubuntu/ bionic-security main restricted",
  "deb http://security.ubuntu.com/ubuntu/ bionic-security universe",
  "deb http://security.ubuntu.com/ubuntu/ bionic-security multiverse"
]
```

##### git dependency
If the `--git` flag is provided with a valid path to a git repository, a git dependency will be added:
```json
{
  "dependencies": [
    {
      "type": "package",
      "source": {
        "type": "git",
        "version": {
          "commit":  "d2c3ccdffd3c5a014891e40a3ed8ba020d00eefd"
         },
        "metadata": {
          "url": "https://github.com/pivotal/deplab.git",
          "refs": ["0.5.0"]
        }
      }
    }
  ]
}
```

You may add multiple git repositories by adding additional git flags:
```bash
./deplab --image <image name> --git <path to git repo> --git <path to another git repo>
```

The output will look like:
```json
{
  "dependencies": [
    {
      "type": "package",
      "source": {
        "type": "git",
        "version": {
          "commit":  "d2c3ccdffd3c5a014891e40a3ed8ba020d00eefd"
         },
        "metadata": {
          "url": "https://github.com/pivotal/deplab.git",
          "refs": ["0.5.0"]
        }
      }
    },
    {
      "type": "package",
      "source": {
        "type": "git",
        "version": {
          "commit":  "d2a3ccdffd3c5a014891e40a3ed8ba020d00eefd"
         },
        "metadata": {
          "url": "https://github.com/pivotal/anotherdeplab.git",
          "refs": ["0.3.0"]
        }
      }
    }
  ]
}
```


##### base
The base image metadata is generated with the following format
```json
  "base": {
    "name": "Ubuntu",
    "version_id": "18.04",
    "version_codename": "bionic"
  }
```

This relies on the `/etc/os-release` file being in the docker container, and `cat` being able to read it. If either are not present all the field will be set to `unknown`.

```json
{
  "name": "unknown",
  "version_id": "unknown",
  "version_codename": "unknown"
}
```

## Testing
Testing requires `go` to be installed.  Please clone this git repository.  Tests can be run with:
```bash
go test ./...
```

## Building

To build for release, please run the following:
```bash
go build -o deplab ./cmd/deplab
```

To build the Concourse task image, please run the following:
```bash
docker build . -f Dockerfile.task
```

## Support

This tool is currently maintained by the Pivotal NavCon team;
@navcon in #navcon-team on Pivotal Slack.

Please reach out to us on Slack first, and then raise a Github issue.
