# hardsub

This is an FFmpeg command generator and runner for burning subs into videos.  
Specifically made to convert mkv files with subtitles included into mp4 files with the subs burned into them.

The idea is to navigate to the folder where you have stored the mkv files and just simply run this command.  
The defaults are set to my own use case. Nearly everything is configurable, so it shouldn't be hard to customize
this to your needs.

I use this software several times a week, without problems. It does what it needs to do because I have a very
specific use case. So this might not fit your specific problem. _You are welcome to make a pull request_ or even
fork the project.

## Work in progress
This is very much a hobby project to scratch an itch. Contributions and suggestions are welcome,
feel free to open an issue.  
I do not pretend this is the best Go code ever written.  
The code was written to automate something boring. It could use some cleaning up, which I'm slowly doing.  
Some of the code is still very much not cleaned up. (especially main.go) _You are welcome to help out_.  
This project is "done" for me. Which means I don't need any more features or nicer CLI output. It does the job
and I'm happy with it.

## Currently working on or thinking about implementing these things:
- Adding functionality for this program to wait for files to be written to a folder and convert them when they are
fully uploaded.
- Searching for the intro part of the video and cutting it out of the hardsubbed file. This will require the user to
create a couple of screenshots so the tool knows what to look for.
- Simple notification support, so you can get notified when the video conversion is done.


## Requirements
You need ffmpeg, ffprobe and mkvtoolnix cli programs installed and on your $PATH.  
There are release versions in the releases page.

Or, you can build it yourself. You need Go installed.  
Clone the repo and run 'go install', this will install the tool into $HOME/go/bin/  (which should be in your $PATH)
