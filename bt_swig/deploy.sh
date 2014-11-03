if [ "$WORK" = "" ]
then
	echo "please run the following:"
	echo ""
	echo "go install -a -work github.com/slimsag/bt_swig"
	echo ""
	echo "Copy the directory printed into this command:"
	echo "export WORK=insert_work_dir_here"
	echo ""
	echo "And re-run this script."
	exit 0
else
	echo "Note: if you want to rebuild the project please run"
	echo "      go install -a -work github.com/slimsag/bt_swig"
	echo "      export WORK=insert_work_dir_here"
fi

DST=$GOPATH/src/github.com/slimsag/bt
rm -rf $DST
mkdir $DST
cp $WORK/github.com/slimsag/bt_swig/_obj/bt_gc.c $DST
cp $WORK/github.com/slimsag/bt_swig/_obj/bt.go $DST
cp $WORK/github.com/slimsag/bt_swig/_obj/bt_wrap.cxx $DST
echo "Deployed to $DST"

