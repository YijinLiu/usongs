udisk-songs
===========
This is a tiny tool to copy songs from hard drive to udisk.
I randomly copies mp3 files and writes a checkpoint file after it's done.
Next time it runs, it reads checkpoint file and copy different mp3 files.
(You'll need to delete the old files from your udisk manually.)
Manual
<pre>
$ cd go
$ GOPATH=`pwd` go install copy_mp3
$ ./bin/copy_mp3 -src-dir=XXX -dst-dir=YYY