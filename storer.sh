#!/bin/bash
#-----------------------------------------------------------------------------#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.
#
#-----------------------------------------------------------------------------#
#-----------------------------------------------------------------------------#
# Script:   /etc/init.d/storer
#
# Brief:    System V init script for storer.
#
#-----------------------------------------------------------------------------#

PREFIX="storer: "
STORER="/usr/sbin/storer"
STORER_PROC="storer"
START_STOP_DAEMON="/sbin/start-stop-daemon"
# According to setting "server.graceful-shutdown-timeout" in lighttpd.conf
WAIT_FOR_STOP_TIMEOUT=60

# When logfile should be rotated we get SIGHUP signals,
# which can cause "start" or "stop" to be interrupted unexpectedly.
# This can also cause "reload" to not work as intended.
# This trap handles SIGHUP to work around this problem.
trap "echo \"[ INFO  ] Received signal SIGHUB\"" SIGHUP

storer_pid_count()
{
    local storer_pid_candidates
    local pid_cmd
    local pid_count=0

    storer_pid_candidates="$(pidof "${STORER_PROC}")"
    for pid in $storer_pid_candidates; do
        pid_cmd="$(cat "/proc/$pid/cmdline" 2>/dev/null | tr -d '\0')"
        if [[ $pid_cmd =~ ${STORER} ]]; then
            pid_count=$((pid_count + 1))
        fi
    done

    echo -n $pid_count
    return 0
}

status_check()
{
    local pid_count=0

    pid_count=$(storer_pid_count)
    if [[ $pid_count -eq 1 ]]; then
        return 0
    elif [[ $pid_count -lt 1 ]]; then
        # No storer process active
        return 1
    else
        # Multiple storer processes active => starting but not started, yet
        return 2
    fi
}

wait_for_start()
{
    local result=0

    echo -n "${PREFIX}starting"
    local try_count=0
    while ! status_check; do
        try_count=$((try_count+1))
        if [[ $try_count -gt 20 ]]; then
            result=1
            log_warning "server does not start within expected time"
            break
        fi
        echo -n "."
        usleep 100000
    done
    echo ""

    return $result
}

start()
{
    local result; result=1
    local try_count; try_count=3
    # START_STOP_DAEMON fails to start a process when SIGHUB (for log rotation) is received,
    # therefor we try again to start lighttpd
    while [[ result -ne 0 && try_count -gt 0 ]]; do
        try_count=$((try_count-1))
        $START_STOP_DAEMON --quiet --start --exec "${STORER}" 
        result=$?
    done
    if [[ $result -ne 0 ]]; then
        log_error "$START_STOP_DAEMON failed to start ${STORER}"
    else
        wait_for_start
        result=$?
    fi

    return $result
}

wait_for_stop()
{
    local result=0

    echo -n "${PREFIX}stopping"
    local try_count=0
    while $START_STOP_DAEMON --stop -t --quiet --exec "${STORER}"; do
        try_count=$((try_count+1))
        if [[ $try_count -gt ${WAIT_FOR_STOP_TIMEOUT} ]]; then
            result=1
            break
        fi
        echo -n "."
        sleep 1
    done
    echo ""

    return $result
}

stop()
{
    local result

    $START_STOP_DAEMON --quiet --stop --oknodo --exec "${STORER}"
    result=$?
    if [[ $result -eq 0 ]]; then
        wait_for_stop
        result=$?
    fi

    return $result
}

# Reload configuration of the service
reload()
{
    local result=1
    if $START_STOP_DAEMON --stop --signal SIGINT --oknodo --quiet --exec "${STORER}"; then
        wait_for_stop
        if start; then
            result=0
        # Startup may fail rarely, try once more in that case
        elif sleep 1 && start; then
            result=0
        fi
    fi

    return $result
}

# main
exec {lock_fd}>/var/lock/storer_init || exit 1
trap 'flock --unlock $lock_fd' EXIT
flock --exclusive $lock_fd

case $1 in

    start)
        log_info "start"
        if start; then
            log_info "start done"
            exit 0
        else
            log_error "could not start storer"
            exit 1
        fi
        ;;

    stop)
        log_info "stop"
        if stop; then
            log_info "stop done"
            exit 0
        else
            log_error "could not stop storer"
            exit 1
        fi
        ;;

    status)
        log_info "status check"
        if status_check; then
            log_info "status is running"
            exit 0
        else
            log_info "status is stopped"
            # This is not really an error case but for an automated status check
            # an other return code than 0 is used here
            exit 1
        fi
        ;;

    restart)
        log_info "restart"
        if stop; then
            log_info "stop done"
            if start; then
                log_info "start done"
                exit 0
            else
                log_error "could not start storer"
                exit 1
            fi
        else
            log_error "could not stop storer"
            exit 1
        fi
        ;;

    reload)
        log_info "reload"
        if reload; then
            log_info "reload done"
            exit 0
        else
            log_error "could not reload storer"
            exit 1
        fi
        ;;

    *)
        echo "Usage: ${0} [start|stop|status|restart|reload]" >&2
        exit 1
        ;;
esac