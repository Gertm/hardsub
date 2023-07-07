# hardsub

This is an FFmpeg command generator and runner for burning subs into videos.  
Specifically made to convert mkv files with subtitles included into mp4 files with the subs burned into them.

The idea is to navigate to the folder where you have stored the mkv files and just simply run this command.  
The defaults are set to my own use case. Nearly everything is configurable, so it shouldn't be hard to customize
this to your needs.

## Work in progress
This is very much a hobby project to scratch an itch. Contributions and suggestions are welcome,
feel free to open an issue.  
I do not pretend this is the best Go code ever written.  
The code was written to automate something boring. It could use some cleaning up, which I'm slowly doing.

## Requirements
You need ffmpeg, ffprobe and mkvtoolnix cli programs installed and on your $PATH.  

I could add goreleaser at some point to get easy release versions of this.

For now, you'll have to build it yourself. You need Go installed.  
Clone the repo and run 'go install', this will install the tool into $HOME/go/bin/  (which should be in your $PATH)
