#!/bin/bash
read -p "Update the live-streamer (live-streamer on pi) y/n ? " -n 1 -r
echo
if [[ $REPLY =~ ^[Nn]$ ]]
then
	echo "Update canceled."
	exit 0

fi
echo "Stop the service"
sudo systemctl stop live-streamer

ZIPDIR="./zips"
CURRDIR="./current"
OLDDIR="./old"

echo "Now starting the process..."

# Make sure dir exits
[ ! -d "$ZIPDIR" ] && mkdir -p "$ZIPDIR"
[ ! -d "$CURRDIR" ] && mkdir -p "$CURRDIR" 
[ ! -d "$OLDDIR" ] && mkdir -p "$OLDDIR" 

# backup the current dir
bckdir=$OLDDIR'/'"$(date +"%Y-%m-%d-%H%M%S")"
echo "Backup dir is: $bckdir"
[ ! -d "$bckdir" ] && mkdir -p "$bckdir" 

mv $CURRDIR'/'*  $bckdir

#zips=$(ls $ZIPDIR)
#echo "$zips"

for file in $ZIPDIR/*.zip ; do 
	fname=$(basename "$file")
	#echo "Name is $fname"
done

echo "Want to unzip $fname"
zippath=$ZIPDIR'/'$fname
destpath=$CURRDIR'/'
echo "The source is $zippath and destination is $destpath"
unzip $zippath -d $destpath

chmod +x $destpath'/'live-streamer.bin

echo "Start the service"
sudo systemctl start live-streamer


echo "That's all folks!"
