## `mirrorswitch`

`mirrorswitch` is used to batchly change your zypper repositories.

eg:

1. if you want to switch all repositories from "openSUSE_Leap_15.2" flavor to "openSUSE_Tumbleweed"

    sudo mirrorswitch -from="openSUSE_Leap_15.2" -to="openSUSE_Tumbleweed"

2. if you want to switch packman repository to its mirror

    sudo mirrorswitch -repo="packman" -from="packman.inode.at" -to="mirrors.hust.edu.cn/packman"
