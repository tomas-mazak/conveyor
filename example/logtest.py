#!/usr/bin/python
# -*- coding: utf-8 -*-

import os
import sys
import time
import signal
import multiprocessing

LOGS = [
    ('access.log', 'Access', 0.5, 10, 0),
    ('error.log', 'Nasty error occured', 2, 30, 0),
    ('random.log', 'Some randomness', 1, 15, 10),
]


def log(args):
    """
    Write and rotate a single log file
    """
    filename, msg, log_period, rotation_period, initial_delay = args
    filename = os.path.join(os.environ.get('LOG_DIRECTORY', '/tmp/logs'), filename)
    time.sleep(initial_delay)
    while True:
        # Write
        with open(filename, "a") as f:
            for _ in range(int(rotation_period//log_period)):
                f.write("{} at {}\n".format(msg, time.time()))
                f.flush()
                time.sleep(log_period)
        # Rotate
        i = 1
        while os.path.exists(filename+".{}".format(i)):
            i += 1
        os.rename(filename, filename+".{}".format(i))


def initializer():
    """Ignore CTRL+C in the worker process."""
    signal.signal(signal.SIGINT, signal.SIG_IGN)


def main():
    try:
        pool = multiprocessing.Pool(processes=len(LOGS), initializer=initializer)
        pool.map_async(log, LOGS)
        # Wait for manual termination
        while True:
            time.sleep(3600)
    except KeyboardInterrupt:
        pool.terminate()
        pool.join()

if __name__ == '__main__':
    main()
