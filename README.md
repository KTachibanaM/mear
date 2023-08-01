# mear
Bring-your-own-cloud, on-demand media encoding

"mear" stands for "Media Encoder At-Request"


## Motivation
For indie developers, media encoding can be expensive. You need to either.

* use a cloud encoding service, which can be expensive (see [pricing comparison](#pricing-comparison))
* keep a powerful VPS instance always-on, which is expensive if you don't fully utilize it
* keep a powerful PC tower at home, which takes up space and does not easily connect to the cloud where your application is deployed

`mear` is an on-demand media encoding tool that spawns VPS instances on-demand (currently supports DigitalOcean), uses it for media encoding, and shuts it down when it's done. By using `mear`, you 

* control the media encoding process because it uses your cloud account
* only pay for what you use
* don't need to keep a powerful machine always-on (whether at home or on the cloud)

For regular users of `ffmpeg` (`mear` uses the mighty `ffmpeg` under the hood), `mear` can be used as a replacement to encode media files without having to hang your machine because `ffmpeg` can take up a lot of your personal computer's resources.


## Overview
`mear` at its core is a Go program that spawns one-off VPS instances and S3 buckets using your cloud credentials, uploads/downloads media files to/from the S3 buckets on your machine and VPS instances, and runs `ffmpeg` on the VPS instances via ssh to encode the media files.

Architecturally it is run on a `host` environment and an `engine` environment. A `host` is where all the cloud resource orchestration happens, e.g., creating VPS instances and S3 buckets. An `engine` is where the media encoding happens, e.g., uploading/downloading media files to/from S3 buckets, downloading and running `ffmpeg`. In reality, the `engine` runs an `agent` binary to accomplish all its tasks so that it's only one ssh execution from `host` to `engine`.

Currently, we distribute two Go binaries for different use cases.

* `mear-host` is an application-facing CLI that accepts a JSON payload, emits log lines as newline-delimited JSON, and terminates with a 0 exit code if media encoding is successful
* `mear` is a user-facing CLI that is a direct replacement for `ffmpeg` users


## Installation
Whether you are an application developer or an `ffmpeg` user, run

```bash
curl -L https://raw.githubusercontent.com/KTachibanaM/mear/master/install.sh | bash
```

This script will install `mear` and `mear-host` to `/usr/local/bin` on Linux and macOS.

You can also put this line in your `Dockerfile` to distribute `mear-host` with your application.


## Usage for application developers
TBD


## Usage for `ffmpeg` users
Once `mear` is installed, you can use it as a direct replacement for `ffmpeg`. For example, to encode a video to 720p, H.264 on DigitalOcean, run

```bash
mear -i test.avi --mear-stack do --mear-agent-timeout 60 --mear-do-ram 16 --mear-do-cpu 8 -vf scale=-1:720 test.mp4
```

Some explanations on the cli usage:

* `-i test.avi` is the input file. It's the same as in `ffmpeg`.
* `--mear-stack do` is the cloud provider stack for media encoding. `do` specifies DigitalOcean.
* `--mear-agent-timeout 60` is the deadline for the `engine` to finish encoding. If encoding doesn't finish within 60 minutes, `host` will terminate the `engine`, and cli will fail.
* `--mear-do-ram 16` is the amount of RAM in GB to use for the `engine` on DigitalOcean. A specific set of RAM/CPU cores combinations are supported. Run `mear` without any arguments to see the supported combinations.
* `--mear-do-cpu 8` is the number of CPU cores to use for the `engine` on DigitalOcean. A specific set of RAM/CPU cores combinations are supported. Run `mear` without any arguments to see the supported combinations.
* `-vf scale=-1:720` is the extra argument you'd pass into `ffmpeg`. `mear` will interpret any arguments without the `--mear` prefix as `ffmpeg` arguments.
* `test.mp4` is the output file. It's the same as in `ffmpeg`.


## Pricing comparison
To encode a 1-hour long, MPEG-4 encoded, 30fps, 9000Kbps, 1080p video to an H.264 equivalent

* [AWS Elemental MediaConvert](https://aws.amazon.com/mediaconvert/) would cost you [$0.9](https://calculator.aws/#/estimate?id=9474477a1f71466e30b55f5de02737da8f756f85)
* [Qencode](https://cloud.qencode.com/pricing) would cost you $0.9
* `mear`, using an `s-8vcpu-16gb` (16GB RAM, 8 CPU cores) Droplet for an hour, would cost you **$0.14**. Additionally, if you exceed the Spaces subscription limit, S3 would cost **$0.06**. Best case `mear` is 6x cheaper.

Some caveats for the comparison.

1. Encoding on DigitalOcean using `mear` is about 4x slower because the compute resource is less powerful than the cloud service equivalent. It is also due to the time taken to upload/download media files from/to S3 (although it is less of a problem if you run `mear-host` on an application running in the cloud because cloud-to-cloud networking is generally faster than cloud-to-residential networking)

2. `mear` on DigitalOcean is probably more expensive for encoding smaller media files (lower resolution, lower bitrate, shorter duration, or images). You should also be aware that since a Droplet's minimal cost is for one hour, you will be charged for one hour even if the media encoding takes less than an hour.


## Development
Clone and [open the project in VSCode Dev Container](https://code.visualstudio.com/docs/devcontainers/containers#_quick-start-open-an-existing-folder-in-a-container).

VSCode should have installed all development dependencies.

Run `make` to compile the Go binaries.

Run `./dev/download-demo-videos.sh` to download demo videos used for testing. Running this command is a prerequisite for running the following two commands.

Run `./dev/test-host.sh` to test the `mear-host` CLI (application-facing). The script encodes two mp4 files into two avi files and saves them into the `./dev` directory.

Run `./dev/test-cli.sh` to test the `mear` CLI (user-facing). The script encodes an mp4 file into an avi file and saves it into the `./dev` directory.

We use docker containers and `minio` in this development environment for the `engine` and S3 buckets.

Run `./dev/clean-dev.sh` to clean up development docker containers and S3 buckets if `mear` fails to clean up properly.
