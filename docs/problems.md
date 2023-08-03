# 1 
Problem:
    after chroot and running the first run all goes well. but when trying to do another run it can't cause they can't find the rootfs to chroot to. after some debugging it turns out that exiting root is not work well. and we stuck in the chroot jail
Suggestions:    
    1. we can open a file descriptor for the system root and save it so then after we done we can break the jail.
    2. don't enter the jail in bash `chroot /path /executable` command can have extra args that can run inside the path without chroot actually. this will be useful in the builder part not in the mounter. cause in builder we actually don't need to mount other directories.