import time
import heapq
import math
from dataclasses import dataclass, field
from collections import deque, defaultdict, Counter
from typing import *

import numpy as np
import itertools
import re

import fire


def main(infile: str):
    print('hi!')

    pts = []
    with open(infile, 'r') as f:
        pt = None
        for line in f.readlines():
            if len(line.strip()) == 0:
                continue

            m = re.match(r'--- scanner ([0-9]*) ---', line.strip())
            if m:
                pt = list()
                pts.append(pt)
                continue
            else:
                pt.append([int(x) for x in line.strip().split(',')] + [1])

    scanners = []
    for i, pt in enumerate(pts):
        scanners.append(Scanner(i, np.array(pt)))

    t0 = time.time()
    s0 = scanners[0]
    dq = deque(scanners[1:])
    while len(dq) > 0:
        s1 = dq.popleft()
        found = False
        for rot in ROTS:
            s1.set_pos((0, 0, 0, 1))
            s1.A[:3, :3] = rot
            if check_overlap(s0, s1):
                s0.union(s1.transformed_pts())
                found = True
                break
        if not found:
            dq.append(s1)
        else:
            print(len(scanners) - len(dq), '/', len(scanners), s0.pts.shape)

    print('1:', s0.pts.shape, time.time() - t0, 's')

    t1 = time.time()
    max_l1 = 0
    for i, s0 in enumerate(scanners):
        for s1 in scanners[i+1:]:
            d = np.abs(s0.A[:, 3] - s1.A[:, 3]).sum()
            if d > max_l1:
                max_l1 = d

    print('2:', max_l1, time.time() - t1, 's')


def check_overlap(s0, s1):
    pts0 = s0.pts
    pts1 = s1.transformed_pts()

    dpts = list()
    for j in range(pts1.shape[0]):
        pt1 = pts1[j, :]
        dx = pts0 - pt1
        dx[:, 3] = 1
        dpts.append(dx)

    dptsar = np.concatenate(dpts, axis=0)
    unq, cnts = np.unique(dptsar, axis=0, return_counts=True)
    max_cnt, max_i = max((v, i) for i, v in enumerate(cnts))
    if max_cnt >= 12:
        s1.set_pos(unq[max_i, :])
        return True

    return False


@dataclass
class Scanner:
    n: int
    pts: np.ndarray
    A: np.ndarray = field(default_factory=lambda: np.eye(4, dtype=int))

    def transformed_pts(self):
        return self.pts.dot(self.A.transpose())

    def set_pos(self, pos: np.ndarray):
        self.A[:, 3] = pos

    def union(self, pts: np.ndarray):
        a = set(tuple(self.pts[i, :]) for i in range(self.pts.shape[0]))
        a.update(tuple(pts[j, :]) for j in range(pts.shape[0]))
        self.pts = np.array(list(a))


ROTS = []
for x, i in itertools.product(range(3), (-1, 1)):
    xv = np.zeros(3, dtype=int)
    xv[x] = i
    for y, k in itertools.product(range(3), (-1, 1)):
        if y == x:
            continue
        yv = np.zeros(3, dtype=int)
        yv[y] = k
        ROTS.append(
            np.stack([xv, yv, np.cross(xv, yv)],
                     axis=0)
        )


if __name__ == '__main__':
    fire.Fire(main)
