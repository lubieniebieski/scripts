#!/bin/bash

# Init variables
isSetM=0
isSetD=0
isSetE=0
isSetCount=0
counter=0

#Display usage info
usage() {

    cat <<EOF

Usage: dng-jpg.sh [-m <path>] [-d <path>] [-e <path>] [-h]

-m: for move   (moves files to <path>/duplicates)
-d: for delete (deletes duplicate files)
-e: for echo   (lists duplicate files)
-h: for help

EOF
  exit 1
}

#Check for parameters
while getopts ":m:d:e:h" opt; do
  case ${opt} in
    m)
        isSetM=1
        let isSetCount="$isSetCount+1"
        arg=${OPTARG}
      echo "Move selected with path:" $arg
      ;;
    d)
        isSetD=1
        let isSetCount="$isSetCount+1"
        arg=${OPTARG}
      echo "Delete selected with path:" $arg
      ;;
    e)
        isSetE=1
        let isSetCount="$isSetCount+1"
        arg=${OPTARG}
      echo "Echo selected with path:" $arg
      ;;
    h)
        let isSetCount="$isSetCount+1"
        usage
      ;;
    \?)
      echo "Invalid option: -$OPTARG" >&2
      usage
      ;;
    :)
      echo "Option -$OPTARG requires a directory argument." >&2
      usage
      ;;
    *)
      usage
      ;;
  esac
done

# If no parameters, show usage help and exit
if test -z "$1"; then
    usage
fi

# If multiple parameters (not counting -a), show usage help and exit
if (($isSetCount > 1)); then
    usage
fi

#Verify directory
if [ ! -d "${arg}" ]; then
  echo "$arg is not a path to a directory." >&2
  usage
fi

#Now set it as a basedir
BASEDIR=$arg
WASTEDIR="$BASEDIR/duplicates/"
if (( $isSetM==1 )); then
    mkdir $WASTEDIR
fi

find "${BASEDIR}" -name '*.RAF' | while read LINE; do
  prefix=${LINE%.RAF}
  if [ -e "$prefix.JPG" ]; then
      let counter="$counter+1"
      if (( $isSetE==1 )); then
          echo "$prefix.JPG"
      fi
      if (( $isSetM==1 )); then
          mv $prefix.JPG $WASTEDIR
      fi
      if (( $isSetD==1 )); then
          rm "${prefix}.JPG"
      fi
  fi

done
echo "$counter files found."
