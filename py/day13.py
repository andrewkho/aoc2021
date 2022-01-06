import math
from dataclasses import dataclass
from collections import deque, defaultdict
from typing import *

import numpy as np
import itertools
import re
import time

import fire


def main(infile: str):
    print('hi!')

    points = list()
    folds = list()
    with open(infile, 'r') as f:
        part2 = False
        for line in f.readlines():
            if len(line.strip()) == 0:
                part2 = True
                continue
            if not part2:
                l, r = line.strip().split(',')
                points.append((int(r), int(l)))
            else:
                _, _, x = line.strip().split(' ')
                l, r = x.strip().split('=')
                folds.append((l, int(r)))

    max_x = max(x for y, x in points) + 1
    max_y = max(y for y, x in points) + 1
    grid = np.zeros((max_y, max_x), dtype=np.int32)

    for pt in points:
        grid[pt] += 1

    t0 = time.time()
    grid = fold(grid, folds[0])

    print(f'1: {(grid > 0).sum()}', time.time() - t0, 's')

    t1 = time.time()
    for fld in folds[1:]:
        grid = fold(grid, fld)

    for i in range(grid.shape[0]):
        l = []
        for j in range(grid.shape[1]):
            l.append('#' if grid[i, j] > 0 else ' ')
        print(''.join(l))

    print(f'2:', time.time() - t1, 's')

def fold(grid, instr):
    axis, loc = instr
    if axis == 'x':
        for di in range(grid.shape[1] - loc):
            grid[:, loc-di] += grid[:, loc+di]
        return grid[:, :loc]
    else:
        for dj in range(grid.shape[0] - loc):
            grid[loc-dj, :] += grid[loc+dj, :]
        return grid[:loc, :]


if __name__ == '__main__':
    fire.Fire(main)
