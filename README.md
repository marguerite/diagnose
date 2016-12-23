## Diagnose tools for openSUSE ##

This project aims to collect the diagnose tools I wrote to help openSUSE freshmen to report issues on forum.
All written in ruby and can be natively ran on openSUSE.

it includes:

#### `instdpkg` AKA "installed packages" ####

Tell the user what packages he/she installed/removed on the specific `date`, since the `time`.

Case: 

A user posted on forum saying "I don't know what happened, I just installed some packages then problem occurs."

Tell him to download and run `instdpkg.rb -d=date -t=time`. it will tell you what "some packages" are.

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

Tell him to doownload and run `pkmswitch100 --all`, problem solved.


