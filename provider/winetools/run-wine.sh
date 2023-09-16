#!/usr/bin/env bash

# Parameter
# 1. image name
# 2. vmid, should be 1-99
# 3. the relative path of the app/game's folder (relate to the mount dir)
# 4. the execute file in the app/game's folder
# 5. window title
# 6. "app"/"game"
# 7. screen width
# 8. screen height 
# 9. target host IP of the RTP protocal
# 10. frame
# 11. vcodec
# 12. use_gpu (none, 0/1/2, all)
# 13.(optional) optional parameters to run the app/game

# Notice
# 1. The supervisor is used inside the container as a daemon, and the supervisord.conf file is used to launch services
# 2. The actual entrance of the container is supervisord.conf
# 3. The system image output number specified is :99, and the same number should be used for video related operations
# 4. Docker command asks for sudo privilege

# Example
# sudo sh run-wine.sh dcwine container_id apppath appfile appname hwkey screenwidth screenheight targethost frame vcodec wineoptions
# sudo sh ./run-wine.sh dcwine 1 /apps/spider sol.exe spider game 480 320 127.0.0.1 30 h.264 

image_name=$1
container_id=$2
apppath=$3
appfile=$4
appname=$5
hwkey=$6
screenwidth=$7
screenheight=$8
targethost=$9
fps=${10}
vcodec=${11}
use_gpu=${12}
wineoptions=${13}

display=:99

if [ ${container_id} -lt 10 ];
then
    vmid_f="0${container_id}"
else
    vmid_f="${container_id}"
fi

conf_pre="$(pwd)"/../winetools/dockertools/containerfiles/supervisord_${vcodec}
if [[ $image_name == *nvidia* ]]; then 
  conf=${conf_pre}_nvidia.conf
else 
  conf=${conf_pre}.conf
fi

container_name="appvm${container_id}"
appdir_name="apps/point${container_id}"

if [ "$use_gpu" = "none" ]; then
    docker run -d --privileged --rm --name ${container_name} \
    -v /etc/localtime:/etc/localtime:ro \
    --mount type=bind,source="$(pwd)"/../winetools/"${appdir_name}",target=/apps \
    --mount type=bind,source=${conf},target=/etc/supervisor/conf.d/supervisord.conf  \
    --env "vmid=${container_id}" \
    --env "apppath=/apps/${apppath}"  \
    --env "appfile=${appfile}" \
    --env "appname=${appname}" \
    --env "hwkey=${hwkey}" \
    --env "screenwidth=${screenwidth}" \
    --env "screenheight=${screenheight}" \
    --env "wineoptions=${wineoptions}" \
    --env "targethost=${targethost}" \
    --env "videoport=1${vmid_f}05" \
    --env "audioport=1${vmid_f}01" \
    --env "inputport=1${vmid_f}09" \
    --env "fps=${fps}" \
    --env "DISPLAY=${display}" \
    --volume "winecfg:/root/.wine" ${image_name} 
    #supervisord
else
    docker run -d --privileged --rm --name ${container_name} \
    -v /etc/localtime:/etc/localtime:ro \
    --mount type=bind,source="$(pwd)"/../winetools/"${appdir_name}",target=/apps \
    --mount type=bind,source=${conf},target=/etc/supervisor/conf.d/supervisord.conf  \
    --gpus "${use_gpu}" \
    --env "vmid=${container_id}" \
    --env "apppath=/apps/${apppath}"  \
    --env "appfile=${appfile}" \
    --env "appname=${appname}" \
    --env "hwkey=${hwkey}" \
    --env "screenwidth=${screenwidth}" \
    --env "screenheight=${screenheight}" \
    --env "wineoptions=${wineoptions}" \
    --env "targethost=${targethost}" \
    --env "videoport=1${vmid_f}05" \
    --env "audioport=1${vmid_f}01" \
    --env "inputport=1${vmid_f}09" \
    --env "fps=${fps}" \
    --env "DISPLAY=${display}" \
    --volume "winecfg:/root/.wine" ${image_name} 
    #supervisord
fi

