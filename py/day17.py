import heapq
import math
from dataclasses import dataclass
from collections import deque, defaultdict, Counter
from typing import *
import time

import numpy as np
import itertools
import re

import fire


def main():
    print('hi!')

    xlim = (257,  286)
    ylim = (-101, -57)

    def hit(pos):
        return (xlim[0] <= pos[0] <= xlim[1]) and (ylim[0] <= pos[1] <= ylim[1])

    def step(pos, vel):
        pos[0] += vel[0]
        pos[1] += vel[1]

        if vel[0] > 0:
            vel[0] -= 1
        elif vel[0] < 0:
            vel[0] += 1
        else:
            vel[0] = 0
        vel[1] -= 1

    def trial(vel):
        pos = [0, 0]
        max_y = pos[1]
        while pos[1] >= ylim[0] and pos[0] <= xlim[1]:
            if pos[1] > max_y:
                max_y = pos[1]
            if hit(pos):
                return pos, max_y
            step(pos, vel)

        return None, max_y

    t0 = time.time()
    max_max_y = 0
    for u in range(1, xlim[1]):
        for v in range(-ylim[0]):
            p, max_y = trial([u, v])
            if p:
                if max_y > max_max_y:
                    max_max_y = max_y

    print('1:', max_max_y, time.time() - t0, 's')

    t1 = time.time()
    hits = 0
    for u in range(1, xlim[1]+1):
        for v in range(ylim[0], -ylim[0]):
            p, max_y = trial([u, v])
            if p:
                hits += 1

    print('2:', hits, time.time() - t1, 's')


if __name__ == '__main__':
    fire.Fire(main)
