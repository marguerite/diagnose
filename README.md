## Diagnose tools for openSUSE ##

This project aims to collect the diagnose tools I wrote to help openSUSE freshmen to report issues on forum.
All written in ruby and can be natively ran on openSUSE.

It includes:

#### `instdpkg` AKA "installed packages" ####

Tell the user what packages he/she installed/removed on the specific `date`, since the `time`.

Case: 

A user posted on forum saying "I don't know what happened, I just installed some packages then problem occurs."

Tell him to download and run `instdpkg.rb -d=date -t=time`. It will tell you what "some packages" are.

	================ Packages altered on 2016-12-22 after 23:26:12 ======================
	|  date   |  time  | oper. | name | version | repo |
	2016-12-22|23:39:36|install|libavformat57|3.2.2-3.3.x86_64|packman
	2016-12-22|23:26:12|install|libgstbasecamerabinsrc-1_0-0|1.10.2-4.7.x86_64|packman
	2016-12-22|23:26:12|install|libgstphotography-1_0-0|1.10.2-4.7.x86_64|packman
	2016-12-22|23:26:12|install|gstreamer-plugins-ugly|1.10.2-4.2.x86_64|packman
	2016-12-22|23:26:13|install|libgstbadvideo-1_0-0|1.10.2-4.7.x86_64|packman
	2016-12-22|23:26:13|install|libgstbadaudio-1_0-0|1.10.2-4.7.x86_64|packman
	2016-12-22|23:26:13|install|libgstadaptivedemux-1_0-0|1.10.2-4.7.x86_64|packman
	2016-12-22|23:26:14|install|gstreamer-plugins-bad|1.10.2-4.7.x86_64|packman

Even forgot the date? "just two weeks ago"? Run `instdpkg.rb -tl` and find out the date.

	2016-06-01
	2016-11-17
	2016-12-23

#### `pkmswitch100` AKA "packman switch 100%?"

Tell the user if he/she has related packages (ffmpeg, vlc and etc) 100% switched from oss's to packman's.

Case:

A user posted on forum saying "I can cut videos using ffmpeg but can't put the cut into mp4 container"
or worse "how to cut using ffmpeg???" and attached many logs that're mostly useless unless you can find
a video that encoded using the same codecs.

In the previous case, `libavformat57` is oss while all others are from packman. If you have "update
from a different repository" disabled which is default and switch packages one by one by yourself, you
are exposed to such cases.

Tell him to doownload and run `pkmswitch100`, problem solved.

NOTE: Always run "sudo zypper ref" first. And this is not a installation tool but a debug tool, 
it will not install the packages that you haven't installed.

There're 3 options: "-ffmpeg", "-vlc", "-gstreamer". by default all of the three will be checked.

	======================= Packages not from Packman =========================
	libavformat57

	FIX: Run 'sudo zypper install libavformat57-3.2.2-3.4.x86_64'.
	======================= Packman Packages need updates =====================
	libavutil55
	ffmpeg
	libavfilter6
	libpostproc54
	libavcodec57
	libswresample2
	libswscale4
	vlc-lang
	gstreamer-plugins-ugly-orig-addon
	libavdevice57
	libavresample3

	FIX: Run 'sudo zypper up libavutil55 ffmpeg libavfilter6 libpostproc54 libavcodec57 libswresample2 libswscale4 vlc-lang gstreamer-plugins-ugly-orig-addon libavdevice57 libavresample3'.

#### `bstmirror` AKA "Best Mirror"

As it says, it will find the best mirror of openSUSE/Packman for you.

The usage is quite easy, `ruby bstmirror.rb -region="North America" -os=422 -file=quick` will get this:

	======================= Best openSUSE Mirror ======================
	http://mirrors.tuna.tsinghua.edu.cn/opensuse/
	======================= Best Packman Mirror =======================
	http://mirror.karneval.cz/pub/linux/packman/

You can omit "-region" option if you want to check worldwide.

If you omit "-os" it will check 42.2 mirrors by default

You can also omit "-file=quick" because by default it's quick, or use "-file=long" option which will download a larger file for more accurate speed test.

NOTE: Mirror lists are obtained from mirrors.opensuse.org and packman.link2linux.de/mirrors. If you want to use some underground mirror, I will not know.
And for some unknown reasons, some mirrors listed on mirrors.opensuse.org do has eg. 42.2 repos which are identified by MirrorBrain as none. For now, I
just take what it provides.

Future Plan: Add a config file so you can add your underground mirrors. Use a spider to visit the repo that are identified to have 0 repo, for more accurate result.
