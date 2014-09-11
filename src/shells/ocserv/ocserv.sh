#!/bin/sh
### BEGIN INIT INFO
# Provides:          ocserv
# Required-Start:    $remote_fs $syslog
# Required-Stop:     $remote_fs $syslog
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
### END INIT INFO
# Copyright Rene Mayrhofer, Gibraltar, 1999
# This script is distibuted under the GPL

PATH=/bin:/usr/bin:/sbin:/usr/sbin:/opt/bin
DAEMON=/opt/bin/ocserv
PIDFILE=/var/run/ocserv.pid
DAEMON_ARGS="-c /etc/ocserv/ocserv.conf"

case "$1" in
start)
if [ ! -r $PIDFILE ]; then
echo -n "Starting OpenConnect VPN Server Daemon: "
start-stop-daemon --start --quiet --pidfile $PIDFILE --exec $DAEMON -- \
$DAEMON_ARGS > /dev/null
echo "ocserv."
else
echo -n "OpenConnect VPN Server is already running.\n\r"
exit 0
fi
;;
stop)
echo -n "Stopping OpenConnect VPN Server Daemon: "
start-stop-daemon --stop --quiet --pidfile $PIDFILE --exec $DAEMON
echo "ocserv."
rm -f $PIDFILE
;;
force-reload|restart)
echo "Restarting OpenConnect VPN Server: "
$0 stop
sleep 1
$0 start
;;
status)
if [ ! -r $PIDFILE ]; then
# no pid file, process doesn't seem to be running correctly
exit 3
fi
PID=`cat $PIDFILE | sed 's/ //g'`
EXE=/proc/$PID/exe
if [ -x "$EXE" ] &&
[ "`ls -l \"$EXE\" | cut -d'>' -f2,2 | cut -d' ' -f2,2`" = \
"$DAEMON" ]; then
# ok, process seems to be running
exit 0
elif [ -r $PIDFILE ]; then
# process not running, but pidfile exists
exit 1
else
# no lock file to check for, so simply return the stopped status
exit 3
fi
;;
*)
echo "Usage: /etc/init.d/ocserv {start|stop|restart|force-reload|status}"
exit 1
;;
esac

exit 0