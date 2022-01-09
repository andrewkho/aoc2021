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


def main(infile: str):
    print('hi!')

    lines = []
    with open(infile, 'r') as f:
        for i, line in enumerate(f.readlines()):
            lines.append([int(x) for x in line.strip()])

    A = np.array(lines)

    t0 = time.time()
    src = 0, 0
    dst = A.shape[0]-1, A.shape[1]-1
    dist, parents = shortest_path(A, src, dst)
    print(dist[dst], time.time() - t0, 's')

    t1 = time.time()
    rows = np.concatenate([
        np.where(A+i > 9, A+i-9, A+i)
        for i in range(5)
    ], axis=0)

    A5 = np.concatenate([
        np.where(rows+i > 9, rows+i-9, rows+i)
        for i in range(5)
    ], axis=1)

    assert (A5 > 9).sum() == 0
    print("A5:", time.time() - t1, 's')

    dst5 = A5.shape[0]-1, A5.shape[1]-1
    dist5, parents5 = shortest_path(A5, src, dst5)
    print(dist5[dst5], time.time() - t1, 's')


def shortest_path(A, src, dst):

    dist = np.ones(A.shape, dtype=int)*np.inf
    dist[src] = 0
    heap = [(dist[src], src)]
    parents = dict()
    while len(heap) > 0:
        v, pos = heapq.heappop(heap)
        if pos == dst:
            break

        i, j = pos
        for di, dj in [(-1, 0), (1, 0), (0, -1), (0, 1)]:
            if not (0 <= i+di < A.shape[0]):
                continue
            if not (0 <= j+dj < A.shape[1]):
                continue
            new_pos = i+di, j+dj
            if dist[i, j] + A[new_pos] < dist[new_pos]:
                parents[new_pos] = i, j
                dist[new_pos] = dist[i, j] + A[new_pos]
                heapq.heappush(heap, (dist[new_pos], new_pos))

    return dist, parents


if __name__ == '__main__':
    fire.Fire(main)
