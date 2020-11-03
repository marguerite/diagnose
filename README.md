## Diagnosis tools for openSUSE ##

This project aims to collect the diagnosis tools I wrote to help openSUSE freshmen to report issues on forum.

It includes:

#### `instdpkg` AKA "installed packages" ####

Tell the user what packages were installed/removed on the specific `date`, since the `time`.

Case:

A user posted on forum saying "I don't know what happened, I just installed some packages then problem occurs."

Tell him to download and run `sudo instdpkg -date=2020-10-20`. It will tell you what "some packages" are.

    ====== Packages modified after 2020-10-20 00:00:00 ======
    time               |action |name            |version     |arch  |repo
    2020-10-20 10:06:43|install|libspiro1       |20200505-1.1|x86_64|repo-oss
    2020-10-20 10:06:44|install|libuninameslist1|20200413-1.2|x86_64|repo-oss
    2020-10-20 10:06:44|install|libwoff2enc1_0_2|1.0.2-3.10  |x86_64|repo-oss
    2020-10-20 10:06:49|install|fontforge       |20200314-3.3|x86_64|repo-oss

Even forgot the date? "just two weeks ago"? Run `sudo instdpkg -timeline` and find out the date.

    2020-10-01 12:34:41 +0000 UTC
    2020-10-02 11:26:26 +0000 UTC
    2020-10-05 14:34:38 +0000 UTC
    2020-10-15 11:48:06 +0000 UTC
    2020-10-20 10:06:41 +0000 UTC
    2020-11-03 11:31:21 +0000 UTC

#### `pkmswitch100` AKA "packman switch 100%?"

Tell the user if related packages (ffmpeg, vlc and gstreamer) were 100% switched from oss to packman.

Case:

A user posted on forum saying "I can cut videos using ffmpeg but can't put the cut into mp4 container"
or worse "how to cut using ffmpeg???" and attached many logs that're mostly useless unless you can find
a video that encoded using the same codecs.

In the previous case, `libavformat57` is oss while all others are from packman. If you have "update
from a different repository" disabled which is default and switch packages one by one by yourself, you
are exposed to such cases.

Tell him to download and run `pkmswitch100`, problem solved.

NOTE: Always run "sudo zypper ref" first. And this is not a installation tool but a debug tool,
it will not install the packages that you haven't installed.

There're 3 options: `-type=ffmpeg`, `-type=vlc`, `-type=gstreamer`. by default all of the three will be checked.

    ====== Packages not installed ======
    gstreamer-plugins-bad-chromaprint libgstplayer-1_0-0 gstreamer-plugins-bad-fluidsynth libgstcodecs-1_0-0 libgstvulkan-1_0-0 libgstinsertbin-1_0-0 gstreamer-transcoder vlc-codecs libgsttranscoder-1_0-0
    FIX: sudo zypper in gstreamer-plugins-bad-chromaprint libgstplayer-1_0-0 gstreamer-plugins-bad-fluidsynth libgstcodecs-1_0-0 libgstvulkan-1_0-0 libgstinsertbin-1_0-0 gstreamer-transcoder vlc-codecs libgsttranscoder-1_0-0 --from packman
    ====== Packages should be updated ASAP ======
    libgstwebrtc-1_0-0 libswscale5 libgstisoff-1_0-0 gstreamer-plugins-ugly gstreamer-plugins-bad-orig-addon libgstbasecamerabinsrc-1_0-0 libavutil56 libavfilter7 libgsturidownloader-1_0-0 gstreamer-plugins-bad libswresample3 libgstcodecparsers-1_0-0 libgstadaptivedemux-1_0-0 libgstbadaudio-1_0-0 libgstsctp-1_0-0 gstreamer-plugins-libav libgstmpegts-1_0-0 libavdevice58 libpostproc55 libavresample4 gstreamer-plugins-ugly-orig-addon libgstwayland-1_0-0 libgstphotography-1_0-0 libavformat58 libavcodec58
    FIX: sudo zypper up libgstwebrtc-1_0-0 libswscale5 libgstisoff-1_0-0 gstreamer-plugins-ugly gstreamer-plugins-bad-orig-addon libgstbasecamerabinsrc-1_0-0 libavutil56 libavfilter7 libgsturidownloader-1_0-0 gstreamer-plugins-bad libswresample3 libgstcodecparsers-1_0-0 libgstadaptivedemux-1_0-0 libgstbadaudio-1_0-0 libgstsctp-1_0-0 gstreamer-plugins-libav libgstmpegts-1_0-0 libavdevice58 libpostproc55 libavresample4 gstreamer-plugins-ugly-orig-addon libgstwayland-1_0-0 libgstphotography-1_0-0 libavformat58 libavcodec58

#### rescue-network

Bring your network up in rescue mode. it can also fix your network when NetworkManager or wickedd is not available.

Case 0: sometimes you need to run `zypper dup` in rescue mode. But you don't have any network connection under rescue mode.

So you connect to your wifi with a Phone and get the gateway, netmask.

Now run `rescue-network -device=wifi -gateway=192.168.31.1 -netmask=255.255.255.0 -essid=MyHomeNetwork --password=12345678`.

And you'll have a running network.

Then do what you want.

Case 1: sometimes your NM/wickedd is not working, but you need to browse the internet to find answers.

You can `sudo systemctl stop network` and run `rescue-network` as above.

Available options:

    -device "wired" or "wifi", when using wired, it will assign a static IP for you. when using wifi, it'll connect to the WIFI with wpa_supplicant first.
    -gateway your router's IP address, usually it's 192.168.1.1
    -netmask your router's netmask, usually it's 255.255.255.0
    -essid your WIFI's name.
    -password your WIFI's password

## rpm-unowned

Check for file not owned by rpm:

    rpm-unowned -dir /usr/lib64

Everything printed are not owned by rpm.

## consistent-repo

Check for packages not installed from a specific repository

    consistent-repo -p <pkg> -r <repo>

The `pkg` is a string that can be separated by ",", specifying the package or package list to be passed to `zypper se`.
`consistent-repo` will search the terms via zypper and check every package returned against the repo.

The `repo` is your `local alias` of an openSUSE repo. eg, the upstream name maybe "openSUSE-Leap-15.0-Oss", but your local
alias is just "oss", use "oss". You can get alias via "zypper lr".
